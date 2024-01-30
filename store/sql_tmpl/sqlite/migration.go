package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/senomas/todo_app/store"
)

var errCtxNoDB = errors.New("no db defined in context")

func init() {
	slog.Debug("Register sql_tmpl.Migrations")
	store.SetupMigrationsImpl(&MigrationsImpl{})
}

type MigrationsImpl struct{}

func (m *MigrationsImpl) Init(ctx context.Context) error {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sql.DB); ok {
		qry := `
      CREATE TABLE IF NOT EXISTS _migration (
        id        INTEGER PRIMARY KEY AUTOINCREMENT,
        filename  TEXT,
        hash      TEXT,
        success   BOOLEAN,
        result    TEXT,
        timestamp DATETIME
      )
    `
		_, err := db.ExecContext(ctx, qry)
		if err != nil {
			slog.Warn("Error insert todo", "qry", qry, "error", err)
			return err
		}
		return err
	}
	return errCtxNoDB
}

func (i *MigrationsImpl) GetMigration(ctx context.Context, filename string) (*store.Migration, error) {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sql.DB); ok {
		qry := `
      SELECT id, filename, hash, success, result, timestamp FROM _migration WHERE filename = $1 
      ORDER BY timestamp DESC LIMIT 1`
		rs, err := db.QueryContext(ctx, qry, filename)
		if err != nil {
			return nil, err
		}
		defer rs.Close()
		if rs.Next() {
			var m store.Migration
			err := rs.Scan(&m.ID, &m.Filename, &m.Hash, &m.Success, &m.Result, &m.Timestamp)
			return &m, err
		}
		return nil, nil
	}
	return nil, errCtxNoDB
}

func (i *MigrationsImpl) AddMigration(ctx context.Context, m *store.Migration) error {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sql.DB); ok {
		qry := `
      INSERT INTO _migration (filename, hash, success, result, timestamp)
      VALUES ($1, $2, $3, $4, $5)
    `
		rs, err := db.ExecContext(ctx, qry, m.Filename, m.Hash, m.Success, m.Result, m.Timestamp)
		if err != nil {
			return err
		}
		m.ID, err = rs.LastInsertId()
		return err
	}
	return errCtxNoDB
}
