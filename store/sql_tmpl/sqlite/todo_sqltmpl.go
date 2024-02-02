package sqlite

import (
	"fmt"
	"strings"

	"github.com/senomas/todo_app/store"
	"github.com/senomas/todo_app/store/sql_tmpl"
)

type TodoStoreTemplateImpl struct{}

func init() {
	sql_tmpl.SetupTodoStoreTemplate(&TodoStoreTemplateImpl{})
}

// InsertTodo implements sql_tmpl.TodoStoreTemplate.
func (s *TodoStoreTemplateImpl) InsertTodo(t *store.Todo) (string, []any) {
	return `INSERT INTO todo (title, completed) VALUES ($1, $2)`, []any{t.Title, t.Completed}
}

// UpdateTodo implements sql_tmpl.TodoStoreTemplate.
func (s *TodoStoreTemplateImpl) UpdateTodo(t *store.Todo) (string, []any) {
	return `UPDATE todo SET title = $1, completed = $2 WHERE id = $3`, []any{t.Title, t.Completed, t.ID}
}

// DeleteTodoByID implements sql_tmpl.TodoStoreTemplate.
func (s *TodoStoreTemplateImpl) DeleteTodoByID(id any) (string, []any) {
	return `DELETE FROM todo WHERE id = $1`, []any{id}
}

// GetTodoByID implements sql_tmpl.TodoStoreTemplate.
func (s *TodoStoreTemplateImpl) GetTodoByID(id any) (string, []any) {
	return `SELECT id, title, completed FROM todo WHERE id = $1`, []any{id}
}

func (s *TodoStoreTemplateImpl) findTodoWhere(filter store.TodoFilter) ([]string, []any) {
	where := []string{}
	args := []any{}

	where, args = filterToString(where, args, "id", filter.ID)
	where, args = filterToString(where, args, "title", filter.Title)
	where, args = filterToString(where, args, "completed", filter.Completed)

	return where, args
}

// FindTodo implements sql_tmpl.TodoStoreTemplate.
func (s *TodoStoreTemplateImpl) FindTodo(filter store.TodoFilter, skip int64, limit int) (string, []any) {
	where, args := s.findTodoWhere(filter)
	sl := ""
	if limit > 0 {
		sl += fmt.Sprintf(" LIMIT %d", limit)
	} else {
		sl += " LIMIT 1000"
	}
	if skip > 0 {
		sl += fmt.Sprintf(" OFFSET %d", skip)
	}
	if len(where) > 0 {
		return `SELECT id, title, completed FROM todo WHERE ` + strings.Join(where, " AND ") + sl, args
	}
	return `SELECT id, title, completed FROM todo` + sl, args
}

// FindTodoTotal implements sql_tmpl.TodoStoreTemplate.
func (s *TodoStoreTemplateImpl) FindTodoTotal(filter store.TodoFilter) (string, []any) {
	where, args := s.findTodoWhere(filter)
	if len(where) > 0 {
		return `SELECT COUNT(*) FROM todo WHERE ` + strings.Join(where, " AND "), args
	}
	return `SELECT COUNT(*) FROM todo`, args
}

// ErrorMapFind implements sql_tmpl.TodoStoreTemplate.
func (*TodoStoreTemplateImpl) ErrorMapFind(err error) error {
	if err.Error() == "sql: no rows in result set" {
		return store.ErrNoData
	}
	return err
}
