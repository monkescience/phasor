// Package server provides a test server factory for integration testing.
package server

import (
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/monkescience/vital"

	instanceapi "phasor/backend/internal/instance"
)

// NewTestServer creates a fully configured test server with the same middleware
// and routing as production. Returns an httptest.Server ready for integration tests.
func NewTestServer(opts ...Option) *httptest.Server {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	router := chi.NewRouter()

	// Apply same middleware as production
	router.Use(vital.Recovery(options.Logger))
	router.Use(vital.RequestLogger(options.Logger))
	router.Use(vital.TraceContext())

	// Wire up handlers exactly like production
	instanceHandler := instanceapi.NewInstanceHandler(options.Version)
	instanceapi.HandlerFromMux(instanceHandler, router)

	// Health endpoints
	healthHandler := vital.NewHealthHandler(
		vital.WithVersion(options.Version),
		vital.WithEnvironment("test"),
	)
	router.Mount("/health", healthHandler)

	return httptest.NewServer(router)
}
