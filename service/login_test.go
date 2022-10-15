package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/uekiGityuto/go_todo_app/entity"
	"github.com/uekiGityuto/go_todo_app/store"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	ErrGetUser := errors.New("failed to get user")
	ErrGenJWT := errors.New("failed to generate JWT")
	pw := "password"
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(pw), 10)
	if err != nil {
		t.Fatal(err)
	}
	wantToken := "token"

	tests := map[string]struct {
		userGetterMoq     *UserGetterMock
		tokenGeneratorMoq *TokenGeneratorMock
		want              string
		wantErr           error
	}{
		"ok": {
			userGetterMoq: &UserGetterMock{
				GetUserFunc: func(ctx context.Context, db store.Queryer, name string) (*entity.User, error) {
					return &entity.User{
						Password: string(hashedPW),
					}, nil
				},
			},
			tokenGeneratorMoq: &TokenGeneratorMock{
				GenerateTokenFunc: func(ctx context.Context, u entity.User) ([]byte, error) {
					return []byte(wantToken), nil
				},
			},
			want:    wantToken,
			wantErr: nil,
		},
		"failed to get user": {
			userGetterMoq: &UserGetterMock{
				GetUserFunc: func(ctx context.Context, db store.Queryer, name string) (*entity.User, error) {
					return nil, ErrGetUser
				},
			},
			tokenGeneratorMoq: nil,
			want:              "",
			wantErr:           ErrGetUser,
		},
		"failed to generate token": {
			userGetterMoq: &UserGetterMock{
				GetUserFunc: func(ctx context.Context, db store.Queryer, name string) (*entity.User, error) {
					return &entity.User{
						Password: string(hashedPW),
					}, nil
				},
			},
			tokenGeneratorMoq: &TokenGeneratorMock{
				GenerateTokenFunc: func(ctx context.Context, u entity.User) ([]byte, error) {
					return nil, ErrGenJWT
				},
			},
			want:    "",
			wantErr: ErrGenJWT,
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			sut := Login{
				DB:             nil, // Repoの中でしか利用しないがRepoはmockにするのでnilで良い
				Repo:           tt.userGetterMoq,
				TokenGenerator: tt.tokenGeneratorMoq,
			}
			ctx := context.Background()
			got, err := sut.Login(ctx, "uekiGityuto", pw)
			if err != nil || tt.wantErr != nil {
				switch {
				case err != nil && tt.wantErr == nil:
					t.Fatalf("unexpected error occurred: %+v", err)
				case err == nil && tt.wantErr != nil:
					t.Errorf("expected error is '%+v', but got error is nil", tt.wantErr)
				case !errors.Is(err, tt.wantErr):
					t.Errorf("expected error is '%+v', but got error is '%+v'", tt.wantErr, err)
				}
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("got differs: (-got +want)\n%s", diff)
			}
		})
	}
}
