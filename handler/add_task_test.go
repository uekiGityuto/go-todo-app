package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/uekiGityuto/go-todo-app/entity"
	"github.com/uekiGityuto/go-todo-app/testutil"
)

func TestAddTask(t *testing.T) {
	t.Parallel() // 他のテストと並行して実行する

	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/add_task/ok_req.json.golden",
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/add_task/ok_rsp.json.golden",
			},
		},
		"badRequest": {
			reqFile: "testdata/add_task/bad_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/add_task/bad_req_rsp.json.golden",
			},
		},
	}

	for n, tt := range tests {
		// ゴルーチンで使用されるttを束縛する（ゴルーチンでは宣言時ではなく利用時のttを参照してしまうため）
		// 参照: https://github.com/golang/go/wiki/TableDrivenTests#parallel-testing
		tt := tt
		t.Run(n, func(t *testing.T) { // サブテストをゴルーチンで実行する
			t.Parallel() // 他のテストと並行して実行する

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)

			moq := &AddTaskServiceMock{}
			moq.AddTaskFunc = func(ctx context.Context, title string) (*entity.Task, error) {
				if tt.want.status == http.StatusOK {
					return &entity.Task{ID: 1}, nil
				}
				return nil, errors.New("error from mock")
			}

			// sutはSystems under Test(テスト対象物)の略
			sut := AddTask{
				Service:   moq,
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(t, resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})
	}
}
