package handler

import (
	"context"
	"net/http"

	"github.com/uekiGityuto/go_todo_app/entity"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . ListTaskService AddTaskService RegisterUserService LoginService LogoutService

type ListTaskService interface {
	ListTasks(ctx context.Context) (entity.Tasks, error)
}

type AddTaskService interface {
	AddTask(ctx context.Context, title string) (*entity.Task, error)
}

type RegisterUserService interface {
	RegisterUser(ctx context.Context, name, password, role string) (*entity.User, error)
}

type LoginService interface {
	Login(ctx context.Context, name, pw string) (string, error)
}

type LogoutService interface {
	Logout(r *http.Request) error
}
