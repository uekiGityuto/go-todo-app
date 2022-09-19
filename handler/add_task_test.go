package handler

import (
	"bytes"
	"github.com/go-playground/validator/v10"
	"github.com/uekiGityuto/go_todo_app/entity"
	"github.com/uekiGityuto/go_todo_app/store"
	"github.com/uekiGityuto/go_todo_app/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
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

			sut := AddTask{
				Store: &store.TaskStore{
					Tasks: map[entity.TaskID]*entity.Task{},
				},
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r) // sutはSystems under Test(テスト対象物)の略

			resp := w.Result()
			testutil.AssertResponse(t,
				resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile),
			)
		})
	}
}
