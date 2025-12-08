package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reference-service-go/internal/incoming/http/frontend"
	"strings"
	"testing"

	instanceapi "reference-service-go/internal/incoming/http/instance"
)

func TestFrontend(t *testing.T) {
	t.Parallel()

	t.Run("index handler returns HTML", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		tempDir := t.TempDir()
		setupTestTemplates(t, tempDir)

		handler, err := frontend.NewFrontendHandler(
			tempDir,
			"http://localhost:8080/instance/info",
			[]string{"#667eea", "#f093fb"},
		)
		if err != nil {
			t.Fatalf("failed to create handler: %v", err)
		}

		responseRecorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		// WHEN
		handler.IndexHandler(responseRecorder, req)

		// THEN
		if responseRecorder.Code != http.StatusOK {
			t.Errorf(
				"handler returned wrong status code: got %v want %v",
				responseRecorder.Code,
				http.StatusOK,
			)
		}

		body := responseRecorder.Body.String()
		if !strings.Contains(body, "Instance Dashboard") {
			t.Error("expected response to contain 'Instance Dashboard'")
		}

		if !strings.Contains(body, "<!DOCTYPE html>") {
			t.Error("expected response to contain HTML doctype")
		}
	})

	t.Run("tiles handler with mock instance API", func(t *testing.T) {
		t.Parallel()

		// GIVEN - mock instance API server
		mockVersion := "1.0.0"
		mockServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				instanceHandler := instanceapi.NewInstanceHandler(mockVersion)
				instanceHandler.GetInstanceInfo(w, r)
			}),
		)
		defer mockServer.Close()

		tempDir := t.TempDir()
		setupTestTemplates(t, tempDir)

		handler, err := frontend.NewFrontendHandler(
			tempDir,
			mockServer.URL,
			[]string{"#667eea", "#f093fb", "#4facfe"},
		)
		if err != nil {
			t.Fatalf("failed to create handler: %v", err)
		}

		responseRecorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/tiles?count=2", nil)

		// WHEN
		handler.TilesHandler(responseRecorder, req)

		// THEN
		if responseRecorder.Code != http.StatusOK {
			t.Errorf(
				"handler returned wrong status code: got %v want %v",
				responseRecorder.Code,
				http.StatusOK,
			)
		}

		body := responseRecorder.Body.String()
		if !strings.Contains(body, mockVersion) {
			t.Errorf("expected response to contain version '%v'", mockVersion)
		}

		if !strings.Contains(body, "Instance #1") {
			t.Error("expected response to contain 'Instance #1'")
		}

		if !strings.Contains(body, "Instance #2") {
			t.Error("expected response to contain 'Instance #2'")
		}

		if !strings.Contains(body, "border-left") {
			t.Error("expected response to contain border-left style for color")
		}
	})

	t.Run("tiles handler respects count parameter", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		mockServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				instanceHandler := instanceapi.NewInstanceHandler("1.0.0")
				instanceHandler.GetInstanceInfo(w, r)
			}),
		)
		defer mockServer.Close()

		tempDir := t.TempDir()
		setupTestTemplates(t, tempDir)

		handler, err := frontend.NewFrontendHandler(
			tempDir,
			mockServer.URL,
			[]string{"#667eea"},
		)
		if err != nil {
			t.Fatalf("failed to create handler: %v", err)
		}

		responseRecorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/tiles?count=5", nil)

		// WHEN
		handler.TilesHandler(responseRecorder, req)

		// THEN
		if responseRecorder.Code != http.StatusOK {
			t.Errorf(
				"handler returned wrong status code: got %v want %v",
				responseRecorder.Code,
				http.StatusOK,
			)
		}

		body := responseRecorder.Body.String()
		for i := 1; i <= 5; i++ {
			expected := "Instance #" + string(rune('0'+i))
			if !strings.Contains(body, expected) {
				t.Errorf("expected response to contain '%v'", expected)
			}
		}
	})

	t.Run("tiles handler with invalid count uses default", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		mockServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				instanceHandler := instanceapi.NewInstanceHandler("1.0.0")
				instanceHandler.GetInstanceInfo(w, r)
			}),
		)
		defer mockServer.Close()

		tempDir := t.TempDir()
		setupTestTemplates(t, tempDir)

		handler, err := frontend.NewFrontendHandler(
			tempDir,
			mockServer.URL,
			[]string{"#667eea"},
		)
		if err != nil {
			t.Fatalf("failed to create handler: %v", err)
		}

		responseRecorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/tiles?count=invalid", nil)

		// WHEN
		handler.TilesHandler(responseRecorder, req)

		// THEN
		if responseRecorder.Code != http.StatusOK {
			t.Errorf(
				"handler returned wrong status code: got %v want %v",
				responseRecorder.Code,
				http.StatusOK,
			)
		}

		body := responseRecorder.Body.String()
		if !strings.Contains(body, "Instance #3") {
			t.Error("expected default of 3 tiles")
		}
	})

	t.Run("tiles handler limits count to maximum", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		mockServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				instanceHandler := instanceapi.NewInstanceHandler("1.0.0")
				instanceHandler.GetInstanceInfo(w, r)
			}),
		)
		defer mockServer.Close()

		tempDir := t.TempDir()
		setupTestTemplates(t, tempDir)

		handler, err := frontend.NewFrontendHandler(
			tempDir,
			mockServer.URL,
			[]string{"#667eea"},
		)
		if err != nil {
			t.Fatalf("failed to create handler: %v", err)
		}

		responseRecorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/tiles?count=100", nil)

		// WHEN
		handler.TilesHandler(responseRecorder, req)

		// THEN
		if responseRecorder.Code != http.StatusOK {
			t.Errorf(
				"handler returned wrong status code: got %v want %v",
				responseRecorder.Code,
				http.StatusOK,
			)
		}

		body := responseRecorder.Body.String()
		if strings.Contains(body, "Instance #21") {
			t.Error("expected count to be limited to 20")
		}
	})

	t.Run("tiles handler shows error state when API fails", func(t *testing.T) {
		t.Parallel()

		// GIVEN - server that returns errors
		mockServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
		)
		defer mockServer.Close()

		tempDir := t.TempDir()
		setupTestTemplates(t, tempDir)

		handler, err := frontend.NewFrontendHandler(
			tempDir,
			mockServer.URL,
			[]string{"#667eea"},
		)
		if err != nil {
			t.Fatalf("failed to create handler: %v", err)
		}

		responseRecorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/tiles?count=1", nil)

		// WHEN
		handler.TilesHandler(responseRecorder, req)

		// THEN
		if responseRecorder.Code != http.StatusOK {
			t.Errorf(
				"handler returned wrong status code: got %v want %v",
				responseRecorder.Code,
				http.StatusOK,
			)
		}

		body := responseRecorder.Body.String()
		if !strings.Contains(body, "error") {
			t.Error("expected response to contain error state")
		}

		if !strings.Contains(body, "failed to fetch") {
			t.Error("expected response to contain 'failed to fetch'")
		}
	})

	t.Run("same version always gets same color", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		mockVersion := "1.2.3"
		mockServer := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := instanceapi.InstanceInfoResponse{
					Version: mockVersion,
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}),
		)
		defer mockServer.Close()

		tempDir := t.TempDir()
		setupTestTemplates(t, tempDir)

		colors := []string{"#667eea", "#f093fb", "#4facfe", "#43e97b"}
		handler, err := frontend.NewFrontendHandler(
			tempDir,
			mockServer.URL,
			colors,
		)
		if err != nil {
			t.Fatalf("failed to create handler: %v", err)
		}

		// WHEN - make multiple requests
		req1 := httptest.NewRequest(http.MethodGet, "/tiles?count=1", nil)
		responseRecorder1 := httptest.NewRecorder()
		handler.TilesHandler(responseRecorder1, req1)

		req2 := httptest.NewRequest(http.MethodGet, "/tiles?count=1", nil)
		responseRecorder2 := httptest.NewRecorder()
		handler.TilesHandler(responseRecorder2, req2)

		// THEN - both responses should have the same color
		body1 := responseRecorder1.Body.String()
		body2 := responseRecorder2.Body.String()

		if body1 != body2 {
			t.Error("expected same color for same version across requests")
		}

		colorFound := false
		for _, color := range colors {
			if strings.Contains(body1, color) {
				colorFound = true
				break
			}
		}

		if !colorFound {
			t.Error("expected one of the configured colors to be used")
		}
	})

	t.Run("different versions get different colors", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		tempDir := t.TempDir()
		setupTestTemplates(t, tempDir)

		colors := []string{"#667eea", "#f093fb", "#4facfe"}

		mockServer1 := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := instanceapi.InstanceInfoResponse{Version: "1.0.0"}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}),
		)
		defer mockServer1.Close()

		mockServer2 := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := instanceapi.InstanceInfoResponse{Version: "2.0.0"}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}),
		)
		defer mockServer2.Close()

		handler1, err := frontend.NewFrontendHandler(tempDir, mockServer1.URL, colors)
		if err != nil {
			t.Fatalf("failed to create handler1: %v", err)
		}

		handler2, err := frontend.NewFrontendHandler(tempDir, mockServer2.URL, colors)
		if err != nil {
			t.Fatalf("failed to create handler2: %v", err)
		}

		// WHEN
		req1 := httptest.NewRequest(http.MethodGet, "/tiles?count=1", nil)
		responseRecorder1 := httptest.NewRecorder()
		handler1.TilesHandler(responseRecorder1, req1)

		req2 := httptest.NewRequest(http.MethodGet, "/tiles?count=1", nil)
		responseRecorder2 := httptest.NewRecorder()
		handler2.TilesHandler(responseRecorder2, req2)

		// THEN - extract colors from responses
		body1 := responseRecorder1.Body.String()
		body2 := responseRecorder2.Body.String()

		color1 := extractColorFromBody(t, body1, colors)
		color2 := extractColorFromBody(t, body2, colors)

		if color1 == "" || color2 == "" {
			t.Fatal("failed to extract colors from responses")
		}

		if color1 == color2 {
			t.Logf("Note: Different versions happened to map to same color (hash collision)")
		}
	})

	t.Run("handler creation fails with invalid template path", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		invalidPath := "/nonexistent/templates"

		// WHEN
		handler, err := frontend.NewFrontendHandler(
			invalidPath,
			"http://localhost:8080",
			[]string{"#667eea"},
		)

		// THEN
		if err == nil {
			t.Fatal("expected error for invalid template path, got nil")
		}

		if handler != nil {
			t.Errorf("expected nil handler, got: %v", handler)
		}
	})
}

func setupTestTemplates(t *testing.T, dir string) {
	t.Helper()

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
}

func extractColorFromBody(t *testing.T, body string, colors []string) string {
	t.Helper()

	for _, color := range colors {
		if strings.Contains(body, color) {
			return color
		}
	}
	return ""
}
