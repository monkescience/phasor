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
	t.Run("load valid config", func(t *testing.T) {
		configContent := `backend_url: "http://localhost:8080/instance/info"
tile_colors:
  - "#667eea"
  - "#f093fb"
  - "#4facfe"
`
		configPath := createTempConfig(t, configContent)

		cfg, err := Load(configPath)

		testastic.NoError(t, err)
		testastic.NotNil(t, cfg)
		testastic.Equal(t, "http://localhost:8080/instance/info", cfg.BackendURL)
		testastic.Len(t, cfg.TileColors, 3)
		testastic.SliceEqual(t, []string{"#667eea", "#f093fb", "#4facfe"}, cfg.TileColors)
	})

	t.Run("backend_url is required", func(t *testing.T) {
		configContent := `tile_colors:
  - "#667eea"
`
		configPath := createTempConfig(t, configContent)

		cfg, err := Load(configPath)

		testastic.ErrorIs(t, ErrBackendURLRequired, err)
		testastic.Nil(t, cfg)
	})

	t.Run("tile_colors is required", func(t *testing.T) {
		configContent := `backend_url: "http://localhost:8080/instance/info"
log_config:
  level: "info"
`
		configPath := createTempConfig(t, configContent)

		cfg, err := Load(configPath)

		testastic.ErrorIs(t, ErrTileColorsRequired, err)
		testastic.Nil(t, cfg)
	})

	t.Run("empty tile colors array returns error", func(t *testing.T) {
		configContent := `backend_url: "http://localhost:8080/instance/info"
tile_colors: []
`
		configPath := createTempConfig(t, configContent)

		cfg, err := Load(configPath)

		testastic.ErrorIs(t, ErrTileColorsRequired, err)
		testastic.Nil(t, cfg)
	})

	t.Run("non-existent file returns error", func(t *testing.T) {
		configPath := "/nonexistent/path/config.yaml"

		cfg, err := Load(configPath)

		testastic.Error(t, err)
		testastic.Nil(t, cfg)
		testastic.Contains(t, err.Error(), "failed to open config file")
	})

	t.Run("invalid YAML returns error", func(t *testing.T) {
		configContent := `invalid: yaml: content: [`
		configPath := createTempConfig(t, configContent)

		cfg, err := Load(configPath)

		testastic.Error(t, err)
		testastic.Nil(t, cfg)
		testastic.Contains(t, err.Error(), "failed to decode config")
	})

	t.Run("load config with many tile colors", func(t *testing.T) {
		configContent := `backend_url: "http://localhost:8080/instance/info"
tile_colors:
  - "#667eea"
  - "#f093fb"
  - "#4facfe"
  - "#43e97b"
  - "#fa709a"
  - "#feca57"
  - "#ff6348"
  - "#1dd1a1"
`
		configPath := createTempConfig(t, configContent)

		cfg, err := Load(configPath)

		testastic.NoError(t, err)
		testastic.NotNil(t, cfg)
		testastic.Len(t, cfg.TileColors, 8)
	})

	t.Run("log config is parsed correctly", func(t *testing.T) {
		configContent := `backend_url: "http://localhost:8080/instance/info"
tile_colors:
  - "#667eea"
log_config:
  level: "debug"
  format: "text"
  add_source: true
`
		configPath := createTempConfig(t, configContent)

		cfg, err := Load(configPath)

		testastic.NoError(t, err)
		testastic.NotNil(t, cfg)
		testastic.Equal(t, "debug", cfg.LogConfig.Level)
		testastic.Equal(t, "text", cfg.LogConfig.Format)
		testastic.True(t, cfg.LogConfig.AddSource)
	})
}
