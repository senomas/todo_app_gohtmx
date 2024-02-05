package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/senomas/gosvc_todo/todo_store"
)

func ListTodoCountHandler(w http.ResponseWriter, r *http.Request) {
	todoStore := todo_store.GetTodoStore()
	_, count, err := todoStore.FindTodo(r.Context(), todo_store.TodoFilter{}, 0, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%d", count)
}

func ListTodoHandler(w http.ResponseWriter, r *http.Request) {
	todoStore := todo_store.GetTodoStore()
	spath := strings.Split(r.URL.EscapedPath(), "/")
	id, idErr := strconv.ParseInt(spath[len(spath)-1], 10, 64)
	if idErr != nil {
		id = -1
	}
	switch r.Method {
	case "POST":
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		v := r.Form["id"]
		var uid int64 = -1
		if len(v) == 1 {
			uid, err = strconv.ParseInt(v[0], 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		}

		var title string
		v = r.Form["title"]
		if len(v) == 1 {
			title = v[0]
		} else {
			http.Error(w, "title should not empty", http.StatusBadRequest)
			return
		}
		if title == "" {
			http.Error(w, "title should not empty", http.StatusBadRequest)
			return
		}
		if strings.Contains(title, "xxx") {
			http.Error(w, "title should not contain xxx", http.StatusBadRequest)
			return
		}
		if uid == -1 {
			slog.Debug("CreateTodoHandler", "title", title)
			_, err = todoStore.CreateTodo(r.Context(), title)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("HX-Trigger", "new-todo")
		} else {
			slog.Debug("UpdateTodoHandler", "id", uid, "title", title)
			err = todoStore.UpdateTodo(r.Context(), todo_store.Todo{
				ID:    uid,
				Title: title,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("HX-Trigger", "update-todo")
		}
	case "DELETE":
		if idErr != nil {
			http.Error(w, idErr.Error(), http.StatusInternalServerError)
			return
		}
		err := todoStore.DeleteTodoByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Trigger", "delete-todo")
	}

	rqry := r.URL.Query()
	filter := todo_store.TodoFilter{}
	filter.ID.Set("id", rqry)
	filter.Title.Set("title", rqry)
	todos, _, err := todoStore.FindTodo(r.Context(), filter, 0, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	component := ListTodo(todos, id)
	component.Render(r.Context(), w)
}
