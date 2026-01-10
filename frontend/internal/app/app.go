package app

import (
	"fmt"
	"log/slog"
	"phasor/frontend/internal/config"
	"phasor/frontend/internal/frontend"
	"phasor/frontend/internal/health"

	"github.com/go-chi/chi/v5"
	"github.com/monkescience/vital"
)

// SetupRouter creates and configures the application router with all middleware and handlers.
func SetupRouter(cfg *config.Config, templatesPath string, logger *slog.Logger) (*chi.Mux, error) {
	router := chi.NewRouter()
	router.Use(vital.Recovery(logger))

	backendChecker, err := health.NewBackendChecker(cfg.BackendURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create backend health checker: %w", err)
	}

	healthHandler := vital.NewHealthHandler(
		vital.WithEnvironment(cfg.Environment),
		vital.WithCheckers(backendChecker),
	)
	router.Mount("/health", healthHandler)

	frontendHandler, err := frontend.NewFrontendHandler(
		templatesPath,
		cfg.BackendURL,
		cfg.TileColors,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create frontend handler: %w", err)
	}

	router.Group(func(r chi.Router) {
		r.Use(vital.TraceContext())
		r.Use(vital.RequestLogger(logger))
		r.Get("/", frontendHandler.IndexHandler)
		r.Get("/tiles", frontendHandler.TilesHandler)
	})

	return router, nil
}
