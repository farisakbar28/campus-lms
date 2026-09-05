package auth

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/farisakbar28/campus-lms/apps/api/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
)

type authIntegrationSuite struct {
	owner *pgx.Conn
	pool  *database.Pool
}

var authSuite *authIntegrationSuite

func TestMain(m *testing.M) {
	ctx := context.Background()
	container, err := postgrescontainer.Run(ctx, "postgres:16.14-alpine3.23",
		postgrescontainer.WithDatabase("auth_session_test"),
		postgrescontainer.WithUsername("owner"),
		postgrescontainer.WithPassword("owner-test-only"),
		postgrescontainer.WithOrderedInitScripts(authIntegrationSQLFiles()...),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start auth integration database: %v\n", err)
		os.Exit(1)
	}

	exitCode := 1
	defer func() {
		if authSuite != nil {
			authSuite.pool.Close()
			_ = authSuite.owner.Close(ctx)
		}
		_ = container.Terminate(ctx)
		os.Exit(exitCode)
	}()

	databaseURL, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "get auth integration connection string: %v\n", err)
		return
	}
	owner, err := connectAuthWithRetry(ctx, databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect auth integration owner: %v\n", err)
		return
	}
	pool, err := database.NewPool(ctx, databaseURL, 0, 4)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect auth integration pool: %v\n", err)
		_ = owner.Close(ctx)
		return
	}

	authSuite = &authIntegrationSuite{owner: owner, pool: pool}
	exitCode = m.Run()
}

func TestAuthSessionCreateIntegration(t *testing.T) {
	ctx := context.Background()
	userID := createAuthUser(t)
	issuedAt := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)
	ttl := 7 * 24 * time.Hour
	service := NewService(authSuite.pool, ttl, WithClock(func() time.Time { return issuedAt }))

	session, err := service.CreateSession(ctx, CreateSessionInput{
		UserID:    userID,
		UserAgent: "integration-agent/1.0",
		IPAddress: net.ParseIP("192.0.2.10"),
	})
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}
	if session.ID == uuid.Nil || session.UserID != userID {
		t.Fatalf("session identity = %s/%s, want non-nil ID and user %s", session.ID, session.UserID, userID)
	}
	if decoded, err := DecodeRefreshCredential(session.RefreshToken); err != nil || len(decoded) != RefreshCredentialBytes {
		t.Fatalf("generated refresh credential decode = %d bytes, error=%v; want %d bytes", len(decoded), err, RefreshCredentialBytes)
	}
	if !session.ExpiresAt.Equal(issuedAt.Add(ttl)) {
		t.Fatalf("expires_at = %s, want %s", session.ExpiresAt, issuedAt.Add(ttl))
	}

	var storedUserID uuid.UUID
	var storedDigest []byte
	var storedIssuedAt, storedExpiresAt time.Time
	var rotatedFromIsNull, lastSeenIsNull bool
	var userAgent, ipAddress string
	err = authSuite.owner.QueryRow(ctx, `
SELECT user_id, refresh_token_hash, issued_at, expires_at,
       rotated_from IS NULL, COALESCE(user_agent, ''),
       COALESCE(host(ip_address)::text, ''), last_seen_at IS NULL
FROM auth_sessions
WHERE id = $1`, session.ID).Scan(
		&storedUserID,
		&storedDigest,
		&storedIssuedAt,
		&storedExpiresAt,
		&rotatedFromIsNull,
		&userAgent,
		&ipAddress,
		&lastSeenIsNull,
	)
	if err != nil {
		t.Fatalf("inspect created auth session: %v", err)
	}
	wantDigest, err := HashEncodedRefreshCredential(session.RefreshToken)
	if err != nil {
		t.Fatalf("hash returned refresh credential: %v", err)
	}
	if storedUserID != userID || !bytesEqual(storedDigest, wantDigest) {
		t.Fatalf("stored identity/digest does not match returned session")
	}
	if len(storedDigest) != sha256.Size {
		t.Fatalf("stored digest length = %d, want %d", len(storedDigest), sha256.Size)
	}
	if bytesEqual(storedDigest, []byte(session.RefreshToken)) {
		t.Fatal("stored digest equals encoded plaintext credential")
	}
	if !storedIssuedAt.Equal(issuedAt) || !storedExpiresAt.Equal(issuedAt.Add(ttl)) {
		t.Fatalf("stored timestamps = %s/%s, want %s/%s", storedIssuedAt, storedExpiresAt, issuedAt, issuedAt.Add(ttl))
	}
	if !rotatedFromIsNull || !lastSeenIsNull || userAgent != "integration-agent/1.0" || !net.ParseIP(ipAddress).Equal(net.ParseIP("192.0.2.10")) {
		t.Fatalf("stored metadata = rotated_from_null=%t last_seen_null=%t user_agent=%q ip=%q", rotatedFromIsNull, lastSeenIsNull, userAgent, ipAddress)
	}
	t.Logf("create_row_found=true digest_bytes=%d plaintext_persisted=false rotated_from_null=%t expiry_match=true user_agent_stored=true ip_metadata_stored=true", len(storedDigest), rotatedFromIsNull)
}

func TestAuthSessionRotateIntegration(t *testing.T) {
	ctx := context.Background()
	userID := createAuthUser(t)
	issuedAt := time.Date(2026, 8, 28, 13, 0, 0, 0, time.UTC)
	ttl := 24 * time.Hour
	service := NewService(authSuite.pool, ttl, WithClock(func() time.Time { return issuedAt }))

	predecessor, err := service.CreateSession(ctx, CreateSessionInput{UserID: userID})
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}
	rotatedAt := issuedAt.Add(time.Hour)
	rotator := NewService(authSuite.pool, ttl, WithClock(func() time.Time { return rotatedAt }))
	child, err := rotator.RotateRefreshToken(ctx, predecessor.RefreshToken)
	if err != nil {
		t.Fatalf("RotateRefreshToken() error = %v", err)
	}
	if child.ID == uuid.Nil || child.ID == predecessor.ID || child.UserID != userID {
		t.Fatalf("child identity = %s/%s, want a new session for %s", child.ID, child.UserID, userID)
	}
	if child.RefreshToken == predecessor.RefreshToken {
		t.Fatal("rotation returned the predecessor credential")
	}
	if child.RotatedFrom == nil || *child.RotatedFrom != predecessor.ID {
		t.Fatalf("child rotated_from = %v, want %s", child.RotatedFrom, predecessor.ID)
	}
	if !child.ExpiresAt.Equal(rotatedAt.Add(ttl)) {
		t.Fatalf("child expires_at = %s, want %s", child.ExpiresAt, rotatedAt.Add(ttl))
	}

	var predecessorRevoked bool
	var predecessorReason string
	var childCount, historyCount int
	var predecessorDigest, childDigest []byte
	if err := authSuite.owner.QueryRow(ctx, `SELECT revoked_at IS NOT NULL, COALESCE(revoked_reason, '') FROM auth_sessions WHERE id = $1`, predecessor.ID).Scan(&predecessorRevoked, &predecessorReason); err != nil {
		t.Fatalf("inspect predecessor state: %v", err)
	}
	if err := authSuite.owner.QueryRow(ctx, `SELECT count(*) FROM auth_sessions WHERE rotated_from = $1`, predecessor.ID).Scan(&childCount); err != nil {
		t.Fatalf("count rotation children: %v", err)
	}
	if err := authSuite.owner.QueryRow(ctx, `SELECT count(*) FROM auth_sessions WHERE user_id = $1`, userID).Scan(&historyCount); err != nil {
		t.Fatalf("count session history: %v", err)
	}
	if err := authSuite.owner.QueryRow(ctx, `SELECT refresh_token_hash FROM auth_sessions WHERE id = $1`, predecessor.ID).Scan(&predecessorDigest); err != nil {
		t.Fatalf("read predecessor digest: %v", err)
	}
	if err := authSuite.owner.QueryRow(ctx, `SELECT refresh_token_hash FROM auth_sessions WHERE id = $1`, child.ID).Scan(&childDigest); err != nil {
		t.Fatalf("read child digest: %v", err)
	}
	if !predecessorRevoked || predecessorReason != "rotated" || childCount != 1 || historyCount != 2 {
		t.Fatalf("rotation state = revoked=%t reason=%q child_count=%d history=%d", predecessorRevoked, predecessorReason, childCount, historyCount)
	}
	if bytesEqual(predecessorDigest, childDigest) {
		t.Fatal("rotation persisted the same digest for predecessor and child")
	}
	t.Logf("rotation_child_count=%d user_history_rows=%d predecessor_revoked=%t same_user=true digests_distinct=true", childCount, historyCount, predecessorRevoked)
}

func TestAuthSessionReuseIntegration(t *testing.T) {
	ctx := context.Background()
	compromisedUser := createAuthUser(t)
	otherUser := createAuthUser(t)
	now := time.Date(2026, 8, 28, 14, 0, 0, 0, time.UTC)
	ttl := 24 * time.Hour
	service := NewService(authSuite.pool, ttl, WithClock(func() time.Time { return now }))

	predecessor, err := service.CreateSession(ctx, CreateSessionInput{UserID: compromisedUser})
	if err != nil {
		t.Fatalf("create predecessor: %v", err)
	}
	otherCompromisedSession, err := service.CreateSession(ctx, CreateSessionInput{UserID: compromisedUser})
	if err != nil {
		t.Fatalf("create second compromised-user session: %v", err)
	}
	otherUserSession, err := service.CreateSession(ctx, CreateSessionInput{UserID: otherUser})
	if err != nil {
		t.Fatalf("create other-user session: %v", err)
	}
	rotated, err := service.RotateRefreshToken(ctx, predecessor.RefreshToken)
	if err != nil {
		t.Fatalf("rotate predecessor: %v", err)
	}

	_, err = service.RotateRefreshToken(ctx, predecessor.RefreshToken)
	if err == nil || !errors.Is(err, ErrAuthenticationFailed) || !errors.Is(err, ErrRefreshReuse) {
		t.Fatalf("reuse error = %v, want generic authentication failure with internal reuse reason", err)
	}
	if err.Error() != ErrAuthenticationFailed.Error() {
		t.Fatalf("reuse error string = %q, want %q", err, ErrAuthenticationFailed)
	}

	var predecessorChildren, userHistory, activeCompromised, activeOther int
	if err := authSuite.owner.QueryRow(ctx, `SELECT count(*) FROM auth_sessions WHERE rotated_from = $1`, predecessor.ID).Scan(&predecessorChildren); err != nil {
		t.Fatalf("count reuse children: %v", err)
	}
	if err := authSuite.owner.QueryRow(ctx, `SELECT count(*) FROM auth_sessions WHERE user_id = $1`, compromisedUser).Scan(&userHistory); err != nil {
		t.Fatalf("count compromised-user history: %v", err)
	}
	if err := authSuite.owner.QueryRow(ctx, `SELECT count(*) FROM auth_sessions WHERE user_id = $1 AND revoked_at IS NULL`, compromisedUser).Scan(&activeCompromised); err != nil {
		t.Fatalf("count compromised-user active sessions: %v", err)
	}
	if err := authSuite.owner.QueryRow(ctx, `SELECT count(*) FROM auth_sessions WHERE user_id = $1 AND revoked_at IS NULL`, otherUser).Scan(&activeOther); err != nil {
		t.Fatalf("count other-user active sessions: %v", err)
	}
	if predecessorChildren != 1 || userHistory != 3 || activeCompromised != 0 || activeOther != 1 {
		t.Fatalf("reuse state = children=%d history=%d compromised_active=%d other_active=%d", predecessorChildren, userHistory, activeCompromised, activeOther)
	}
	assertSessionRevoked(t, otherCompromisedSession.ID)
	assertSessionRevoked(t, rotated.ID)
	assertSessionUnrevoked(t, otherUserSession.ID)
	t.Logf("reuse_child_count=%d compromised_user_history_rows=%d compromised_user_unrevoked=%d other_user_unrevoked=%d", predecessorChildren, userHistory, activeCompromised, activeOther)
}

func TestAuthSessionManualRevokeIsNotReuseIntegration(t *testing.T) {
	ctx := context.Background()
	userID := createAuthUser(t)
	now := time.Date(2026, 8, 28, 15, 0, 0, 0, time.UTC)
	service := NewService(authSuite.pool, 24*time.Hour, WithClock(func() time.Time { return now }))

	revoked, err := service.CreateSession(ctx, CreateSessionInput{UserID: userID})
	if err != nil {
		t.Fatalf("create session to revoke: %v", err)
	}
	active, err := service.CreateSession(ctx, CreateSessionInput{UserID: userID})
	if err != nil {
		t.Fatalf("create active session: %v", err)
	}
	if err := service.RevokeSession(ctx, revoked.ID, "logout"); err != nil {
		t.Fatalf("RevokeSession() error = %v", err)
	}
	_, err = service.RotateRefreshToken(ctx, revoked.RefreshToken)
	if err == nil || !errors.Is(err, ErrAuthenticationFailed) || !errors.Is(err, ErrRevokedCredential) || errors.Is(err, ErrRefreshReuse) {
		t.Fatalf("manually revoked credential error = %v, want revoked-only authentication failure", err)
	}

	var childCount, activeCount int
	if err := authSuite.owner.QueryRow(ctx, `SELECT count(*) FROM auth_sessions WHERE rotated_from = $1`, revoked.ID).Scan(&childCount); err != nil {
		t.Fatalf("count manual-revoke children: %v", err)
	}
	if err := authSuite.owner.QueryRow(ctx, `SELECT count(*) FROM auth_sessions WHERE user_id = $1 AND revoked_at IS NULL`, userID).Scan(&activeCount); err != nil {
		t.Fatalf("count active sessions after manual revoke: %v", err)
	}
	if childCount != 0 || activeCount != 1 {
		t.Fatalf("manual revoke state = children=%d active=%d", childCount, activeCount)
	}
	assertSessionUnrevoked(t, active.ID)
	t.Logf("manual_revoke_child_count=%d same_user_unrevoked=%d reuse_classification=false", childCount, activeCount)
}

func TestAuthSessionExpiryBoundaryIntegration(t *testing.T) {
	ctx := context.Background()
	issuedAt := time.Date(2026, 8, 28, 16, 0, 0, 0, time.UTC)
	ttl := time.Hour
	tests := []struct {
		name      string
		rotateAt  time.Time
		wantValid bool
	}{
		{name: "just before expiry", rotateAt: issuedAt.Add(ttl - time.Nanosecond), wantValid: true},
		{name: "exact expiry", rotateAt: issuedAt.Add(ttl), wantValid: false},
		{name: "after expiry", rotateAt: issuedAt.Add(ttl + time.Nanosecond), wantValid: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userID := createAuthUser(t)
			creator := NewService(authSuite.pool, ttl, WithClock(func() time.Time { return issuedAt }))
			created, err := creator.CreateSession(ctx, CreateSessionInput{UserID: userID})
			if err != nil {
				t.Fatalf("create session: %v", err)
			}
			rotator := NewService(authSuite.pool, ttl, WithClock(func() time.Time { return test.rotateAt }))
			rotated, err := rotator.RotateRefreshToken(ctx, created.RefreshToken)
			if test.wantValid {
				if err != nil || rotated.ID == uuid.Nil {
					t.Fatalf("rotation error = %v, result=%+v; want success", err, rotated)
				}
				t.Log("expiry_boundary=just_before result=rotated")
				return
			}
			if err == nil || !errors.Is(err, ErrAuthenticationFailed) || !errors.Is(err, ErrExpiredCredential) {
				t.Fatalf("expiry error = %v, want expired authentication failure", err)
			}
			var childCount, historyCount int
			if err := authSuite.owner.QueryRow(ctx, `SELECT count(*) FROM auth_sessions WHERE rotated_from = $1`, created.ID).Scan(&childCount); err != nil {
				t.Fatalf("count expired-session children: %v", err)
			}
			if err := authSuite.owner.QueryRow(ctx, `SELECT count(*) FROM auth_sessions WHERE user_id = $1`, userID).Scan(&historyCount); err != nil {
				t.Fatalf("count expired-session history: %v", err)
			}
			if childCount != 0 || historyCount != 1 {
				t.Fatalf("expired state = children=%d history=%d, want 0/1", childCount, historyCount)
			}
			t.Logf("expiry_boundary=%s child_count=%d history_rows=%d result=authentication_failure", test.name, childCount, historyCount)
		})
	}
}

func TestAuthSessionConcurrentRotationIntegration(t *testing.T) {
	ctx := context.Background()
	compromisedUser := createAuthUser(t)
	otherUser := createAuthUser(t)
	now := time.Date(2026, 8, 28, 17, 0, 0, 0, time.UTC)
	ttl := 24 * time.Hour
	creator := NewService(authSuite.pool, ttl, WithClock(func() time.Time { return now }))
	predecessor, err := creator.CreateSession(ctx, CreateSessionInput{UserID: compromisedUser})
	if err != nil {
		t.Fatalf("create concurrent predecessor: %v", err)
	}
	otherUserSession, err := creator.CreateSession(ctx, CreateSessionInput{UserID: otherUser})
	if err != nil {
		t.Fatalf("create concurrent control session: %v", err)
	}

	start := make(chan struct{})
	results := make(chan struct {
		session Session
		err     error
	}, 2)
	var waitGroup sync.WaitGroup
	for range 2 {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			<-start
			service := NewService(authSuite.pool, ttl, WithClock(func() time.Time { return now }))
			session, err := service.RotateRefreshToken(ctx, predecessor.RefreshToken)
			results <- struct {
				session Session
				err     error
			}{session: session, err: err}
		}()
	}
	close(start)
	waitGroup.Wait()
	close(results)

	successes := 0
	reuseFailures := 0
	for result := range results {
		if result.err == nil {
			successes++
			continue
		}
		if errors.Is(result.err, ErrAuthenticationFailed) && errors.Is(result.err, ErrRefreshReuse) {
			reuseFailures++
			continue
		}
		t.Fatalf("concurrent rotation error = %v, want reuse authentication failure for loser", result.err)
	}

	var childCount, activeCompromised, activeOther, historyCount int
	if err := authSuite.owner.QueryRow(ctx, `SELECT count(*) FROM auth_sessions WHERE rotated_from = $1`, predecessor.ID).Scan(&childCount); err != nil {
		t.Fatalf("count concurrent children: %v", err)
	}
	if err := authSuite.owner.QueryRow(ctx, `SELECT count(*) FROM auth_sessions WHERE user_id = $1 AND revoked_at IS NULL`, compromisedUser).Scan(&activeCompromised); err != nil {
		t.Fatalf("count concurrent compromised-user active sessions: %v", err)
	}
	if err := authSuite.owner.QueryRow(ctx, `SELECT count(*) FROM auth_sessions WHERE user_id = $1 AND revoked_at IS NULL`, otherUser).Scan(&activeOther); err != nil {
		t.Fatalf("count concurrent other-user active sessions: %v", err)
	}
	if err := authSuite.owner.QueryRow(ctx, `SELECT count(*) FROM auth_sessions WHERE user_id = $1`, compromisedUser).Scan(&historyCount); err != nil {
		t.Fatalf("count concurrent history: %v", err)
	}
	if successes != 1 || reuseFailures != 1 || childCount != 1 || activeCompromised != 0 || activeOther != 1 || historyCount != 2 {
		t.Fatalf("concurrent state = successes=%d reuse_failures=%d children=%d compromised_active=%d other_active=%d history=%d", successes, reuseFailures, childCount, activeCompromised, activeOther, historyCount)
	}
	assertSessionUnrevoked(t, otherUserSession.ID)
	t.Logf("successful_rotation_calls=%d reuse_failures=%d child_count=%d compromised_user_unrevoked=%d compromised_user_history_rows=%d other_user_unrevoked=%d", successes, reuseFailures, childCount, activeCompromised, historyCount, activeOther)
}

func TestAuthSessionRevokeOneAndAllIntegration(t *testing.T) {
	ctx := context.Background()
	userID := createAuthUser(t)
	otherUser := createAuthUser(t)
	now := time.Date(2026, 8, 28, 18, 0, 0, 0, time.UTC)
	service := NewService(authSuite.pool, 24*time.Hour, WithClock(func() time.Time { return now }))

	first, err := service.CreateSession(ctx, CreateSessionInput{UserID: userID})
	if err != nil {
		t.Fatalf("create first revoke session: %v", err)
	}
	second, err := service.CreateSession(ctx, CreateSessionInput{UserID: userID})
	if err != nil {
		t.Fatalf("create second revoke session: %v", err)
	}
	third, err := service.CreateSession(ctx, CreateSessionInput{UserID: userID})
	if err != nil {
		t.Fatalf("create third revoke session: %v", err)
	}
	other, err := service.CreateSession(ctx, CreateSessionInput{UserID: otherUser})
	if err != nil {
		t.Fatalf("create revoke control session: %v", err)
	}

	if err := service.RevokeSession(ctx, first.ID, "logout"); err != nil {
		t.Fatalf("first RevokeSession() error = %v", err)
	}
	if err := service.RevokeSession(ctx, first.ID, "different-reason"); err != nil {
		t.Fatalf("second RevokeSession() error = %v", err)
	}
	assertSessionRevokedWithReason(t, first.ID, "logout")

	if err := service.RevokeAllUserSessions(ctx, userID, "security-event"); err != nil {
		t.Fatalf("RevokeAllUserSessions() error = %v", err)
	}
	var userHistory, activeUser, activeOther int
	if err := authSuite.owner.QueryRow(ctx, `SELECT count(*) FROM auth_sessions WHERE user_id = $1`, userID).Scan(&userHistory); err != nil {
		t.Fatalf("count revoke-all history: %v", err)
	}
	if err := authSuite.owner.QueryRow(ctx, `SELECT count(*) FROM auth_sessions WHERE user_id = $1 AND revoked_at IS NULL`, userID).Scan(&activeUser); err != nil {
		t.Fatalf("count revoke-all active user sessions: %v", err)
	}
	if err := authSuite.owner.QueryRow(ctx, `SELECT count(*) FROM auth_sessions WHERE user_id = $1 AND revoked_at IS NULL`, otherUser).Scan(&activeOther); err != nil {
		t.Fatalf("count revoke-all active other sessions: %v", err)
	}
	if userHistory != 3 || activeUser != 0 || activeOther != 1 {
		t.Fatalf("revoke-all state = history=%d active_user=%d active_other=%d", userHistory, activeUser, activeOther)
	}
	assertSessionRevokedWithReason(t, second.ID, "security-event")
	assertSessionRevokedWithReason(t, third.ID, "security-event")
	assertSessionUnrevoked(t, other.ID)
	t.Logf("revoke_one_idempotent=true user_history_rows=%d user_unrevoked=%d other_user_unrevoked=%d", userHistory, activeUser, activeOther)
}

func createAuthUser(t *testing.T) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	userID := uuid.New()
	_, err := authSuite.owner.Exec(ctx, `
INSERT INTO users (id, email, display_name, status, created_at)
VALUES ($1, $2, $3, 'active', $4)`, userID, userID.String()+"@auth-session.test", "Auth Session Test User", time.Date(2026, 8, 28, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("create auth test user: %v", err)
	}

	return userID
}

func assertSessionRevoked(t *testing.T, sessionID uuid.UUID) {
	t.Helper()
	var revoked bool
	if err := authSuite.owner.QueryRow(context.Background(), `SELECT revoked_at IS NOT NULL FROM auth_sessions WHERE id = $1`, sessionID).Scan(&revoked); err != nil {
		t.Fatalf("inspect revoked session %s: %v", sessionID, err)
	}
	if !revoked {
		t.Errorf("session %s is not revoked", sessionID)
	}
}

func assertSessionRevokedWithReason(t *testing.T, sessionID uuid.UUID, wantReason string) {
	t.Helper()
	var revoked bool
	var reason string
	if err := authSuite.owner.QueryRow(context.Background(), `SELECT revoked_at IS NOT NULL, COALESCE(revoked_reason, '') FROM auth_sessions WHERE id = $1`, sessionID).Scan(&revoked, &reason); err != nil {
		t.Fatalf("inspect revoked session %s: %v", sessionID, err)
	}
	if !revoked || reason != wantReason {
		t.Errorf("session %s state = revoked=%t reason=%q, want revoked with %q", sessionID, revoked, reason, wantReason)
	}
}

func assertSessionUnrevoked(t *testing.T, sessionID uuid.UUID) {
	t.Helper()
	var unrevoked bool
	if err := authSuite.owner.QueryRow(context.Background(), `SELECT revoked_at IS NULL FROM auth_sessions WHERE id = $1`, sessionID).Scan(&unrevoked); err != nil {
		t.Fatalf("inspect active session %s: %v", sessionID, err)
	}
	if !unrevoked {
		t.Errorf("session %s is unexpectedly revoked", sessionID)
	}
}

func bytesEqual(left, right []byte) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func authIntegrationSQLFiles() []string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}
	apiRoot := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", ".."))
	return []string{
		filepath.Join(apiRoot, "migrations", "0001_tenant_identity_schema.up.sql"),
		filepath.Join(apiRoot, "migrations", "0002_academic_core_schema.up.sql"),
		filepath.Join(apiRoot, "migrations", "0003_auth_membership_schema.up.sql"),
		filepath.Join(apiRoot, "migrations", "0004_academic_term_time_range_check.up.sql"),
		filepath.Join(apiRoot, "migrations", "0005_enrollments_active_student_lookup_index.up.sql"),
		filepath.Join(apiRoot, "migrations", "0006_auth_sessions_schema.up.sql"),
	}
}

func connectAuthWithRetry(ctx context.Context, connectionURL string) (*pgx.Conn, error) {
	deadline := time.NewTimer(30 * time.Second)
	defer deadline.Stop()
	for {
		connection, err := pgx.Connect(ctx, connectionURL)
		if err == nil {
			return connection, nil
		}
		select {
		case <-deadline.C:
			return nil, err
		case <-time.After(250 * time.Millisecond):
		}
	}
}
