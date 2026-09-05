package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/farisakbar28/campus-lms/apps/api/internal/auth"
)

var errNilAccessTokenVerifier = errors.New("access token verifier is required")

const (
	bearerScheme        = "Bearer"
	bearerSchemeLower   = "bearer"
	bearerRealm         = "campus-lms-api"
	bearerChallenge     = bearerScheme + ` realm="` + bearerRealm + `"`
	unauthenticatedCode = "unauthenticated"
	unauthenticatedText = "authentication is required"
)

// NewBearerMiddleware constructs a trusted authentication boundary around a
// downstream handler. Tenant authorization and Principal creation happen in a
// later layer.
func NewBearerMiddleware(verifier auth.AccessTokenVerifier) (func(http.Handler) http.Handler, error) {
	if verifier == nil {
		return nil, errNilAccessTokenVerifier
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			encoded, ok := parseBearerCredential(request.Header.Values("Authorization"))
			if !ok {
				writeUnauthenticated(writer)
				return
			}

			identity, err := verifier.Verify(encoded)
			if err != nil {
				writeUnauthenticated(writer)
				return
			}

			next.ServeHTTP(writer, request.WithContext(WithAuthIdentity(request.Context(), identity)))
		})
	}, nil
}

// parseBearerCredential validates only the HTTP Bearer credential syntax. JWT
// structure and signature validation belong to AccessTokenVerifier.Verify.
func parseBearerCredential(values []string) (string, bool) {
	if len(values) != 1 {
		return "", false
	}

	fieldValue := strings.Trim(values[0], " \t")
	separator := strings.IndexByte(fieldValue, ' ')
	if separator <= 0 || !isBearerScheme(fieldValue[:separator]) {
		return "", false
	}

	credential := strings.TrimLeft(fieldValue[separator:], " ")
	if !isBearerToken68(credential) {
		return "", false
	}

	return credential, true
}

func isBearerScheme(value string) bool {
	if len(value) != len(bearerScheme) {
		return false
	}

	for index := range value {
		character := value[index]
		if character >= 'A' && character <= 'Z' {
			character += 'a' - 'A'
		}
		if character != bearerSchemeLower[index] {
			return false
		}
	}

	return true
}

func isBearerToken68(value string) bool {
	if value == "" {
		return false
	}

	seenTokenCharacter := false
	seenPadding := false
	for index := range value {
		character := value[index]
		switch {
		case isBearerTokenCharacter(character):
			if seenPadding {
				return false
			}
			seenTokenCharacter = true
		case character == '=':
			if !seenTokenCharacter {
				return false
			}
			seenPadding = true
		default:
			return false
		}
	}

	return seenTokenCharacter
}

func isBearerTokenCharacter(character byte) bool {
	return (character >= 'A' && character <= 'Z') ||
		(character >= 'a' && character <= 'z') ||
		(character >= '0' && character <= '9') ||
		strings.ContainsRune("-._~+/", rune(character))
}

type unauthenticatedResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeUnauthenticated(writer http.ResponseWriter) {
	writer.Header().Set("WWW-Authenticate", bearerChallenge)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(writer).Encode(unauthenticatedResponse{
		Code:    unauthenticatedCode,
		Message: unauthenticatedText,
	})
}
