package instanceapi

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"time"
)

// InstanceHandler handles instance information requests.
type InstanceHandler struct {
	version   string
	startTime time.Time
}

// NewInstanceHandler creates a new instance handler with the specified version.
func NewInstanceHandler(version string) *InstanceHandler {
	return &InstanceHandler{
		version:   version,
		startTime: time.Now(),
	}
}

// GetInstanceInfo returns information about the running instance including version,
// hostname, uptime, and Go version.
func (h *InstanceHandler) GetInstanceInfo(writer http.ResponseWriter, _ *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	uptime := time.Since(h.startTime)

	response := InstanceInfoResponse{
		Version:   h.version,
		Hostname:  hostname,
		Uptime:    uptime.String(),
		GoVersion: runtime.Version(),
		Timestamp: time.Now(),
	}

	writer.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		http.Error(writer, "failed to encode response", http.StatusInternalServerError)

		return
	}
}
