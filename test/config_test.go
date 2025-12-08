package test

import (
	"os"
	"path/filepath"
	"reference-service-go/internal/config"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("load config with valid tile colors", func(t *testing.T) {
		// GIVEN
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.yaml")
		configContent := `tile_colors:
  - "#667eea"
  - "#f093fb"
  - "#4facfe"
`
		if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		t.Setenv("VERSION", "1.0.0")

		// WHEN
		cfg, err := config.Load(configPath)
		// THEN
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cfg.Version != "1.0.0" {
			t.Errorf("expected version '1.0.0', got: %v", cfg.Version)
		}

		if len(cfg.TileColors) != 3 {
			t.Errorf("expected 3 tile colors, got: %d", len(cfg.TileColors))
		}

		expectedColors := []string{"#667eea", "#f093fb", "#4facfe"}
		for i, expected := range expectedColors {
			if cfg.TileColors[i] != expected {
				t.Errorf("expected color[%d] to be %v, got: %v", i, expected, cfg.TileColors[i])
			}
		}
	})

	t.Run("load config without VERSION env var", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.yaml")
		configContent := `tile_colors:
  - "#667eea"
`
		if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		// WHEN
		cfg, err := config.Load(configPath)

		// THEN
		if err == nil {
			t.Fatal("expected error for missing VERSION env var, got nil")
		}

		if cfg != nil {
			t.Errorf("expected nil config, got: %v", cfg)
		}

		expectedErr := "VERSION environment variable is required"
		if err.Error() != expectedErr {
			t.Errorf("expected error '%v', got: '%v'", expectedErr, err.Error())
		}
	})

	t.Run("load config without tile colors", func(t *testing.T) {
		// GIVEN
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.yaml")
		configContent := `log_config:
  level: "info"
  format: "json"
  add_source: false
`
		if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		t.Setenv("VERSION", "1.0.0")

		// WHEN
		cfg, err := config.Load(configPath)

		// THEN
		if err == nil {
			t.Fatal("expected error for missing tile_colors, got nil")
		}

		if cfg != nil {
			t.Errorf("expected nil config, got: %v", cfg)
		}

		expectedErr := "tile_colors must be configured in the config file"
		if err.Error() != expectedErr {
			t.Errorf("expected error '%v', got: '%v'", expectedErr, err.Error())
		}
	})

	t.Run("load config with empty tile colors array", func(t *testing.T) {
		// GIVEN
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.yaml")
		configContent := `tile_colors: []
`
		if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		t.Setenv("VERSION", "2.0.0")

		// WHEN
		cfg, err := config.Load(configPath)

		// THEN
		if err == nil {
			t.Fatal("expected error for empty tile_colors array, got nil")
		}

		if cfg != nil {
			t.Errorf("expected nil config, got: %v", cfg)
		}

		expectedErr := "tile_colors must be configured in the config file"
		if err.Error() != expectedErr {
			t.Errorf("expected error '%v', got: '%v'", expectedErr, err.Error())
		}
	})

	t.Run("load config with non-existent file", func(t *testing.T) {
		// GIVEN
		configPath := "/nonexistent/path/config.yaml"
		t.Setenv("VERSION", "1.0.0")

		// WHEN
		cfg, err := config.Load(configPath)

		// THEN
		if err == nil {
			t.Fatal("expected error for non-existent file, got nil")
		}

		if cfg != nil {
			t.Errorf("expected nil config, got: %v", cfg)
		}
	})

	t.Run("version cannot be set in config file", func(t *testing.T) {
		// GIVEN
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.yaml")
		configContent := `version: "should-be-ignored"
tile_colors:
  - "#667eea"
`
		if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		t.Setenv("VERSION", "env-version")

		// WHEN
		cfg, err := config.Load(configPath)
		// THEN
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cfg.Version != "env-version" {
			t.Errorf("expected version from env 'env-version', got: %v", cfg.Version)
		}
	})

	t.Run("load config with many tile colors", func(t *testing.T) {
		// GIVEN
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.yaml")
		configContent := `tile_colors:
  - "#667eea"
  - "#f093fb"
  - "#4facfe"
  - "#43e97b"
  - "#fa709a"
  - "#feca57"
  - "#ff6348"
  - "#1dd1a1"
`
		if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		t.Setenv("VERSION", "1.0.0")

		// WHEN
		cfg, err := config.Load(configPath)
		// THEN
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(cfg.TileColors) != 8 {
			t.Errorf("expected 8 tile colors, got: %d", len(cfg.TileColors))
		}
	})
}
