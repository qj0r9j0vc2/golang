package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func blocking(w http.ResponseWriter, r *http.Request) {
	select {}
}

func TestBlocking(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(blocking))
	_, _ = http.Get(ts.URL)
	t.Fatal("client did not block")

}

func TestBlockingWithTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(blocking))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatal(err)
		}
		return
	}

	_ = resp.Body.Close()
}
