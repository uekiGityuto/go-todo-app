package auth

import (
	"context"

	"github.com/uekiGityuto/go-todo-app/entity"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . Store
type Store interface {
	Save(ctx context.Context, key string, userID entity.UserID) error
	Load(ctx context.Context, key string) (entity.UserID, error)
	Delete(ctx context.Context, key string) error
}
