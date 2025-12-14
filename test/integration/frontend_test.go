package integration

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/monkescience/testastic"

	backendserver "phasor/backend/pkg/server"
	frontendserver "phasor/frontend/pkg/server"
)

func TestFrontendHandler(t *testing.T) {
	t.Parallel()

	t.Run("index page returns HTML", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a frontend server connected to a backend
		backend := backendserver.NewTestServer(
			backendserver.WithTestLogger(t),
			backendserver.WithVersion("1.0.0"),
		)
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithTestLogger(t),
			frontendserver.WithBackendURL(backend.URL+"/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
		)
		testastic.NoError(t, err)
		defer frontend.Close()

		// WHEN: requesting the index page
		resp, err := http.Get(frontend.URL + "/")

		// THEN: response contains HTML with dashboard content
		testastic.NoError(t, err)
		defer resp.Body.Close()
		testastic.Equal(t, http.StatusOK, resp.StatusCode)

		body := readBody(t, resp)
		testastic.Contains(t, body, "<!DOCTYPE html>")
		testastic.Contains(t, body, "Instance Dashboard")
	})

	t.Run("health endpoint responds OK", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a frontend server
		backend := backendserver.NewTestServer(backendserver.WithTestLogger(t))
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithTestLogger(t),
			frontendserver.WithBackendURL(backend.URL+"/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
		)
		testastic.NoError(t, err)
		defer frontend.Close()

		// WHEN: requesting the health endpoint
		resp, err := http.Get(frontend.URL + "/health/live")

		// THEN: response status is OK
		testastic.NoError(t, err)
		defer resp.Body.Close()
		testastic.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestFrontendTiles(t *testing.T) {
	t.Parallel()

	t.Run("tiles endpoint fetches from backend", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a frontend server with configured tile colors
		backend := backendserver.NewTestServer(
			backendserver.WithTestLogger(t),
			backendserver.WithVersion("2.0.0"),
		)
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithTestLogger(t),
			frontendserver.WithBackendURL(backend.URL+"/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
			frontendserver.WithTileColors([]string{"#667eea", "#f093fb"}),
		)
		testastic.NoError(t, err)
		defer frontend.Close()

		// WHEN: requesting tiles with count=2
		resp, err := http.Get(frontend.URL + "/tiles?count=2")

		// THEN: response contains two tiles with backend version
		testastic.NoError(t, err)
		defer resp.Body.Close()
		testastic.Equal(t, http.StatusOK, resp.StatusCode)

		body := readBody(t, resp)
		testastic.Contains(t, body, "2.0.0")
		testastic.Contains(t, body, "Instance #1")
		testastic.Contains(t, body, "Instance #2")
		testastic.Contains(t, body, "border-left")
	})

	t.Run("tile count parameter is respected", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a frontend server
		backend := backendserver.NewTestServer(
			backendserver.WithTestLogger(t),
			backendserver.WithVersion("1.0.0"),
		)
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithTestLogger(t),
			frontendserver.WithBackendURL(backend.URL+"/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
		)
		testastic.NoError(t, err)
		defer frontend.Close()

		// WHEN: requesting tiles with count=5
		resp, err := http.Get(frontend.URL + "/tiles?count=5")

		// THEN: response contains exactly 5 tiles
		testastic.NoError(t, err)
		defer resp.Body.Close()
		testastic.Equal(t, http.StatusOK, resp.StatusCode)

		body := readBody(t, resp)
		testastic.Contains(t, body, "Instance #1")
		testastic.Contains(t, body, "Instance #2")
		testastic.Contains(t, body, "Instance #3")
		testastic.Contains(t, body, "Instance #4")
		testastic.Contains(t, body, "Instance #5")
	})

	t.Run("invalid count uses default of 3", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a frontend server
		backend := backendserver.NewTestServer(backendserver.WithTestLogger(t))
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithTestLogger(t),
			frontendserver.WithBackendURL(backend.URL+"/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
		)
		testastic.NoError(t, err)
		defer frontend.Close()

		// WHEN: requesting tiles with invalid count parameter
		resp, err := http.Get(frontend.URL + "/tiles?count=invalid")

		// THEN: response contains default 3 tiles
		testastic.NoError(t, err)
		defer resp.Body.Close()
		testastic.Equal(t, http.StatusOK, resp.StatusCode)

		body := readBody(t, resp)
		testastic.Contains(t, body, "Instance #3")
	})

	t.Run("count is limited to maximum of 20", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a frontend server
		backend := backendserver.NewTestServer(backendserver.WithTestLogger(t))
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithTestLogger(t),
			frontendserver.WithBackendURL(backend.URL+"/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
		)
		testastic.NoError(t, err)
		defer frontend.Close()

		// WHEN: requesting tiles with count exceeding maximum
		resp, err := http.Get(frontend.URL + "/tiles?count=100")

		// THEN: response is limited to maximum of 20 tiles
		testastic.NoError(t, err)
		defer resp.Body.Close()
		testastic.Equal(t, http.StatusOK, resp.StatusCode)

		body := readBody(t, resp)
		testastic.NotContains(t, body, "Instance #21")
	})

	t.Run("handles backend failure gracefully", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a frontend server with unreachable backend
		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithTestLogger(t),
			frontendserver.WithBackendURL("http://localhost:59999/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
		)
		testastic.NoError(t, err)
		defer frontend.Close()

		// WHEN: requesting tiles
		resp, err := http.Get(frontend.URL + "/tiles?count=1")

		// THEN: response shows error state gracefully
		testastic.NoError(t, err)
		defer resp.Body.Close()
		testastic.Equal(t, http.StatusOK, resp.StatusCode)

		body := readBody(t, resp)
		testastic.Contains(t, body, "error")
	})
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()

	buf := new(strings.Builder)

	_, err := io.Copy(buf, resp.Body)
	testastic.NoError(t, err)

	return buf.String()
}
