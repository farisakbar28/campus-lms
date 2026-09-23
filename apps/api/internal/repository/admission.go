package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrAuthenticationState = errors.New("authentication state is inactive")
	ErrTenantAdmission     = errors.New("tenant selection is unavailable")
)

const activeAuthenticationSQL = `
SELECT EXISTS (
    SELECT 1
    FROM users AS u
    JOIN auth_sessions AS s ON s.user_id = u.id
    WHERE u.id = $1
      AND u.status = 'active'
      AND s.id = $2
      AND s.revoked_at IS NULL
      AND s.expires_at > $3
)`

const activeMembershipSQL = `
SELECT m.id
FROM tenants AS t
JOIN memberships AS m ON m.tenant_id = t.id
WHERE t.id = $1
  AND t.status = 'active'
  AND t.suspended_at IS NULL
  AND m.user_id = $2
  AND m.status = 'active'`

// AdmissionService resolves one verified token identity and one untrusted URL
// tenant selector in a tenant-scoped transaction. No tenant data is returned
// unless the complete admission check succeeds.
type AdmissionService struct {
	transactions TenantTransactioner
	clock        func() time.Time
}

func NewAdmissionService(transactions TenantTransactioner) AdmissionService {
	return AdmissionService{transactions: transactions, clock: time.Now}
}

func (s AdmissionService) Admit(ctx context.Context, userID, sessionID, tenantID uuid.UUID) (uuid.UUID, error) {
	if userID == uuid.Nil || sessionID == uuid.Nil {
		return uuid.Nil, ErrAuthenticationState
	}
	if tenantID == uuid.Nil {
		return uuid.Nil, ErrTenantAdmission
	}
	if s.transactions == nil || s.clock == nil {
		return uuid.Nil, fmt.Errorf("admit tenant: missing dependency")
	}

	var membershipID uuid.UUID
	err := s.transactions.WithTenantTx(ctx, tenantID.String(), func(tx pgx.Tx) error {
		var active bool
		if err := tx.QueryRow(ctx, activeAuthenticationSQL, userID, sessionID, s.clock()).Scan(&active); err != nil {
			return classifyDatabaseError("resolve active user and session", err)
		}
		if !active {
			return ErrAuthenticationState
		}

		if err := tx.QueryRow(ctx, activeMembershipSQL, tenantID, userID).Scan(&membershipID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrTenantAdmission
			}
			return classifyDatabaseError("resolve active tenant membership", err)
		}
		if membershipID == uuid.Nil {
			return ErrTenantAdmission
		}
		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	return membershipID, nil
}
