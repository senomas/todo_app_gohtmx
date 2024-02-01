package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/senomas/todo_app/handler"
	"github.com/senomas/todo_app/handler/todo"
	"github.com/senomas/todo_app/store"
	_ "github.com/senomas/todo_app/store/sql_tmpl"
	_ "github.com/senomas/todo_app/store/sql_tmpl/sqlite"
)

func main() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path := filepath.Dir(ex)
	assets := filepath.Join(path, "assets")

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx := context.WithValue(context.Background(), store.StoreCtxDB, db)
	_, err = store.Migrate(ctx)
	if err != nil {
		panic(err)
	}
	todoStore := store.GetTodoStore()
	for i := 1; i <= 4; i++ {
		_, err := todoStore.CreateTodo(ctx, "Todo "+strconv.Itoa(i))
		if err != nil {
			panic(err)
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, assets+"/index.html")
	})
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assets))))
	http.HandleFunc("/api/todo", handler.Handle(handler.Mux{
		Get:  todo.ListTodoHandler,
		Post: todo.CreateTodoHandler,
	}))
	server := &http.Server{Addr: ":8080", BaseContext: func(net.Listener) context.Context {
		return ctx
	}}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
