package database

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestClassifyTransactionError(t *testing.T) {
	tests := []struct {
		name            string
		err             error
		wantCanceled    bool
		wantDeadline    bool
		wantUnavailable bool
	}{
		{"context canceled", context.Canceled, true, false, false},
		{"context deadline exceeded", context.DeadlineExceeded, false, true, false},
		{"ordinary unexpected error", errors.New("transaction state mismatch"), false, false, false},
		{"PostgreSQL connection error", &pgconn.PgError{Code: "08006", Message: "connection failure"}, false, false, true},
		{"pgconn connect error", &pgconn.ConnectError{}, false, false, true},
		{"network operation error", &net.OpError{Op: "read", Net: "tcp", Err: errors.New("connection reset")}, false, false, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := classifyTransactionError("begin tenant transaction", test.err)
			if got := errors.Is(err, context.Canceled); got != test.wantCanceled {
				t.Errorf("context canceled = %t, want %t", got, test.wantCanceled)
			}
			if got := errors.Is(err, context.DeadlineExceeded); got != test.wantDeadline {
				t.Errorf("context deadline exceeded = %t, want %t", got, test.wantDeadline)
			}
			if got := errors.Is(err, ErrUnavailable); got != test.wantUnavailable {
				t.Errorf("database unavailable = %t, want %t", got, test.wantUnavailable)
			}
		})
	}
}
