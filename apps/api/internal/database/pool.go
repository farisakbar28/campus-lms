// Package database owns PostgreSQL pool and tenant transaction lifecycle.
package database

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const startupPingTimeout = 5 * time.Second

// ErrUnavailable marks database connectivity failures safe for HTTP classification.
var ErrUnavailable = errors.New("database unavailable")

// Pool is one reusable PostgreSQL connection pool.
type Pool struct {
	pool *pgxpool.Pool
}

// NewPool creates and verifies a pool before the HTTP server accepts traffic.
func NewPool(ctx context.Context, databaseURL string, minConns, maxConns int32) (*Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database pool configuration: %w", err)
	}
	poolConfig.MinConns = minConns
	poolConfig.MaxConns = maxConns

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("%w: create database pool: %w", ErrUnavailable, err)
	}

	databasePool := &Pool{pool: pool}
	pingContext, cancel := context.WithTimeout(ctx, startupPingTimeout)
	defer cancel()
	if err := databasePool.Ping(pingContext); err != nil {
		pool.Close()
		return nil, fmt.Errorf("startup database ping: %w", err)
	}

	return databasePool, nil
}

// Ping checks the existing pool without allocating a new connection pool.
func (p *Pool) Ping(ctx context.Context) error {
	if err := p.pool.Ping(ctx); err != nil {
		return fmt.Errorf("%w: ping database: %w", ErrUnavailable, err)
	}

	return nil
}

// Close releases every connection after HTTP shutdown has drained requests.
func (p *Pool) Close() {
	p.pool.Close()
}

// WithTenantTx scopes a trusted tenant ID to exactly one transaction.
func (p *Pool) WithTenantTx(ctx context.Context, tenantID string, fn func(pgx.Tx) error) (err error) {
	return p.withTx(ctx, "tenant", func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, "SELECT set_config('app.tenant_id', $1, true)", tenantID); err != nil {
			return fmt.Errorf("set tenant transaction context: %w", err)
		}
		return nil
	}, fn)
}

// WithTx runs work in a transaction without setting tenant context.
// Global tables such as auth_sessions must use this boundary instead of
// WithTenantTx because their authorization is not tenant-scoped.
func (p *Pool) WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
	return p.withTx(ctx, "global", nil, fn)
}

func (p *Pool) withTx(ctx context.Context, scope string, setup func(pgx.Tx) error, fn func(pgx.Tx) error) (err error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return classifyTransactionError("begin "+scope+" transaction", err)
	}

	committed := false
	defer func() {
		if committed {
			return
		}
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			err = errors.Join(err, fmt.Errorf("rollback %s transaction: %w", scope, rollbackErr))
		}
	}()

	if setup != nil {
		if err = setup(tx); err != nil {
			return err
		}
	}
	if err = fn(tx); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return classifyTransactionError("commit "+scope+" transaction", err)
	}

	committed = true
	return nil
}

func classifyTransactionError(operation string, err error) error {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%s: %w", operation, err)
	}

	if isTransactionUnavailable(err) {
		return fmt.Errorf("%w: %s: %w", ErrUnavailable, operation, err)
	}

	return fmt.Errorf("%s: %w", operation, err)
}

func isTransactionUnavailable(err error) bool {
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
