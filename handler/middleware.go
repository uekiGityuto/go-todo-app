package handler

import (
	"fmt"
	"net/http"

	"github.com/uekiGityuto/go-todo-app/auth"
)

func AuthMiddleware(j *auth.JWTer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req, err := j.FillContext(r)
			if err != nil {
				RespondJSON(r.Context(), w, ErrResponse{
					Message: "not find auth info",
					Details: []string{err.Error()},
				}, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}

// AdminMiddleware contextにユーザ情報が含まれていることを前提としているため、AuthMiddlewareより後に適用すること。
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !auth.IsAdmin(r.Context()) {
			RespondJSON(r.Context(), w, ErrResponse{
				Message: "not admin",
			}, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				RespondJSON(r.Context(), w, ErrResponse{
					Message: "error occurred",
					Details: []string{fmt.Sprintf("%v", err)},
				}, http.StatusUnauthorized)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
