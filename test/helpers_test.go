package test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// testdataDir returns the absolute path to the testdata directory.
func testdataDir(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to get caller information")
	}
	return filepath.Join(filepath.Dir(filename), "testdata")
}

// templatesDir returns the absolute path to test templates.
func templatesDir(t *testing.T) string {
	t.Helper()
	return filepath.Join(testdataDir(t), "templates")
}

// expectedDir returns the absolute path to expected JSON files.
func expectedDir(t *testing.T) string {
	t.Helper()
	return filepath.Join(testdataDir(t), "expected")
}

// createTempConfig creates a temporary config file with the given content.
func createTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	return path
}

// defaultColors returns a standard set of test colors.
func defaultColors() []string {
	return []string{"#667eea", "#f093fb", "#4facfe"}
}
