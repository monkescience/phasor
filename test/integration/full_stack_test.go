package integration_test

import (
	"net/http"
	"strings"
	"testing"

	backendserver "phasor/backend/pkg/server"
	frontendserver "phasor/frontend/pkg/server"

	"github.com/monkescience/testastic"
)

func TestFullStackFlow(t *testing.T) {
	t.Parallel()

	t.Run("frontend fetches tiles from backend", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a full stack with frontend and backend servers
		backend := backendserver.NewTestServer(
			backendserver.WithTestLogger(t),
			backendserver.WithVersion("2.0.0"),
		)
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithTestLogger(t),
			frontendserver.WithBackendURL(backend.URL+"/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
			frontendserver.WithTileColors([]string{"#667eea", "#f093fb", "#4facfe"}),
		)
		testastic.NoError(t, err)

		defer frontend.Close()

		// WHEN: requesting tiles from frontend
		resp := httpGet(t, frontend.URL+"/tiles?count=3")
		defer resp.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: response contains tiles with backend version
		testastic.Equal(t, http.StatusOK, resp.StatusCode)

		body := readBody(t, resp)
		testastic.Contains(t, body, "2.0.0")
		testastic.Contains(t, body, "class=\"tile\"")
	})

	t.Run("same version always gets same color", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a frontend with multiple configured colors
		backend := backendserver.NewTestServer(
			backendserver.WithTestLogger(t),
			backendserver.WithVersion("1.2.3"),
		)
		defer backend.Close()

		colors := []string{"#667eea", "#f093fb", "#4facfe", "#43e97b"}

		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithTestLogger(t),
			frontendserver.WithBackendURL(backend.URL+"/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
			frontendserver.WithTileColors(colors),
		)
		testastic.NoError(t, err)

		defer frontend.Close()

		// WHEN: requesting tiles twice
		resp1 := httpGet(t, frontend.URL+"/tiles?count=1")
		body1 := readBody(t, resp1)
		resp1.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		resp2 := httpGet(t, frontend.URL+"/tiles?count=1")
		body2 := readBody(t, resp2)
		resp2.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: same version produces same color in both responses
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

		// GIVEN: a frontend with custom color palette
		backend := backendserver.NewTestServer(
			backendserver.WithTestLogger(t),
			backendserver.WithVersion("test-version"),
		)
		defer backend.Close()

		customColors := []string{"#ff0000", "#00ff00", "#0000ff"}

		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithTestLogger(t),
			frontendserver.WithBackendURL(backend.URL+"/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
			frontendserver.WithTileColors(customColors),
		)
		testastic.NoError(t, err)

		defer frontend.Close()

		// WHEN: requesting a tile
		resp := httpGet(t, frontend.URL+"/tiles?count=1")
		defer resp.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: tile uses one of the configured colors
		testastic.Equal(t, http.StatusOK, resp.StatusCode)

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
