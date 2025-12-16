package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	// ErrVersionRequired is returned when the VERSION environment variable is not set.
	ErrVersionRequired = errors.New("VERSION environment variable is required")
	// ErrConfigPathNotAbsolute is returned when the config file path is not absolute.
	ErrConfigPathNotAbsolute = errors.New("config file path must be absolute")
)

// Config holds the backend application configuration.
type Config struct {
	Version   string `yaml:"-"` // Version must be set via VERSION environment variable only
	LogConfig struct {
		Level     string `yaml:"level"`      // Log level (debug, info, warn, error)
		Format    string `yaml:"format"`     // Log format (json, text)
		AddSource bool   `yaml:"add_source"` // Include source file and line number
	} `yaml:"log_config"`
}

// Load reads configuration from the specified YAML file and environment variables.
// The VERSION environment variable is required and must be set; it cannot be configured via the config file.
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

	cfg.Version = os.Getenv("VERSION")
	if cfg.Version == "" {
		return nil, ErrVersionRequired
	}

	return &cfg, nil
}
