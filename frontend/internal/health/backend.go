package health

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/monkescience/vital"
)

const (
	healthCheckTimeout = 2 * time.Second
)

// BackendChecker checks the health of the backend service.
type BackendChecker struct {
	client    *http.Client
	healthURL string
}

// NewBackendChecker creates a new backend health checker from the backend URL.
// It derives the health endpoint by using the base URL with /health/ready path.
func NewBackendChecker(backendURL string) (*BackendChecker, error) {
	parsed, err := url.Parse(backendURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse backend URL: %w", err)
	}

	healthURL := fmt.Sprintf("%s://%s/health/ready", parsed.Scheme, parsed.Host)

	return &BackendChecker{
		client: &http.Client{
			Timeout: healthCheckTimeout,
		},
		healthURL: healthURL,
	}, nil
}

// Name returns the name of this health check.
func (c *BackendChecker) Name() string {
	return "backend"
}

// Check performs a health check against the backend service.
func (c *BackendChecker) Check(ctx context.Context) (vital.Status, string) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.healthURL, nil)
	if err != nil {
		return vital.StatusError, fmt.Sprintf("failed to create request: %v", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return vital.StatusError, fmt.Sprintf("failed to reach backend: %v", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return vital.StatusError, fmt.Sprintf("backend returned status %d", resp.StatusCode)
	}

	return vital.StatusOK, ""
}
