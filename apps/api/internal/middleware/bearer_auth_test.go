package middleware

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/farisakbar28/campus-lms/apps/api/internal/auth"
	"github.com/google/uuid"
)

const (
	testBearerSecretOne  = "phase3b-test-signing-secret-one-32-bytes"
	testBearerSecretTwo  = "phase3b-test-signing-secret-two-32-bytes"
	wantBearerChallenge  = `Bearer realm="campus-lms-api"`
	wantUnauthorizedJSON = `{"code":"unauthenticated","message":"authentication is required"}
`
)

type fakeAccessTokenVerifier struct {
	identity auth.AuthIdentity
	err      error
	calls    int
	encoded  []string
}

func (verifier *fakeAccessTokenVerifier) Verify(encoded string) (auth.AuthIdentity, error) {
	verifier.calls++
	verifier.encoded = append(verifier.encoded, encoded)
	return verifier.identity, verifier.err
}

func TestNewBearerMiddlewareRejectsNilVerifier(t *testing.T) {
	middleware, err := NewBearerMiddleware(nil)
	if err == nil {
		t.Fatal("NewBearerMiddleware() error = nil, want an error")
	}
	if middleware != nil {
		t.Fatal("NewBearerMiddleware() middleware is non-nil after construction failure")
	}
}

func TestBearerMiddlewareHeaderContract(t *testing.T) {
	identity := testAuthIdentity()
	tests := []struct {
		name           string
		values         []string
		wantStatus     int
		wantCalls      int
		wantDownstream int
		wantToken      string
	}{
		{name: "absent", wantStatus: http.StatusUnauthorized},
		{name: "blank", values: []string{""}, wantStatus: http.StatusUnauthorized},
		{name: "wrong scheme", values: []string{"Basic token"}, wantStatus: http.StatusUnauthorized},
		{name: "no credential", values: []string{"Bearer"}, wantStatus: http.StatusUnauthorized},
		{name: "only separator", values: []string{"Bearer   "}, wantStatus: http.StatusUnauthorized},
		{name: "case insensitive scheme", values: []string{"bEaReR token"}, wantStatus: http.StatusOK, wantCalls: 1, wantDownstream: 1, wantToken: "token"},
		{name: "one SP", values: []string{"Bearer token"}, wantStatus: http.StatusOK, wantCalls: 1, wantDownstream: 1, wantToken: "token"},
		{name: "multiple SP", values: []string{"Bearer  token"}, wantStatus: http.StatusOK, wantCalls: 1, wantDownstream: 1, wantToken: "token"},
		{name: "outer OWS", values: []string{" \tBeArEr  token\t "}, wantStatus: http.StatusOK, wantCalls: 1, wantDownstream: 1, wantToken: "token"},
		{name: "token68 punctuation", values: []string{"Bearer a0-._~+/"}, wantStatus: http.StatusOK, wantCalls: 1, wantDownstream: 1, wantToken: "a0-._~+/"},
		{name: "trailing padding", values: []string{"Bearer token=="}, wantStatus: http.StatusOK, wantCalls: 1, wantDownstream: 1, wantToken: "token=="},
		{name: "HTAB separator", values: []string{"Bearer\ttoken"}, wantStatus: http.StatusUnauthorized},
		{name: "internal credential whitespace", values: []string{"Bearer to ken"}, wantStatus: http.StatusUnauthorized},
		{name: "extra field", values: []string{"Bearer token extra"}, wantStatus: http.StatusUnauthorized},
		{name: "non-trailing padding", values: []string{"Bearer token=extra"}, wantStatus: http.StatusUnauthorized},
		{name: "duplicate field values", values: []string{"Bearer token", "Bearer token"}, wantStatus: http.StatusUnauthorized},
		{name: "combined values", values: []string{"Bearer token,Bearer other"}, wantStatus: http.StatusUnauthorized},
		{name: "invalid character", values: []string{"Bearer token#"}, wantStatus: http.StatusUnauthorized},
		{name: "valid syntax rejected by verifier", values: []string{"Bearer not-a-jwt"}, wantStatus: http.StatusUnauthorized, wantCalls: 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			verifier := &fakeAccessTokenVerifier{identity: identity, err: errors.New("fake verification failure")}
			if test.wantStatus == http.StatusOK {
				verifier.err = nil
			}
			bearer, err := NewBearerMiddleware(verifier)
			if err != nil {
				t.Fatalf("NewBearerMiddleware() error = %v", err)
			}

			downstreamCalls := 0
			handler := bearer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				downstreamCalls++
				_, _ = io.WriteString(writer, "ok")
			}))
			request := httptest.NewRequest(http.MethodGet, "/resource", nil)
			for _, value := range test.values {
				request.Header.Add("Authorization", value)
			}
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if verifier.calls != test.wantCalls {
				t.Fatalf("verifier calls = %d, want %d", verifier.calls, test.wantCalls)
			}
			if downstreamCalls != test.wantDownstream {
				t.Fatalf("downstream calls = %d, want %d", downstreamCalls, test.wantDownstream)
			}
			if test.wantStatus == http.StatusUnauthorized {
				assertGenericBearerChallenge(t, response)
			}
			if test.wantToken != "" {
				if len(verifier.encoded) != 1 || verifier.encoded[0] != test.wantToken {
					t.Fatalf("verifier did not receive the parsed credential")
				}
			}
		})
	}
}

func TestBearerMiddlewareReturnsExactUnauthorizedResponse(t *testing.T) {
	verifier := &fakeAccessTokenVerifier{err: errors.New("internal verification detail")}
	bearer, err := NewBearerMiddleware(verifier)
	if err != nil {
		t.Fatalf("NewBearerMiddleware() error = %v", err)
	}

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/resource", nil)
	request.Header.Set("Authorization", "Bearer not-a-jwt")
	bearer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("downstream must not run after verifier failure")
	})).ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
	assertGenericBearerChallenge(t, response)
	if strings.Contains(response.Body.String(), "internal verification detail") || strings.Contains(response.Body.String(), "not-a-jwt") {
		t.Fatal("response exposed internal verification detail or token")
	}
}

func TestBearerMiddlewarePropagatesCanonicalIdentityWithoutPrincipal(t *testing.T) {
	trusted := testAuthIdentity()
	attacker := auth.AuthIdentity{
		UserID:    uuid.MustParse("cccccccc-cccc-4ccc-8ccc-cccccccccccc"),
		SessionID: uuid.MustParse("dddddddd-dddd-4ddd-8ddd-dddddddddddd"),
	}
	verifier := &fakeAccessTokenVerifier{identity: trusted}
	bearer, err := NewBearerMiddleware(verifier)
	if err != nil {
		t.Fatalf("NewBearerMiddleware() error = %v", err)
	}

	var gotIdentity auth.AuthIdentity
	var identityPresent bool
	var principalPresent bool
	downstreamCalls := 0
	handler := bearer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		downstreamCalls++
		gotIdentity, identityPresent = AuthIdentityFromContext(request.Context())
		_, principalPresent = PrincipalFromContext(request.Context())
		_, _ = io.WriteString(writer, "ok")
	}))
	request := httptest.NewRequest(
		http.MethodPost,
		"/resource?user_id="+attacker.UserID.String()+"&session_id="+attacker.SessionID.String()+"&tenant_id="+attacker.UserID.String(),
		strings.NewReader(`{"user_id":"`+attacker.UserID.String()+`","session_id":"`+attacker.SessionID.String()+`","tenant_id":"`+attacker.UserID.String()+`"}`),
	)
	request.Header.Set("Authorization", "Bearer client-token")
	request.Header.Set("X-User-ID", attacker.UserID.String())
	request.Header.Set("X-Session-ID", attacker.SessionID.String())
	request.Header.Set("X-Tenant-ID", attacker.UserID.String())
	request.AddCookie(&http.Cookie{Name: "user_id", Value: attacker.UserID.String()})
	request.AddCookie(&http.Cookie{Name: "session_id", Value: attacker.SessionID.String()})
	request.AddCookie(&http.Cookie{Name: "tenant_id", Value: attacker.UserID.String()})
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if downstreamCalls != 1 {
		t.Fatalf("downstream calls = %d, want 1", downstreamCalls)
	}
	if verifier.calls != 1 {
		t.Fatalf("verifier calls = %d, want 1", verifier.calls)
	}
	if !identityPresent {
		t.Fatal("context identity was absent")
	}
	if gotIdentity.UserID == attacker.UserID || gotIdentity.SessionID == attacker.SessionID {
		t.Fatal("context identity accepted a client-supplied identity")
	}
	if gotIdentity != trusted {
		t.Fatalf("context identity did not exactly match verifier output")
	}
	if principalPresent {
		t.Fatal("Bearer authentication must not create middleware.Principal")
	}
}

func TestAuthIdentityContextDistinguishesAbsentFromZeroValue(t *testing.T) {
	identity, present := AuthIdentityFromContext(httptest.NewRequest(http.MethodGet, "/", nil).Context())
	if present {
		t.Fatal("AuthIdentityFromContext() present = true for absent identity")
	}
	if identity != (auth.AuthIdentity{}) {
		t.Fatal("absent identity returned a non-zero value")
	}

	contextWithZeroIdentity := WithAuthIdentity(httptest.NewRequest(http.MethodGet, "/", nil).Context(), auth.AuthIdentity{})
	identity, present = AuthIdentityFromContext(contextWithZeroIdentity)
	if !present {
		t.Fatal("AuthIdentityFromContext() present = false for explicitly stored zero identity")
	}
	if identity != (auth.AuthIdentity{}) {
		t.Fatal("explicit zero identity was changed")
	}
}

func TestBearerMiddlewareIgnoresRequestIdentityInputsWhenAuthorizationIsAbsent(t *testing.T) {
	verifier := &fakeAccessTokenVerifier{identity: testAuthIdentity()}
	bearer, err := NewBearerMiddleware(verifier)
	if err != nil {
		t.Fatalf("NewBearerMiddleware() error = %v", err)
	}
	downstreamCalls := 0
	handler := bearer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		downstreamCalls++
	}))
	request := httptest.NewRequest(
		http.MethodPost,
		"/resource?user_id="+testAuthIdentity().UserID.String()+"&tenant_id="+testAuthIdentity().UserID.String(),
		strings.NewReader(`{"user_id":"`+testAuthIdentity().UserID.String()+`","tenant_id":"`+testAuthIdentity().UserID.String()+`"}`),
	)
	request.Header.Set("X-User-ID", testAuthIdentity().UserID.String())
	request.Header.Set("X-Tenant-ID", testAuthIdentity().UserID.String())
	request.AddCookie(&http.Cookie{Name: "user_id", Value: testAuthIdentity().UserID.String()})
	request.AddCookie(&http.Cookie{Name: "tenant_id", Value: testAuthIdentity().UserID.String()})
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
	assertGenericBearerChallenge(t, response)
	if verifier.calls != 0 {
		t.Fatalf("verifier calls = %d, want 0", verifier.calls)
	}
	if downstreamCalls != 0 {
		t.Fatalf("downstream calls = %d, want 0", downstreamCalls)
	}
}

func TestBearerMiddlewareWithAccessTokenManager(t *testing.T) {
	issuedAt := time.Date(2026, time.August, 28, 12, 0, 0, 0, time.UTC)
	currentTime := issuedAt
	manager, err := auth.NewAccessTokenManager(
		[]byte(testBearerSecretOne),
		15*time.Minute,
		auth.WithAccessTokenClock(func() time.Time { return currentTime }),
	)
	if err != nil {
		t.Fatalf("NewAccessTokenManager() error = %v", err)
	}
	userID := uuid.MustParse("11111111-1111-4111-8111-111111111111")
	sessionID := uuid.MustParse("22222222-2222-4222-8222-222222222222")
	accessToken, err := manager.Issue(userID, sessionID)
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}

	tests := []struct {
		name       string
		verifier   auth.AccessTokenVerifier
		token      string
		wantStatus int
	}{
		{name: "valid issued access JWT", verifier: manager, token: accessToken.Token, wantStatus: http.StatusOK},
		{name: "malformed JWT", verifier: manager, token: "not-a-jwt", wantStatus: http.StatusUnauthorized},
	}

	currentTime = issuedAt.Add(16 * time.Minute)
	tests = append(tests, struct {
		name       string
		verifier   auth.AccessTokenVerifier
		token      string
		wantStatus int
	}{name: "expired JWT", verifier: manager, token: accessToken.Token, wantStatus: http.StatusUnauthorized})

	alternateManager, err := auth.NewAccessTokenManager(
		[]byte(testBearerSecretTwo),
		15*time.Minute,
		auth.WithAccessTokenClock(func() time.Time { return issuedAt }),
	)
	if err != nil {
		t.Fatalf("NewAccessTokenManager() alternate error = %v", err)
	}
	alternateToken, err := alternateManager.Issue(userID, sessionID)
	if err != nil {
		t.Fatalf("Issue() alternate error = %v", err)
	}
	tests = append(tests, struct {
		name       string
		verifier   auth.AccessTokenVerifier
		token      string
		wantStatus int
	}{name: "alternate signing key", verifier: manager, token: alternateToken.Token, wantStatus: http.StatusUnauthorized})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.name == "expired JWT" {
				currentTime = issuedAt.Add(16 * time.Minute)
			} else {
				currentTime = issuedAt
			}
			downstreamCalls := 0
			bearer, err := NewBearerMiddleware(test.verifier)
			if err != nil {
				t.Fatalf("NewBearerMiddleware() error = %v", err)
			}
			handler := bearer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				downstreamCalls++
				_, _ = io.WriteString(writer, "ok")
			}))
			request := httptest.NewRequest(http.MethodGet, "/resource", nil)
			request.Header.Set("Authorization", "Bearer "+test.token)
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if downstreamCalls != 1 && test.wantStatus == http.StatusOK {
				t.Fatalf("downstream calls = %d, want 1", downstreamCalls)
			}
			if downstreamCalls != 0 && test.wantStatus == http.StatusUnauthorized {
				t.Fatalf("downstream calls = %d, want 0", downstreamCalls)
			}
			if test.wantStatus == http.StatusUnauthorized {
				assertGenericBearerChallenge(t, response)
				if strings.Contains(response.Body.String(), test.token) || strings.Contains(response.Body.String(), testBearerSecretOne) || strings.Contains(response.Body.String(), testBearerSecretTwo) {
					t.Fatal("response exposed token or signing material")
				}
			}
		})
	}
}

func assertGenericBearerChallenge(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()

	if response.Header().Get("WWW-Authenticate") != wantBearerChallenge {
		t.Fatalf("WWW-Authenticate = %q, want %q", response.Header().Get("WWW-Authenticate"), wantBearerChallenge)
	}
	if response.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", response.Header().Get("Content-Type"))
	}
	if response.Body.String() != wantUnauthorizedJSON {
		t.Fatalf("response body = %q, want exact generic JSON", response.Body.String())
	}
}

func testAuthIdentity() auth.AuthIdentity {
	return auth.AuthIdentity{
		UserID:    uuid.MustParse("aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"),
		SessionID: uuid.MustParse("bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"),
	}
}
