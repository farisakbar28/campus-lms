package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

const RefreshCredentialBytes = 32

// GenerateRefreshCredential creates a production refresh credential using the
// standard library's cryptographically secure random source.
func GenerateRefreshCredential() (string, error) {
	return GenerateRefreshCredentialFrom(rand.Reader)
}

// GenerateRefreshCredentialFrom is injectable so tests can control the byte
// source without making production randomness deterministic.
func GenerateRefreshCredentialFrom(source io.Reader) (string, error) {
	credential, _, err := generateRefreshCredential(source)
	return credential, err
}

func generateRefreshCredential(source io.Reader) (string, []byte, error) {
	if source == nil {
		return "", nil, fmt.Errorf("generate refresh credential: %w", ErrInvalidCredential)
	}

	raw := make([]byte, RefreshCredentialBytes)
	if _, err := io.ReadFull(source, raw); err != nil {
		return "", nil, fmt.Errorf("generate refresh credential: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(raw), raw, nil
}

// DecodeRefreshCredential validates the transport representation before any
// digest is calculated. Canonical re-encoding rejects padded/non-canonical
// representations even when a decoder could accept them.
func DecodeRefreshCredential(encoded string) ([]byte, error) {
	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil || len(raw) != RefreshCredentialBytes || base64.RawURLEncoding.EncodeToString(raw) != encoded {
		return nil, ErrInvalidCredential
	}

	return raw, nil
}

// HashRefreshCredential returns the deterministic digest persisted in
// auth_sessions.refresh_token_hash. The caller must pass decoded credential
// bytes, never the encoded transport string.
func HashRefreshCredential(raw []byte) []byte {
	digest := sha256.Sum256(raw)
	return digest[:]
}

// HashEncodedRefreshCredential validates and hashes one transport credential.
func HashEncodedRefreshCredential(encoded string) ([]byte, error) {
	raw, err := DecodeRefreshCredential(encoded)
	if err != nil {
		return nil, err
	}

	return HashRefreshCredential(raw), nil
}

// IsInvalidCredential reports whether an error is caused by malformed token
// input. It is useful for internal tests and future transport mapping.
func IsInvalidCredential(err error) bool {
	return errors.Is(err, ErrInvalidCredential)
}
