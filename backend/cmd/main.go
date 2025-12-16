package main

import (
	"flag"
	"log"
	"phasor/backend/internal/app"
	"phasor/backend/internal/config"

	"github.com/monkescience/vital"
)

const serverPort = 8080

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

	router := app.SetupRouter(cfg, logger)

	vital.NewServer(router, vital.WithPort(serverPort), vital.WithLogger(logger)).Run()
}
