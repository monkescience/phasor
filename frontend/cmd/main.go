package main

import (
	"flag"
	"log"
	"path/filepath"
	"phasor/frontend/internal/app"
	"phasor/frontend/internal/config"

	"github.com/monkescience/vital"
)

const serverPort = 8081

func main() {
	configPath := flag.String("config", "/config/config.yaml", "Path to the configuration file")

	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger, err := app.SetupLogger(app.LogConfig{
		Level:     cfg.LogConfig.Level,
		Format:    cfg.LogConfig.Format,
		AddSource: cfg.LogConfig.AddSource,
	})
	if err != nil {
		log.Fatalf("failed to setup logger: %v", err)
	}

	templatesPath := filepath.Join("frontend", "internal", "frontend", "templates")

	router, err := app.SetupRouter(
		app.WithBackendURL(cfg.BackendURL),
		app.WithTileColors(cfg.TileColors),
		app.WithTemplatesPath(templatesPath),
		app.WithEnvironment("production"),
		app.WithLogger(logger),
	)
	if err != nil {
		log.Fatalf("failed to setup router: %v", err)
	}

	vital.NewServer(router, vital.WithPort(serverPort), vital.WithLogger(logger)).Run()
}
