package middleware

import (
	"context"

	"github.com/farisakbar28/campus-lms/apps/api/internal/auth"
)

type authIdentityContextKey struct{}

// WithAuthIdentity attaches verifier-produced identity for trusted middleware
// and internal tests. Client-controlled request data must never call this
// helper directly.
func WithAuthIdentity(ctx context.Context, identity auth.AuthIdentity) context.Context {
	return context.WithValue(ctx, authIdentityContextKey{}, identity)
}

// AuthIdentityFromContext returns the canonical authenticated identity, if a
// trusted middleware has populated it.
func AuthIdentityFromContext(ctx context.Context) (auth.AuthIdentity, bool) {
	identity, ok := ctx.Value(authIdentityContextKey{}).(auth.AuthIdentity)
	return identity, ok
}
