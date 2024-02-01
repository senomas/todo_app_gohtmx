package handler

import (
	"net/http"
)

type Mux struct {
	Get  func(http.ResponseWriter, *http.Request)
	Post func(http.ResponseWriter, *http.Request)
}

func Handle(m Mux) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if m.Get != nil {
				m.Get(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case http.MethodPost:
			if m.Post != nil {
				m.Post(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
