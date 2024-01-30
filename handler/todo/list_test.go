package todo_test

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/senomas/todo_app/handler/todo"
	"github.com/senomas/todo_app/store"
	_ "github.com/senomas/todo_app/store/sql_tmpl"
	_ "github.com/senomas/todo_app/store/sql_tmpl/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestListHandler(t *testing.T) {
	t.Log("TestListHandler")
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(log)

	todoStore := store.GetTodoStore()
	assert.NotNil(t, todoStore, "todo store should not be nil")

	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err, "sql open should not error")
	defer db.Close()

	initCtx := func() (context.Context, context.CancelFunc) {
		return context.WithTimeout(
			context.WithValue(context.Background(), store.StoreCtxDB, db),
			time.Millisecond*200,
		)
	}

	t.Run("Init", func(t *testing.T) {
		ctx, cancel := initCtx()
		defer cancel()

		c, err := store.Migrate(ctx)
		assert.NoError(t, err, "migrate should not error")
		assert.EqualValues(t, 1, c, "migrate rows")

		err = todoStore.Init(ctx)
		assert.NoError(t, err, "todo store init should not error")
	})

	t.Run("Init data", func(t *testing.T) {
		ctx, cancel := initCtx()
		defer cancel()

		for i := 1; i <= 4; i++ {
			todo, err := todoStore.CreateTodo(ctx, "Todo "+strconv.Itoa(i))
			assert.NoError(t, err, "todo store create todo should not error")
			assert.EqualValues(t, todo, &store.Todo{ID: int64(i), Title: "Todo " + strconv.Itoa(i), Completed: false}, "todo should be equal")
		}
	})

	t.Run("List", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/todo", nil)
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(
			context.WithValue(context.Background(), store.StoreCtxDB, db),
		)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todo.ListHandler)

		handler.ServeHTTP(rr, req)

		assert.EqualValuesf(t, http.StatusOK, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)

		body := rr.Body.String()
		ebody := `<h1>TODO List</h1><ul>`
		for i := 1; i <= 4; i++ {
			ebody += fmt.Sprintf("<li><a href=\"/todo/%d\">Todo %d</a></li>", i, i)
		}
		ebody += `</ul>`
		assert.EqualValues(t, ebody, body, "handler returned unexpected body")
	})

	t.Run("List title like", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/todo?title.like="+url.QueryEscape("%3%"), nil)
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(
			context.WithValue(context.Background(), store.StoreCtxDB, db),
		)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todo.ListHandler)

		handler.ServeHTTP(rr, req)

		assert.EqualValuesf(t, http.StatusOK, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)

		body := rr.Body.String()
		ebody := `<h1>TODO List</h1><ul>`
		ebody += "<li><a href=\"/todo/3\">Todo 3</a></li>"
		ebody += `</ul>`
		assert.EqualValues(t, ebody, body, "handler returned unexpected body")
	})
}
