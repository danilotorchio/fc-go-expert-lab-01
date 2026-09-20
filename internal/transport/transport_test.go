package transport

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestRouter(t *testing.T) {
	h := NewRouter().Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.PathValue("cep")))
	}))

	tests := []struct {
		method, path string
		status       int
		body         string
	}{
		{http.MethodGet, "/weather/01001000", http.StatusOK, "01001000"},
		{http.MethodPost, "/weather/01001000", http.StatusMethodNotAllowed, ""},
		{http.MethodGet, "/", http.StatusNotFound, ""},
	}
	for _, tt := range tests {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, nil))
		if rec.Code != tt.status {
			t.Errorf("%s %s = %d, want %d", tt.method, tt.path, rec.Code, tt.status)
		}
		if tt.body != "" && rec.Body.String() != tt.body {
			t.Errorf("%s %s body = %q, want %q", tt.method, tt.path, rec.Body.String(), tt.body)
		}
	}
}

func TestServerStart(t *testing.T) {
	t.Run("shuts down gracefully on context cancel", func(t *testing.T) {
		port := freePort(t)
		srv := NewServer(&NewServerOpts{
			Port:    port,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}),
		})

		ctx, cancel := context.WithCancel(t.Context())
		done := make(chan error, 1)
		go func() { done <- srv.Start(ctx, time.Second) }()

		waitUntilServing(t, port)
		cancel()

		select {
		case err := <-done:
			if err != nil {
				t.Errorf("Start() = %v, want nil", err)
			}
		case <-time.After(3 * time.Second):
			t.Fatal("Start() did not return after cancel")
		}
	})

	t.Run("returns error when port is in use", func(t *testing.T) {
		l, err := net.Listen("tcp", ":0")
		if err != nil {
			t.Fatal(err)
		}
		defer l.Close()

		srv := NewServer(&NewServerOpts{Port: strconv.Itoa(l.Addr().(*net.TCPAddr).Port)})
		if err := srv.Start(t.Context(), time.Second); err == nil {
			t.Error("Start() on busy port: want error")
		}
	})
}

func freePort(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
}

func waitUntilServing(t *testing.T, port string) {
	t.Helper()
	for range 50 {
		if resp, err := http.Get("http://localhost:" + port); err == nil {
			resp.Body.Close()
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("server did not start")
}
