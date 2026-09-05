package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const accessTokenTestSecret = "01234567890123456789012345678901"

func TestAccessTokenIssueAndVerifyRoundTrip(t *testing.T) {
	now := time.Date(2026, 8, 28, 19, 0, 0, 0, time.UTC)
	userID := uuid.MustParse("11111111-1111-4111-8111-111111111111")
	sessionID := uuid.MustParse("22222222-2222-4222-8222-222222222222")
	manager := newTestAccessTokenManager(t, now, 15*time.Minute)

	issued, err := manager.Issue(userID, sessionID)
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}
	if issued.Token == "" {
		t.Fatal("Issue() returned an empty token")
	}
	if !issued.IssuedAt.Equal(now) {
		t.Errorf("IssuedAt = %s, want %s", issued.IssuedAt, now)
	}
	if !issued.ExpiresAt.Equal(now.Add(15 * time.Minute)) {
		t.Errorf("ExpiresAt = %s, want %s", issued.ExpiresAt, now.Add(15*time.Minute))
	}

	identity, err := manager.Verify(issued.Token)
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if identity.UserID != userID || identity.SessionID != sessionID {
		t.Fatalf("identity = %+v, want user %s session %s", identity, userID, sessionID)
	}

	t.Logf("round_trip=true access_ttl=%s identity_claims=user_id_session_id", 15*time.Minute)
}

func TestAccessTokenIssueUsesNumericDateSecondPrecision(t *testing.T) {
	clockNow := time.Date(2026, 8, 28, 19, 0, 0, 123456789, time.FixedZone("test", 7*60*60))
	wantIssuedAt := clockNow.UTC().Truncate(time.Second)
	manager := newTestAccessTokenManager(t, clockNow, 15*time.Minute)

	issued, err := manager.Issue(testAccessTokenUserID(), testAccessTokenSessionID())
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}
	if !issued.IssuedAt.Equal(wantIssuedAt) {
		t.Errorf("IssuedAt = %s, want second-precision %s", issued.IssuedAt, wantIssuedAt)
	}
	if !issued.ExpiresAt.Equal(wantIssuedAt.Add(15 * time.Minute)) {
		t.Errorf("ExpiresAt = %s, want %s", issued.ExpiresAt, wantIssuedAt.Add(15*time.Minute))
	}
}

func TestAccessTokenIssueRejectsNilIdentity(t *testing.T) {
	manager := newTestAccessTokenManager(t, time.Date(2026, 8, 28, 19, 0, 0, 0, time.UTC), 15*time.Minute)

	tests := []struct {
		name       string
		userID     uuid.UUID
		sessionID  uuid.UUID
		wantReason error
	}{
		{name: "nil user", userID: uuid.Nil, sessionID: testAccessTokenSessionID(), wantReason: ErrInvalidUserID},
		{name: "nil session", userID: testAccessTokenUserID(), sessionID: uuid.Nil, wantReason: ErrInvalidSessionID},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			issued, err := manager.Issue(test.userID, test.sessionID)
			if !errors.Is(err, test.wantReason) {
				t.Fatalf("Issue() error = %v, want %v", err, test.wantReason)
			}
			if issued.Token != "" || !issued.IssuedAt.IsZero() || !issued.ExpiresAt.IsZero() {
				t.Fatalf("Issue() returned output for invalid identity: %+v", issued)
			}
		})
	}
}

func TestAccessTokenExpiryBoundaries(t *testing.T) {
	now := time.Date(2026, 8, 28, 19, 0, 0, 0, time.UTC)
	currentTime := now
	manager, err := NewAccessTokenManager(
		[]byte(accessTokenTestSecret),
		15*time.Minute,
		WithAccessTokenClock(func() time.Time { return currentTime }),
	)
	if err != nil {
		t.Fatalf("NewAccessTokenManager() error = %v", err)
	}
	issued, err := manager.Issue(testAccessTokenUserID(), testAccessTokenSessionID())
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}

	tests := []struct {
		name       string
		verifyTime time.Time
		wantReason error
	}{
		{name: "before expiration", verifyTime: issued.ExpiresAt.Add(-time.Second)},
		{name: "at expiration", verifyTime: issued.ExpiresAt, wantReason: ErrExpiredAccessToken},
		{name: "after expiration", verifyTime: issued.ExpiresAt.Add(time.Second), wantReason: ErrExpiredAccessToken},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			currentTime = test.verifyTime
			identity, verifyErr := manager.Verify(issued.Token)
			if test.wantReason == nil {
				if verifyErr != nil {
					t.Fatalf("Verify() error = %v, want nil", verifyErr)
				}
				if identity.UserID != testAccessTokenUserID() || identity.SessionID != testAccessTokenSessionID() {
					t.Fatalf("identity = %+v, want issued identity", identity)
				}
				return
			}

			requireGenericAccessTokenFailure(t, verifyErr, test.wantReason)
		})
	}

	t.Log("expiry_semantics=now_greater_than_or_equal_to_exp_invalid positive_leeway=false")
}

func TestAccessTokenRejectsInvalidSignature(t *testing.T) {
	manager := newTestAccessTokenManager(t, testAccessTokenNow(), 15*time.Minute)
	issued, err := manager.Issue(testAccessTokenUserID(), testAccessTokenSessionID())
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}

	tests := []struct {
		name  string
		token string
	}{
		{name: "tampered signature", token: tamperAccessTokenSignature(issued.Token)},
		{name: "another hs256 key", token: issueAccessTokenWithKey(t, jwt.SigningMethodHS256, []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), validAccessTokenClaims(testAccessTokenNow()))},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, verifyErr := manager.Verify(test.token)
			requireGenericAccessTokenFailure(t, verifyErr)
		})
	}
}

func TestAccessTokenRejectsWrongAlgorithm(t *testing.T) {
	manager := newTestAccessTokenManager(t, testAccessTokenNow(), 15*time.Minute)
	token := issueAccessTokenWithKey(t, jwt.SigningMethodHS384, []byte(accessTokenTestSecret), validAccessTokenClaims(testAccessTokenNow()))

	_, verifyErr := manager.Verify(token)
	requireGenericAccessTokenFailure(t, verifyErr, ErrInvalidAccessToken)
}

func TestAccessTokenRejectsNoneAlgorithm(t *testing.T) {
	manager := newTestAccessTokenManager(t, testAccessTokenNow(), 15*time.Minute)
	token := issueAccessTokenWithKey(t, jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType, validAccessTokenClaims(testAccessTokenNow()))

	_, verifyErr := manager.Verify(token)
	requireGenericAccessTokenFailure(t, verifyErr, ErrInvalidAccessToken)
}

func TestAccessTokenRejectsWrongIssuerAndAudience(t *testing.T) {
	manager := newTestAccessTokenManager(t, testAccessTokenNow(), 15*time.Minute)

	tests := []struct {
		name   string
		mutate func(*AccessTokenClaims)
	}{
		{name: "wrong issuer", mutate: func(claims *AccessTokenClaims) { claims.Issuer = "another-service" }},
		{name: "wrong audience", mutate: func(claims *AccessTokenClaims) { claims.Audience = jwt.ClaimStrings{"another-service"} }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			claims := validAccessTokenClaims(testAccessTokenNow())
			test.mutate(&claims)
			token := issueAccessTokenWithKey(t, jwt.SigningMethodHS256, []byte(accessTokenTestSecret), claims)
			_, verifyErr := manager.Verify(token)
			requireGenericAccessTokenFailure(t, verifyErr, ErrInvalidAccessToken)
		})
	}
}

func TestAccessTokenRejectsMissingRequiredClaims(t *testing.T) {
	manager := newTestAccessTokenManager(t, testAccessTokenNow(), 15*time.Minute)

	tests := []struct {
		name   string
		mutate func(*AccessTokenClaims)
	}{
		{name: "missing sub", mutate: func(claims *AccessTokenClaims) { claims.Subject = "" }},
		{name: "missing sid", mutate: func(claims *AccessTokenClaims) { claims.SessionID = "" }},
		{name: "missing iat", mutate: func(claims *AccessTokenClaims) { claims.IssuedAt = nil }},
		{name: "missing exp", mutate: func(claims *AccessTokenClaims) { claims.ExpiresAt = nil }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			claims := validAccessTokenClaims(testAccessTokenNow())
			test.mutate(&claims)
			token := issueAccessTokenWithKey(t, jwt.SigningMethodHS256, []byte(accessTokenTestSecret), claims)
			_, verifyErr := manager.Verify(token)
			requireGenericAccessTokenFailure(t, verifyErr, ErrMissingAccessTokenClaim)
		})
	}
}

func TestAccessTokenRejectsMalformedUUIDClaims(t *testing.T) {
	manager := newTestAccessTokenManager(t, testAccessTokenNow(), 15*time.Minute)

	tests := []struct {
		name   string
		mutate func(*AccessTokenClaims)
	}{
		{name: "malformed sub", mutate: func(claims *AccessTokenClaims) { claims.Subject = "not-a-uuid" }},
		{name: "malformed sid", mutate: func(claims *AccessTokenClaims) { claims.SessionID = "not-a-uuid" }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			claims := validAccessTokenClaims(testAccessTokenNow())
			test.mutate(&claims)
			token := issueAccessTokenWithKey(t, jwt.SigningMethodHS256, []byte(accessTokenTestSecret), claims)
			_, verifyErr := manager.Verify(token)
			requireGenericAccessTokenFailure(t, verifyErr, ErrInvalidAccessTokenClaim)
		})
	}
}

func TestAccessTokenRejectsNilUUIDClaims(t *testing.T) {
	manager := newTestAccessTokenManager(t, testAccessTokenNow(), 15*time.Minute)

	tests := []struct {
		name   string
		mutate func(*AccessTokenClaims)
	}{
		{name: "nil sub", mutate: func(claims *AccessTokenClaims) { claims.Subject = uuid.Nil.String() }},
		{name: "nil sid", mutate: func(claims *AccessTokenClaims) { claims.SessionID = uuid.Nil.String() }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			claims := validAccessTokenClaims(testAccessTokenNow())
			test.mutate(&claims)
			token := issueAccessTokenWithKey(t, jwt.SigningMethodHS256, []byte(accessTokenTestSecret), claims)
			_, verifyErr := manager.Verify(token)
			requireGenericAccessTokenFailure(t, verifyErr, ErrInvalidAccessTokenClaim)
		})
	}
}

func TestAccessTokenClaimsAreMinimized(t *testing.T) {
	manager := newTestAccessTokenManager(t, testAccessTokenNow(), 15*time.Minute)
	issued, err := manager.Issue(testAccessTokenUserID(), testAccessTokenSessionID())
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}

	parts := strings.Split(issued.Token, ".")
	if len(parts) != 3 {
		t.Fatalf("issued token segments = %d, want 3", len(parts))
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("decode test payload: %v", err)
	}
	var claims map[string]json.RawMessage
	if err := json.Unmarshal(payload, &claims); err != nil {
		t.Fatalf("decode test claims: %v", err)
	}

	for _, required := range []string{"sub", "sid", "iat", "exp", "iss", "aud"} {
		if _, ok := claims[required]; !ok {
			t.Errorf("issued claims missing required key %q", required)
		}
	}
	for _, forbidden := range []string{
		"tenant_id", "tenant_role", "membership_role", "course_role", "role", "roles",
		"permissions", "email", "display_name", "profile", "refresh_token", "refresh_token_hash",
	} {
		if _, ok := claims[forbidden]; ok {
			t.Errorf("issued claims contain forbidden key %q", forbidden)
		}
	}

	t.Log("claims=iss_sub_sid_aud_iat_exp authorization_and_refresh_data=false")
}

func TestAccessTokenErrorsDoNotExposeSecretOrToken(t *testing.T) {
	manager := newTestAccessTokenManager(t, testAccessTokenNow(), 15*time.Minute)
	issued, err := manager.Issue(testAccessTokenUserID(), testAccessTokenSessionID())
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}

	_, verifyErr := manager.Verify(tamperAccessTokenSignature(issued.Token))
	if verifyErr == nil {
		t.Fatal("Verify() error = nil, want invalid signature error")
	}
	if verifyErr.Error() != ErrAuthenticationFailed.Error() {
		t.Fatalf("Verify() error = %q, want generic authentication failure", verifyErr)
	}
	if strings.Contains(verifyErr.Error(), accessTokenTestSecret) || strings.Contains(verifyErr.Error(), issued.Token) {
		t.Fatal("Verify() error exposed secret or token")
	}

	shortSecret := "secret-value"
	_, constructionErr := NewAccessTokenManager([]byte(shortSecret), 15*time.Minute)
	if constructionErr == nil || strings.Contains(constructionErr.Error(), shortSecret) {
		t.Fatalf("constructor error = %v, secret value must not be exposed", constructionErr)
	}
}

func TestNewAccessTokenManagerRejectsInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name       string
		secret     []byte
		accessTTL  time.Duration
		wantReason error
	}{
		{name: "missing secret", secret: nil, accessTTL: 15 * time.Minute, wantReason: ErrInvalidAccessTokenSecret},
		{name: "blank secret", secret: []byte("                                "), accessTTL: 15 * time.Minute, wantReason: ErrInvalidAccessTokenSecret},
		{name: "short secret", secret: []byte("short"), accessTTL: 15 * time.Minute, wantReason: ErrInvalidAccessTokenSecret},
		{name: "zero ttl", secret: []byte(accessTokenTestSecret), accessTTL: 0, wantReason: ErrInvalidAccessTokenTTL},
		{name: "negative ttl", secret: []byte(accessTokenTestSecret), accessTTL: -time.Second, wantReason: ErrInvalidAccessTokenTTL},
		{name: "over maximum ttl", secret: []byte(accessTokenTestSecret), accessTTL: time.Hour + time.Second, wantReason: ErrInvalidAccessTokenTTL},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewAccessTokenManager(test.secret, test.accessTTL)
			if !errors.Is(err, test.wantReason) {
				t.Fatalf("NewAccessTokenManager() error = %v, want %v", err, test.wantReason)
			}
		})
	}
}

func newTestAccessTokenManager(t *testing.T, now time.Time, accessTTL time.Duration) *AccessTokenManager {
	t.Helper()
	currentTime := now
	manager, err := NewAccessTokenManager(
		[]byte(accessTokenTestSecret),
		accessTTL,
		WithAccessTokenClock(func() time.Time { return currentTime }),
	)
	if err != nil {
		t.Fatalf("NewAccessTokenManager() error = %v", err)
	}
	return manager
}

func testAccessTokenNow() time.Time {
	return time.Date(2026, 8, 28, 19, 0, 0, 0, time.UTC)
}

func testAccessTokenUserID() uuid.UUID {
	return uuid.MustParse("11111111-1111-4111-8111-111111111111")
}

func testAccessTokenSessionID() uuid.UUID {
	return uuid.MustParse("22222222-2222-4222-8222-222222222222")
}

func validAccessTokenClaims(now time.Time) AccessTokenClaims {
	issuedAt := now.UTC().Truncate(time.Second)
	return AccessTokenClaims{
		SessionID: testAccessTokenSessionID().String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    AccessTokenIssuerName,
			Subject:   testAccessTokenUserID().String(),
			Audience:  jwt.ClaimStrings{AccessTokenAudienceName},
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(15 * time.Minute)),
		},
	}
}

func issueAccessTokenWithKey(t *testing.T, method jwt.SigningMethod, key any, claims AccessTokenClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(method, claims)
	encoded, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("sign test access token: %v", err)
	}
	return encoded
}

func tamperAccessTokenSignature(encoded string) string {
	last := encoded[len(encoded)-1]
	if last == 'A' {
		return encoded[:len(encoded)-1] + "B"
	}
	return encoded[:len(encoded)-1] + "A"
}

func requireGenericAccessTokenFailure(t *testing.T, err error, reasons ...error) {
	t.Helper()
	if err == nil {
		t.Fatal("Verify() error = nil, want authentication failure")
	}
	if err.Error() != ErrAuthenticationFailed.Error() {
		t.Fatalf("Verify() error = %q, want generic authentication failure", err)
	}
	if !errors.Is(err, ErrAuthenticationFailed) {
		t.Fatal("Verify() error does not match ErrAuthenticationFailed")
	}
	for _, reason := range reasons {
		if !errors.Is(err, reason) {
			t.Errorf("Verify() error = %v, want internal reason %v", err, reason)
		}
	}
}
