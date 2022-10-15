package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/uekiGityuto/go_todo_app/testutil"
)

func TestLogout_ServeHTTP(t *testing.T) {
	moq := &LogoutServiceMock{}
	moq.LogoutFunc = func(r *http.Request) error {
		return nil
	}
	sut := Logout{Service: moq}

	wantStatus := http.StatusOK
	wantRspFile := "testdata/logout/ok_rsp.json.golden"
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/logout", nil)
	sut.ServeHTTP(w, r)

	resp := w.Result()
	testutil.AssertResponse(t, resp, wantStatus, testutil.LoadFile(t, wantRspFile))
}
