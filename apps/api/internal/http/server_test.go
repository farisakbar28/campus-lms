package http

import (
	"context"
	"errors"
	"io"
	"log/slog"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/farisakbar28/campus-lms/apps/api/internal/domain"
)

type readyStub struct{ err error }

func (stub readyStub) Ping(context.Context) error { return stub.err }

type rosterStub struct{}

func (rosterStub) AuthorizedRoster(context.Context, string, string, string) (domain.Roster, error) {
	return domain.Roster{}, nil
}

func TestHealthEndpointsReturnOK(t *testing.T) {
	server := NewServer(":8080", slog.New(slog.NewTextHandler(io.Discard, nil)), readyStub{}, rosterStub{})

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
	server := NewServer(":8080", slog.New(slog.NewTextHandler(io.Discard, nil)), readyStub{}, rosterStub{})

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

func TestReadinessReturnsServiceUnavailableWhenDatabaseCannotBeReached(t *testing.T) {
	server := NewServer(":8080", slog.New(slog.NewTextHandler(io.Discard, nil)), readyStub{err: errors.New("database unavailable")}, rosterStub{})
	request := httptest.NewRequest(nethttp.MethodGet, "/readyz", nil)
	response := httptest.NewRecorder()

	server.Handler.ServeHTTP(response, request)

	if response.Code != nethttp.StatusServiceUnavailable {
		t.Errorf("status = %d, want %d", response.Code, nethttp.StatusServiceUnavailable)
	}
}
