package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	instanceapi "reference-service-go/internal/incoming/http/instance"
)

func TestInstanceAPI(t *testing.T) {
	t.Parallel()

	t.Run("get instance info returns version", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		version := "1.2.3"
		handler := instanceapi.NewInstanceHandler(version)
		responseRecorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/instance/info", nil)

		// WHEN
		handler.GetInstanceInfo(responseRecorder, req)

		// THEN
		if responseRecorder.Code != http.StatusOK {
			t.Errorf(
				"handler returned wrong status code: got %v want %v",
				responseRecorder.Code,
				http.StatusOK,
			)
		}

		var response instanceapi.InstanceInfoResponse
		if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response.Version != version {
			t.Errorf("expected version %v, got %v", version, response.Version)
		}

		if response.Hostname == "" {
			t.Error("expected hostname to be set")
		}

		if response.GoVersion == "" {
			t.Error("expected go version to be set")
		}

		if response.Uptime == "" {
			t.Error("expected uptime to be set")
		}

		if response.Timestamp.IsZero() {
			t.Error("expected timestamp to be set")
		}
	})

	t.Run("get instance info with different version", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		version := "2.0.0-beta"
		handler := instanceapi.NewInstanceHandler(version)
		responseRecorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/instance/info", nil)

		// WHEN
		handler.GetInstanceInfo(responseRecorder, req)

		// THEN
		if responseRecorder.Code != http.StatusOK {
			t.Errorf(
				"handler returned wrong status code: got %v want %v",
				responseRecorder.Code,
				http.StatusOK,
			)
		}

		var response instanceapi.InstanceInfoResponse
		if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response.Version != version {
			t.Errorf("expected version %v, got %v", version, response.Version)
		}
	})

	t.Run("get instance info returns consistent hostname", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		handler := instanceapi.NewInstanceHandler("1.0.0")

		// WHEN - make two requests
		req1 := httptest.NewRequest(http.MethodGet, "/instance/info", nil)
		responseRecorder1 := httptest.NewRecorder()
		handler.GetInstanceInfo(responseRecorder1, req1)

		req2 := httptest.NewRequest(http.MethodGet, "/instance/info", nil)
		responseRecorder2 := httptest.NewRecorder()
		handler.GetInstanceInfo(responseRecorder2, req2)

		// THEN - hostname should be the same
		var response1 instanceapi.InstanceInfoResponse
		if err := json.NewDecoder(responseRecorder1.Body).Decode(&response1); err != nil {
			t.Fatalf("failed to decode response1: %v", err)
		}

		var response2 instanceapi.InstanceInfoResponse
		if err := json.NewDecoder(responseRecorder2.Body).Decode(&response2); err != nil {
			t.Fatalf("failed to decode response2: %v", err)
		}

		if response1.Hostname != response2.Hostname {
			t.Errorf(
				"expected consistent hostname, got %v and %v",
				response1.Hostname,
				response2.Hostname,
			)
		}
	})

	t.Run("get instance info uptime increases", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		handler := instanceapi.NewInstanceHandler("1.0.0")

		// WHEN - make first request
		req1 := httptest.NewRequest(http.MethodGet, "/instance/info", nil)
		responseRecorder1 := httptest.NewRecorder()
		handler.GetInstanceInfo(responseRecorder1, req1)

		var response1 instanceapi.InstanceInfoResponse
		if err := json.NewDecoder(responseRecorder1.Body).Decode(&response1); err != nil {
			t.Fatalf("failed to decode response1: %v", err)
		}

		time.Sleep(100 * time.Millisecond)

		// WHEN - make second request after delay
		req2 := httptest.NewRequest(http.MethodGet, "/instance/info", nil)
		responseRecorder2 := httptest.NewRecorder()
		handler.GetInstanceInfo(responseRecorder2, req2)

		var response2 instanceapi.InstanceInfoResponse
		if err := json.NewDecoder(responseRecorder2.Body).Decode(&response2); err != nil {
			t.Fatalf("failed to decode response2: %v", err)
		}

		// THEN - uptime should be different (second should be greater)
		if response1.Uptime == response2.Uptime {
			t.Errorf(
				"expected uptime to increase, got same value: %v",
				response1.Uptime,
			)
		}
	})

	t.Run("get instance info timestamp is recent", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		handler := instanceapi.NewInstanceHandler("1.0.0")
		beforeRequest := time.Now()

		// WHEN
		req := httptest.NewRequest(http.MethodGet, "/instance/info", nil)
		responseRecorder := httptest.NewRecorder()
		handler.GetInstanceInfo(responseRecorder, req)

		afterRequest := time.Now()

		// THEN
		var response instanceapi.InstanceInfoResponse
		if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response.Timestamp.Before(beforeRequest) || response.Timestamp.After(afterRequest) {
			t.Errorf(
				"expected timestamp between %v and %v, got %v",
				beforeRequest,
				afterRequest,
				response.Timestamp,
			)
		}
	})

	t.Run("get instance info with empty version", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		handler := instanceapi.NewInstanceHandler("")
		responseRecorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/instance/info", nil)

		// WHEN
		handler.GetInstanceInfo(responseRecorder, req)

		// THEN
		if responseRecorder.Code != http.StatusOK {
			t.Errorf(
				"handler returned wrong status code: got %v want %v",
				responseRecorder.Code,
				http.StatusOK,
			)
		}

		var response instanceapi.InstanceInfoResponse
		if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response.Version != "" {
			t.Errorf("expected empty version, got %v", response.Version)
		}
	})

	t.Run("get instance info content type is JSON", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		handler := instanceapi.NewInstanceHandler("1.0.0")
		responseRecorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/instance/info", nil)

		// WHEN
		handler.GetInstanceInfo(responseRecorder, req)

		// THEN
		contentType := responseRecorder.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf(
				"expected Content-Type 'application/json', got %v",
				contentType,
			)
		}
	})
}
