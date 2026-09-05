// Package main is the entrypoint for the campus-lms API.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/farisakbar28/campus-lms/apps/api/internal/config"
	"github.com/farisakbar28/campus-lms/apps/api/internal/database"
	"github.com/farisakbar28/campus-lms/apps/api/internal/healthcheck"
	apphttp "github.com/farisakbar28/campus-lms/apps/api/internal/http"
	"github.com/farisakbar28/campus-lms/apps/api/internal/repository"
)

func main() {
	healthcheckMode := flag.Bool("healthcheck", false, "probe the local health endpoint and exit")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		bootstrapLogger().Error("load configuration", "error", err)
		os.Exit(1)
	}
	if *healthcheckMode {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		endpoint := fmt.Sprintf("http://127.0.0.1:%d/healthz", cfg.Port)
		if err := healthcheck.Probe(ctx, endpoint); err != nil {
			bootstrapLogger().Error("healthcheck failed", "error", err)
			os.Exit(1)
		}

		return
	}

	logger := newLogger(cfg)
	startupContext, cancelStartup := context.WithTimeout(context.Background(), 10*time.Second)
	databasePool, err := database.NewPool(startupContext, cfg.DatabaseURL, cfg.DBMinConns, cfg.DBMaxConns)
	cancelStartup()
	if err != nil {
		logger.Error("connect database", "error", err)
		os.Exit(1)
	}
	defer databasePool.Close()

	rosterService := repository.NewRosterService(databasePool)
	server := apphttp.NewServer(cfg.Address(), logger, databasePool, rosterService)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("HTTP server started", "address", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, apphttp.ErrServerClosed) {
			logger.Error("HTTP server stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		logger.Info("shutting down gracefully", "signal", ctx.Err())

		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("graceful shutdown failed", "error", err)
			os.Exit(1)
		}
	}
}

func newLogger(cfg config.Config) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel})).With(
		"service", "campus-api",
		"env", cfg.Env,
		"version", buildVersion(),
	)
}

func bootstrapLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, nil)).With(
		"service", "campus-api",
		"env", "unknown",
		"version", buildVersion(),
	)
}

func buildVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}

	return "devel"
}
