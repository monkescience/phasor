// Package server provides a test server factory for integration testing.
package server

import (
	"net/http/httptest"
	"phasor/backend/internal/app"
	"phasor/backend/internal/config"
)

// NewTestServer creates a fully configured test server with the same middleware
// and routing as production. Returns an httptest.Server ready for integration tests.
func NewTestServer(opts ...Option) *httptest.Server {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	cfg := &config.Config{
		Version:     options.Version,
		Environment: "test",
	}

	router := app.SetupRouterWithHostname(cfg, options.Logger, options.GetHostname)

	return httptest.NewServer(router)
}
