package integration

import (
	"net/http"
	"strings"
	"testing"

	"github.com/monkescience/testastic"

	backendserver "phasor/backend/pkg/server"
	frontendserver "phasor/frontend/pkg/server"
)

func TestFullStackFlow(t *testing.T) {
	t.Parallel()

	t.Run("frontend fetches tiles from backend", func(t *testing.T) {
		t.Parallel()

		// Given
		backend := backendserver.NewTestServer(
			backendserver.WithVersion("2.0.0"),
		)
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithBackendURL(backend.URL+"/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
			frontendserver.WithTileColors([]string{"#667eea", "#f093fb", "#4facfe"}),
		)
		testastic.NoError(t, err)
		defer frontend.Close()

		// When
		resp, err := http.Get(frontend.URL + "/tiles?count=3")

		// Then
		testastic.NoError(t, err)
		defer resp.Body.Close()
		testastic.Equal(t, http.StatusOK, resp.StatusCode)

		body := readBody(t, resp)
		testastic.Contains(t, body, "2.0.0")
		testastic.Contains(t, body, "Instance #1")
		testastic.Contains(t, body, "Instance #2")
		testastic.Contains(t, body, "Instance #3")
	})

	t.Run("same version always gets same color", func(t *testing.T) {
		t.Parallel()

		// Given
		backend := backendserver.NewTestServer(
			backendserver.WithVersion("1.2.3"),
		)
		defer backend.Close()

		colors := []string{"#667eea", "#f093fb", "#4facfe", "#43e97b"}
		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithBackendURL(backend.URL+"/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
			frontendserver.WithTileColors(colors),
		)
		testastic.NoError(t, err)
		defer frontend.Close()

		// When
		resp1, err := http.Get(frontend.URL + "/tiles?count=1")
		testastic.NoError(t, err)
		body1 := readBody(t, resp1)
		resp1.Body.Close()

		resp2, err := http.Get(frontend.URL + "/tiles?count=1")
		testastic.NoError(t, err)
		body2 := readBody(t, resp2)
		resp2.Body.Close()

		// Then - check color consistency (uptime changes between requests, so we check color only)
		var color1, color2 string
		for _, color := range colors {
			if strings.Contains(body1, color) {
				color1 = color
			}
			if strings.Contains(body2, color) {
				color2 = color
			}
		}
		testastic.NotEmpty(t, color1)
		testastic.Equal(t, color1, color2)
	})

	t.Run("color is derived from configured colors", func(t *testing.T) {
		t.Parallel()

		// Given
		backend := backendserver.NewTestServer(
			backendserver.WithVersion("test-version"),
		)
		defer backend.Close()

		customColors := []string{"#ff0000", "#00ff00", "#0000ff"}
		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithBackendURL(backend.URL+"/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
			frontendserver.WithTileColors(customColors),
		)
		testastic.NoError(t, err)
		defer frontend.Close()

		// When
		resp, err := http.Get(frontend.URL + "/tiles?count=1")

		// Then
		testastic.NoError(t, err)
		defer resp.Body.Close()

		body := readBody(t, resp)
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
