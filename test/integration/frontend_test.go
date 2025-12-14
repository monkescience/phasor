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

		// Given
		backend := backendserver.NewTestServer(
			backendserver.WithVersion("1.0.0"),
		)
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithBackendURL(backend.URL+"/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
		)
		testastic.NoError(t, err)
		defer frontend.Close()

		// When
		resp, err := http.Get(frontend.URL + "/")

		// Then
		testastic.NoError(t, err)
		defer resp.Body.Close()
		testastic.Equal(t, http.StatusOK, resp.StatusCode)

		body := readBody(t, resp)
		testastic.Contains(t, body, "<!DOCTYPE html>")
		testastic.Contains(t, body, "Instance Dashboard")
	})

	t.Run("health endpoint responds OK", func(t *testing.T) {
		t.Parallel()

		// Given
		backend := backendserver.NewTestServer()
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithBackendURL(backend.URL+"/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
		)
		testastic.NoError(t, err)
		defer frontend.Close()

		// When
		resp, err := http.Get(frontend.URL + "/health/live")

		// Then
		testastic.NoError(t, err)
		defer resp.Body.Close()
		testastic.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestFrontendTiles(t *testing.T) {
	t.Parallel()

	t.Run("tiles endpoint fetches from backend", func(t *testing.T) {
		t.Parallel()

		// Given
		backend := backendserver.NewTestServer(
			backendserver.WithVersion("2.0.0"),
		)
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithBackendURL(backend.URL+"/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
			frontendserver.WithTileColors([]string{"#667eea", "#f093fb"}),
		)
		testastic.NoError(t, err)
		defer frontend.Close()

		// When
		resp, err := http.Get(frontend.URL + "/tiles?count=2")

		// Then
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

		// Given
		backend := backendserver.NewTestServer(
			backendserver.WithVersion("1.0.0"),
		)
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithBackendURL(backend.URL+"/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
		)
		testastic.NoError(t, err)
		defer frontend.Close()

		// When
		resp, err := http.Get(frontend.URL + "/tiles?count=5")

		// Then
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

		// Given
		backend := backendserver.NewTestServer()
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithBackendURL(backend.URL+"/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
		)
		testastic.NoError(t, err)
		defer frontend.Close()

		// When
		resp, err := http.Get(frontend.URL + "/tiles?count=invalid")

		// Then
		testastic.NoError(t, err)
		defer resp.Body.Close()
		testastic.Equal(t, http.StatusOK, resp.StatusCode)

		body := readBody(t, resp)
		testastic.Contains(t, body, "Instance #3")
	})

	t.Run("count is limited to maximum of 20", func(t *testing.T) {
		t.Parallel()

		// Given
		backend := backendserver.NewTestServer()
		defer backend.Close()

		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithBackendURL(backend.URL+"/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
		)
		testastic.NoError(t, err)
		defer frontend.Close()

		// When
		resp, err := http.Get(frontend.URL + "/tiles?count=100")

		// Then
		testastic.NoError(t, err)
		defer resp.Body.Close()
		testastic.Equal(t, http.StatusOK, resp.StatusCode)

		body := readBody(t, resp)
		testastic.NotContains(t, body, "Instance #21")
	})

	t.Run("handles backend failure gracefully", func(t *testing.T) {
		t.Parallel()

		// Given - no backend server running
		frontend, err := frontendserver.NewTestServer(
			frontendserver.WithBackendURL("http://localhost:59999/instance/info"),
			frontendserver.WithTemplatesPath(templatesPath()),
		)
		testastic.NoError(t, err)
		defer frontend.Close()

		// When
		resp, err := http.Get(frontend.URL + "/tiles?count=1")

		// Then
		testastic.NoError(t, err)
		defer resp.Body.Close()
		testastic.Equal(t, http.StatusOK, resp.StatusCode)

		body := readBody(t, resp)
		testastic.Contains(t, body, "error") // Error shown in version field
	})
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	buf := new(strings.Builder)
	_, err := io.Copy(buf, resp.Body)
	testastic.NoError(t, err)
	return buf.String()
}
