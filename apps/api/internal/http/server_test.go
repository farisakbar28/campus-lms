package http

import (
	"io"
	"log/slog"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpointsReturnOK(t *testing.T) {
	server := NewServer(":8080", slog.New(slog.NewTextHandler(io.Discard, nil)))

	for _, path := range []string{"/healthz", "/readyz"} {
		request := httptest.NewRequest(nethttp.MethodGet, path, nil)
		response := httptest.NewRecorder()

		server.Handler.ServeHTTP(response, request)

		if response.Code != nethttp.StatusOK {
			t.Errorf("%s status = %d, want %d", path, response.Code, nethttp.StatusOK)
		}
		if got := response.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
			t.Errorf("%s Content-Type = %q, want application/json", path, got)
		}
	}
}

func TestServerSetsAllTimeouts(t *testing.T) {
	server := NewServer(":8080", slog.New(slog.NewTextHandler(io.Discard, nil)))

	if server.ReadHeaderTimeout == 0 {
		t.Error("ReadHeaderTimeout = 0, want non-zero")
	}
	if server.ReadTimeout == 0 {
		t.Error("ReadTimeout = 0, want non-zero")
	}
	if server.WriteTimeout == 0 {
		t.Error("WriteTimeout = 0, want non-zero")
	}
	if server.IdleTimeout == 0 {
		t.Error("IdleTimeout = 0, want non-zero")
	}
}
