package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/farisakbar28/campus-lms/apps/api/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// GlobalTransactioner is deliberately narrower than the database pool. It
// prevents the auth service from reaching around the global transaction
// boundary or accidentally choosing tenant-scoped work.
type GlobalTransactioner interface {
	WithTx(context.Context, func(pgx.Tx) error) error
}

type SessionIDSource func() (uuid.UUID, error)
type Clock func() time.Time
type Option func(*Service)

// WithClock replaces the wall clock for deterministic lifecycle tests.
func WithClock(clock Clock) Option {
	return func(service *Service) {
		if clock != nil {
			service.clock = clock
		}
	}
}

// WithTokenSource replaces crypto/rand.Reader for deterministic unit tests.
func WithTokenSource(source io.Reader) Option {
	return func(service *Service) {
		if source != nil {
			service.tokenSource = source
		}
	}
}

// WithSessionIDSource replaces UUID generation for deterministic tests.
func WithSessionIDSource(source SessionIDSource) Option {
	return func(service *Service) {
		if source != nil {
			service.sessionIDSource = source
		}
	}
}

// Service owns auth session security and lifecycle semantics. PostgreSQL
// statements remain in repository.AuthSessionRepository.
type Service struct {
	transactions    GlobalTransactioner
	sessions        repository.AuthSessionRepository
	refreshTTL      time.Duration
	clock           Clock
	tokenSource     io.Reader
	sessionIDSource SessionIDSource
}

// NewService constructs the internal refresh-session service. The configured
// TTL is validated at operation time so configuration and service tests can
// exercise the same failure path without a second constructor error type.
func NewService(transactions GlobalTransactioner, refreshTTL time.Duration, options ...Option) *Service {
	service := &Service{
		transactions: transactions,
		refreshTTL:   refreshTTL,
		clock:        time.Now,
		tokenSource:  rand.Reader,
		sessionIDSource: func() (uuid.UUID, error) {
			return uuid.NewRandom()
		},
	}
	for _, option := range options {
		if option != nil {
			option(service)
		}
	}

	return service
}

type CreateSessionInput struct {
	UserID    uuid.UUID
	UserAgent string
	IPAddress net.IP
}

type Session struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	RefreshToken string
	IssuedAt     time.Time
	ExpiresAt    time.Time
	RotatedFrom  *uuid.UUID
}

// CreateSession issues one independent refresh session. It never revokes an
// existing session and only returns the plaintext token after the insert has
// committed successfully.
func (service *Service) CreateSession(ctx context.Context, input CreateSessionInput) (Session, error) {
	if input.UserID == uuid.Nil {
		return Session{}, ErrInvalidUserID
	}
	if err := service.validate(); err != nil {
		return Session{}, err
	}

	sessionID, err := service.sessionIDSource()
	if err != nil {
		return Session{}, fmt.Errorf("generate auth session id: %w", err)
	}
	if sessionID == uuid.Nil {
		return Session{}, fmt.Errorf("generate auth session id: %w", ErrInvalidSessionID)
	}

	refreshToken, rawCredential, err := generateRefreshCredential(service.tokenSource)
	if err != nil {
		return Session{}, err
	}
	issuedAt := service.clock()
	expiresAt := calculateSessionExpiry(issuedAt, service.refreshTTL)
	persisted := repository.AuthSession{
		ID:               sessionID,
		UserID:           input.UserID,
		RefreshTokenHash: HashRefreshCredential(rawCredential),
		IssuedAt:         issuedAt,
		ExpiresAt:        expiresAt,
		UserAgent:        input.UserAgent,
		IPAddress:        input.IPAddress,
	}

	if err := service.transactions.WithTx(ctx, func(tx pgx.Tx) error {
		return service.sessions.Insert(ctx, tx, persisted)
	}); err != nil {
		return Session{}, err
	}

	return Session{
		ID:           sessionID,
		UserID:       input.UserID,
		RefreshToken: refreshToken,
		IssuedAt:     issuedAt,
		ExpiresAt:    expiresAt,
	}, nil
}

// RotateRefreshToken atomically consumes one valid predecessor and creates a
// child session. A reuse branch commits revocation before returning a generic
// authentication failure, so the compromise response cannot roll back its
// own containment work.
func (service *Service) RotateRefreshToken(ctx context.Context, refreshToken string) (Session, error) {
	digest, err := HashEncodedRefreshCredential(refreshToken)
	if err != nil {
		return Session{}, newAuthenticationError(ErrInvalidCredential)
	}
	if err := service.validate(); err != nil {
		return Session{}, err
	}

	var rotated Session
	var authenticationFailure error
	err = service.transactions.WithTx(ctx, func(tx pgx.Tx) error {
		predecessor, err := service.sessions.FindByRefreshTokenHashForUpdate(ctx, tx, digest)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				authenticationFailure = newAuthenticationError(ErrInvalidCredential)
				return nil
			}
			return err
		}

		hasRotationChild, err := service.sessions.HasRotationChild(ctx, tx, predecessor.ID)
		if err != nil {
			return err
		}

		now := service.clock()
		if hasRotationChild {
			if _, err := service.sessions.RevokeAllUserSessions(ctx, tx, predecessor.UserID, now, "refresh_token_reuse"); err != nil {
				return err
			}
			authenticationFailure = newAuthenticationError(ErrRefreshReuse)
			return nil
		}
		if predecessor.RevokedAt != nil {
			authenticationFailure = newAuthenticationError(ErrRevokedCredential)
			return nil
		}
		if !now.Before(predecessor.ExpiresAt) {
			authenticationFailure = newAuthenticationError(ErrExpiredCredential)
			return nil
		}

		childID, err := service.sessionIDSource()
		if err != nil {
			return fmt.Errorf("generate rotated auth session id: %w", err)
		}
		if childID == uuid.Nil {
			return fmt.Errorf("generate rotated auth session id: %w", ErrInvalidSessionID)
		}
		childToken, childRawCredential, err := generateRefreshCredential(service.tokenSource)
		if err != nil {
			return err
		}
		childIssuedAt := now
		childExpiresAt := calculateSessionExpiry(childIssuedAt, service.refreshTTL)
		child := repository.AuthSession{
			ID:               childID,
			UserID:           predecessor.UserID,
			RefreshTokenHash: HashRefreshCredential(childRawCredential),
			IssuedAt:         childIssuedAt,
			ExpiresAt:        childExpiresAt,
			RotatedFrom:      &predecessor.ID,
			UserAgent:        "",
			IPAddress:        nil,
		}
		if err := service.sessions.Insert(ctx, tx, child); err != nil {
			return err
		}
		if _, err := service.sessions.RevokeSession(ctx, tx, predecessor.ID, now, "rotated"); err != nil {
			return err
		}

		rotated = Session{
			ID:           childID,
			UserID:       predecessor.UserID,
			RefreshToken: childToken,
			IssuedAt:     childIssuedAt,
			ExpiresAt:    childExpiresAt,
			RotatedFrom:  &predecessor.ID,
		}
		return nil
	})
	if err != nil {
		return Session{}, err
	}
	if authenticationFailure != nil {
		return Session{}, authenticationFailure
	}

	return rotated, nil
}

// RevokeSession preserves one session row and makes repeated calls harmless.
func (service *Service) RevokeSession(ctx context.Context, sessionID uuid.UUID, reason string) error {
	if sessionID == uuid.Nil {
		return ErrInvalidSessionID
	}
	if err := service.validate(); err != nil {
		return err
	}

	revokedAt := service.clock()
	revokedReason := normalizeRevocationReason(reason)
	return service.transactions.WithTx(ctx, func(tx pgx.Tx) error {
		updated, err := service.sessions.RevokeSession(ctx, tx, sessionID, revokedAt, revokedReason)
		if err != nil {
			return err
		}
		if !updated {
			return ErrSessionNotFound
		}
		return nil
	})
}

// RevokeAllUserSessions revokes only unrevoked sessions belonging to one user
// and retains all historical rows.
func (service *Service) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID, reason string) error {
	if userID == uuid.Nil {
		return ErrInvalidUserID
	}
	if err := service.validate(); err != nil {
		return err
	}

	revokedAt := service.clock()
	revokedReason := normalizeRevocationReason(reason)
	return service.transactions.WithTx(ctx, func(tx pgx.Tx) error {
		_, err := service.sessions.RevokeAllUserSessions(ctx, tx, userID, revokedAt, revokedReason)
		return err
	})
}

func (service *Service) validate() error {
	if service.transactions == nil || service.clock == nil || service.tokenSource == nil || service.sessionIDSource == nil {
		return fmt.Errorf("%w: missing dependency", ErrInvalidServiceConfig)
	}
	if service.refreshTTL <= 0 {
		return ErrInvalidRefreshTTL
	}

	return nil
}

func normalizeRevocationReason(reason string) string {
	if strings.TrimSpace(reason) == "" {
		return "manual"
	}

	return reason
}

func calculateSessionExpiry(issuedAt time.Time, refreshTTL time.Duration) time.Time {
	return issuedAt.Add(refreshTTL)
}
