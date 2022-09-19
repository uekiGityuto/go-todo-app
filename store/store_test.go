package store

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/google/go-cmp/cmp"
	"github.com/uekiGityuto/go_todo_app/clock"
	"github.com/uekiGityuto/go_todo_app/entity"

	"github.com/uekiGityuto/go_todo_app/testutil"
)

func prepareTasks(ctx context.Context, t *testing.T, con Execer) entity.Tasks {
	t.Helper()
	// 一度綺麗にしておく
	if _, err := con.ExecContext(ctx, "DELETE FROM task;"); err != nil {
		t.Logf("failed to initialize task: %v", err)
	}
	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			Title:    "want task 1",
			Status:   entity.TaskStatusTodo,
			Created:  c.Now(),
			Modified: c.Now(),
		}, {
			Title:    "want task 2",
			Status:   entity.TaskStatusTodo,
			Created:  c.Now(),
			Modified: c.Now(),
		}, {
			Title:    "want task 3",
			Status:   entity.TaskStatusDone,
			Created:  c.Now(),
			Modified: c.Now(),
		},
	}

	result, err := con.ExecContext(ctx,
		`INSERT INTO task (title, status, created, modified) VALUES
                                                       (?, ?, ?, ?),
                                                       (?, ?, ?, ?),
                                                       (?, ?, ?, ?);`,
		wants[0].Title, wants[0].Status, wants[0].Created, wants[0].Modified,
		wants[1].Title, wants[1].Status, wants[1].Created, wants[1].Modified,
		wants[2].Title, wants[2].Status, wants[2].Created, wants[2].Modified,
	)
	if err != nil {
		t.Fatal(err)
	}
	// 複数のレコードを作成した時の戻り値は発行されたIDの中で一番小さなIDになることに注意
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	wants[0].ID = entity.TaskID(id)
	wants[1].ID = entity.TaskID(id + 1)
	wants[2].ID = entity.TaskID(id + 2)
	return wants
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

	wants := prepareTasks(ctx, t, tx)

	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx)
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
		WithArgs(okTask.Title, okTask.Status, c.Now(), c.Now()).
		WillReturnResult(sqlmock.NewResult(wantID, 1))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if err := r.AddTask(ctx, xdb, okTask); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}
