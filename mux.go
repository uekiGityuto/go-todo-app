package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/uekiGityuto/go_todo_app/handler"
	"github.com/uekiGityuto/go_todo_app/store"
	"net/http"
)

func NewMux() http.Handler {
	mux := chi.NewRouter()

	mux.HandleFunc("/health",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=uft-8")
			_, _ = w.Write([]byte(`{"status": "ok"}`))
		})

	v := validator.New()
	at := &handler.AddTask{Store: store.Tasks, Validator: v}
	mux.Post("/tasks", at.ServeHTTP)

	lt := &handler.ListTask{Store: store.Tasks}
	mux.Get("/tasks", lt.ServeHTTP)

	return mux
}
