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

	"github.com/farisakbar28/campus-lms/apps/api/internal/auth"
	"github.com/farisakbar28/campus-lms/apps/api/internal/domain"
	"github.com/farisakbar28/campus-lms/apps/api/internal/repository"
	"github.com/google/uuid"
)

type readyStub struct{ err error }

func (stub readyStub) Ping(context.Context) error { return stub.err }

type rosterStub struct{}

func (rosterStub) AuthorizedRoster(context.Context, string, string, string) (domain.Roster, error) {
	return domain.Roster{}, nil
}

type verifierStub struct{}

func (verifierStub) Verify(string) (auth.AuthIdentity, error) {
	return auth.AuthIdentity{}, errors.New("invalid test token")
}

type admitterStub struct{}

func (admitterStub) Admit(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) (uuid.UUID, error) {
	return uuid.Nil, errors.New("admission is unavailable")
}

type acceptingVerifier struct {
	identity auth.AuthIdentity
	err      error
}

func (stub acceptingVerifier) Verify(string) (auth.AuthIdentity, error) {
	return stub.identity, stub.err
}

type acceptingAdmitter struct {
	membershipID uuid.UUID
	err          error
	userID       uuid.UUID
	sessionID    uuid.UUID
	tenantID     uuid.UUID
	calls        int
}

func (stub *acceptingAdmitter) Admit(_ context.Context, userID, sessionID, tenantID uuid.UUID) (uuid.UUID, error) {
	stub.calls++
	stub.userID = userID
	stub.sessionID = sessionID
	stub.tenantID = tenantID
	return stub.membershipID, stub.err
}

type capturingRoster struct {
	roster     domain.Roster
	userID     string
	tenantID   string
	offeringID string
	calls      int
}

func (stub *capturingRoster) AuthorizedRoster(_ context.Context, tenantID, userID, offeringID string) (domain.Roster, error) {
	stub.calls++
	stub.tenantID = tenantID
	stub.userID = userID
	stub.offeringID = offeringID
	return stub.roster, nil
}

func newServerForTest(t *testing.T, readiness readinessChecker) *nethttp.Server {
	t.Helper()
	server, err := NewServer(":8080", slog.New(slog.NewTextHandler(io.Discard, nil)), readiness, rosterStub{}, verifierStub{}, admitterStub{})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}
	return server
}

func TestHealthEndpointsReturnOK(t *testing.T) {
	server := newServerForTest(t, readyStub{})

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
	server := newServerForTest(t, readyStub{})

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
	server := newServerForTest(t, readyStub{err: errors.New("database unavailable")})
	request := httptest.NewRequest(nethttp.MethodGet, "/readyz", nil)
	response := httptest.NewRecorder()

	server.Handler.ServeHTTP(response, request)

	if response.Code != nethttp.StatusServiceUnavailable {
		t.Errorf("status = %d, want %d", response.Code, nethttp.StatusServiceUnavailable)
	}
}

func TestServerFailsClosedWithoutAuthenticationDependencies(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	if server, err := NewServer(":8080", logger, readyStub{}, rosterStub{}, nil, admitterStub{}); err == nil || server != nil {
		t.Fatal("NewServer() accepted a missing verifier")
	}
	if server, err := NewServer(":8080", logger, readyStub{}, rosterStub{}, verifierStub{}, nil); err == nil || server != nil {
		t.Fatal("NewServer() accepted a missing admission service")
	}
}

func TestLegacyRosterRouteIsNotRegistered(t *testing.T) {
	server := newServerForTest(t, readyStub{})
	request := httptest.NewRequest(nethttp.MethodGet, "/course-offerings/2888c021-06ae-73da-2f57-884c1dd5d059/participants", nil)
	response := httptest.NewRecorder()
	server.Handler.ServeHTTP(response, request)
	if response.Code != nethttp.StatusNotFound {
		t.Fatalf("legacy route status = %d, want 404", response.Code)
	}
}

func TestProtectedRosterRouteComposesTokenTenantAdmissionAndPrincipal(t *testing.T) {
	userID := uuid.MustParse("4ec42919-bc1a-17bc-10b0-d75b8343dff8")
	sessionID := uuid.MustParse("7bd7c5f5-f2ec-48cd-9259-c9f7e3f20510")
	tenantID := uuid.MustParse("19cd4773-2aeb-d614-028f-e21bf9b73d0c")
	membershipID := uuid.MustParse("1ec2ad6e-a42a-4bc8-bf6a-c0a3f79c5cf2")
	offeringID := "2888c021-06ae-73da-2f57-884c1dd5d059"
	roster := &capturingRoster{roster: domain.Roster{Offering: domain.CourseOffering{ID: offeringID}}}
	admitter := &acceptingAdmitter{membershipID: membershipID}
	server, err := NewServer(
		":8080",
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		readyStub{}, roster,
		acceptingVerifier{identity: auth.AuthIdentity{UserID: userID, SessionID: sessionID}}, admitter,
	)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	request := httptest.NewRequest(nethttp.MethodGet, "/tenants/"+tenantID.String()+"/course-offerings/"+offeringID+"/participants", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()
	server.Handler.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", response.Code, response.Body.String())
	}
	if admitter.calls != 1 || admitter.userID != userID || admitter.sessionID != sessionID || admitter.tenantID != tenantID {
		t.Fatalf("admission = calls=%d user=%s session=%s tenant=%s", admitter.calls, admitter.userID, admitter.sessionID, admitter.tenantID)
	}
	if roster.calls != 1 || roster.tenantID != tenantID.String() || roster.userID != userID.String() || roster.offeringID != offeringID {
		t.Fatalf("roster = calls=%d tenant=%s user=%s offering=%s", roster.calls, roster.tenantID, roster.userID, roster.offeringID)
	}
}

func TestProtectedRosterRouteFailsBeforeAdmissionForInvalidCredentialsOrTenant(t *testing.T) {
	userID := uuid.MustParse("4ec42919-bc1a-17bc-10b0-d75b8343dff8")
	sessionID := uuid.MustParse("7bd7c5f5-f2ec-48cd-9259-c9f7e3f20510")
	admitter := &acceptingAdmitter{membershipID: uuid.New()}
	newServer := func(verifier acceptingVerifier) *nethttp.Server {
		server, err := NewServer(":8080", slog.New(slog.NewTextHandler(io.Discard, nil)), readyStub{}, &capturingRoster{}, verifier, admitter)
		if err != nil {
			t.Fatalf("NewServer() error = %v", err)
		}
		return server
	}

	t.Run("invalid credentials", func(t *testing.T) {
		server := newServer(acceptingVerifier{err: errors.New("invalid token")})
		request := httptest.NewRequest(nethttp.MethodGet, "/tenants/19cd4773-2aeb-d614-028f-e21bf9b73d0c/course-offerings/2888c021-06ae-73da-2f57-884c1dd5d059/participants", nil)
		request.Header.Set("Authorization", "Bearer test-token")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		if response.Code != nethttp.StatusUnauthorized || admitter.calls != 0 {
			t.Fatalf("status=%d admission_calls=%d, want 401 and 0", response.Code, admitter.calls)
		}
	})

	t.Run("malformed tenant selector", func(t *testing.T) {
		server := newServer(acceptingVerifier{identity: auth.AuthIdentity{UserID: userID, SessionID: sessionID}})
		request := httptest.NewRequest(nethttp.MethodGet, "/tenants/not-a-uuid/course-offerings/2888c021-06ae-73da-2f57-884c1dd5d059/participants", nil)
		request.Header.Set("Authorization", "Bearer test-token")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		if response.Code != nethttp.StatusBadRequest || admitter.calls != 0 {
			t.Fatalf("status=%d admission_calls=%d, want 400 and 0", response.Code, admitter.calls)
		}
	})
}

func TestProtectedRosterRouteMapsUnknownTenantWithoutLeakingState(t *testing.T) {
	admitter := &acceptingAdmitter{err: repository.ErrTenantAdmission}
	server, err := NewServer(
		":8080", slog.New(slog.NewTextHandler(io.Discard, nil)), readyStub{}, &capturingRoster{},
		acceptingVerifier{identity: auth.AuthIdentity{UserID: uuid.New(), SessionID: uuid.New()}}, admitter,
	)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}
	request := httptest.NewRequest(nethttp.MethodGet, "/tenants/19cd4773-2aeb-d614-028f-e21bf9b73d0c/course-offerings/2888c021-06ae-73da-2f57-884c1dd5d059/participants", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()
	server.Handler.ServeHTTP(response, request)
	if response.Code != nethttp.StatusNotFound || strings.Contains(response.Body.String(), "membership") {
		t.Fatalf("status=%d body=%q, want indistinguishable 404", response.Code, response.Body.String())
	}
}
