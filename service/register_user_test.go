package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/go-cmp/cmp"
	"github.com/uekiGityuto/go_todo_app/entity"
	"github.com/uekiGityuto/go_todo_app/store"
)

func TestRegisterUser(t *testing.T) {
	created := time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC)
	modified := time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC)
	user := "テスト太郎"
	password := "password"
	role := "user"
	err := errors.New("error in repository")

	tests := map[string]struct {
		moq           *UserRegisterMock
		wantWithoutPW *entity.User
		wantPW        []byte
		wantErr       error
	}{
		"ok": {
			moq: &UserRegisterMock{
				RegisterUserFunc: func(ctx context.Context, db store.Execer, u *entity.User) error {
					u.ID = 1
					u.Created = created
					u.Modified = modified
					return nil
				},
			},
			wantWithoutPW: &entity.User{
				ID:       1,
				Name:     user,
				Role:     role,
				Created:  created,
				Modified: modified,
			},
			wantPW:  []byte(password),
			wantErr: nil,
		},
		"error": {
			moq: &UserRegisterMock{
				RegisterUserFunc: func(ctx context.Context, db store.Execer, t *entity.User) error {
					return err
				},
			},
			wantWithoutPW: nil,
			wantPW:        nil,
			wantErr:       fmt.Errorf("failed to register: %w", err),
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			sut := &RegisterUser{
				DB:   nil, // Repoの中でしか利用しないがRepoはmockにするのでnilで良い
				Repo: tt.moq,
			}
			ctx := context.Background()
			got, err := sut.RegisterUser(ctx, user, password, role)
			if err != nil || tt.wantErr != nil {
				switch {
				case err != nil && tt.wantErr == nil:
					t.Fatalf("unexpected error occurred: %+v", err)
				case err == nil && tt.wantErr != nil:
					t.Errorf("expected error is '%+v', but got error is nil", tt.wantErr)
				case err.Error() != tt.wantErr.Error():
					t.Errorf("expected error is '%+v', but got error is '%+v'", tt.wantErr, err)
				default:
					return
				}
			}
			gotHashedPW := got.Password
			got.Password = ""
			if diff := cmp.Diff(got, tt.wantWithoutPW); diff != "" {
				t.Errorf("got differs: (-got +want)\n%s", diff)
			}
			if bcrypt.CompareHashAndPassword([]byte(gotHashedPW), tt.wantPW) != nil {
				t.Errorf("failed to compare hash(%+v) and password(%+v)", got.Password, string(tt.wantPW))
			}
		})
	}
}
