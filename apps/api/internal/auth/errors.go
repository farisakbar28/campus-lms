package auth

import (
	"errors"

	"github.com/farisakbar28/campus-lms/apps/api/internal/database"
)

var (
	// ErrAuthenticationFailed is the only public-safe classification for a
	// refresh credential failure. The concrete reason remains available to
	// internal callers through errors.Is.
	ErrAuthenticationFailed    = errors.New("authentication failed")
	ErrInvalidCredential       = errors.New("invalid refresh credential")
	ErrExpiredCredential       = errors.New("expired refresh credential")
	ErrRevokedCredential       = errors.New("revoked refresh credential")
	ErrRefreshReuse            = errors.New("refresh credential reuse detected")
	ErrInvalidAccessToken      = errors.New("invalid access token")
	ErrExpiredAccessToken      = errors.New("expired access token")
	ErrMissingAccessTokenClaim = errors.New("missing access token claim")
	ErrInvalidAccessTokenClaim = errors.New("invalid access token claim")

	ErrInvalidUserID            = errors.New("invalid user id")
	ErrInvalidSessionID         = errors.New("invalid session id")
	ErrInvalidRefreshTTL        = errors.New("invalid refresh ttl")
	ErrInvalidAccessTokenSecret = errors.New("invalid access token secret")
	ErrInvalidAccessTokenTTL    = errors.New("invalid access token ttl")
	ErrInvalidAccessTokenConfig = errors.New("invalid access token configuration")
	ErrSessionNotFound          = errors.New("auth session not found")
	ErrInvalidServiceConfig     = errors.New("invalid auth service configuration")
	ErrDatabaseUnavailable      = database.ErrUnavailable
)

// AuthenticationError keeps detailed internal classification while exposing
// one generic Error string suitable for a future HTTP response.
type AuthenticationError struct {
	reason error
}

func (e *AuthenticationError) Error() string {
	return ErrAuthenticationFailed.Error()
}

func (e *AuthenticationError) Unwrap() error {
	return e.reason
}

func (e *AuthenticationError) Is(target error) bool {
	return target == ErrAuthenticationFailed || errors.Is(e.reason, target)
}

func newAuthenticationError(reason error) error {
	return &AuthenticationError{reason: reason}
}

// IsAuthenticationFailure reports the safe classification for future
// transport code without exposing the internal failure reason.
func IsAuthenticationFailure(err error) bool {
	return errors.Is(err, ErrAuthenticationFailed)
}
