package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/farisakbar28/campus-lms/apps/api/internal/database"
	"github.com/farisakbar28/campus-lms/apps/api/internal/repository"
	"github.com/google/uuid"
)

var errNilTenantAdmitter = errors.New("tenant admission dependency is required")

type TenantAdmitter interface {
	Admit(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) (uuid.UUID, error)
}

// NewTenantAdmissionMiddleware turns a verified token identity and an
// untrusted URL selector into a trusted Principal only after database checks.
func NewTenantAdmissionMiddleware(admitter TenantAdmitter, logger *slog.Logger) (func(http.Handler) http.Handler, error) {
	if admitter == nil || logger == nil {
		return nil, errNilTenantAdmitter
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			identity, ok := AuthIdentityFromContext(request.Context())
			if !ok || identity.UserID == uuid.Nil || identity.SessionID == uuid.Nil {
				writeUnauthenticated(writer)
				return
			}

			selector := request.PathValue("tenant_id")
			tenantID, err := uuid.Parse(selector)
			if err != nil || len(selector) != 36 || !strings.EqualFold(selector, tenantID.String()) {
				writeAdmissionError(writer, http.StatusBadRequest, "invalid_tenant_id", "tenant ID must be a UUID")
				return
			}

			membershipID, err := admitter.Admit(request.Context(), identity.UserID, identity.SessionID, tenantID)
			switch {
			case errors.Is(err, repository.ErrAuthenticationState):
				writeUnauthenticated(writer)
				return
			case errors.Is(err, repository.ErrTenantAdmission):
				writeAdmissionError(writer, http.StatusNotFound, "tenant_not_found", "tenant was not found")
				return
			case errors.Is(err, database.ErrUnavailable):
				logger.Error("tenant admission dependency unavailable", "error", err)
				writeAdmissionError(writer, http.StatusServiceUnavailable, "service_unavailable", "service is temporarily unavailable")
				return
			case err != nil:
				logger.Error("tenant admission failed", "error", err)
				writeAdmissionError(writer, http.StatusInternalServerError, "internal_error", "internal server error")
				return
			case membershipID == uuid.Nil:
				writeAdmissionError(writer, http.StatusNotFound, "tenant_not_found", "tenant was not found")
				return
			}

			principal := Principal{
				UserID:       identity.UserID,
				SessionID:    identity.SessionID,
				TenantID:     tenantID,
				MembershipID: membershipID,
			}
			next.ServeHTTP(writer, request.WithContext(WithPrincipal(request.Context(), principal)))
		})
	}, nil
}

func writeAdmissionError(writer http.ResponseWriter, status int, code, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(unauthenticatedResponse{Code: code, Message: message})
}
