package store

import "context"

type TodoStore interface {
	Init(ctx context.Context) error
	CreateTodo(ctx context.Context, title string) (*Todo, error)
	UpdateTodo(ctx context.Context, todo Todo) error
	DeleteTodoByID(ctx context.Context, id int64) error

	GetTodoByID(ctx context.Context, id int64) (*Todo, error)
	FindTodo(ctx context.Context, filter TodoFilter, skip int64, count int) ([]*Todo, int64, error)
}

var todoStoreImpl TodoStore

func SetupTodoStoreImplementation(s TodoStore) {
	todoStoreImpl = s
}

func GetTodoStore() TodoStore {
	if todoStoreImpl == nil {
		panic("todo store not initialized")
	}
	return todoStoreImpl
}
