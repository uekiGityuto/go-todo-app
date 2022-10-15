package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogout(t *testing.T) {
	moq := &UserIDDeleterMock{}
	moq.DeleteUserIDFunc = func(r *http.Request) error {
		return nil
	}

	sut := Logout{UserIDDeleter: moq}
	req := httptest.NewRequest(
		http.MethodGet,
		`https://github.com/uekiGityuto`,
		nil,
	)
	err := sut.Logout(req)
	if err != nil {
		t.Errorf("want no error, but got error: %+v", err)
	}
}
