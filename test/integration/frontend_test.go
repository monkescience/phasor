package integration_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	backendserver "phasor/backend/testutil"
	frontendserver "phasor/frontend/testutil"

	"github.com/monkescience/testastic"
)

var defaultTileColors = []string{"#667eea", "#f093fb", "#4facfe", "#43e97b", "#fa709a", "#feca57", "#ff6348", "#1dd1a1"}

func TestFrontendHandler(t *testing.T) {
	t.Parallel()

	t.Run("index page returns HTML", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a frontend server connected to a backend
		backend := backendserver.NewTestServer("1.0.0", backendserver.NewTestLogger(t))
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			backend.URL+"/instance/info",
			defaultTileColors,
			templatesPath(),
			frontendserver.NewTestLogger(t),
		)
		testastic.NoError(t, err)

		defer frontend.Close()

		// WHEN: requesting the index page
		resp := httpGet(t, frontend.URL+"/")
		defer resp.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: response matches expected HTML structure
		testastic.Equal(t, http.StatusOK, resp.StatusCode)
		testastic.AssertHTML(t, testdataPath("frontend_index", "expected_response.html"), resp.Body)
	})

	t.Run("health live endpoint responds OK", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a frontend server
		backend := backendserver.NewTestServer("test-version", backendserver.NewTestLogger(t))
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			backend.URL+"/instance/info",
			defaultTileColors,
			templatesPath(),
			frontendserver.NewTestLogger(t),
		)
		testastic.NoError(t, err)

		defer frontend.Close()

		// WHEN: requesting the live health endpoint
		resp := httpGet(t, frontend.URL+"/health/live")
		defer resp.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: response matches expected JSON structure
		testastic.Equal(t, http.StatusOK, resp.StatusCode)
		testastic.AssertJSON(t, testdataPath("frontend_health_live", "expected_response.json"), resp.Body)
	})

	t.Run("health ready endpoint responds OK", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a frontend server
		backend := backendserver.NewTestServer("test-version", backendserver.NewTestLogger(t))
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			backend.URL+"/instance/info",
			defaultTileColors,
			templatesPath(),
			frontendserver.NewTestLogger(t),
		)
		testastic.NoError(t, err)

		defer frontend.Close()

		// WHEN: requesting the ready health endpoint
		resp := httpGet(t, frontend.URL+"/health/ready")
		defer resp.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: response matches expected JSON structure
		testastic.Equal(t, http.StatusOK, resp.StatusCode)
		testastic.AssertJSON(t, testdataPath("frontend_health_ready", "expected_response.json"), resp.Body)
	})
}

func TestFrontendTiles(t *testing.T) {
	t.Parallel()

	t.Run("tiles endpoint fetches from backend", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a frontend server with configured tile colors
		backend := backendserver.NewTestServer("2.0.0", backendserver.NewTestLogger(t))
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			backend.URL+"/instance/info",
			[]string{"#667eea", "#f093fb"},
			templatesPath(),
			frontendserver.NewTestLogger(t),
		)
		testastic.NoError(t, err)

		defer frontend.Close()

		// WHEN: requesting tiles with count=2
		resp := httpGet(t, frontend.URL+"/tiles?count=2")
		defer resp.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: response matches expected HTML structure
		testastic.Equal(t, http.StatusOK, resp.StatusCode)
		testastic.AssertHTML(t, testdataPath("frontend_tiles_count_2", "expected_response.html"), resp.Body)
	})

	t.Run("tile count parameter is respected", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a frontend server
		backend := backendserver.NewTestServer("1.0.0", backendserver.NewTestLogger(t))
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			backend.URL+"/instance/info",
			defaultTileColors,
			templatesPath(),
			frontendserver.NewTestLogger(t),
		)
		testastic.NoError(t, err)

		defer frontend.Close()

		// WHEN: requesting tiles with count=5
		resp := httpGet(t, frontend.URL+"/tiles?count=5")
		defer resp.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: response matches expected HTML structure with 5 tiles
		testastic.Equal(t, http.StatusOK, resp.StatusCode)
		testastic.AssertHTML(t, testdataPath("frontend_tiles_count_5", "expected_response.html"), resp.Body)
	})

	t.Run("invalid count uses default of 3", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a frontend server
		backend := backendserver.NewTestServer("test-version", backendserver.NewTestLogger(t))
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			backend.URL+"/instance/info",
			defaultTileColors,
			templatesPath(),
			frontendserver.NewTestLogger(t),
		)
		testastic.NoError(t, err)

		defer frontend.Close()

		// WHEN: requesting tiles with invalid count parameter
		resp := httpGet(t, frontend.URL+"/tiles?count=invalid")
		defer resp.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: response contains default 3 tiles
		testastic.Equal(t, http.StatusOK, resp.StatusCode)

		body := readBody(t, resp)
		// Count the number of tiles in the response (should be 3)
		tileCount := strings.Count(body, "class=\"tile\"")
		testastic.Equal(t, 3, tileCount)
	})

	t.Run("count is limited to maximum of 20", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a frontend server
		backend := backendserver.NewTestServer("test-version", backendserver.NewTestLogger(t))
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			backend.URL+"/instance/info",
			defaultTileColors,
			templatesPath(),
			frontendserver.NewTestLogger(t),
		)
		testastic.NoError(t, err)

		defer frontend.Close()

		// WHEN: requesting tiles with count exceeding maximum
		resp := httpGet(t, frontend.URL+"/tiles?count=100")
		defer resp.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: response is limited to maximum of 20 tiles
		testastic.Equal(t, http.StatusOK, resp.StatusCode)

		body := readBody(t, resp)
		testastic.NotContains(t, body, "Instance #21")
	})

	t.Run("handles backend failure gracefully", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a frontend server with unreachable backend
		frontend, err := frontendserver.NewTestServer(
			"http://localhost:59999/instance/info",
			defaultTileColors,
			templatesPath(),
			frontendserver.NewTestLogger(t),
		)
		testastic.NoError(t, err)

		defer frontend.Close()

		// WHEN: requesting tiles
		resp := httpGet(t, frontend.URL+"/tiles?count=1")
		defer resp.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: response shows error state gracefully
		testastic.Equal(t, http.StatusOK, resp.StatusCode)
		testastic.AssertHTML(t, testdataPath("frontend_tiles_error", "expected_response.html"), resp.Body)
	})
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()

	buf := new(strings.Builder)

	_, err := io.Copy(buf, resp.Body)
	testastic.NoError(t, err)

	return buf.String()
}
