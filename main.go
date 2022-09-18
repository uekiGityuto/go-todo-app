package main

import (
	"context"
	"fmt"
	"github.com/uekiGityuto/go_todo_app/config"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.Serve(l); err != nil && err != http.ErrServerClosed { // ErrServerClosedは正常にシャットダウンされたことを示す
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// キャンセルされるまでここでブロックされ、cancelされたらWebサーバをGraceful Shutdownする
	// （ただ、キャンセルするような処理を実装していないので、テスト以外では以下は実行されない）
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	// Go関数で実行した関数の戻り値(err)を返す
	return eg.Wait()
}
