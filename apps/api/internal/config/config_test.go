package config

import (
	"log/slog"
	"strings"
	"testing"
	"time"
)

const validJWTSecret = "01234567890123456789012345678901"

func TestLoad(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_PORT", "8080")
	t.Setenv("APP_LOG_LEVEL", "warn")
	t.Setenv("APP_SHUTDOWN_TIMEOUT", "3s")
	t.Setenv("JWT_SECRET", validJWTSecret)
	t.Setenv("JWT_ACCESS_TTL", "15m")
	t.Setenv("JWT_REFRESH_TTL", "168h")
	t.Setenv("DATABASE_URL", "postgres://user:password@postgres:5432/campus_lms?sslmode=disable")
	t.Setenv("DB_MAX_CONNS", "10")
	t.Setenv("DB_MIN_CONNS", "2")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Env != "test" {
		t.Errorf("Env = %q, want %q", cfg.Env, "test")
	}
	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want 8080", cfg.Port)
	}
	if cfg.LogLevel != slog.LevelWarn {
		t.Errorf("LogLevel = %v, want %v", cfg.LogLevel, slog.LevelWarn)
	}
	if cfg.ShutdownTimeout != 3*time.Second {
		t.Errorf("ShutdownTimeout = %v, want 3s", cfg.ShutdownTimeout)
	}
	if string(cfg.JWTSecret) != validJWTSecret {
		t.Errorf("JWTSecret length = %d, want %d", len(cfg.JWTSecret), len(validJWTSecret))
	}
	if cfg.AccessTTL != 15*time.Minute {
		t.Errorf("AccessTTL = %v, want 15m", cfg.AccessTTL)
	}
	if cfg.RefreshTTL != 168*time.Hour {
		t.Errorf("RefreshTTL = %v, want 168h", cfg.RefreshTTL)
	}
	if cfg.DBMaxConns != 10 || cfg.DBMinConns != 2 {
		t.Errorf("database connections = min %d max %d, want min 2 max 10", cfg.DBMinConns, cfg.DBMaxConns)
	}
}

func TestLoadRequiresTLSForProductionDatabase(t *testing.T) {
	tests := []struct {
		name      string
		database  string
		wantError string
	}{
		{
			name:      "missing sslmode",
			database:  "postgres://db.example/campus_lms",
			wantError: "production requires sslmode",
		},
		{
			name:      "disabled",
			database:  "postgres://db.example/campus_lms?sslmode=disable",
			wantError: "production requires sslmode",
		},
		{
			name:      "prefer",
			database:  "postgres://db.example/campus_lms?sslmode=prefer",
			wantError: "production requires sslmode",
		},
		{
			name:      "unknown",
			database:  "postgres://db.example/campus_lms?sslmode=unknown",
			wantError: "production requires sslmode",
		},
		{
			name:     "require",
			database: "postgres://db.example/campus_lms?sslmode=require",
		},
		{
			name:     "verify-ca",
			database: "postgres://db.example/campus_lms?sslmode=verify-ca",
		},
		{
			name:     "verify-full",
			database: "postgres://db.example/campus_lms?sslmode=verify-full",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setValidEnvironment(t)
			t.Setenv("APP_ENV", "production")
			t.Setenv("DATABASE_URL", test.database)

			_, err := Load()
			if test.wantError == "" {
				if err != nil {
					t.Fatalf("Load() error = %v, want nil", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("Load() error = %v, want substring %q", err, test.wantError)
			}
			if strings.Contains(err.Error(), "password") || strings.Contains(err.Error(), "db.example") {
				t.Errorf("Load() error exposes database URL details: %v", err)
			}
		})
	}
}

func TestLoadAcceptsOnlyKnownEnvironmentValues(t *testing.T) {
	for _, environment := range []string{"development", "test", "production"} {
		t.Run("accepts "+environment, func(t *testing.T) {
			setValidEnvironment(t)
			t.Setenv("APP_ENV", environment)
			if environment == "production" {
				t.Setenv("DATABASE_URL", "postgres://db.example/campus_lms?sslmode=require")
			}

			if _, err := Load(); err != nil {
				t.Fatalf("Load() error = %v, want nil", err)
			}
		})
	}

	for _, environment := range []string{"prod", "Production", "staging", "production-with-typo"} {
		t.Run("rejects "+environment, func(t *testing.T) {
			setValidEnvironment(t)
			t.Setenv("APP_ENV", environment)
			// A non-production TLS mode must not be accepted merely because an
			// environment typo avoids the exact production string.
			t.Setenv("DATABASE_URL", "postgres://db.example/campus_lms?sslmode=disable")

			_, err := Load()
			if err == nil || !strings.Contains(err.Error(), "validate APP_ENV") {
				t.Fatalf("Load() error = %v, want APP_ENV validation failure", err)
			}
			if strings.Contains(err.Error(), "db.example") {
				t.Errorf("Load() error exposes database URL details: %v", err)
			}
		})
	}
}

func TestLoadFailsWhenRequiredVariableIsMissing(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_PORT", "8080")
	t.Setenv("APP_LOG_LEVEL", "info")
	t.Setenv("APP_SHUTDOWN_TIMEOUT", "")
	t.Setenv("JWT_SECRET", validJWTSecret)
	t.Setenv("JWT_ACCESS_TTL", "15m")
	t.Setenv("JWT_REFRESH_TTL", "168h")
	t.Setenv("DATABASE_URL", "postgres://user:password@postgres:5432/campus_lms?sslmode=disable")
	t.Setenv("DB_MAX_CONNS", "10")
	t.Setenv("DB_MIN_CONNS", "2")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want missing required variable error")
	}
	if !strings.Contains(err.Error(), "APP_SHUTDOWN_TIMEOUT") {
		t.Errorf("Load() error = %q, want APP_SHUTDOWN_TIMEOUT", err)
	}
}

func TestLoadValidatesAccessTokenConfiguration(t *testing.T) {
	tests := []struct {
		name       string
		secret     string
		accessTTL  string
		want       string
		wantAccess bool
	}{
		{name: "missing secret", secret: "", accessTTL: "15m", want: "required environment variable JWT_SECRET is missing"},
		{name: "blank secret", secret: "   ", accessTTL: "15m", want: "required environment variable JWT_SECRET is missing"},
		{name: "too short secret", secret: "short", accessTTL: "15m", want: "validate JWT_SECRET"},
		{name: "missing access ttl", secret: validJWTSecret, accessTTL: "", want: "required environment variable JWT_ACCESS_TTL is missing"},
		{name: "malformed access ttl", secret: validJWTSecret, accessTTL: "not-a-duration", want: "parse JWT_ACCESS_TTL"},
		{name: "zero access ttl", secret: validJWTSecret, accessTTL: "0s", want: "JWT_ACCESS_TTL: must be greater than zero"},
		{name: "negative access ttl", secret: validJWTSecret, accessTTL: "-1s", want: "JWT_ACCESS_TTL: must be greater than zero"},
		{name: "over maximum access ttl", secret: validJWTSecret, accessTTL: "1h1s", want: "JWT_ACCESS_TTL: must not exceed 1h0m0s"},
		{name: "valid 15m", secret: validJWTSecret, accessTTL: "15m", wantAccess: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setValidEnvironment(t)
			t.Setenv("JWT_SECRET", test.secret)
			t.Setenv("JWT_ACCESS_TTL", test.accessTTL)

			cfg, err := Load()
			if test.wantAccess {
				if err != nil {
					t.Fatalf("Load() error = %v, want nil", err)
				}
				if cfg.AccessTTL != 15*time.Minute {
					t.Errorf("AccessTTL = %v, want 15m", cfg.AccessTTL)
				}
				return
			}

			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Load() error = %v, want substring %q", err, test.want)
			}
			if test.secret != "" && strings.Contains(err.Error(), test.secret) {
				t.Errorf("Load() error exposes secret value")
			}
		})
	}
}

func TestLoadValidatesRefreshTTL(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "missing", value: "", want: "required environment variable JWT_REFRESH_TTL is missing"},
		{name: "malformed", value: "not-a-duration", want: "parse JWT_REFRESH_TTL"},
		{name: "zero", value: "0s", want: "JWT_REFRESH_TTL: must be greater than zero"},
		{name: "negative", value: "-1h", want: "JWT_REFRESH_TTL: must be greater than zero"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setValidEnvironment(t)
			t.Setenv("JWT_REFRESH_TTL", test.value)

			_, err := Load()
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Load() error = %v, want substring %q", err, test.want)
			}
		})
	}
}

func TestLoadRejectsInvalidDatabaseConnectionBounds(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_PORT", "8080")
	t.Setenv("APP_LOG_LEVEL", "info")
	t.Setenv("APP_SHUTDOWN_TIMEOUT", "3s")
	t.Setenv("JWT_SECRET", validJWTSecret)
	t.Setenv("JWT_ACCESS_TTL", "15m")
	t.Setenv("JWT_REFRESH_TTL", "168h")
	t.Setenv("DATABASE_URL", "postgres://user:password@postgres:5432/campus_lms?sslmode=disable")
	t.Setenv("DB_MAX_CONNS", "2")
	t.Setenv("DB_MIN_CONNS", "3")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "DB_MIN_CONNS") {
		t.Fatalf("Load() error = %v, want DB_MIN_CONNS validation error", err)
	}
}

func setValidEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_PORT", "8080")
	t.Setenv("APP_LOG_LEVEL", "info")
	t.Setenv("APP_SHUTDOWN_TIMEOUT", "3s")
	t.Setenv("JWT_SECRET", validJWTSecret)
	t.Setenv("JWT_ACCESS_TTL", "15m")
	t.Setenv("JWT_REFRESH_TTL", "168h")
	t.Setenv("DATABASE_URL", "postgres://user:password@postgres:5432/campus_lms?sslmode=disable")
	t.Setenv("DB_MAX_CONNS", "10")
	t.Setenv("DB_MIN_CONNS", "2")
}
