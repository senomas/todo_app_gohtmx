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
	return `INSERT INTO todo (title, completed) VALUES ($1, false)`, []any{t.Title}
}

// UpdateTodo implements sql_tmpl.TodoStoreTemplate.
func (s *TodoStoreTemplateImpl) UpdateTodo(t *store.Todo) (string, []any) {
	return `UPDATE todo SET title = $2, completed = $3 WHERE id = $1`, []any{t.ID, t.Title, t.Completed}
}

// DeleteTodoByID implements sql_tmpl.TodoStoreTemplate.
func (s *TodoStoreTemplateImpl) DeleteTodoByID(id any) (string, []any) {
	return `DELETE FROM todo WHERE id=?`, []any{id}
}

// GetTodoByID implements sql_tmpl.TodoStoreTemplate.
func (s *TodoStoreTemplateImpl) GetTodoByID(id any) (string, []any) {
	return `SELECT id, title, completed FROM todo WHERE id = $1`, []any{id}
}

func filterToString(where []string, args []any, field string, filter any) ([]string, []any) {
	switch f := filter.(type) {
	case store.FilterInt64:
		switch f.Op {
		case store.OP_NOP:
		case store.OP_EQ:
			args = append(args, f.Value)
			where = append(where, fmt.Sprintf("%s = $%d", field, len(args)))
		default:
			panic(fmt.Sprintf("unknown filter int64 op: %+v", f))
		}
	case store.FilterBool:
		switch f.Op {
		case store.OP_NOP:
		case store.OP_EQ:
			args = append(args, f.Value)
			where = append(where, fmt.Sprintf("%s = $%d", field, len(args)))
		default:
			panic(fmt.Sprintf("unknown filter bool op: %+v", f))
		}
	case store.FilterString:
		switch f.Op {
		case store.OP_NOP:
		case store.OP_EQ:
			args = append(args, f.Value)
			where = append(where, fmt.Sprintf("%s = $%d", field, len(args)))
		case store.OP_LIKE:
			args = append(args, f.Value)
			where = append(where, fmt.Sprintf("%s LIKE $%d", field, len(args)))
		default:
			panic(fmt.Sprintf("unknown filter string op: %+v", f))
		}
	default:
		panic(fmt.Sprintf("unknown filter type: %+v", f))
	}
	return where, args
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
func (*TodoStoreTemplateImpl) FindTodo(store.TodoFilter, int64, int) (string, []any) {
	panic("unimplemented")
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
