package service

import (
	"context"
	"errors"
	"fmt"
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
	err := errors.New("error in repository")

	tests := map[string]struct {
		moq     *TaskAdderMock
		want    *entity.Task
		wantErr error
	}{
		"ok": {
			moq: &TaskAdderMock{
				AddTaskFunc: func(ctx context.Context, db store.Execer, t *entity.Task) error {
					t.ID = 1
					t.Created = created
					t.Modified = modified
					return nil
				},
			},
			want: &entity.Task{
				ID:       1,
				Title:    title,
				Status:   entity.TaskStatusTodo,
				Created:  created,
				Modified: modified,
			},
			wantErr: nil,
		},
		"error": {
			moq: &TaskAdderMock{
				AddTaskFunc: func(ctx context.Context, db store.Execer, t *entity.Task) error {
					return err
				},
			},
			want:    nil,
			wantErr: fmt.Errorf("failed to register: %w", err),
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			sut := &AddTask{
				DB:   nil,
				Repo: tt.moq,
			}
			ctx := context.Background()
			got, err := sut.AddTask(ctx, title)
			if err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Fatalf("unexpected error occurred: %+v", err)
				}
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("got differs: (-got +want)\n%s", diff)
			}
		})
	}
}
