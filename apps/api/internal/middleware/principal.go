// Package middleware contains trusted request identity plumbing.
package middleware

import (
	"context"

	"github.com/google/uuid"
)

// Principal is populated only after trusted user, session, tenant, and
// membership admission. It does not carry role or object authority.
type Principal struct {
	UserID       uuid.UUID
	SessionID    uuid.UUID
	TenantID     uuid.UUID
	MembershipID uuid.UUID
}

type principalContextKey struct{}

// WithPrincipal is for trusted middleware and internal tests, never client input.
func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, principalContextKey{}, principal)
}

// PrincipalFromContext returns the authenticated principal, if middleware set one.
func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(principalContextKey{}).(Principal)
	return principal, ok
}
