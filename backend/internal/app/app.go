package app

import (
	"log/slog"
	"os"
	"phasor/backend/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/monkescience/vital"

	instanceapi "phasor/backend/internal/instance"
)

// SetupRouter creates and configures the application router with all middleware and handlers.
func SetupRouter(cfg *config.Config, logger *slog.Logger) *chi.Mux {
	return SetupRouterWithHostname(cfg, logger, systemHostname)
}

// SetupRouterWithHostname creates and configures the application router with a custom hostname function.
// This is primarily useful for testing with deterministic hostnames.
func SetupRouterWithHostname(cfg *config.Config, logger *slog.Logger, getHostname instanceapi.HostnameFunc) *chi.Mux {
	router := chi.NewRouter()
	router.Use(vital.Recovery(logger))

	healthHandler := vital.NewHealthHandler(
		vital.WithVersion(cfg.Version),
		vital.WithEnvironment(cfg.Environment),
	)
	router.Mount("/health", healthHandler)

	router.Group(func(r chi.Router) {
		r.Use(vital.TraceContext())
		r.Use(vital.RequestLogger(logger))

		instanceHandler := instanceapi.NewInstanceHandler(cfg.Version, getHostname)
		instanceapi.HandlerFromMux(instanceHandler, r)
	})

	return router
}

// systemHostname returns the system hostname or "unknown" if it cannot be determined.
func systemHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}

	return hostname
}
