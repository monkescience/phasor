package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ErrVersionRequired is returned when the VERSION environment variable is not set.
var ErrVersionRequired = errors.New("VERSION environment variable is required")

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

	cfg.Version = os.Getenv("VERSION")
	if cfg.Version == "" {
		return nil, ErrVersionRequired
	}

	return &cfg, nil
}
