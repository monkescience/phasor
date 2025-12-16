// Package server provides a test server factory for integration testing.
package server

import (
	"net/http/httptest"
	"phasor/backend/internal/app"
)

// NewTestServer creates a fully configured test server with the same middleware
// and routing as production. Returns an httptest.Server ready for integration tests.
func NewTestServer(opts ...Option) *httptest.Server {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	router := app.SetupRouter(
		app.WithVersion(options.Version),
		app.WithHostnameFunc(options.GetHostname),
		app.WithEnvironment("test"),
		app.WithLogger(options.Logger),
	)

	return httptest.NewServer(router)
}
