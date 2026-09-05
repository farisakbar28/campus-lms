package auth

import (
	"testing"
	"time"
)

func TestCalculateSessionExpiry(t *testing.T) {
	issuedAt := time.Date(2026, 8, 28, 19, 0, 0, 0, time.UTC)
	ttl := 168 * time.Hour
	want := issuedAt.Add(ttl)
	got := calculateSessionExpiry(issuedAt, ttl)
	if !got.Equal(want) {
		t.Fatalf("session expiry = %s, want %s", got, want)
	}
	t.Logf("expiry_calculation=issued_at_plus_ttl ttl=%s", ttl)
}
