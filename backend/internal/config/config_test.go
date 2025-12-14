package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/monkescience/testastic"
)

func createTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	return path
}

func TestConfig(t *testing.T) {
	t.Run("load valid config with VERSION env", func(t *testing.T) {
		configContent := `log_config:
  level: "info"
  format: "json"
  add_source: false
`
		configPath := createTempConfig(t, configContent)
		t.Setenv("VERSION", "1.0.0")

		cfg, err := Load(configPath)

		testastic.NoError(t, err)
		testastic.NotNil(t, cfg)
		testastic.Equal(t, "1.0.0", cfg.Version)
		testastic.Equal(t, "info", cfg.LogConfig.Level)
		testastic.Equal(t, "json", cfg.LogConfig.Format)
		testastic.False(t, cfg.LogConfig.AddSource)
	})

	t.Run("missing VERSION env returns error", func(t *testing.T) {
		configContent := `log_config:
  level: "info"
`
		configPath := createTempConfig(t, configContent)
		// Don't set VERSION env var

		cfg, err := Load(configPath)

		testastic.ErrorIs(t, ErrVersionRequired, err)
		testastic.Nil(t, cfg)
	})

	t.Run("invalid YAML returns error", func(t *testing.T) {
		configPath := createTempConfig(t, "invalid: yaml: content: [")
		t.Setenv("VERSION", "1.0.0")

		cfg, err := Load(configPath)

		testastic.Error(t, err)
		testastic.Nil(t, cfg)
		testastic.Contains(t, err.Error(), "failed to decode config")
	})

	t.Run("non-existent file returns error", func(t *testing.T) {
		configPath := "/nonexistent/path/config.yaml"
		t.Setenv("VERSION", "1.0.0")

		cfg, err := Load(configPath)

		testastic.Error(t, err)
		testastic.Nil(t, cfg)
		testastic.Contains(t, err.Error(), "failed to open config file")
	})

	t.Run("log config is parsed correctly", func(t *testing.T) {
		configContent := `log_config:
  level: "debug"
  format: "text"
  add_source: true
`
		configPath := createTempConfig(t, configContent)
		t.Setenv("VERSION", "2.0.0")

		cfg, err := Load(configPath)

		testastic.NoError(t, err)
		testastic.NotNil(t, cfg)
		testastic.Equal(t, "debug", cfg.LogConfig.Level)
		testastic.Equal(t, "text", cfg.LogConfig.Format)
		testastic.True(t, cfg.LogConfig.AddSource)
	})

	t.Run("VERSION in yaml is ignored", func(t *testing.T) {
		configContent := `version: "should-be-ignored"
log_config:
  level: "info"
`
		configPath := createTempConfig(t, configContent)
		t.Setenv("VERSION", "env-version")

		cfg, err := Load(configPath)

		testastic.NoError(t, err)
		testastic.NotNil(t, cfg)
		testastic.Equal(t, "env-version", cfg.Version)
	})
}
