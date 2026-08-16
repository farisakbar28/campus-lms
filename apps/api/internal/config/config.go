// Package config loads the API configuration from environment variables.
package config

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config contains the environment-backed settings required by the Week 1 API.
type Config struct {
	Env             string
	Port            uint16
	LogLevel        slog.Level
	ShutdownTimeout time.Duration
}

// Load reads the configuration needed by the currently implemented API.
// Future feature settings are added only when their feature exists, so unused
// credentials never prevent the Week 1 health service from starting.
func Load() (Config, error) {
	env, err := required("APP_ENV")
	if err != nil {
		return Config{}, err
	}

	portValue, err := required("APP_PORT")
	if err != nil {
		return Config{}, err
	}
	port, err := parsePort(portValue)
	if err != nil {
		return Config{}, fmt.Errorf("parse APP_PORT: %w", err)
	}

	logLevelValue, err := required("APP_LOG_LEVEL")
	if err != nil {
		return Config{}, err
	}
	var logLevel slog.Level
	if err := logLevel.UnmarshalText([]byte(logLevelValue)); err != nil {
		return Config{}, fmt.Errorf("parse APP_LOG_LEVEL: %w", err)
	}

	shutdownTimeoutValue, err := required("APP_SHUTDOWN_TIMEOUT")
	if err != nil {
		return Config{}, err
	}
	shutdownTimeout, err := time.ParseDuration(shutdownTimeoutValue)
	if err != nil {
		return Config{}, fmt.Errorf("parse APP_SHUTDOWN_TIMEOUT: %w", err)
	}
	if shutdownTimeout <= 0 {
		return Config{}, fmt.Errorf("parse APP_SHUTDOWN_TIMEOUT: must be greater than zero")
	}

	return Config{
		Env:             env,
		Port:            port,
		LogLevel:        logLevel,
		ShutdownTimeout: shutdownTimeout,
	}, nil
}

// Address returns the loopback-agnostic address used by the HTTP server.
func (c Config) Address() string {
	return net.JoinHostPort("", strconv.Itoa(int(c.Port)))
}

func required(name string) (string, error) {
	value, ok := os.LookupEnv(name)
	if !ok || strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("required environment variable %s is missing", name)
	}

	return value, nil
}

func parsePort(value string) (uint16, error) {
	port, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return 0, err
	}
	if port == 0 {
		return 0, fmt.Errorf("must be between 1 and 65535")
	}

	return uint16(port), nil
}
