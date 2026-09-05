// Package http provides the API's HTTP transport.
package http

import (
	"context"
	"log/slog"
	nethttp "net/http"
	"time"
)

type readinessChecker interface {
	Ping(context.Context) error
}

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
)

// ErrServerClosed is re-exported to keep main's transport dependency explicit.
var ErrServerClosed = nethttp.ErrServerClosed

// NewServer creates the API server with timeouts that bound slow clients.
func NewServer(address string, logger *slog.Logger, readiness readinessChecker, roster rosterReader) *nethttp.Server {
	mux := nethttp.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz(logger))
	mux.HandleFunc("GET /readyz", readyz(logger, readiness))
	mux.HandleFunc("GET /course-offerings/{id}/participants", courseOfferingParticipants(roster, logger))

	return &nethttp.Server{
		Addr:              address,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
}

func healthz(logger *slog.Logger) nethttp.HandlerFunc {
	return func(response nethttp.ResponseWriter, _ *nethttp.Request) {
		logger.Info("liveness check passed")
		writeStatus(response, logger)
	}
}

func readyz(logger *slog.Logger, readiness readinessChecker) nethttp.HandlerFunc {
	return func(response nethttp.ResponseWriter, request *nethttp.Request) {
		readinessContext, cancel := context.WithTimeout(request.Context(), 2*time.Second)
		defer cancel()

		if err := readiness.Ping(readinessContext); err != nil {
			logger.Error("readiness check failed", "error", err)
			writeError(response, nethttp.StatusServiceUnavailable, "service_unavailable", "service is temporarily unavailable")
			return
		}

		logger.Info("readiness check passed")
		writeStatus(response, logger)
	}
}

func writeStatus(response nethttp.ResponseWriter, logger *slog.Logger) {
	writeJSON(response, nethttp.StatusOK, map[string]string{"status": "ok"})
}
