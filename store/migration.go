package store

import (
	"bufio"
	"context"
	"crypto/sha512"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Migration struct {
	Timestamp time.Time
	Filename  string
	Hash      string
	Result    string
	ID        int64
	Success   bool
}

type Migrations interface {
	Init(context.Context) error
	GetMigration(ctx context.Context, filename string) (*Migration, error)
	AddMigration(ctx context.Context, m *Migration) error
}

var migrationsImpl Migrations

func SetupMigrationsImpl(m Migrations) {
	migrationsImpl = m
}

func calcHash(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("error reading %s: %v", filename, err)
	}
	defer f.Close()
	hasher := sha512.New()
	p := make([]byte, 1024)
	for {
		n, err := f.Read(p)
		hasher.Write(p[:n])
		if err == io.EOF {
			return fmt.Sprintf("%x", hasher.Sum(nil)), nil
		} else if err != nil {
			return "", err
		}
	}
}

func Migrate(ctx context.Context) (int, error) {
	if db, ok := ctx.Value(StoreCtxDB).(*sql.DB); ok {
		migrationsImpl.Init(ctx)
		mpath := os.Getenv("MIGRATIONS_PATH")
		if mpath == "" {
			ex, err := os.Executable()
			if err != nil {
				panic(err)
			}
			path := filepath.Dir(ex)
			mpath = filepath.Join(path, "migrations")
		}
		slog.Debug("MIGRATIONS_PATH", "path", mpath)
		files, err := os.ReadDir(mpath)
		if err != nil {
			return -1, err
		}

		migrate := func(db *sql.DB, qry string) error {
			_, err := db.ExecContext(ctx, qry)
			return err
		}
		count := 0
		for _, file := range files {
			mig, err := migrationsImpl.GetMigration(ctx, file.Name())
			if err != nil {
				return -1, err
			}

			hash, err := calcHash(path.Join(mpath, file.Name()))
			if err != nil {
				return -1, err
			}

			if mig != nil && mig.Success && mig.Hash == hash {
				fmt.Printf("SKIP file: %s\n", file.Name())
				continue
			}
			fmt.Printf("file: %s - %s %v\n", file.Name(), hash, mig)

			m := &Migration{
				Filename: file.Name(),
				Hash:     hash,
				Result:   "",
				Success:  false,
			}
			defer func() {
				str, _ := json.MarshalIndent(m, "", "  ")
				fmt.Printf("MIGRATE %s\n", string(str))
			}()
			f, err := os.Open(path.Join(mpath, file.Name()))
			if err != nil {
				return -1, fmt.Errorf("error reading %s: %v", file.Name(), err)
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			qry := ""
			for scanner.Scan() {
				ln := scanner.Text()
				qry = fmt.Sprintf("%s%s\n", qry, ln)
				if strings.HasSuffix(strings.TrimSpace(ln), ";") {
					err := migrate(db, qry)
					if err != nil {
						m.Result = fmt.Sprintf("%s%s\nERROR: %v\n", m.Result, qry, err)
						return -1, fmt.Errorf("error migrating %s: [%s]\n%v", file.Name(), qry, err)
					} else {
						m.Result = fmt.Sprintf("%s%s\n", m.Result, qry)
						count++
					}
					qry = ""
				}
			}
			if strings.TrimSpace(qry) != "" {
				err := migrate(db, qry)
				if err != nil {
					return -1, fmt.Errorf("error migrating %s: [%s]\n%v", file.Name(), qry, err)
				}
			}
			m.Success = true
			m.Timestamp = time.Now()
			err = migrationsImpl.AddMigration(ctx, m)
			if err != nil {
				return -1, err
			}
		}
		return count, nil
	}
	return -1, errors.New("context not have db")
}
