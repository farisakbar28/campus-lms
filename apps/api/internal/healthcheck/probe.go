// Package healthcheck probes the API liveness endpoint from inside its container.
package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const timeout = 2 * time.Second

// Probe performs a bounded GET request and returns an error for any non-200 response.
func Probe(ctx context.Context, endpoint string) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("create healthcheck request: %w", err)
	}

	client := &http.Client{Timeout: timeout}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("request healthcheck endpoint: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("healthcheck endpoint returned status %d", response.StatusCode)
	}

	return nil
}
