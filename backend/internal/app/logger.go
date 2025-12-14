package app

import (
	"fmt"
	"log/slog"

	"github.com/monkescience/vital"
)

// LogConfig holds logging configuration.
type LogConfig struct {
	Level     string
	Format    string
	AddSource bool
}

// SetupLogger creates a configured slog.Logger using vital's handler.
// It also sets the logger as the default slog logger.
func SetupLogger(cfg LogConfig) (*slog.Logger, error) {
	vitalConfig := vital.LogConfig{
		Level:     cfg.Level,
		Format:    cfg.Format,
		AddSource: cfg.AddSource,
	}

	handler, err := vital.NewHandlerFromConfig(vitalConfig, vital.WithBuiltinKeys())
	if err != nil {
		return nil, fmt.Errorf("failed to create logger handler: %w", err)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger, nil
}
