package service

import (
	"context"
	"testing"
	"time"

	"github.com/uekiGityuto/go_todo_app/store"

	"github.com/google/go-cmp/cmp"

	"github.com/uekiGityuto/go_todo_app/entity"
)

func TestAddTask(t *testing.T) {
	created := time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC)
	modified := time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC)
	title := "テストタスク"
	want := &entity.Task{
		ID:       1,
		Title:    title,
		Status:   entity.TaskStatusTodo,
		Created:  created,
		Modified: modified,
	}

	moq := &TaskAdderMock{}
	moq.AddTaskFunc = func(ctx context.Context, db store.Execer, t *entity.Task) error {
		t.ID = 1
		t.Created = created
		t.Modified = modified
		return nil
	}

	sut := &AddTask{
		DB:   nil,
		Repo: moq,
	}

	ctx := context.Background()
	got, err := sut.AddTask(ctx, title)
	if err != nil {
		t.Fatalf("failed to add task: %+v", err)
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("got differs: (-got +want)\n%s", diff)
	}
}
