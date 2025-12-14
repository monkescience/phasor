package integration

import (
	"path/filepath"
	"runtime"
)

// templatesPath returns the path to test templates directory.
func templatesPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "testdata", "templates")
}
