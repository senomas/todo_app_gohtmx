package store_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"strconv"
	"testing"
	"time"

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

	t.Run("GetTodoByID", func(t *testing.T) {
		ctx, cancel := initCtx()
		defer cancel()

		todo, err := todoStore.GetTodoByID(ctx, 1)
		assert.NoError(t, err, "todo store get todo by id should not error")
		assert.EqualValues(t, todo, &store.Todo{ID: 1, Title: "Todo 1", Completed: false}, "todo should be equal")
	})

	t.Run("GetTodoByID not-found", func(t *testing.T) {
		ctx, cancel := initCtx()
		defer cancel()

		_, err := todoStore.GetTodoByID(ctx, 2)
		assert.ErrorContains(t, err, "sql: no data")
	})

	t.Run("CreateTodo 2", func(t *testing.T) {
		ctx, cancel := initCtx()
		defer cancel()

		for i := 2; i <= 10; i++ {
			todo, err := todoStore.CreateTodo(ctx, "Todo "+strconv.Itoa(i))
			assert.NoError(t, err, "todo store create todo should not error")
			assert.EqualValues(t, todo, &store.Todo{ID: int64(i), Title: "Todo " + strconv.Itoa(i), Completed: false}, "todo should be equal")
		}
	})

	t.Run("FindTodo", func(t *testing.T) {
		t.Skip("TODO: Fix this test")
		ctx, cancel := initCtx()
		defer cancel()

		todos, total, err := todoStore.FindTodo(ctx, store.TodoFilter{}, 1, 0)
		assert.NoError(t, err, "todo store find todo should not error")
		assert.EqualValues(t, total, 10, "total should be equal")
		assert.EqualValues(t, todos, []*store.Todo{
			{ID: 1, Title: "Todo 1", Completed: false},
			{ID: 2, Title: "Todo 2", Completed: false},
			{ID: 3, Title: "Todo 3", Completed: false},
			{ID: 4, Title: "Todo 4", Completed: false},
		}, "todos should be equal")
	})
}
