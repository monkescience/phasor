package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/monkescience/testastic"

	backendserver "phasor/backend/pkg/server"
)

func TestBackendInstanceAPI(t *testing.T) {
	t.Parallel()

	t.Run("returns instance info with version", func(t *testing.T) {
		t.Parallel()

		// Given
		server := backendserver.NewTestServer(
			backendserver.WithVersion("1.2.3"),
		)
		defer server.Close()

		// When
		resp, err := http.Get(server.URL + "/instance/info")

		// Then
		testastic.NoError(t, err)
		defer resp.Body.Close()
		testastic.Equal(t, http.StatusOK, resp.StatusCode)
		testastic.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		var info map[string]any
		err = json.NewDecoder(resp.Body).Decode(&info)
		testastic.NoError(t, err)
		testastic.Equal(t, "1.2.3", info["version"])
		testastic.NotEmpty(t, info["hostname"])
		testastic.NotEmpty(t, info["go_version"])
		testastic.NotEmpty(t, info["uptime"])
		testastic.NotEmpty(t, info["timestamp"])
	})

	t.Run("returns consistent hostname across requests", func(t *testing.T) {
		t.Parallel()

		// Given
		server := backendserver.NewTestServer(
			backendserver.WithVersion("1.0.0"),
		)
		defer server.Close()

		// When
		resp1, err := http.Get(server.URL + "/instance/info")
		testastic.NoError(t, err)
		var info1 map[string]any
		json.NewDecoder(resp1.Body).Decode(&info1)
		resp1.Body.Close()

		resp2, err := http.Get(server.URL + "/instance/info")
		testastic.NoError(t, err)
		var info2 map[string]any
		json.NewDecoder(resp2.Body).Decode(&info2)
		resp2.Body.Close()

		// Then
		testastic.Equal(t, info1["hostname"], info2["hostname"])
	})

	t.Run("health live endpoint responds OK", func(t *testing.T) {
		t.Parallel()

		// Given
		server := backendserver.NewTestServer()
		defer server.Close()

		// When
		resp, err := http.Get(server.URL + "/health/live")

		// Then
		testastic.NoError(t, err)
		defer resp.Body.Close()
		testastic.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("health ready endpoint responds OK", func(t *testing.T) {
		t.Parallel()

		// Given
		server := backendserver.NewTestServer()
		defer server.Close()

		// When
		resp, err := http.Get(server.URL + "/health/ready")

		// Then
		testastic.NoError(t, err)
		defer resp.Body.Close()
		testastic.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
