// Package middleware contains trusted request identity plumbing.
package middleware

import "context"

// Principal is populated only by trusted authentication middleware in a future loop.
type Principal struct {
	TenantID string
	UserID   string
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
