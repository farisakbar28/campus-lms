package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	// AccessTokenIssuerName identifies this API as the access-token issuer.
	AccessTokenIssuerName = "campus-lms-api"
	// AccessTokenAudienceName identifies the API that accepts these tokens.
	AccessTokenAudienceName = "campus-lms-api"

	minimumAccessTokenSecretBytes = 32
	maximumAccessTokenTTL         = time.Hour
)

// AuthIdentity is the trusted identity carried across the authentication
// boundary. Tenant authorization is resolved separately from this identity.
type AuthIdentity struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
}

// AccessTokenClaims is the complete claim set issued by this API. It contains
// identity and protocol claims only; tenant and authorization data are not
// token claims.
type AccessTokenClaims struct {
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}

// AccessToken is the result of issuing one short-lived access token.
type AccessToken struct {
	Token     string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// AccessTokenIssuer issues access tokens for an authenticated user/session.
type AccessTokenIssuer interface {
	Issue(userID, sessionID uuid.UUID) (AccessToken, error)
}

// AccessTokenVerifier verifies access tokens and returns their trusted
// authentication identity.
type AccessTokenVerifier interface {
	Verify(encoded string) (AuthIdentity, error)
}

// AccessTokenOption customizes construction of an AccessTokenManager.
type AccessTokenOption func(*AccessTokenManager)

// WithAccessTokenClock replaces the wall clock for deterministic tests.
func WithAccessTokenClock(clock Clock) AccessTokenOption {
	return func(manager *AccessTokenManager) {
		if clock != nil {
			manager.clock = clock
		}
	}
}

// AccessTokenManager signs and verifies short-lived HS256 access tokens.
// Configuration is copied and validated at construction so issuance does not
// reread environment variables or depend on database/network state.
type AccessTokenManager struct {
	signingKey []byte
	accessTTL  time.Duration
	clock      Clock
}

// NewAccessTokenManager constructs an access-token issuer/verifier from
// validated key material and a bounded access-token lifetime.
func NewAccessTokenManager(signingKey []byte, accessTTL time.Duration, options ...AccessTokenOption) (*AccessTokenManager, error) {
	manager := &AccessTokenManager{
		signingKey: append([]byte(nil), signingKey...),
		accessTTL:  accessTTL,
		clock:      time.Now,
	}
	for _, option := range options {
		if option != nil {
			option(manager)
		}
	}

	if err := manager.validate(); err != nil {
		return nil, err
	}

	return manager, nil
}

// Issue creates a signed access token for one user and auth session.
func (manager *AccessTokenManager) Issue(userID, sessionID uuid.UUID) (AccessToken, error) {
	if err := manager.validate(); err != nil {
		return AccessToken{}, err
	}
	if userID == uuid.Nil {
		return AccessToken{}, ErrInvalidUserID
	}
	if sessionID == uuid.Nil {
		return AccessToken{}, ErrInvalidSessionID
	}

	issuedAt := manager.clock().UTC().Truncate(time.Second)
	expiresAt := issuedAt.Add(manager.accessTTL)
	claims := AccessTokenClaims{
		SessionID: sessionID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    AccessTokenIssuerName,
			Subject:   userID.String(),
			Audience:  jwt.ClaimStrings{AccessTokenAudienceName},
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encoded, err := token.SignedString(manager.signingKey)
	if err != nil {
		return AccessToken{}, ErrInvalidAccessTokenConfig
	}

	return AccessToken{
		Token:     encoded,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}, nil
}

// Verify authenticates one encoded access token and returns only its user and
// auth-session identity. All token parsing failures use a safe generic error.
func (manager *AccessTokenManager) Verify(encoded string) (AuthIdentity, error) {
	if err := manager.validate(); err != nil {
		return AuthIdentity{}, err
	}
	if strings.TrimSpace(encoded) == "" {
		return AuthIdentity{}, newAuthenticationError(ErrInvalidAccessToken)
	}

	var claims AccessTokenClaims
	parsed, err := jwt.ParseWithClaims(
		encoded,
		&claims,
		func(token *jwt.Token) (any, error) {
			if token.Method == nil || token.Method != jwt.SigningMethodHS256 {
				return nil, ErrInvalidAccessToken
			}
			return manager.signingKey, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithIssuer(AccessTokenIssuerName),
		jwt.WithAudience(AccessTokenAudienceName),
		jwt.WithTimeFunc(manager.clock),
	)
	if err != nil {
		return AuthIdentity{}, classifyAccessTokenError(err)
	}
	if parsed == nil || !parsed.Valid {
		return AuthIdentity{}, newAuthenticationError(ErrInvalidAccessToken)
	}

	userID, err := parseAccessTokenIdentityUUID(claims.Subject)
	if err != nil {
		return AuthIdentity{}, newAuthenticationError(err)
	}
	sessionID, err := parseAccessTokenIdentityUUID(claims.SessionID)
	if err != nil {
		return AuthIdentity{}, newAuthenticationError(err)
	}

	return AuthIdentity{UserID: userID, SessionID: sessionID}, nil
}

// Validate enforces the application-specific claims required for an access
// token. jwt.WithIssuedAt validates iat when present, so presence is checked
// here explicitly along with the custom identity claims.
func (claims AccessTokenClaims) Validate() error {
	switch {
	case strings.TrimSpace(claims.Subject) == "":
		return ErrMissingAccessTokenClaim
	case strings.TrimSpace(claims.SessionID) == "":
		return ErrMissingAccessTokenClaim
	case claims.IssuedAt == nil:
		return ErrMissingAccessTokenClaim
	case claims.ExpiresAt == nil:
		return ErrMissingAccessTokenClaim
	}
	if _, err := parseAccessTokenIdentityUUID(claims.Subject); err != nil {
		return err
	}
	if _, err := parseAccessTokenIdentityUUID(claims.SessionID); err != nil {
		return err
	}

	return nil
}

func (manager *AccessTokenManager) validate() error {
	if manager == nil || len(manager.signingKey) < minimumAccessTokenSecretBytes || strings.TrimSpace(string(manager.signingKey)) == "" {
		return ErrInvalidAccessTokenSecret
	}
	if manager.accessTTL <= 0 || manager.accessTTL > maximumAccessTokenTTL {
		return ErrInvalidAccessTokenTTL
	}
	if manager.clock == nil {
		return ErrInvalidAccessTokenConfig
	}

	return nil
}

func parseAccessTokenIdentityUUID(value string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil || parsed == uuid.Nil {
		return uuid.Nil, ErrInvalidAccessTokenClaim
	}

	return parsed, nil
}

func classifyAccessTokenError(err error) error {
	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return newAuthenticationError(ErrExpiredAccessToken)
	case errors.Is(err, jwt.ErrTokenRequiredClaimMissing), errors.Is(err, ErrMissingAccessTokenClaim):
		return newAuthenticationError(ErrMissingAccessTokenClaim)
	case errors.Is(err, ErrInvalidAccessTokenClaim):
		return newAuthenticationError(ErrInvalidAccessTokenClaim)
	default:
		return newAuthenticationError(ErrInvalidAccessToken)
	}
}
