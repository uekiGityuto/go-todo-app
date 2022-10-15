package store

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/uekiGityuto/go_todo_app/testutil"

	"github.com/uekiGityuto/go_todo_app/config"
)

func TestNew(t *testing.T) {
	type args struct {
		ctx context.Context
		cfg *config.Config
	}
	// DB接続に成功するようなconfigを設定する。
	correctCfg := testutil.NewConfig(t)

	// DB接続に失敗するようなconfigを設定する。
	incorrectCfg := testutil.NewConfig(t)
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
			gotDB, gotFunc, gotErr := New(tt.args.ctx, tt.args.cfg)
			if gotErr != nil || tt.isWantErr == true {
				switch {
				case gotErr != nil && tt.isWantErr == false:
					t.Fatalf("unexpected error occurred: %+v", gotErr)
				case gotErr == nil && tt.isWantErr == true:
					t.Error("expected error, but got error is nil")
				default:
					fmt.Printf("error occurred as expected: %+v\n", gotErr)
					return
				}
			}

			// 接続出来ればOK
			ctx, cancel := context.WithTimeout(tt.args.ctx, 2*time.Second)
			defer cancel()
			if err := gotDB.PingContext(ctx); err != nil {
				t.Errorf("failed to connect DB: %+v", err)
			}
			// コネクションクローズ
			gotFunc()
			// クローズしたので接続出来なければOK
			if err := gotDB.PingContext(tt.args.ctx); err == nil {
				t.Errorf("failed to close DB connection: %+v", err)
			}
		})
	}
}
