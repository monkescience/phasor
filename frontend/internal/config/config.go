package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	// ErrTileColorsRequired is returned when tile_colors is not configured in the config file.
	ErrTileColorsRequired = errors.New("tile_colors must be configured in the config file")
	// ErrBackendURLRequired is returned when backend_url is not configured.
	ErrBackendURLRequired = errors.New("backend_url must be configured in the config file")
)

// Config holds the frontend application configuration.
type Config struct {
	BackendURL string   `yaml:"backend_url"` // URL of the backend service
	TileColors []string `yaml:"tile_colors"` // Available colors for instance tiles based on version
	LogConfig  struct {
		Level     string `yaml:"level"`      // Log level (debug, info, warn, error)
		Format    string `yaml:"format"`     // Log format (json, text)
		AddSource bool   `yaml:"add_source"` // Include source file and line number
	} `yaml:"log_config"`
}

// Load reads configuration from the specified YAML file.
func Load(path string) (*Config, error) {
	//nolint:gosec // Config file path is expected to be provided by trusted deployment configuration
	configFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	//nolint:noinlineerr,wsl // Defer close pattern is idiomatic for resource cleanup
	defer func() {
		if closeErr := configFile.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to close config file: %w", closeErr))
		}
	}()

	var cfg Config

	decoder := yaml.NewDecoder(configFile)

	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	if cfg.BackendURL == "" {
		return nil, ErrBackendURLRequired
	}

	if len(cfg.TileColors) == 0 {
		return nil, ErrTileColorsRequired
	}

	return &cfg, nil
}
