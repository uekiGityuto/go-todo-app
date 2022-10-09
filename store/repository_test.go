package store

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/uekiGityuto/go_todo_app/config"
)

func TestNew(t *testing.T) {
	type args struct {
		ctx context.Context
		cfg *config.Config
	}
	// DBに接続するためのconfigを設定する。
	// config.New()を使わずに*config.Configに一つずつ値を設定した方が疎結合になるので良いかもしれないが、
	// DBに接続するためのデフォルト値を*config.Configで管理しているので、それを利用した方が良いと思い、config.New()を実行する形にした。
	correctCfg, err := config.New()
	if err != nil {
		t.Fatalf("cannot create config: %v", err)
	}
	if _, defined := os.LookupEnv("CI"); defined {
		correctCfg.DBPort = 3306 // CIで実行する場合はCIで実行するようのポート番号に上書きする
	}

	incorrectCfg, err := config.New()
	if err != nil {
		t.Fatalf("cannot create config: %v", err)
	}
	incorrectCfg.DBPort = 99999 // 正しくないポート番号を設定

	tests := map[string]struct {
		args      args
		isWantErr bool
	}{
		"ok": {
			args: args{
				ctx: context.Background(),
				cfg: correctCfg,
			},
			isWantErr: false,
		},
		"error": {
			args: args{
				ctx: context.Background(),
				cfg: incorrectCfg,
			},
			isWantErr: true,
		},
	}
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			gotDB, gotFunc, err := New(tt.args.ctx, tt.args.cfg)
			if err != nil || tt.isWantErr == true {
				if err != nil && tt.isWantErr == false {
					t.Fatalf("unexpected error occurred: %+v", err)
				} else if err == nil && tt.isWantErr == true {
					t.Error("expected error, but got error is nil")
				} else {
					// 期待通りにエラーになった場合はこの時点でテスト成功とする
					fmt.Printf("error occurred as expected: %+v\n", err)
					return
				}
			}

			// 接続出来ればOK
			ctx, cancel := context.WithTimeout(tt.args.ctx, 2*time.Second)
			defer cancel()
			if err = gotDB.PingContext(ctx); err != nil {
				t.Errorf("failed to connect DB: %+v", err)
			}
			// コネクションクローズ
			gotFunc()
			// クローズしたので接続出来なければOK
			if err = gotDB.PingContext(tt.args.ctx); err == nil {
				t.Errorf("failed to close DB connection: %+v", err)
			}
		})
	}
}
