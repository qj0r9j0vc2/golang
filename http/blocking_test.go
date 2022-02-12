package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func blocking(w http.ResponseWriter, r *http.Request) {
	select {}
}

func TestBlocking(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(blocking))
	_, _ = http.Get(ts.URL)
	t.Fatal("client did not block")

}
