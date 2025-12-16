// Package testutil provides test utilities for integration testing.
package testutil

import (
	"log/slog"
	"net/http/httptest"
	"phasor/backend/internal/app"
	"phasor/backend/internal/config"
)

// NewTestServer creates a fully configured test server with the same middleware
// and routing as production. Returns an httptest.Server ready for integration tests.
func NewTestServer(version string, logger *slog.Logger) *httptest.Server {
	cfg := &config.Config{
		Version:     version,
		Environment: "test",
	}

	router := app.SetupRouter(cfg, logger)

	return httptest.NewServer(router)
}
