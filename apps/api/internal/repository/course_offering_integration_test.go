package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/farisakbar28/campus-lms/apps/api/internal/database"
	"github.com/farisakbar28/campus-lms/apps/api/internal/domain"
	"github.com/jackc/pgx/v5"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	tenantAID       = "19cd4773-2aeb-d614-028f-e21bf9b73d0c"
	tenantBID       = "f7338b92-8230-c6fb-0bdf-34bd142e6357"
	offeringAID     = "2888c021-06ae-73da-2f57-884c1dd5d059"
	offeringBID     = "e898a0dd-a8ac-8e9c-e13a-e4e4cd69a168"
	testRole        = "roster_test_app"
	fixtureNoStaff  = "a1000000-0000-0000-0000-000000000001"
	fixtureTA       = "a1000000-0000-0000-0000-000000000002"
	fixtureInactive = "a1000000-0000-0000-0000-000000000003"
	fixtureMember   = "a1000000-0000-0000-0000-000000000004"
	fixtureRole     = "a1000000-0000-0000-0000-000000000005"
	fixtureLead     = "a1000000-0000-0000-0000-000000000006"
	smallOfferingID = "a2000000-0000-0000-0000-000000000001"
)

type integrationSuite struct {
	owner        *pgx.Conn
	appPool      *database.Pool
	appURL       string
	instructorA  string
	instructorB  string
	smallStudent string
}

var repositorySuite *integrationSuite

func TestMain(m *testing.M) {
	ctx := context.Background()
	container, err := postgrescontainer.Run(ctx, "postgres:16.14-alpine3.23",
		postgrescontainer.WithDatabase("campus_lms"),
		postgrescontainer.WithUsername("owner"),
		postgrescontainer.WithPassword("owner-test-only"),
		postgrescontainer.WithOrderedInitScripts(integrationSQLFiles()...),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start repository integration database: %v\n", err)
		os.Exit(1)
	}

	code := 1
	defer func() {
		if repositorySuite != nil {
			repositorySuite.appPool.Close()
			_ = repositorySuite.owner.Close(ctx)
		}
		_ = container.Terminate(ctx)
		os.Exit(code)
	}()

	ownerURL, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "get owner connection: %v\n", err)
		return
	}
	owner, err := connectWithRetry(ctx, ownerURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect owner: %v\n", err)
		return
	}

	appURL, err := bootstrapApplicationRole(ctx, owner, ownerURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bootstrap application role: %v\n", err)
		return
	}
	if err := bootstrapAuthorizationFixtures(ctx, owner); err != nil {
		fmt.Fprintf(os.Stderr, "bootstrap authorization fixtures: %v\n", err)
		return
	}

	appPool, err := database.NewPool(ctx, appURL, 0, 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect application pool: %v\n", err)
		return
	}

	suite := &integrationSuite{owner: owner, appPool: appPool, appURL: appURL}
	if err := owner.QueryRow(ctx, `SELECT cs.user_id FROM course_staff cs WHERE cs.tenant_id = $1::uuid AND cs.course_offering_id = $2::uuid`, tenantAID, offeringAID).Scan(&suite.instructorA); err != nil {
		fmt.Fprintf(os.Stderr, "load tenant A fixture: %v\n", err)
		return
	}
	if err := owner.QueryRow(ctx, `SELECT cs.user_id FROM course_staff cs WHERE cs.tenant_id = $1::uuid AND cs.course_offering_id = $2::uuid`, tenantBID, offeringBID).Scan(&suite.instructorB); err != nil {
		fmt.Fprintf(os.Stderr, "load tenant B fixture: %v\n", err)
		return
	}
	if err := owner.QueryRow(ctx, `SELECT student_user_id FROM enrollments WHERE course_offering_id = $1::uuid`, smallOfferingID).Scan(&suite.smallStudent); err != nil {
		fmt.Fprintf(os.Stderr, "load small-roster student fixture: %v\n", err)
		return
	}
	repositorySuite = suite
	code = m.Run()
}

func TestAuthorizedRosterAndAuthorizationBoundaries(t *testing.T) {
	tests := []struct {
		name      string
		tenantID  string
		userID    string
		offering  string
		wantCount int
		wantError error
	}{
		{"authorized seeded instructor", tenantAID, repositorySuite.instructorA, offeringAID, 50, nil},
		{"cross tenant offering is hidden", tenantAID, repositorySuite.instructorA, offeringBID, 0, domain.ErrNotFound},
		{"same tenant no staff is hidden", tenantAID, fixtureNoStaff, offeringAID, 0, domain.ErrNotFound},
		{"teaching assistant is hidden", tenantAID, fixtureTA, offeringAID, 0, domain.ErrNotFound},
		{"inactive course staff is hidden", tenantAID, fixtureInactive, offeringAID, 0, domain.ErrNotFound},
		{"inactive membership is hidden", tenantAID, fixtureMember, offeringAID, 0, domain.ErrNotFound},
		{"revoked lecturer role is hidden", tenantAID, fixtureRole, offeringAID, 0, domain.ErrNotFound},
		{"lead instructor is authorized", tenantAID, fixtureLead, offeringAID, 50, nil},
		{"unknown offering is hidden", tenantAID, repositorySuite.instructorA, "ffffffff-ffff-ffff-ffff-ffffffffffff", 0, domain.ErrNotFound},
	}

	service := NewRosterService(repositorySuite.appPool)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			roster, err := service.AuthorizedRoster(context.Background(), test.tenantID, test.userID, test.offering)
			if test.wantError != nil {
				if !errors.Is(err, test.wantError) {
					t.Fatalf("error = %v, want %v", err, test.wantError)
				}
				return
			}
			if err != nil {
				t.Fatalf("AuthorizedRoster() error = %v", err)
			}
			if len(roster.Participants) != test.wantCount {
				t.Errorf("participants = %d, want %d", len(roster.Participants), test.wantCount)
			}
		})
	}
}

func TestTenantTransactionAndRLSRole(t *testing.T) {
	ctx := context.Background()
	var superuser, bypassRLS, ownsTenantTables bool
	if err := repositorySuite.owner.QueryRow(ctx, `
SELECT r.rolsuper, r.rolbypassrls,
       EXISTS (
           SELECT 1 FROM pg_class c
           JOIN pg_namespace n ON n.oid = c.relnamespace
           WHERE n.nspname = 'public'
             AND c.relname IN ('course_offerings', 'courses', 'academic_terms', 'course_staff', 'memberships', 'membership_roles', 'enrollments')
             AND c.relowner = r.oid
       )
FROM pg_roles r WHERE r.rolname = $1`, testRole).Scan(&superuser, &bypassRLS, &ownsTenantTables); err != nil {
		t.Fatalf("inspect application role: %v", err)
	}
	if superuser || bypassRLS || ownsTenantTables {
		t.Fatalf("application role properties: superuser=%t bypassrls=%t ownsTenantTables=%t", superuser, bypassRLS, ownsTenantTables)
	}

	appConnection, err := pgx.Connect(ctx, repositorySuite.appURL)
	if err != nil {
		t.Fatalf("connect test application role: %v", err)
	}
	defer func() {
		if err := appConnection.Close(ctx); err != nil {
			t.Errorf("close test application connection: %v", err)
		}
	}()

	tx, err := appConnection.Begin(ctx)
	if err != nil {
		t.Fatalf("begin RLS control transaction: %v", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			t.Errorf("rollback RLS control transaction: %v", err)
		}
	}()
	if _, err := tx.Exec(ctx, "SELECT set_config('app.tenant_id', $1, true)", tenantAID); err != nil {
		t.Fatalf("set tenant A context: %v", err)
	}
	var crossTenantRows int
	if err := tx.QueryRow(ctx, `SELECT count(*) FROM course_offerings WHERE id = $1::uuid`, offeringBID).Scan(&crossTenantRows); err != nil {
		t.Fatalf("query cross tenant RLS control: %v", err)
	}
	if crossTenantRows != 0 {
		t.Fatalf("RLS control returned %d tenant B rows for tenant A", crossTenantRows)
	}
}

func TestSequentialTenantTransactionsAndFixedDataQueryCount(t *testing.T) {
	ctx := context.Background()
	repository := CourseOfferingRepository{}
	assertCount := func(tenantID, userID, offeringID string, wantParticipants int) (domain.Roster, int) {
		t.Helper()
		count := 0
		backendPID := 0
		var roster domain.Roster
		err := repositorySuite.appPool.WithTenantTx(ctx, tenantID, func(tx pgx.Tx) error {
			if err := tx.QueryRow(ctx, `SELECT pg_backend_pid()`).Scan(&backendPID); err != nil {
				return fmt.Errorf("read physical backend PID: %w", err)
			}
			rosterResult, err := repository.AuthorizedRoster(ctx, countingQuerier{Querier: tx, count: &count}, tenantID, userID, offeringID)
			roster = rosterResult
			return err
		})
		if err != nil {
			t.Fatalf("authorized roster: %v", err)
		}
		if count != 2 {
			t.Errorf("repository data query count = %d, want 2", count)
		}
		if len(roster.Participants) != wantParticipants {
			t.Errorf("participants = %d, want %d", len(roster.Participants), wantParticipants)
		}
		t.Logf("tenant %s offering %s: repository data query count=%d, participants=%d", tenantID, offeringID, count, len(roster.Participants))
		return roster, backendPID
	}

	tenantARoster, tenantABackendPID := assertCount(tenantAID, repositorySuite.instructorA, offeringAID, 50)
	tenantBRoster, tenantBBackendPID := assertCount(tenantBID, repositorySuite.instructorB, offeringBID, 50)
	smallRoster, _ := assertCount(tenantAID, fixtureNoStaff, smallOfferingID, 1)
	if tenantABackendPID != tenantBBackendPID {
		t.Fatalf("physical backend PID changed between sequential tenants: tenant A=%d tenant B=%d", tenantABackendPID, tenantBBackendPID)
	}
	t.Logf("sequential tenant transactions reused PostgreSQL backend PID %d", tenantABackendPID)
	if tenantARoster.Offering.ID != offeringAID || tenantBRoster.Offering.ID != offeringBID {
		t.Fatalf("sequential rosters returned unexpected offerings: tenant A=%s tenant B=%s", tenantARoster.Offering.ID, tenantBRoster.Offering.ID)
	}
	tenantAStudents := make(map[string]struct{}, len(tenantARoster.Participants))
	for _, participant := range tenantARoster.Participants {
		tenantAStudents[participant.StudentUserID] = struct{}{}
	}
	for _, participant := range tenantBRoster.Participants {
		if _, found := tenantAStudents[participant.StudentUserID]; found {
			t.Fatalf("tenant B roster leaked tenant A participant %s", participant.StudentUserID)
		}
	}
	if len(smallRoster.Participants) != 1 || smallRoster.Participants[0].StudentUserID != repositorySuite.smallStudent {
		t.Fatalf("small roster participant = %#v, want seeded student %s", smallRoster.Participants, repositorySuite.smallStudent)
	}
}

type countingQuerier struct {
	Querier
	count *int
}

func (q countingQuerier) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	*q.count++
	return q.Querier.QueryRow(ctx, sql, args...)
}

func (q countingQuerier) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	*q.count++
	return q.Querier.Query(ctx, sql, args...)
}

func integrationSQLFiles() []string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("resolve integration test source path")
	}
	apiRoot := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", ".."))
	return []string{
		filepath.Join(apiRoot, "migrations", "0001_tenant_identity_schema.up.sql"),
		filepath.Join(apiRoot, "migrations", "0002_academic_core_schema.up.sql"),
		filepath.Join(apiRoot, "migrations", "0003_auth_membership_schema.up.sql"),
		filepath.Join(apiRoot, "migrations", "0004_academic_term_time_range_check.up.sql"),
		filepath.Join(apiRoot, "migrations", "0005_enrollments_active_student_lookup_index.up.sql"),
		filepath.Join(apiRoot, "testdata", "seed.sql"),
	}
}

func bootstrapApplicationRole(ctx context.Context, owner *pgx.Conn, ownerURL string) (string, error) {
	password, err := randomHex(24)
	if err != nil {
		return "", fmt.Errorf("generate test application password: %w", err)
	}
	if _, err := owner.Exec(ctx, fmt.Sprintf("CREATE ROLE %s LOGIN NOSUPERUSER NOBYPASSRLS PASSWORD '%s'", testRole, password)); err != nil {
		return "", fmt.Errorf("create test application role: %w", err)
	}
	if _, err := owner.Exec(ctx, `GRANT USAGE ON SCHEMA public TO roster_test_app`); err != nil {
		return "", fmt.Errorf("grant schema usage to test application role: %w", err)
	}
	if _, err := owner.Exec(ctx, `GRANT SELECT ON course_offerings, courses, academic_terms, course_staff, memberships, membership_roles, enrollments, users TO roster_test_app`); err != nil {
		return "", fmt.Errorf("grant test application role: %w", err)
	}
	parsed, err := url.Parse(ownerURL)
	if err != nil {
		return "", fmt.Errorf("parse owner URL: %w", err)
	}
	parsed.User = url.UserPassword(testRole, password)
	return parsed.String(), nil
}

func bootstrapAuthorizationFixtures(ctx context.Context, owner *pgx.Conn) error {
	users := []struct{ id, status string }{
		{fixtureNoStaff, "active"}, {fixtureTA, "active"}, {fixtureInactive, "active"}, {fixtureMember, "inactive"}, {fixtureRole, "active"}, {fixtureLead, "active"},
	}
	for _, user := range users {
		if _, err := owner.Exec(ctx, `INSERT INTO users (id, email, display_name, status, created_at) VALUES ($1::uuid, $2, $2, 'active', now())`, user.id, user.id+"@test.invalid"); err != nil {
			return err
		}
		if _, err := owner.Exec(ctx, `INSERT INTO memberships (id, tenant_id, user_id, status, joined_at) VALUES (md5($1)::uuid, $2::uuid, $1::uuid, $3, now())`, user.id, tenantAID, user.status); err != nil {
			return err
		}
		if _, err := owner.Exec(ctx, `INSERT INTO membership_roles (id, tenant_id, membership_id, role, granted_at, revoked_at) VALUES (md5('role-' || $1)::uuid, $2::uuid, md5($1)::uuid, 'lecturer', now(), CASE WHEN $1 = $3 THEN now() ELSE NULL END)`, user.id, tenantAID, fixtureRole); err != nil {
			return err
		}
	}
	for _, staff := range []struct {
		userID, role string
		active       bool
	}{{fixtureTA, "teaching_assistant", true}, {fixtureInactive, "instructor", false}, {fixtureMember, "instructor", true}, {fixtureRole, "instructor", true}, {fixtureLead, "lead_instructor", true}} {
		if _, err := owner.Exec(ctx, `INSERT INTO course_staff (id, tenant_id, course_offering_id, user_id, role, source, active) VALUES (md5('staff-' || $1)::uuid, $2::uuid, $3::uuid, $1::uuid, $4, 'test', $5)`, staff.userID, tenantAID, offeringAID, staff.role, staff.active); err != nil {
			return err
		}
	}
	var courseID, termID string
	if err := owner.QueryRow(ctx, `SELECT course_id, academic_term_id FROM course_offerings WHERE id = $1::uuid`, offeringAID).Scan(&courseID, &termID); err != nil {
		return err
	}
	if _, err := owner.Exec(ctx, `INSERT INTO course_offerings (id, tenant_id, external_id, course_id, academic_term_id, display_name, lms_status, created_at) VALUES ($1::uuid, $2::uuid, 'test-small-roster', $3::uuid, $4::uuid, 'Test small roster', 'published', now())`, smallOfferingID, tenantAID, courseID, termID); err != nil {
		return err
	}
	if _, err := owner.Exec(ctx, `INSERT INTO course_staff (id, tenant_id, course_offering_id, user_id, role, source, active) VALUES (md5('staff-small')::uuid, $1::uuid, $2::uuid, $3::uuid, 'instructor', 'test', true)`, tenantAID, smallOfferingID, fixtureNoStaff); err != nil {
		return err
	}
	var seededStudentID string
	if err := owner.QueryRow(ctx, `
SELECT m.user_id
FROM memberships AS m
JOIN membership_roles AS mr
  ON mr.tenant_id = m.tenant_id
 AND mr.membership_id = m.id
WHERE m.tenant_id = $1::uuid
  AND m.status = 'active'
  AND mr.role = 'student'
  AND mr.revoked_at IS NULL
ORDER BY m.user_id
LIMIT 1`, tenantAID).Scan(&seededStudentID); err != nil {
		return fmt.Errorf("select seeded small-roster student: %w", err)
	}
	var hasActiveStudentRole bool
	if err := owner.QueryRow(ctx, `
SELECT EXISTS (
    SELECT 1
    FROM memberships AS m
    JOIN membership_roles AS mr
      ON mr.tenant_id = m.tenant_id
     AND mr.membership_id = m.id
    WHERE m.tenant_id = $1::uuid
      AND m.user_id = $2::uuid
      AND m.status = 'active'
      AND mr.role = 'student'
      AND mr.revoked_at IS NULL
)`, tenantAID, seededStudentID).Scan(&hasActiveStudentRole); err != nil {
		return fmt.Errorf("verify seeded small-roster student role: %w", err)
	}
	if !hasActiveStudentRole {
		return errors.New("selected small-roster participant has no active student role")
	}
	if _, err := owner.Exec(ctx, `INSERT INTO enrollments (id, tenant_id, course_offering_id, student_user_id, external_id, status, enrolled_at) VALUES (md5('enrollment-small')::uuid, $1::uuid, $2::uuid, $3::uuid, 'test-small-enrollment', 'active', now())`, tenantAID, smallOfferingID, seededStudentID); err != nil {
		return err
	}
	return nil
}

func randomHex(bytes int) (string, error) {
	buffer := make([]byte, bytes)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}

func connectWithRetry(ctx context.Context, connectionURL string) (*pgx.Conn, error) {
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
