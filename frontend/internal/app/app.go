package app

import (
	"errors"
	"fmt"
	"phasor/frontend/internal/frontend"

	"github.com/go-chi/chi/v5"
	"github.com/monkescience/vital"
)

// ErrTemplatesPathRequired is returned when TemplatesPath is not provided.
var ErrTemplatesPathRequired = errors.New("TemplatesPath is required: use WithTemplatesPath option")

// SetupRouter creates and configures the application router with all middleware and handlers.
func SetupRouter(opts ...Option) (*chi.Mux, error) {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	if options.TemplatesPath == "" {
		return nil, ErrTemplatesPathRequired
	}

	router := chi.NewRouter()
	router.Use(vital.Recovery(options.Logger))
	router.Use(vital.RequestLogger(options.Logger))
	router.Use(vital.TraceContext())

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

	healthHandler := vital.NewHealthHandler(
		vital.WithEnvironment(options.Environment),
	)
	router.Mount("/health", healthHandler)

	return router, nil
}
