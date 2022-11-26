package main

import (
	"context"
	"net/http"

	"github.com/uekiGityuto/go-todo-app/auth"

	"github.com/uekiGityuto/go-todo-app/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/uekiGityuto/go-todo-app/clock"
	"github.com/uekiGityuto/go-todo-app/config"
	"github.com/uekiGityuto/go-todo-app/handler"
	"github.com/uekiGityuto/go-todo-app/store"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()

	mux.HandleFunc("/health",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=uft-8")
			_, _ = w.Write([]byte(`{"status": "ok"}`))
		})

	v := validator.New()
	clocker := clock.RealClocker{}

	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}

	r := &store.Repository{Clocker: clocker}
	redisClient, err := store.NewKVS(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	jwter, err := auth.NewJWTer(redisClient, clocker)
	if err != nil {
		return nil, cleanup, err
	}

	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: r},
		Validator: v,
	}
	mux.Post("/register", ru.ServeHTTP)

	lin := &handler.Login{
		Service: &service.Login{
			DB:             db,
			Repo:           r,
			TokenGenerator: jwter,
		},
		Validator: v,
	}
	mux.Post("/login", lin.ServeHTTP)

	at := &handler.AddTask{
		Service:   &service.AddTask{DB: db, Repo: r},
		Validator: v,
	}
	lt := &handler.ListTask{
		Service: &service.ListTask{DB: db, Repo: r},
	}
	mux.Route("/tasks", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter))
		r.Post("/", at.ServeHTTP)
		r.Get("/", lt.ServeHTTP)
	})

	mux.Route("/admin", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter), handler.AdminMiddleware)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=uft-8")
			_, _ = w.Write([]byte(`{"message": "admin only"}`))
		})
	})

	lout := &handler.Logout{
		Service: &service.Logout{UserIDDeleter: jwter},
	}
	mux.Route("/logout", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter))
		r.Get("/", lout.ServeHTTP)
	})

	return mux, cleanup, err
}
