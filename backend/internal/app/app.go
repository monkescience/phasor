package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/monkescience/vital"

	instanceapi "phasor/backend/internal/instance"
)

// SetupRouter creates and configures the application router with all middleware and handlers.
func SetupRouter(opts ...Option) *chi.Mux {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	router := chi.NewRouter()
	router.Use(vital.Recovery(options.Logger))
	router.Use(vital.RequestLogger(options.Logger))
	router.Use(vital.TraceContext())

	instanceHandler := instanceapi.NewInstanceHandler(options.Version, options.GetHostname)
	instanceapi.HandlerFromMux(instanceHandler, router)

	healthHandler := vital.NewHealthHandler(
		vital.WithVersion(options.Version),
		vital.WithEnvironment(options.Environment),
	)
	router.Mount("/health", healthHandler)

	return router
}
