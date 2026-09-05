// Package config loads the API configuration from environment variables.
package config

import (
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config contains the environment-backed settings required by the API.
type Config struct {
	Env             string
	Port            uint16
	LogLevel        slog.Level
	ShutdownTimeout time.Duration
	JWTSecret       []byte
	AccessTTL       time.Duration
	// RefreshTTL keeps the historical JWT_REFRESH_TTL environment name while
	// configuring the opaque refresh-session credential lifetime.
	RefreshTTL  time.Duration
	DatabaseURL string
	DBMaxConns  int32
	DBMinConns  int32
}

const (
	minimumJWTSecretBytes = 32
	maximumAccessTTL      = time.Hour
)

// Load reads the configuration needed by the currently implemented API.
func Load() (Config, error) {
	env, err := required("APP_ENV")
	if err != nil {
		return Config{}, err
	}
	if err := validateEnvironment(env); err != nil {
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

	jwtSecretValue, err := required("JWT_SECRET")
	if err != nil {
		return Config{}, err
	}
	jwtSecret, err := parseJWTSecret(jwtSecretValue)
	if err != nil {
		return Config{}, err
	}

	accessTTLValue, err := required("JWT_ACCESS_TTL")
	if err != nil {
		return Config{}, err
	}
	accessTTL, err := parseAccessTokenTTL(accessTTLValue)
	if err != nil {
		return Config{}, err
	}

	refreshTTLValue, err := required("JWT_REFRESH_TTL")
	if err != nil {
		return Config{}, err
	}
	refreshTTL, err := parsePositiveDuration("JWT_REFRESH_TTL", refreshTTLValue)
	if err != nil {
		return Config{}, err
	}

	databaseURL, err := required("DATABASE_URL")
	if err != nil {
		return Config{}, err
	}
	if err := validateDatabaseURL(databaseURL); err != nil {
		return Config{}, err
	}
	if env == "production" {
		if err := validateProductionDatabaseURL(databaseURL); err != nil {
			return Config{}, err
		}
	}

	dbMaxConnsValue, err := required("DB_MAX_CONNS")
	if err != nil {
		return Config{}, err
	}
	dbMaxConns, err := parseConnections("DB_MAX_CONNS", dbMaxConnsValue, false)
	if err != nil {
		return Config{}, err
	}

	dbMinConnsValue, err := required("DB_MIN_CONNS")
	if err != nil {
		return Config{}, err
	}
	dbMinConns, err := parseConnections("DB_MIN_CONNS", dbMinConnsValue, true)
	if err != nil {
		return Config{}, err
	}
	if dbMinConns > dbMaxConns {
		return Config{}, fmt.Errorf("validate database connections: DB_MIN_CONNS must not exceed DB_MAX_CONNS")
	}

	return Config{
		Env:             env,
		Port:            port,
		LogLevel:        logLevel,
		ShutdownTimeout: shutdownTimeout,
		JWTSecret:       jwtSecret,
		AccessTTL:       accessTTL,
		RefreshTTL:      refreshTTL,
		DatabaseURL:     databaseURL,
		DBMaxConns:      dbMaxConns,
		DBMinConns:      dbMinConns,
	}, nil
}

func validateEnvironment(value string) error {
	switch value {
	case "development", "test", "production":
		return nil
	default:
		return fmt.Errorf("validate APP_ENV: must be development, test, or production")
	}
}

func parseJWTSecret(value string) ([]byte, error) {
	secret := []byte(value)
	if len(secret) < minimumJWTSecretBytes {
		return nil, fmt.Errorf("validate JWT_SECRET: must be at least %d bytes", minimumJWTSecretBytes)
	}

	return secret, nil
}

func parseAccessTokenTTL(value string) (time.Duration, error) {
	duration, err := parsePositiveDuration("JWT_ACCESS_TTL", value)
	if err != nil {
		return 0, err
	}
	if duration > maximumAccessTTL {
		return 0, fmt.Errorf("parse JWT_ACCESS_TTL: must not exceed %s", maximumAccessTTL)
	}

	return duration, nil
}

func parsePositiveDuration(name, value string) (time.Duration, error) {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", name, err)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("parse %s: must be greater than zero", name)
	}

	return duration, nil
}

func validateDatabaseURL(value string) error {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("parse DATABASE_URL: must be a valid connection URL")
	}

	return nil
}

func validateProductionDatabaseURL(value string) error {
	parsed, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("parse DATABASE_URL: %w", err)
	}

	switch strings.ToLower(parsed.Query().Get("sslmode")) {
	case "require", "verify-ca", "verify-full":
		return nil
	default:
		return fmt.Errorf("parse DATABASE_URL: production requires sslmode=require, verify-ca, or verify-full")
	}
}

func parseConnections(name, value string, allowZero bool) (int32, error) {
	connections, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", name, err)
	}
	if connections < 0 || (!allowZero && connections == 0) {
		return 0, fmt.Errorf("parse %s: must be %s", name, map[bool]string{true: "zero or greater", false: "greater than zero"}[allowZero])
	}

	return int32(connections), nil
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
