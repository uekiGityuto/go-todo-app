package entity

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestUser_ComparePassword(t *testing.T) {
	pw := "password"
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(pw), 10)
	if err != nil {
		t.Fatal(err)
	}
	u := &User{
		Password: string(hashedPW),
	}
	if err := u.ComparePassword(pw); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}

func TestUser_ComparePassword_NG(t *testing.T) {
	pw := "password"
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(pw), 10)
	if err != nil {
		t.Fatal(err)
	}
	u := &User{
		Password: string(hashedPW),
	}
	if err := u.ComparePassword(pw + "test"); err == nil {
		t.Error("want error, but got nil")
	}
}
