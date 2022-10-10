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

func TestRepository_RegisterUser_WhenDuplicate(t *testing.T) {
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

	// 一度綺麗にしておく
	if _, err := tx.ExecContext(ctx, "DELETE FROM user;"); err != nil {
		t.Logf("failed to initialize user: %v", err)
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

func prepareUser(ctx context.Context, t *testing.T, con Execer, u *entity.User) *entity.User {
	t.Helper()
	// 一度綺麗にしておく
	if _, err := con.ExecContext(ctx, "DELETE FROM user;"); err != nil {
		t.Logf("failed to initialize user: %v", err)
	}

	sql := `INSERT INTO user (name, password, role, created, modified) VALUES (?, ?, ?, ?, ?);`
	result, err := con.ExecContext(ctx, sql, u.Name, u.Password, u.Role, u.Created, u.Modified)
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	u.ID = entity.UserID(id)

	return u
}

func TestRepository_GetUser(t *testing.T) {
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

	name := "uekiGityuto"
	user := &entity.User{
		Name:     name,
		Password: "password",
		Role:     "admin",
		Created:  clock.FixedClocker{}.Now(),
		Modified: clock.FixedClocker{}.Now(),
	}
	want := prepareUser(ctx, t, tx, user)

	sut := &Repository{}
	got, err := sut.GetUser(ctx, tx, name)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := cmp.Diff(got, want); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}
