package instanceapi

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"
)

// HostnameFunc is a function that returns the hostname.
type HostnameFunc func() string

// InstanceHandler handles instance information requests.
type InstanceHandler struct {
	version     string
	getHostname HostnameFunc
	startTime   time.Time
}

// NewInstanceHandler creates a new instance handler with the specified version and hostname function.
func NewInstanceHandler(version string, getHostname HostnameFunc) *InstanceHandler {
	return &InstanceHandler{
		version:     version,
		getHostname: getHostname,
		startTime:   time.Now(),
	}
}

// GetInstanceInfo returns information about the running instance including version,
// hostname, uptime, and Go version.
func (h *InstanceHandler) GetInstanceInfo(writer http.ResponseWriter, _ *http.Request) {
	hostname := h.getHostname()
	uptime := time.Since(h.startTime)

	response := InstanceInfoResponse{
		Version:   h.version,
		Hostname:  hostname,
		Uptime:    uptime.String(),
		GoVersion: runtime.Version(),
		Timestamp: time.Now(),
	}

	writer.Header().Set("Content-Type", "application/json")

	encodeErr := json.NewEncoder(writer).Encode(response)
	if encodeErr != nil {
		http.Error(writer, "failed to encode response", http.StatusInternalServerError)

		return
	}
}
