package store_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/senomas/todo_app/store"
	_ "github.com/senomas/todo_app/store/sql_tmpl"
	_ "github.com/senomas/todo_app/store/sql_tmpl/sqlite"
	"github.com/stretchr/testify/assert"
)

func mustEncode(t *testing.T, v interface{}) string {
	jv, err := json.MarshalIndent(v, "", "  ")
	assert.NoError(t, err, "encode json should not error")
	return string(jv)
}

func TestTodoStore(t *testing.T) {
	t.Log("Test TodoStore")
	todoStore := store.GetTodoStore()
	assert.NotNil(t, todoStore, "todo store should not be nil")

	os.Setenv("MIGRATIONS_PATH", "../migrations")

	db, err := sqlx.Open("sqlite3", ":memory:")
	assert.NoError(t, err, "sqlx open should not error")
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

	t.Run("Init 2", func(t *testing.T) {
		ctx, cancel := initCtx()
		defer cancel()

		c, err := store.Migrate(ctx)
		assert.NoError(t, err, "migrate should not error")
		assert.EqualValues(t, 0, c, "migrate rows")
	})

	t.Run("CreateTodo", func(t *testing.T) {
		ctx, cancel := initCtx()
		defer cancel()

		todo, err := todoStore.CreateTodo(ctx, "Todo 1")
		assert.NoError(t, err, "todo store create todo should not error")
		assert.EqualValues(t, todo, &store.Todo{ID: 1, Title: "Todo 1", Completed: false}, "todo should be equal")
	})
}
