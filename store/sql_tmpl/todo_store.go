package sql_tmpl

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/senomas/todo_app/store"
)

var (
	errCtxNoDB     = errors.New("no db defined in context")
	errNoStoreTmpl = errors.New("todo store template not initialized")
)

type TodoStoreTmpl struct{}

type TodoStoreTemplate interface {
	InsertTodo(t *store.Todo) string
	UpdateTodo(t *store.Todo) string
	DeleteTodoByID() string

	GetTodoByID() string
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

func (t *TodoStoreTmpl) CreateTodo(ctx context.Context, title string) (*store.Todo, error) {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sqlx.DB); ok {
		todo := store.Todo{Title: title}
		qry := todoStoreTemplateImpl.InsertTodo(&todo)
		rs, err := db.NamedExecContext(ctx, qry, &todo)
		if err != nil {
			slog.Warn("Error insert todo", "qry", qry, "error", err)
			return nil, err
		}
		todo.ID, err = rs.LastInsertId()
		return &todo, err
	}
	return nil, errCtxNoDB
}

func (t *TodoStoreTmpl) UpdateTodo(ctx context.Context, todo store.Todo) error {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sqlx.DB); ok {
		qry := todoStoreTemplateImpl.UpdateTodo(&todo)
		db.ExecContext(ctx, qry, todo)
	}
	return errCtxNoDB
}

func (t *TodoStoreTmpl) DeleteTodoByID(ctx context.Context, id int64) error {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sqlx.DB); ok {
		qry := todoStoreTemplateImpl.DeleteTodoByID()
		db.ExecContext(ctx, qry, id)
	}
	return errCtxNoDB
}

func (t *TodoStoreTmpl) GetTodoByID(ctx context.Context, id int64) (*store.Todo, error) {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sqlx.DB); ok {
		todo := store.Todo{}
		qry := todoStoreTemplateImpl.GetTodoByID()
		err := db.GetContext(ctx, &todo, qry, id)
		return &todo, err
	}
	return nil, errCtxNoDB
}
