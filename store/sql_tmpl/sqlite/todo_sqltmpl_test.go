package sqlite_test

import (
	"testing"

	"github.com/senomas/todo_app/store"
	"github.com/senomas/todo_app/store/sql_tmpl"
	"github.com/senomas/todo_app/store/sql_tmpl/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestTodoSqlTemplateSqlite(t *testing.T) {
	var tl sql_tmpl.TodoStoreTemplate = &sqlite.TodoStoreTemplateImpl{}
	assert.NotNil(t, tl, "TodoStoreTemplateImpl not nil")

	t.Run("FindTodoTotal no filter", func(t *testing.T) {
		qry, args := tl.FindTodoTotal(store.TodoFilter{})
		assert.Equal(t, `SELECT COUNT(*) FROM todo`, qry)
		assert.EqualValues(t, len(args), 0, "no args")
	})

	t.Run("FindTodoTotal where id =", func(t *testing.T) {
		qry, args := tl.FindTodoTotal(store.TodoFilter{ID: store.FilterInt64{Value: 100, Op: store.OP_EQ}})
		assert.Equal(t, `SELECT COUNT(*) FROM todo WHERE id = $1`, qry)
		assert.EqualValues(t, args, []any{int64(100)}, "1 args")
	})

	t.Run("FindTodoTotal where title =", func(t *testing.T) {
		qry, args := tl.FindTodoTotal(store.TodoFilter{Title: store.FilterString{Value: "foo", Op: store.OP_EQ}})
		assert.Equal(t, `SELECT COUNT(*) FROM todo WHERE title = $1`, qry)
		assert.EqualValues(t, args, []any{"foo"}, "1 args")
	})

	t.Run("FindTodoTotal where title like", func(t *testing.T) {
		qry, args := tl.FindTodoTotal(store.TodoFilter{Title: store.FilterString{Value: "foo", Op: store.OP_LIKE}})
		assert.Equal(t, `SELECT COUNT(*) FROM todo WHERE title LIKE $1`, qry)
		assert.EqualValues(t, args, []any{"foo"}, "1 args")
	})
}
