// Package repository contains explicit PostgreSQL queries for LMS data.
package repository

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/farisakbar28/campus-lms/apps/api/internal/database"
	"github.com/farisakbar28/campus-lms/apps/api/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const authorizedOfferingSQL = `
SELECT
    o.id,
    COALESCE(o.external_id, ''),
    COALESCE(o.external_section_code, ''),
    COALESCE(o.display_name, ''),
    o.lms_status,
    c.id,
    COALESCE(c.code, ''),
    COALESCE(c.name, ''),
    t.id,
    COALESCE(t.code, ''),
    COALESCE(t.name, '')
FROM course_offerings AS o
JOIN courses AS c
    ON c.id = o.course_id
   AND c.tenant_id = o.tenant_id
JOIN academic_terms AS t
    ON t.id = o.academic_term_id
   AND t.tenant_id = o.tenant_id
JOIN course_staff AS cs
    ON cs.course_offering_id = o.id
   AND cs.tenant_id = o.tenant_id
JOIN memberships AS m
    ON m.tenant_id = o.tenant_id
   AND m.user_id = cs.user_id
JOIN membership_roles AS mr
    ON mr.tenant_id = m.tenant_id
   AND mr.membership_id = m.id
WHERE o.id = $1::uuid
  AND o.tenant_id = $2::uuid
  AND cs.user_id = $3::uuid
  AND cs.active = true
  AND cs.role IN ('instructor', 'lead_instructor')
  AND m.status = 'active'
  AND mr.role = 'lecturer'
  AND mr.revoked_at IS NULL`

const participantsSQL = `
SELECT
    e.student_user_id,
    e.status,
    u.display_name
FROM enrollments AS e
JOIN users AS u ON u.id = e.student_user_id
WHERE e.course_offering_id = $1::uuid
  AND e.tenant_id = $2::uuid
ORDER BY u.display_name, e.student_user_id`

// Querier is the narrow data-query surface needed by the roster repository.
type Querier interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

// TenantTransactioner keeps transaction ownership explicit outside HTTP transport.
type TenantTransactioner interface {
	WithTenantTx(context.Context, string, func(pgx.Tx) error) error
}

// CourseOfferingRepository executes the roster's two explicit data queries.
type CourseOfferingRepository struct{}

// RosterService groups transaction scope and repository work for the HTTP layer.
type RosterService struct {
	transactions TenantTransactioner
	repository   CourseOfferingRepository
}

func NewRosterService(transactions TenantTransactioner) RosterService {
	return RosterService{transactions: transactions, repository: CourseOfferingRepository{}}
}

func (s RosterService) AuthorizedRoster(ctx context.Context, tenantID, userID, offeringID string) (roster domain.Roster, err error) {
	err = s.transactions.WithTenantTx(ctx, tenantID, func(tx pgx.Tx) error {
		roster, err = s.repository.AuthorizedRoster(ctx, tx, tenantID, userID, offeringID)
		return err
	})
	if err != nil {
		return domain.Roster{}, err
	}

	return roster, nil
}

// AuthorizedRoster returns an offering only when all tenant and course authority checks pass.
func (CourseOfferingRepository) AuthorizedRoster(ctx context.Context, queries Querier, tenantID, userID, offeringID string) (domain.Roster, error) {
	var offering domain.CourseOffering
	err := queries.QueryRow(ctx, authorizedOfferingSQL, offeringID, tenantID, userID).Scan(
		&offering.ID,
		&offering.ExternalID,
		&offering.SectionCode,
		&offering.DisplayName,
		&offering.LMSStatus,
		&offering.Course.ID,
		&offering.Course.Code,
		&offering.Course.Name,
		&offering.Term.ID,
		&offering.Term.Code,
		&offering.Term.Name,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Roster{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Roster{}, classifyDatabaseError("query authorized course offering", err)
	}

	rows, err := queries.Query(ctx, participantsSQL, offeringID, tenantID)
	if err != nil {
		return domain.Roster{}, classifyDatabaseError("query course offering participants", err)
	}
	defer rows.Close()

	participants := make([]domain.Participant, 0)
	for rows.Next() {
		var participant domain.Participant
		if err := rows.Scan(&participant.StudentUserID, &participant.EnrollmentStatus, &participant.DisplayName); err != nil {
			return domain.Roster{}, classifyDatabaseError("scan course offering participant", err)
		}
		participants = append(participants, participant)
	}
	if err := rows.Err(); err != nil {
		return domain.Roster{}, classifyDatabaseError("iterate course offering participants", err)
	}

	return domain.Roster{Offering: offering, Participants: participants}, nil
}

func classifyDatabaseError(operation string, err error) error {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%s: %w", operation, err)
	}

	if isDatabaseUnavailable(err) {
		return fmt.Errorf("%w: %s: %w", database.ErrUnavailable, operation, err)
	}

	return fmt.Errorf("%s: %w", operation, err)
}

func isDatabaseUnavailable(err error) bool {
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) {
		return strings.HasPrefix(postgresError.Code, "08")
	}

	var connectError *pgconn.ConnectError
	if errors.As(err, &connectError) {
		return true
	}

	var networkError *net.OpError
	return errors.As(err, &networkError)
}
