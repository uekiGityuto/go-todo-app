package testutil

import (
	"os"
	"testing"

	"github.com/uekiGityuto/go_todo_app/config"
)

func NewConfig(t *testing.T) *config.Config {
	t.Helper()

	// テスト用にconfigを設定する。
	// config.New()を使わずに*config.Configに一つずつ値を設定した方が疎結合になるので良いかもしれないが、
	// DBに接続するためのデフォルト値を*config.Configで管理しているので、それを利用した方が良いと思い、config.New()を実行する形にした。
	cfg, err := config.New()
	if err != nil {
		t.Fatalf("cannot create config: %v", err)
	}
	// CIで実行する場合はCIで実行するようのポート番号に上書きする
	if _, defined := os.LookupEnv("CI"); defined {
		cfg.DBPort = 3306
	}

	return cfg
}
