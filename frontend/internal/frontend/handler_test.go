package frontend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/monkescience/testastic"
)

func setupTestTemplates(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	indexTemplate := `<!DOCTYPE html>
<html>
<head><title>Instance Dashboard</title></head>
<body>
<h1>Instance Dashboard</h1>
<p>Count: {{.Count}}</p>
</body>
</html>`

	tilesTemplate := `{{range .Instances}}
<div class="tile" style="border-left: 6px solid {{.Color}};">
    <h3>Instance #{{.Index}}</h3>
    <div>Version: {{.Info.Version}}</div>
    <div>Hostname: {{.Info.Hostname}}</div>
    <div>Uptime: {{.Info.Uptime}}</div>
</div>
{{end}}`

	if err := os.WriteFile(filepath.Join(dir, "index.gohtml"), []byte(indexTemplate), 0o644); err != nil {
		t.Fatalf("failed to write index template: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dir, "tiles.gohtml"), []byte(tilesTemplate), 0o644); err != nil {
		t.Fatalf("failed to write tiles template: %v", err)
	}

	return dir
}

// mockInstanceResponse represents the backend instance info response for testing.
type mockInstanceResponse struct {
	Version   string `json:"version"`
	Hostname  string `json:"hostname"`
	GoVersion string `json:"go_version"`
	Uptime    string `json:"uptime"`
}

func startMockBackend(t *testing.T, version string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := mockInstanceResponse{
			Version:   version,
			Hostname:  "test-host",
			GoVersion: "go1.23",
			Uptime:    "1h0m0s",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
}

func TestFrontendHandler(t *testing.T) {
	t.Parallel()

	t.Run("index handler returns HTML", func(t *testing.T) {
		t.Parallel()

		templatesDir := setupTestTemplates(t)
		handler, err := NewFrontendHandler(templatesDir, "http://localhost:8080", []string{"#667eea"})
		testastic.NoError(t, err)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		handler.IndexHandler(rec, req)

		testastic.Equal(t, http.StatusOK, rec.Code)
		testastic.Contains(t, rec.Body.String(), "Instance Dashboard")
		testastic.Contains(t, rec.Body.String(), "<!DOCTYPE html>")
	})

	t.Run("handler creation fails with invalid template path", func(t *testing.T) {
		t.Parallel()

		handler, err := NewFrontendHandler(
			"/nonexistent/templates",
			"http://localhost:8080",
			[]string{"#667eea"},
		)

		testastic.Error(t, err)
		testastic.Nil(t, handler)
	})
}

func TestFrontendBackendIntegration(t *testing.T) {
	t.Parallel()

	t.Run("tiles handler fetches from backend", func(t *testing.T) {
		t.Parallel()

		backend := startMockBackend(t, "2.0.0")
		defer backend.Close()

		templatesDir := setupTestTemplates(t)
		handler, err := NewFrontendHandler(templatesDir, backend.URL, []string{"#667eea", "#f093fb"})
		testastic.NoError(t, err)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/tiles?count=2", nil)
		handler.TilesHandler(rec, req)

		testastic.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()
		testastic.Contains(t, body, "2.0.0")
		testastic.Contains(t, body, "Instance #1")
		testastic.Contains(t, body, "Instance #2")
		testastic.Contains(t, body, "border-left")
	})

	t.Run("tile count parameter is respected", func(t *testing.T) {
		t.Parallel()

		backend := startMockBackend(t, "1.0.0")
		defer backend.Close()

		templatesDir := setupTestTemplates(t)
		handler, err := NewFrontendHandler(templatesDir, backend.URL, []string{"#667eea"})
		testastic.NoError(t, err)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/tiles?count=5", nil)
		handler.TilesHandler(rec, req)

		testastic.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()
		for i := 1; i <= 5; i++ {
			testastic.Contains(t, body, fmt.Sprintf("Instance #%d", i))
		}
	})

	t.Run("invalid count uses default of 3", func(t *testing.T) {
		t.Parallel()

		backend := startMockBackend(t, "1.0.0")
		defer backend.Close()

		templatesDir := setupTestTemplates(t)
		handler, err := NewFrontendHandler(templatesDir, backend.URL, []string{"#667eea"})
		testastic.NoError(t, err)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/tiles?count=invalid", nil)
		handler.TilesHandler(rec, req)

		testastic.Equal(t, http.StatusOK, rec.Code)
		testastic.Contains(t, rec.Body.String(), "Instance #3")
	})

	t.Run("count is limited to maximum of 20", func(t *testing.T) {
		t.Parallel()

		backend := startMockBackend(t, "1.0.0")
		defer backend.Close()

		templatesDir := setupTestTemplates(t)
		handler, err := NewFrontendHandler(templatesDir, backend.URL, []string{"#667eea"})
		testastic.NoError(t, err)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/tiles?count=100", nil)
		handler.TilesHandler(rec, req)

		testastic.Equal(t, http.StatusOK, rec.Code)
		testastic.NotContains(t, rec.Body.String(), "Instance #21")
	})

	t.Run("shows error when backend fails", func(t *testing.T) {
		t.Parallel()

		errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer errorServer.Close()

		templatesDir := setupTestTemplates(t)
		handler, err := NewFrontendHandler(templatesDir, errorServer.URL, []string{"#667eea"})
		testastic.NoError(t, err)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/tiles?count=1", nil)
		handler.TilesHandler(rec, req)

		testastic.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()
		testastic.Contains(t, body, "error")
		testastic.Contains(t, body, "failed to fetch")
	})
}

func TestColorConsistency(t *testing.T) {
	t.Parallel()

	t.Run("same version always gets same color", func(t *testing.T) {
		t.Parallel()

		backend := startMockBackend(t, "1.2.3")
		defer backend.Close()

		colors := []string{"#667eea", "#f093fb", "#4facfe", "#43e97b"}
		templatesDir := setupTestTemplates(t)
		handler, err := NewFrontendHandler(templatesDir, backend.URL, colors)
		testastic.NoError(t, err)

		rec1 := httptest.NewRecorder()
		handler.TilesHandler(rec1, httptest.NewRequest(http.MethodGet, "/tiles?count=1", nil))

		rec2 := httptest.NewRecorder()
		handler.TilesHandler(rec2, httptest.NewRequest(http.MethodGet, "/tiles?count=1", nil))

		testastic.Equal(t, rec1.Body.String(), rec2.Body.String())

		colorFound := false
		for _, color := range colors {
			if strings.Contains(rec1.Body.String(), color) {
				colorFound = true
				break
			}
		}
		testastic.True(t, colorFound)
	})

	t.Run("color is derived from configured colors", func(t *testing.T) {
		t.Parallel()

		backend := startMockBackend(t, "test-version")
		defer backend.Close()

		customColors := []string{"#ff0000", "#00ff00", "#0000ff"}
		templatesDir := setupTestTemplates(t)
		handler, err := NewFrontendHandler(templatesDir, backend.URL, customColors)
		testastic.NoError(t, err)

		rec := httptest.NewRecorder()
		handler.TilesHandler(rec, httptest.NewRequest(http.MethodGet, "/tiles?count=1", nil))

		body := rec.Body.String()
		hasConfiguredColor := false
		for _, color := range customColors {
			if strings.Contains(body, color) {
				hasConfiguredColor = true
				break
			}
		}
		testastic.True(t, hasConfiguredColor)
	})
}
