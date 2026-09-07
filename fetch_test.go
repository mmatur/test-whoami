package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchHandler(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("hello from upstream"))
	}))
	defer upstream.Close()

	req := httptest.NewRequest(http.MethodGet, "/fetch?url="+upstream.URL, nil)
	rec := httptest.NewRecorder()

	fetchHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if body := rec.Body.String(); body != "hello from upstream" {
		t.Fatalf("unexpected body: %q", body)
	}
}
