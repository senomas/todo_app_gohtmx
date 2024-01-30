package sql_tmpl

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/senomas/todo_app/store"
)

var (
	errCtxNoDB     = errors.New("no db defined in context")
	errNoStoreTmpl = errors.New("todo store template not initialized")
)

type TodoStoreTmpl struct{}

type TodoStoreTemplate interface {
	InsertTodo(t *store.Todo) (string, []any)
	UpdateTodo(t *store.Todo) (string, []any)
	DeleteTodoByID(id any) (string, []any)

	GetTodoByID(id any) (string, []any)
	FindTodo(store.TodoFilter, int64, int) (string, []any)
	FindTodoTotal(store.TodoFilter) (string, []any)

	ErrorMapFind(error) error
}

func init() {
	slog.Debug("Register sql_tmpl.TodoStore")
	store.SetupTodoStoreImplementation(&TodoStoreTmpl{})
}

var todoStoreTemplateImpl TodoStoreTemplate

func SetupTodoStoreTemplate(t TodoStoreTemplate) {
	todoStoreTemplateImpl = t
}

func (t *TodoStoreTmpl) Init(ctx context.Context) error {
	if todoStoreTemplateImpl == nil {
		return errNoStoreTmpl
	}
	return nil
}

// CreateTodo implements store.TodoStore.
func (t *TodoStoreTmpl) CreateTodo(ctx context.Context, title string) (*store.Todo, error) {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sql.DB); ok {
		todo := store.Todo{Title: title}
		qry, args := todoStoreTemplateImpl.InsertTodo(&todo)
		slog.Debug("CreateTodo", "qry", qry, "args", &store.JsonLogValue{V: args})
		rs, err := db.ExecContext(ctx, qry, args...)
		if err != nil {
			slog.Warn("Error insert todo", "qry", qry, "error", err)
			return nil, err
		}
		todo.ID, err = rs.LastInsertId()
		return &todo, err
	}
	return nil, errCtxNoDB
}

// UpdateTodo implements store.TodoStore.
func (t *TodoStoreTmpl) UpdateTodo(ctx context.Context, todo store.Todo) error {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sql.DB); ok {
		qry, args := todoStoreTemplateImpl.UpdateTodo(&todo)
		slog.Debug("UpdateTodo", "qry", qry, "args", &store.JsonLogValue{V: args})
		db.ExecContext(ctx, qry, args...)
	}
	return errCtxNoDB
}

// DeleteTodoByID implements store.TodoStore.
func (t *TodoStoreTmpl) DeleteTodoByID(ctx context.Context, id int64) error {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sql.DB); ok {
		qry, args := todoStoreTemplateImpl.DeleteTodoByID(id)
		slog.Debug("DeleteTodoByID", "qry", qry, "args", &store.JsonLogValue{V: args})
		db.ExecContext(ctx, qry, args...)
	}
	return errCtxNoDB
}

// GetTodoByID implements store.TodoStore.
func (t *TodoStoreTmpl) GetTodoByID(ctx context.Context, id int64) (*store.Todo, error) {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sql.DB); ok {
		todo := store.Todo{}
		qry, args := todoStoreTemplateImpl.GetTodoByID(id)
		slog.Debug("GetTodoByID", "qry", qry, "args", &store.JsonLogValue{V: args})
		err := db.QueryRowContext(ctx, qry, args...).Scan(&todo.ID, &todo.Title, &todo.Completed)
		if err != nil {
			err = todoStoreTemplateImpl.ErrorMapFind(err)
		}
		return &todo, err
	}
	return nil, errCtxNoDB
}

// FindTodo implements store.TodoStore.
func (*TodoStoreTmpl) FindTodo(ctx context.Context, filter store.TodoFilter, skip int64, count int) ([]*store.Todo, int64, error) {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sql.DB); ok {
		total := int64(0)
		qry, args := todoStoreTemplateImpl.FindTodoTotal(filter)
		slog.Debug("FindTodoTotal", "qry", qry, "args", &store.JsonLogValue{V: args})
		err := db.QueryRowContext(ctx, qry, args...).Scan(&total)
		if err != nil {
			err = todoStoreTemplateImpl.ErrorMapFind(err)
			return nil, total, err
		}
		qry, args = todoStoreTemplateImpl.FindTodo(filter, skip, count)
		slog.Debug("FindTodo", "qry", qry, "args", &store.JsonLogValue{V: args})
		rows, err := db.QueryContext(ctx, qry, args...)
		if err != nil {
			err = todoStoreTemplateImpl.ErrorMapFind(err)
			return nil, total, err
		}
		defer rows.Close()
		todos := []*store.Todo{}
		for rows.Next() {
			todo := store.Todo{}
			err = rows.Scan(&todo.ID, &todo.Title, &todo.Completed)
			if err != nil {
				err = todoStoreTemplateImpl.ErrorMapFind(err)
				return nil, total, err
			}
			todos = append(todos, &todo)
		}
		return todos, total, nil
	}
	return nil, 0, errCtxNoDB
}
