package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/uekiGityuto/go_todo_app/testutil"
)

func TestNewMux_health(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)

	ctx := context.Background()
	cfg := testutil.NewConfig(t)
	gotHandler, _, gotErr := NewMux(ctx, cfg)
	if gotErr != nil {
		t.Fatalf("failed to create mux: %+v", gotErr)
	}

	gotHandler.ServeHTTP(w, r)
	resp := w.Result()
	t.Cleanup(func() {
		_ = resp.Body.Close()
	})

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want status code %d, but %d", http.StatusOK, resp.StatusCode)
	}
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body %v", err)
	}

	want := `{"status": "ok"}`
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
}
