package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/uekiGityuto/go_todo_app/auth"

	"github.com/google/go-cmp/cmp"

	"github.com/uekiGityuto/go_todo_app/entity"
	"github.com/uekiGityuto/go_todo_app/store"
)

func TestListTasks(t *testing.T) {
	userID := entity.UserID(1)
	tasks := []*entity.Task{
		{
			ID:       entity.TaskID(1),
			UserID:   userID,
			Title:    "テストタスク1",
			Status:   entity.TaskStatusTodo,
			Created:  time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC),
			Modified: time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC),
		}, {
			ID:       entity.TaskID(2),
			UserID:   userID,
			Title:    "テストタスク2",
			Status:   entity.TaskStatusDoing,
			Created:  time.Date(2022, 5, 11, 12, 34, 56, 0, time.UTC),
			Modified: time.Date(2022, 5, 12, 12, 34, 56, 0, time.UTC),
		},
	}
	err := errors.New("error in repository")

	tests := map[string]struct {
		moq     *TaskListerMock
		want    entity.Tasks
		wantErr error
	}{
		"ok": {
			moq: &TaskListerMock{
				ListTasksFunc: func(ctx context.Context, db store.Queryer, id entity.UserID) (entity.Tasks, error) {
					return tasks, nil
				},
			},
			want:    tasks,
			wantErr: nil,
		},
		"error": {
			moq: &TaskListerMock{
				ListTasksFunc: func(ctx context.Context, db store.Queryer, id entity.UserID) (entity.Tasks, error) {
					return nil, errors.New("error in repository")
				},
			},
			want:    nil,
			wantErr: fmt.Errorf("failed to list: %w", err),
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			sut := ListTask{
				DB:   nil, // Repoの中でしか利用しないがRepoはmockにするのでnilで良い
				Repo: tt.moq,
			}
			ctx := context.Background()
			ctx = auth.SetUserID(ctx, userID) // 本当はauth.GetUserIDをモックに置き換えられるようにした方が良い気がする。
			got, err := sut.ListTasks(ctx)
			if err != nil || tt.wantErr != nil {
				switch {
				case err != nil && tt.wantErr == nil:
					t.Fatalf("unexpected error occurred: %+v", err)
				case err == nil && tt.wantErr != nil:
					t.Errorf("expected error is '%+v', but got error is nil", tt.wantErr)
				case err.Error() != tt.wantErr.Error():
					t.Errorf("expected error is '%+v', but got error is '%+v'", tt.wantErr, err)
				}
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("got differs: (-got +want)\n%s", diff)
			}
		})
	}
}
