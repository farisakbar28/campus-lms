package repository

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// AuthSession is the persistence shape used by the auth service. The
// plaintext refresh credential is intentionally absent from this type.
type AuthSession struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	RefreshTokenHash []byte
	IssuedAt         time.Time
	ExpiresAt        time.Time
	RotatedFrom      *uuid.UUID
	RevokedAt        *time.Time
	RevokedReason    *string
	UserAgent        string
	IPAddress        net.IP
	LastSeenAt       *time.Time
}

type AuthSessionRepository struct{}

const insertAuthSessionSQL = `
INSERT INTO auth_sessions (
    id,
    user_id,
    refresh_token_hash,
    issued_at,
    expires_at,
    rotated_from,
    revoked_at,
    revoked_reason,
    user_agent,
    ip_address,
    last_seen_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

const findAuthSessionByHashForUpdateSQL = `
SELECT
    id,
    user_id,
    refresh_token_hash,
    issued_at,
    expires_at,
    rotated_from,
    revoked_at,
    revoked_reason
FROM auth_sessions
WHERE refresh_token_hash = $1
FOR UPDATE`

const hasAuthSessionRotationChildSQL = `
SELECT EXISTS (
    SELECT 1
    FROM auth_sessions
    WHERE rotated_from = $1
)`

const revokeAuthSessionSQL = `
UPDATE auth_sessions
SET
    revoked_at = COALESCE(revoked_at, $2),
    revoked_reason = COALESCE(revoked_reason, $3),
    last_seen_at = CASE
        WHEN revoked_at IS NULL THEN GREATEST(COALESCE(last_seen_at, issued_at), $2)
        ELSE last_seen_at
    END
WHERE id = $1`

const revokeAllUserAuthSessionsSQL = `
UPDATE auth_sessions
SET
    revoked_at = $2,
    revoked_reason = $3,
    last_seen_at = GREATEST(COALESCE(last_seen_at, issued_at), $2)
WHERE user_id = $1
  AND revoked_at IS NULL`

func (AuthSessionRepository) Insert(ctx context.Context, tx pgx.Tx, session AuthSession) error {
	var userAgent any
	if session.UserAgent != "" {
		userAgent = session.UserAgent
	}
	var ipAddress any
	if len(session.IPAddress) > 0 {
		ipAddress = session.IPAddress
	}

	_, err := tx.Exec(ctx, insertAuthSessionSQL,
		session.ID,
		session.UserID,
		session.RefreshTokenHash,
		session.IssuedAt,
		session.ExpiresAt,
		session.RotatedFrom,
		session.RevokedAt,
		session.RevokedReason,
		userAgent,
		ipAddress,
		session.LastSeenAt,
	)
	if err != nil {
		return classifyDatabaseError("insert auth session", err)
	}

	return nil
}

// FindByRefreshTokenHashForUpdate locks the predecessor before service code
// evaluates lineage or attempts to insert a rotation child.
func (AuthSessionRepository) FindByRefreshTokenHashForUpdate(ctx context.Context, tx pgx.Tx, digest []byte) (AuthSession, error) {
	var session AuthSession
	err := tx.QueryRow(ctx, findAuthSessionByHashForUpdateSQL, digest).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshTokenHash,
		&session.IssuedAt,
		&session.ExpiresAt,
		&session.RotatedFrom,
		&session.RevokedAt,
		&session.RevokedReason,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return AuthSession{}, pgx.ErrNoRows
	}
	if err != nil {
		return AuthSession{}, classifyDatabaseError("find auth session by refresh digest", err)
	}

	return session, nil
}

func (AuthSessionRepository) HasRotationChild(ctx context.Context, tx pgx.Tx, predecessorID uuid.UUID) (bool, error) {
	var exists bool
	if err := tx.QueryRow(ctx, hasAuthSessionRotationChildSQL, predecessorID).Scan(&exists); err != nil {
		return false, classifyDatabaseError("check auth session rotation lineage", err)
	}

	return exists, nil
}

func (AuthSessionRepository) RevokeSession(ctx context.Context, tx pgx.Tx, sessionID uuid.UUID, revokedAt time.Time, reason string) (bool, error) {
	tag, err := tx.Exec(ctx, revokeAuthSessionSQL, sessionID, revokedAt, reason)
	if err != nil {
		return false, classifyDatabaseError("revoke auth session", err)
	}

	return tag.RowsAffected() == 1, nil
}

func (AuthSessionRepository) RevokeAllUserSessions(ctx context.Context, tx pgx.Tx, userID uuid.UUID, revokedAt time.Time, reason string) (int64, error) {
	tag, err := tx.Exec(ctx, revokeAllUserAuthSessionsSQL, userID, revokedAt, reason)
	if err != nil {
		return 0, classifyDatabaseError("revoke all auth sessions for user", err)
	}

	return tag.RowsAffected(), nil
}
