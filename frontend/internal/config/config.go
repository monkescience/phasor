package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	// ErrTileColorsRequired is returned when tile_colors is not configured in the config file.
	ErrTileColorsRequired = errors.New("tile_colors must be configured in the config file")
	// ErrBackendURLRequired is returned when backend_url is not configured.
	ErrBackendURLRequired = errors.New("backend_url must be configured in the config file")
	// ErrConfigPathNotAbsolute is returned when the config file path is not absolute.
	ErrConfigPathNotAbsolute = errors.New("config file path must be absolute")
	// ErrEnvironmentRequired is returned when environment is not configured in the config file.
	ErrEnvironmentRequired = errors.New("environment must be configured in the config file")
)

// Config holds the frontend application configuration.
type Config struct {
	BackendURL  string   `yaml:"backend_url"` // URL of the backend service
	Environment string   `yaml:"environment"` // Environment name (e.g., local, dev, staging, prod)
	TileColors  []string `yaml:"tile_colors"` // Colors for instance tiles
	LogConfig   struct {
		Level     string `yaml:"level"`      // Log level (debug, info, warn, error)
		Format    string `yaml:"format"`     // Log format (json, text)
		AddSource bool   `yaml:"add_source"` // Include source file and line number
	} `yaml:"log_config"`
}

// Load reads configuration from the specified YAML file.
func Load(path string) (*Config, error) {
	cleanPath := filepath.Clean(path)
	if !filepath.IsAbs(cleanPath) {
		return nil, fmt.Errorf("%w: %s", ErrConfigPathNotAbsolute, path)
	}

	configFile, err := os.Open(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	defer func() {
		closeErr := configFile.Close()
		if closeErr != nil {
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

	if cfg.Environment == "" {
		return nil, ErrEnvironmentRequired
	}

	if len(cfg.TileColors) == 0 {
		return nil, ErrTileColorsRequired
	}

	return &cfg, nil
}
