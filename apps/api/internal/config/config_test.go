package config

import (
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_PORT", "8080")
	t.Setenv("APP_LOG_LEVEL", "warn")
	t.Setenv("APP_SHUTDOWN_TIMEOUT", "3s")

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
}

func TestLoadFailsWhenRequiredVariableIsMissing(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_PORT", "8080")
	t.Setenv("APP_LOG_LEVEL", "info")
	t.Setenv("APP_SHUTDOWN_TIMEOUT", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want missing required variable error")
	}
	if !strings.Contains(err.Error(), "APP_SHUTDOWN_TIMEOUT") {
		t.Errorf("Load() error = %q, want APP_SHUTDOWN_TIMEOUT", err)
	}
}
