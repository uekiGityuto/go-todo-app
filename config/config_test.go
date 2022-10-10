package config

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	wantPort := 3333
	t.Setenv("PORT", fmt.Sprint(wantPort))

	got, err := New()
	if err != nil {
		t.Fatalf("cannot create config: %v", err)
	}

	// 環境変数の値がセットされるか確認
	if got.Port != wantPort {
		t.Errorf("want %d, but %d", wantPort, got.Port)
	}
	// 環境変数の値がないときにデフォルトの値がセットされるか確認
	wantEnv := "dev"
	if got.Env != wantEnv {
		t.Errorf("want %s, but %s", wantEnv, got.Env)
	}
	wantDBHost := "127.0.0.1"
	if got.DBHost != wantDBHost {
		t.Errorf("want %s, but %s", wantDBHost, got.DBHost)
	}
	wantDBPort := 33306
	if got.DBPort != wantDBPort {
		t.Errorf("want %d, but %d", wantDBPort, got.DBPort)
	}
	wantDBUser := "todo"
	if got.DBUser != wantDBUser {
		t.Errorf("want %s, but %s", wantDBUser, got.DBHost)
	}
	wantDBPassword := "todo"
	if got.DBPassword != wantDBPassword {
		t.Errorf("want %s, but %s", wantDBPassword, got.DBHost)
	}
	wantDBName := "todo"
	if got.DBName != wantDBName {
		t.Errorf("want %s, but %s", wantDBName, got.DBHost)
	}
	wantRedisHost := "127.0.0.1"
	if got.RedisHost != wantRedisHost {
		t.Errorf("want %s, but %s", wantRedisHost, got.RedisHost)
	}
	wantRedisPort := 36379
	if got.RedisPort != wantRedisPort {
		t.Errorf("want %d, but %d", wantRedisPort, got.RedisPort)
	}
}
