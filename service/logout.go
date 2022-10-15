package service

import (
	"net/http"
)

type Logout struct {
	UserIDDeleter UserIDDeleter
}

func (l *Logout) Logout(r *http.Request) error {
	return l.UserIDDeleter.DeleteUserID(r)
}
