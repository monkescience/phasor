package instanceapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/monkescience/testastic"
)

func TestInstanceHandler(t *testing.T) {
	t.Parallel()

	t.Run("returns instance info with all fields", func(t *testing.T) {
		t.Parallel()

		handler := NewInstanceHandler("1.2.3")
		server := httptest.NewServer(http.HandlerFunc(handler.GetInstanceInfo))
		defer server.Close()

		resp, err := http.Get(server.URL)
		testastic.NoError(t, err)
		defer resp.Body.Close()

		testastic.Equal(t, http.StatusOK, resp.StatusCode)
		testastic.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		var response InstanceInfoResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		testastic.NoError(t, err)

		testastic.Equal(t, "1.2.3", response.Version)
		testastic.NotEmpty(t, response.Hostname)
		testastic.NotEmpty(t, response.GoVersion)
		testastic.NotEmpty(t, response.Uptime)
		testastic.False(t, response.Timestamp.IsZero())
	})

	t.Run("returns consistent hostname across requests", func(t *testing.T) {
		t.Parallel()

		handler := NewInstanceHandler("1.0.0")
		server := httptest.NewServer(http.HandlerFunc(handler.GetInstanceInfo))
		defer server.Close()

		resp1, err := http.Get(server.URL)
		testastic.NoError(t, err)
		var response1 InstanceInfoResponse
		json.NewDecoder(resp1.Body).Decode(&response1)
		resp1.Body.Close()

		resp2, err := http.Get(server.URL)
		testastic.NoError(t, err)
		var response2 InstanceInfoResponse
		json.NewDecoder(resp2.Body).Decode(&response2)
		resp2.Body.Close()

		testastic.Equal(t, response1.Hostname, response2.Hostname)
	})

	t.Run("uptime increases over time", func(t *testing.T) {
		t.Parallel()

		handler := NewInstanceHandler("1.0.0")
		server := httptest.NewServer(http.HandlerFunc(handler.GetInstanceInfo))
		defer server.Close()

		resp1, _ := http.Get(server.URL)
		var response1 InstanceInfoResponse
		json.NewDecoder(resp1.Body).Decode(&response1)
		resp1.Body.Close()

		time.Sleep(100 * time.Millisecond)

		resp2, _ := http.Get(server.URL)
		var response2 InstanceInfoResponse
		json.NewDecoder(resp2.Body).Decode(&response2)
		resp2.Body.Close()

		testastic.NotEqual(t, response1.Uptime, response2.Uptime)
	})

	t.Run("timestamp is within request window", func(t *testing.T) {
		t.Parallel()

		handler := NewInstanceHandler("1.0.0")
		server := httptest.NewServer(http.HandlerFunc(handler.GetInstanceInfo))
		defer server.Close()

		before := time.Now()
		resp, err := http.Get(server.URL)
		after := time.Now()
		testastic.NoError(t, err)
		defer resp.Body.Close()

		var response InstanceInfoResponse
		json.NewDecoder(resp.Body).Decode(&response)

		testastic.True(t, response.Timestamp.After(before) || response.Timestamp.Equal(before))
		testastic.True(t, response.Timestamp.Before(after) || response.Timestamp.Equal(after))
	})

	t.Run("handles empty version", func(t *testing.T) {
		t.Parallel()

		handler := NewInstanceHandler("")
		server := httptest.NewServer(http.HandlerFunc(handler.GetInstanceInfo))
		defer server.Close()

		resp, err := http.Get(server.URL)
		testastic.NoError(t, err)
		defer resp.Body.Close()

		testastic.Equal(t, http.StatusOK, resp.StatusCode)

		var response InstanceInfoResponse
		json.NewDecoder(resp.Body).Decode(&response)
		testastic.Empty(t, response.Version)
	})
}
