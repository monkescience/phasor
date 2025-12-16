package integration_test

import (
	"net/http"
	"testing"

	backendserver "phasor/backend/testutil"

	"github.com/monkescience/testastic"
)

func TestBackendInstanceAPI(t *testing.T) {
	t.Parallel()

	t.Run("returns instance info with version", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a backend server with version 1.2.3
		server := backendserver.NewTestServer("1.2.3", backendserver.NewTestLogger(t))
		defer server.Close()

		// WHEN: requesting instance info
		resp := httpGet(t, server.URL+"/instance/info")
		defer resp.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: response matches expected JSON structure
		testastic.Equal(t, http.StatusOK, resp.StatusCode)
		testastic.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		testastic.AssertJSON(t, testdataPath("backend_instance_info", "expected_response.json"), resp.Body)
	})

	t.Run("returns consistent hostname across requests", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a backend server
		server := backendserver.NewTestServer("1.0.0", backendserver.NewTestLogger(t))
		defer server.Close()

		// WHEN: requesting instance info twice
		resp1 := httpGet(t, server.URL+"/instance/info")
		defer resp1.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		resp2 := httpGet(t, server.URL+"/instance/info")
		defer resp2.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: both responses match expected JSON structure
		testastic.AssertJSON(t, testdataPath("backend_consistent_hostname", "expected_response.json"), resp1.Body)
		testastic.AssertJSON(t, testdataPath("backend_consistent_hostname", "expected_response.json"), resp2.Body)
	})

	t.Run("health live endpoint responds OK", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a backend server
		server := backendserver.NewTestServer("test-version", backendserver.NewTestLogger(t))
		defer server.Close()

		// WHEN: requesting the live health endpoint
		resp := httpGet(t, server.URL+"/health/live")
		defer resp.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: response matches expected JSON structure
		testastic.Equal(t, http.StatusOK, resp.StatusCode)
		testastic.AssertJSON(t, testdataPath("backend_health_live", "expected_response.json"), resp.Body)
	})

	t.Run("health ready endpoint responds OK", func(t *testing.T) {
		t.Parallel()

		// GIVEN: a backend server
		server := backendserver.NewTestServer("test-version", backendserver.NewTestLogger(t))
		defer server.Close()

		// WHEN: requesting the ready health endpoint
		resp := httpGet(t, server.URL+"/health/ready")
		defer resp.Body.Close() //nolint:errcheck // Ignoring close error in test cleanup.

		// THEN: response matches expected JSON structure
		testastic.Equal(t, http.StatusOK, resp.StatusCode)
		testastic.AssertJSON(t, testdataPath("backend_health_ready", "expected_response.json"), resp.Body)
	})
}
