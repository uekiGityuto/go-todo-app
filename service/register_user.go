package service

import (
	"context"
	"fmt"

	"github.com/uekiGityuto/go_todo_app/entity"
	"github.com/uekiGityuto/go_todo_app/store"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	DB   store.Execer
	Repo UserRegister
}

func (r *RegisterUser) RegisterUser(ctx context.Context, name, password, role string) (*entity.User, error) {
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}
	user := &entity.User{
		Name:     name,
		Password: string(hashedPW),
		Role:     role,
	}
	if err := r.Repo.RegisterUser(ctx, r.DB, user); err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return user, nil
}
