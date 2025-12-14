// Package server provides a test server factory for integration testing.
package server

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"phasor/frontend/internal/frontend"

	"github.com/go-chi/chi/v5"
	"github.com/monkescience/vital"
)

// ErrTemplatesPathRequired is returned when TemplatesPath is not provided.
var ErrTemplatesPathRequired = errors.New("TemplatesPath is required: use WithTemplatesPath option")

// NewTestServer creates a fully configured test server with the same middleware
// and routing as production. Returns an httptest.Server ready for integration tests.
// TemplatesPath must be provided via WithTemplatesPath option.
func NewTestServer(opts ...Option) (*httptest.Server, error) {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	if options.TemplatesPath == "" {
		return nil, ErrTemplatesPathRequired
	}

	router := chi.NewRouter()

	// Apply same middleware as production
	router.Use(vital.Recovery(options.Logger))
	router.Use(vital.RequestLogger(options.Logger))
	router.Use(vital.TraceContext())

	// Wire up frontend handler exactly like production
	frontendHandler, err := frontend.NewFrontendHandler(
		options.TemplatesPath,
		options.BackendURL,
		options.TileColors,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create frontend handler: %w", err)
	}

	router.Get("/", frontendHandler.IndexHandler)
	router.Get("/tiles", frontendHandler.TilesHandler)

	// Health endpoints
	healthHandler := vital.NewHealthHandler(
		vital.WithEnvironment("test"),
	)
	router.Mount("/health", healthHandler)

	return httptest.NewServer(router), nil
}
