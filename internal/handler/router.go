package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Post("/register", Register)
	r.Post("/login", Login)
	r.Get("/items", ListItems)
	r.Post("/items", CreateItem)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 not found", http.StatusNotFound)
	})

	return r
}
