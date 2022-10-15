package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/uekiGityuto/go_todo_app/auth"

	"github.com/uekiGityuto/go_todo_app/store"

	"github.com/google/go-cmp/cmp"

	"github.com/uekiGityuto/go_todo_app/entity"
)

func TestAddTask(t *testing.T) {
	created := time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC)
	modified := time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC)
	title := "テストタスク"
	userID := entity.UserID(1)
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
				UserID:   userID,
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
				DB:   nil, // Repoの中でしか利用しないがRepoはmockにするのでnilで良い
				Repo: tt.moq,
			}
			ctx := context.Background()
			ctx = auth.SetUserID(ctx, userID) // 本当はauth.GetUserIDをモックに置き換えられるようにした方が良い気がする。
			got, err := sut.AddTask(ctx, title)
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
