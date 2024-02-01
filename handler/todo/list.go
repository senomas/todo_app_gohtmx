package todo

import (
	"net/http"

	"github.com/senomas/todo_app/store"
)

func ListTodoHandler(w http.ResponseWriter, r *http.Request) {
	todoStore := store.GetTodoStore()
	rqry := r.URL.Query()
	filter := store.TodoFilter{}
	filter.ID.Set("id", rqry)
	filter.Title.Set("title", rqry)
	todos, _, err := todoStore.FindTodo(r.Context(), filter, 0, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	component := list(todos)
	component.Render(r.Context(), w)
}
