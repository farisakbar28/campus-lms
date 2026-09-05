package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	tenantAIDForCleanupTest = "19cd4773-2aeb-d614-028f-e21bf9b73d0c"
	tenantBIDForCleanupTest = "f7338b92-8230-c6fb-0bdf-34bd142e6357"
)

func TestWithTenantTxClearsTransactionLocalTenantContextOnReusedConnection(t *testing.T) {
	ctx := context.Background()
	container, err := postgrescontainer.Run(ctx, "postgres:16.14-alpine3.23",
		postgrescontainer.WithDatabase("tenant_context_test"),
		postgrescontainer.WithUsername("owner"),
		postgrescontainer.WithPassword("owner-test-only"),
	)
	if err != nil {
		t.Fatalf("start PostgreSQL container: %v", err)
	}
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Errorf("terminate PostgreSQL container: %v", err)
		}
	})

	databaseURL, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("get PostgreSQL connection string: %v", err)
	}
	pool, err := newSingleConnectionTestPool(ctx, databaseURL)
	if err != nil {
		t.Fatalf("create single-connection test pool: %v", err)
	}
	t.Cleanup(pool.Close)

	tenantABackendPID := backendPIDInTenantTransaction(t, pool, tenantAIDForCleanupTest)
	t.Logf("tenant_a_backend_pid=%d", tenantABackendPID)

	connection, err := pool.pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("acquire pooled connection after tenant A commit: %v", err)
	}
	var postTenantABackendPID int
	if err := connection.QueryRow(ctx, "SELECT pg_backend_pid()").Scan(&postTenantABackendPID); err != nil {
		connection.Release()
		t.Fatalf("read backend PID after tenant A commit: %v", err)
	}
	var postTenantASetting *string
	if err := connection.QueryRow(ctx, "SELECT current_setting('app.tenant_id', true)").Scan(&postTenantASetting); err != nil {
		connection.Release()
		t.Fatalf("read tenant setting after tenant A commit: %v", err)
	}
	connection.Release()

	if postTenantABackendPID != tenantABackendPID {
		t.Fatalf("post-tenant-A backend PID = %d, want tenant A PID %d", postTenantABackendPID, tenantABackendPID)
	}
	if postTenantASetting != nil && *postTenantASetting == tenantAIDForCleanupTest {
		t.Fatalf("post-tenant-A tenant setting retained tenant A ID %q", tenantAIDForCleanupTest)
	}
	if postTenantASetting != nil && *postTenantASetting != "" {
		t.Fatalf("post-tenant-A tenant setting = %q, want unset or empty", *postTenantASetting)
	}
	t.Logf("post_tenant_a_backend_pid=%d", postTenantABackendPID)
	t.Log("post_tenant_a_setting=unset")

	tenantBBackendPID := backendPIDInTenantTransaction(t, pool, tenantBIDForCleanupTest)
	if tenantBBackendPID != tenantABackendPID {
		t.Fatalf("tenant B backend PID = %d, want tenant A PID %d", tenantBBackendPID, tenantABackendPID)
	}
	t.Logf("tenant_b_backend_pid=%d", tenantBBackendPID)
	t.Log("same_physical_connection=true")

	if err := pool.WithTx(ctx, func(tx pgx.Tx) error {
		var tenantSetting *string
		if err := tx.QueryRow(ctx, "SELECT current_setting('app.tenant_id', true)").Scan(&tenantSetting); err != nil {
			return fmt.Errorf("read global transaction tenant setting: %w", err)
		}
		if tenantSetting != nil && *tenantSetting != "" {
			return fmt.Errorf("global transaction unexpectedly set app.tenant_id=%q", *tenantSetting)
		}
		return nil
	}); err != nil {
		t.Fatalf("global transaction: %v", err)
	}
	t.Log("global_transaction_tenant_setting=unset")
}

func newSingleConnectionTestPool(ctx context.Context, databaseURL string) (*Pool, error) {
	deadline := time.Now().Add(10 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		pool, err := NewPool(ctx, databaseURL, 0, 1)
		if err == nil {
			return pool, nil
		}
		lastErr = err
		time.Sleep(100 * time.Millisecond)
	}
	return nil, fmt.Errorf("connect single-connection test pool before deadline: %w", lastErr)
}

func backendPIDInTenantTransaction(t *testing.T, pool *Pool, tenantID string) int {
	t.Helper()
	backendPID := 0
	err := pool.WithTenantTx(context.Background(), tenantID, func(tx pgx.Tx) error {
		if err := tx.QueryRow(context.Background(), "SELECT pg_backend_pid()").Scan(&backendPID); err != nil {
			return fmt.Errorf("read backend PID: %w", err)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("run tenant transaction for %s: %v", tenantID, err)
	}
	return backendPID
}
