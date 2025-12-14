// Package server provides a test server factory for integration testing.
package server

import (
	"fmt"
	"net/http/httptest"
	"phasor/frontend/internal/app"
)

// NewTestServer creates a fully configured test server with the same middleware
// and routing as production. Returns an httptest.Server ready for integration tests.
// TemplatesPath must be provided via WithTemplatesPath option.
func NewTestServer(opts ...Option) (*httptest.Server, error) {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	router, err := app.SetupRouter(
		app.WithBackendURL(options.BackendURL),
		app.WithTileColors(options.TileColors),
		app.WithTemplatesPath(options.TemplatesPath),
		app.WithEnvironment("test"),
		app.WithLogger(options.Logger),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to setup router: %w", err)
	}

	return httptest.NewServer(router), nil
}
