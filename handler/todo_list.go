package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/senomas/todo_app/store"
)

func ListTodoHandler(w http.ResponseWriter, r *http.Request) {
	todoStore := store.GetTodoStore()
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		title := r.Form["title"][0]
		slog.Info("CreateTodoHandler", "title", title)
		if title == "" {
			http.Error(w, "title should not empty", http.StatusBadRequest)
			return
		}
		if strings.Contains(title, "xxx") {
			http.Error(w, "title should not contain xxx", http.StatusBadRequest)
			return
		}
		_, err = todoStore.CreateTodo(r.Context(), title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	rqry := r.URL.Query()
	filter := store.TodoFilter{}
	filter.ID.Set("id", rqry)
	filter.Title.Set("title", rqry)
	todos, _, err := todoStore.FindTodo(r.Context(), filter, 0, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	component := ListTodo(todos)
	component.Render(r.Context(), w)
}
