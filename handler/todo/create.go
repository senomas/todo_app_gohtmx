package todo

import (
	"log/slog"
	"net/http"

	"github.com/senomas/todo_app/store"
)

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	todoStore := store.GetTodoStore()

	title := r.Form["title"][0]
	slog.Info("CreateTodoHandler", "title", title)
	_, err = todoStore.CreateTodo(r.Context(), title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	ListTodoHandler(w, r)
}
