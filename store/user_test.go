package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/uekiGityuto/go_todo_app/testutil"

	"github.com/google/go-cmp/cmp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/uekiGityuto/go_todo_app/clock"
	"github.com/uekiGityuto/go_todo_app/entity"
)

func TestRepository_RegisterUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c := clock.FixedClocker{}
	var wantID int64 = 20
	okUser := &entity.User{
		Name:     "uekiGityuto",
		Password: "password",
		Role:     "user",
		Created:  c.Now(),
		Modified: c.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	mock.ExpectExec(`INSERT INTO user`).
		WithArgs(okUser.Name, okUser.Password, okUser.Role, c.Now(), c.Now()).
		WillReturnResult(sqlmock.NewResult(wantID, 1))

	xdb := sqlx.NewDb(db, "mysql")
	sut := &Repository{Clocker: c}
	if err := sut.RegisterUser(ctx, xdb, okUser); err != nil {
		t.Errorf("want no error, but got '%+v'", err)
	}
	if diff := cmp.Diff(okUser.ID, entity.UserID(wantID)); diff != "" {
		t.Errorf("got differs: (-got +want)\n%s", diff)
	}
}

func TestRepository_RegisterUserWhenDuplicate(t *testing.T) {
	ctx := context.Background()
	// トランザクションを張ることでこのテストケースの中だけのテーブル状態にする
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	// このテストケースが終わったらテーブルの状態を元に戻す
	t.Cleanup(func() {
		_ = tx.Rollback()
	})
	if err != nil {
		t.Fatal(err)
	}

	c := clock.FixedClocker{}
	user := &entity.User{
		Name:     "uekiGityuto",
		Password: "password",
		Role:     "user",
		Created:  c.Now(),
		Modified: c.Now(),
	}

	sut := &Repository{
		Clocker: c,
	}
	if err := sut.RegisterUser(ctx, tx, user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 同じNameでユーザ登録する
	wantErr := fmt.Errorf("cannot create same user: %w", ErrAlreadyEntry)
	err = sut.RegisterUser(ctx, tx, user)
	if err == nil {
		t.Errorf("expected error is '%+v', but got error is nil", wantErr)
	} else if err.Error() != wantErr.Error() {
		t.Errorf("expected error is '%+v', but got error is '%+v'", wantErr, err)
	}
}
