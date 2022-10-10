package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/uekiGityuto/go_todo_app/testutil/fixture"

	"github.com/jmoiron/sqlx"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/google/go-cmp/cmp"
	"github.com/uekiGityuto/go_todo_app/clock"
	"github.com/uekiGityuto/go_todo_app/entity"

	"github.com/uekiGityuto/go_todo_app/testutil"
)

func prepareUser(ctx context.Context, t *testing.T, con Execer) entity.UserID {
	t.Helper()

	u := fixture.User(nil)
	sql := `INSERT INTO user (name, password, role, created, modified) VALUES (?, ?, ?, ?, ?);`
	result, err := con.ExecContext(ctx, sql, u.Name, u.Password, u.Role, u.Created, u.Modified)
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("failed to got user_id: %v", err)
	}
	return entity.UserID(id)
}

func prepareTasks(ctx context.Context, t *testing.T, con Execer) (entity.UserID, entity.Tasks) {
	t.Helper()

	// 一度綺麗にしておく
	if _, err := con.ExecContext(ctx, "DELETE FROM user;"); err != nil {
		t.Logf("failed to initialize user: %v", err)
	}
	if _, err := con.ExecContext(ctx, "DELETE FROM task;"); err != nil {
		t.Logf("failed to initialize task: %v", err)
	}

	userID := prepareUser(ctx, t, con)
	otherUserID := prepareUser(ctx, t, con)

	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			UserID:   userID,
			Title:    "want task 1",
			Status:   entity.TaskStatusTodo,
			Created:  c.Now(),
			Modified: c.Now(),
		}, {
			UserID:   userID,
			Title:    "want task 2",
			Status:   entity.TaskStatusTodo,
			Created:  c.Now(),
			Modified: c.Now(),
		},
	}

	tasks := entity.Tasks{
		wants[0],
		{
			UserID:   otherUserID,
			Title:    "want task 3",
			Status:   entity.TaskStatusDone,
			Created:  c.Now(),
			Modified: c.Now(),
		},
		wants[1],
	}

	fmt.Printf("user_id: %d, tasks[0].UserID: %d\n", userID, tasks[0].UserID)

	result, err := con.ExecContext(ctx,
		`INSERT INTO task (user_id, title, status, created, modified) VALUES
                                                       (?, ?, ?, ?, ?),
                                                       (?, ?, ?, ?, ?),
                                                       (?, ?, ?, ?, ?);`,
		tasks[0].UserID, tasks[0].Title, tasks[0].Status, tasks[0].Created, tasks[0].Modified,
		tasks[1].UserID, tasks[1].Title, tasks[1].Status, tasks[1].Created, tasks[1].Modified,
		tasks[2].UserID, tasks[2].Title, tasks[2].Status, tasks[2].Created, tasks[2].Modified,
	)
	if err != nil {
		t.Fatal(err)
	}
	// 複数のレコードを作成した時の戻り値は発行されたIDの中で一番小さなIDになることに注意
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	tasks[0].ID = entity.TaskID(id)
	tasks[1].ID = entity.TaskID(id + 1)
	tasks[2].ID = entity.TaskID(id + 2)
	return userID, wants
}

func TestRepository_ListTasks(t *testing.T) {
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

	wantUserID, wants := prepareTasks(ctx, t, tx)

	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx, wantUserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := cmp.Diff(gots, wants); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

func TestRepository_AddTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c := clock.FixedClocker{}
	var wantID int64 = 20
	okTask := &entity.Task{
		UserID:   1,
		Title:    "ok task",
		Status:   entity.TaskStatusTodo,
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
	//mock.ExpectExec(`INSERT INTO task \(title, status, created, modified\) VALUES \(\?, \?, \?, \?\)`).
	mock.ExpectExec(`INSERT INTO task`).
		WithArgs(okTask.UserID, okTask.Title, okTask.Status, c.Now(), c.Now()).
		WillReturnResult(sqlmock.NewResult(wantID, 1))

	xdb := sqlx.NewDb(db, "mysql")
	sut := &Repository{Clocker: c}
	if err := sut.AddTask(ctx, xdb, okTask); err != nil {
		t.Errorf("want no error, but got '%+v'", err)
	}
	if diff := cmp.Diff(okTask.ID, entity.TaskID(wantID)); diff != "" {
		t.Errorf("got differs: (-got +want)\n%s", diff)
	}
}
