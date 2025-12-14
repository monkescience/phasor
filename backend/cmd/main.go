package main

import (
	"flag"
	"log"
	"log/slog"
	"phasor/backend/internal/config"
	instanceapi "phasor/backend/internal/instance"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/monkescience/vital"
)

const (
	serverPort         = 8080
	serverReadTimeout  = 10 * time.Second
	serverWriteTimeout = 10 * time.Second
	serverIdleTimeout  = 120 * time.Second
	shutdownTimeout    = 20 * time.Second
)

func setupLogger(cfg *config.Config) *slog.Logger {
	logConfig := vital.LogConfig{
		Level:     cfg.LogConfig.Level,
		Format:    cfg.LogConfig.Format,
		AddSource: cfg.LogConfig.AddSource,
	}

	handler, err := vital.NewHandlerFromConfig(logConfig, vital.WithBuiltinKeys())
	if err != nil {
		log.Fatalf("failed to create logger handler: %v", err)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

func setupRouter(logger *slog.Logger, cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()

	// Add vital middleware
	router.Use(vital.Recovery(logger))
	router.Use(vital.RequestLogger(logger))
	router.Use(vital.TraceContext())

	// Instance API handler
	instanceHandler := instanceapi.NewInstanceHandler(cfg.Version)
	instanceapi.HandlerFromMux(instanceHandler, router)

	// Add vital health endpoints
	healthHandler := vital.NewHealthHandler(
		vital.WithVersion(cfg.Version),
		vital.WithEnvironment("production"),
	)
	router.Mount("/health", healthHandler)

	return router
}

func main() {
	configPath := flag.String("config", "/config/config.yaml", "Path to the configuration file")

	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := setupLogger(cfg)

	router := setupRouter(logger, cfg)

	// Create vital server with configuration options
	server := vital.NewServer(
		router,
		vital.WithPort(serverPort),
		vital.WithReadTimeout(serverReadTimeout),
		vital.WithWriteTimeout(serverWriteTimeout),
		vital.WithIdleTimeout(serverIdleTimeout),
		vital.WithShutdownTimeout(shutdownTimeout),
		vital.WithLogger(logger),
	)

	// Run server with graceful shutdown
	server.Run()
}
