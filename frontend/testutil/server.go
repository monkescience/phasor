// Package testutil provides test utilities for integration testing.
package testutil

import (
	"fmt"
	"log/slog"
	"net/http/httptest"
	"phasor/frontend/internal/app"
	"phasor/frontend/internal/config"
)

// NewTestServer creates a fully configured test server with the same middleware
// and routing as production. Returns an httptest.Server ready for integration tests.
func NewTestServer(
	backendURL string,
	tileColors []string,
	templatesPath string,
	logger *slog.Logger,
) (*httptest.Server, error) {
	cfg := &config.Config{
		BackendURL:  backendURL,
		Environment: "test",
		TileColors:  tileColors,
	}

	router, err := app.SetupRouter(cfg, templatesPath, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to setup router: %w", err)
	}

	return httptest.NewServer(router), nil
}
