// Package http provides the API's HTTP transport.
package http

import (
	"encoding/json"
	"log/slog"
	nethttp "net/http"
	"time"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
)

// ErrServerClosed is re-exported to keep main's transport dependency explicit.
var ErrServerClosed = nethttp.ErrServerClosed

// NewServer creates the Week 1 API server with timeouts that bound slow clients.
func NewServer(address string, logger *slog.Logger) *nethttp.Server {
	mux := nethttp.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz(logger))
	mux.HandleFunc("GET /readyz", readyz(logger))

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

func readyz(logger *slog.Logger) nethttp.HandlerFunc {
	return func(response nethttp.ResponseWriter, _ *nethttp.Request) {
		// Dependency checks start in Week 3; until then a live process is ready.
		logger.Info("readiness check passed")
		writeStatus(response, logger)
	}
}

func writeStatus(response nethttp.ResponseWriter, logger *slog.Logger) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(nethttp.StatusOK)
	if err := json.NewEncoder(response).Encode(map[string]string{"status": "ok"}); err != nil {
		// A response write failure cannot be returned after headers are sent.
		logger.Error("write health response", "error", err)
	}
}
