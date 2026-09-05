package repository

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/farisakbar28/campus-lms/apps/api/internal/database"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestClassifyDatabaseError(t *testing.T) {
	tests := []struct {
		name            string
		err             error
		wantUnavailable bool
	}{
		{"ordinary unexpected error", errors.New("scan destination mismatch"), false},
		{"PostgreSQL data error", &pgconn.PgError{Code: "22P02", Message: "invalid input syntax"}, false},
		{"PostgreSQL connection error", &pgconn.PgError{Code: "08006", Message: "connection failure"}, true},
		{"pgconn connect error", &pgconn.ConnectError{}, true},
		{"network operation error", &net.OpError{Op: "read", Net: "tcp", Err: errors.New("connection reset")}, true},
		{"context canceled", context.Canceled, false},
		{"context deadline exceeded", context.DeadlineExceeded, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := classifyDatabaseError("test operation", test.err)
			if got := errors.Is(err, database.ErrUnavailable); got != test.wantUnavailable {
				t.Errorf("database unavailable = %t, want %t", got, test.wantUnavailable)
			}
		})
	}
}
