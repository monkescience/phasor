package app

import (
	"fmt"
	"log/slog"
	"phasor/frontend/internal/config"
	"phasor/frontend/internal/frontend"

	"github.com/go-chi/chi/v5"
	"github.com/monkescience/vital"
)

// SetupRouter creates and configures the application router with all middleware and handlers.
func SetupRouter(cfg *config.Config, templatesPath string, logger *slog.Logger) (*chi.Mux, error) {
	router := chi.NewRouter()
	router.Use(vital.Recovery(logger))
	router.Use(vital.RequestLogger(logger))
	router.Use(vital.TraceContext())

	frontendHandler, err := frontend.NewFrontendHandler(
		templatesPath,
		cfg.BackendURL,
		cfg.TileColors,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create frontend handler: %w", err)
	}

	router.Get("/", frontendHandler.IndexHandler)
	router.Get("/tiles", frontendHandler.TilesHandler)

	healthHandler := vital.NewHealthHandler(
		vital.WithEnvironment(cfg.Environment),
	)
	router.Mount("/health", healthHandler)

	return router, nil
}
