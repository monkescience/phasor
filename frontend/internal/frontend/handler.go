package frontend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

const (
	defaultTileCount         = 3
	maxTileCount             = 20
	httpClientTimeout        = 5 * time.Second
	defaultFallbackColor     = "#667eea"
	httpRequestTimeout       = 3 * time.Second
	maxTileColorIndex        = int64(1<<31 - 1) // Max safe int32 value
	transportMaxIdleConns    = 10
	transportIdleConnTimeout = 30 * time.Second
	transportMaxIdlePerHost  = 2
)

// ErrUnexpectedStatusCode is returned when the instance API returns a non-200 status code.
var ErrUnexpectedStatusCode = errors.New("unexpected status code from instance API")

// InstanceInfoResponse represents the response from the backend instance API.
type InstanceInfoResponse struct {
	// GoVersion Go runtime version
	GoVersion string `json:"go_version"`

	// Hostname Instance hostname
	Hostname string `json:"hostname"`

	// Timestamp Current server timestamp
	Timestamp time.Time `json:"timestamp"`

	// Uptime Human-readable process uptime
	Uptime string `json:"uptime"`

	// Version Application version
	Version string `json:"version"`
}

// FrontendHandler handles frontend HTTP requests for the web UI.
type FrontendHandler struct {
	templates      *template.Template
	instanceClient *http.Client
	instanceURL    string
	tileColors     []string
}

// InstanceTileData represents data for a single instance tile in the UI.
type InstanceTileData struct {
	Index         int
	Info          InstanceInfoResponse
	Color         string
	HostnameColor string
}

// TilesData holds the collection of instance tiles to render.
type TilesData struct {
	Instances []InstanceTileData
}

// IndexData contains data for rendering the index page.
type IndexData struct {
	Count int
}

// NewFrontendHandler creates a new frontend handler with the specified templates path,
// instance API URL, and tile colors.
func NewFrontendHandler(
	templatesPath, instanceURL string,
	tileColors []string,
) (*FrontendHandler, error) {
	tmpl, err := template.ParseGlob(filepath.Join(templatesPath, "*.gohtml"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &FrontendHandler{
		templates: tmpl,
		//nolint:exhaustruct // Other http.Client fields use sensible defaults
		instanceClient: &http.Client{
			Timeout: httpClientTimeout,
			//nolint:exhaustruct // Other http.Transport fields use sensible defaults
			Transport: &http.Transport{
				MaxIdleConns:        transportMaxIdleConns,
				IdleConnTimeout:     transportIdleConnTimeout,
				DisableCompression:  false,
				DisableKeepAlives:   false,
				MaxIdleConnsPerHost: transportMaxIdlePerHost,
			},
		},
		instanceURL: instanceURL,
		tileColors:  tileColors,
	}, nil
}

// IndexHandler serves the main index page with the default tile count.
func (h *FrontendHandler) IndexHandler(writer http.ResponseWriter, _ *http.Request) {
	data := IndexData{
		Count: defaultTileCount,
	}

	err := h.templates.ExecuteTemplate(writer, "index.gohtml", data)
	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("failed to render template: %v", err),
			http.StatusInternalServerError,
		)

		return
	}
}

// TilesHandler renders instance tiles based on the count query parameter.
func (h *FrontendHandler) TilesHandler(writer http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	count := defaultTileCount

	if countStr != "" {
		//nolint:noinlineerr // Inline error check is clearer for optional parameter parsing
		if parsedCount, err := strconv.Atoi(countStr); err == nil && parsedCount > 0 &&
			parsedCount <= maxTileCount {
			count = parsedCount
		}
	}

	instances := make([]InstanceTileData, count)
	//nolint:varnamelen // 'i' is idiomatic for loop index
	for i := range count {
		info, err := h.fetchInstanceInfo(req.Context())
		if err != nil {
			info = InstanceInfoResponse{
				Version:   "error",
				Hostname:  "failed to fetch",
				Uptime:    "N/A",
				GoVersion: "N/A",
				Timestamp: time.Now(),
			}
		}

		instances[i] = InstanceTileData{
			Index:         i + 1,
			Info:          info,
			Color:         h.getColorForVersion(info.Version),
			HostnameColor: h.getColorForHostname(info.Hostname),
		}
	}

	data := TilesData{
		Instances: instances,
	}

	err := h.templates.ExecuteTemplate(writer, "tiles.gohtml", data)
	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("failed to render tiles: %v", err),
			http.StatusInternalServerError,
		)

		return
	}
}

func (h *FrontendHandler) fetchInstanceInfo(
	ctx context.Context,
) (InstanceInfoResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, httpRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.instanceURL, nil)
	if err != nil {
		return InstanceInfoResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := h.instanceClient.Do(req)
	if err != nil {
		return InstanceInfoResponse{}, fmt.Errorf(
			"failed to fetch instance info: %w",
			err,
		)
	}

	//nolint:noinlineerr,wsl // Defer close pattern is idiomatic for resource cleanup
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to close response body: %w", closeErr))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return InstanceInfoResponse{}, fmt.Errorf(
			"%w: %d",
			ErrUnexpectedStatusCode,
			resp.StatusCode,
		)
	}

	var info InstanceInfoResponse

	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return InstanceInfoResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return info, nil
}

// getColorForVersion returns a color from the configured tile colors based on the version string.
// Uses a hash function to deterministically select a color.
func (h *FrontendHandler) getColorForVersion(version string) string {
	if len(h.tileColors) == 0 {
		return defaultFallbackColor
	}

	hasher := fnv.New32a()
	//nolint:noinlineerr // Hash.Write error is extremely unlikely and non-critical
	if _, err := hasher.Write([]byte(version)); err != nil {
		return defaultFallbackColor
	}

	hashValue := min(int64(hasher.Sum32()), maxTileColorIndex)
	index := hashValue % int64(len(h.tileColors))

	return h.tileColors[index]
}

// getColorForHostname returns a color from the configured tile colors based on the hostname string.
// Uses a hash function to deterministically select a color.
func (h *FrontendHandler) getColorForHostname(hostname string) string {
	if len(h.tileColors) == 0 {
		return defaultFallbackColor
	}

	hasher := fnv.New32a()
	//nolint:noinlineerr // Hash.Write error is extremely unlikely and non-critical
	if _, err := hasher.Write([]byte(hostname)); err != nil {
		return defaultFallbackColor
	}

	hashValue := min(int64(hasher.Sum32()), maxTileColorIndex)
	index := hashValue % int64(len(h.tileColors))

	return h.tileColors[index]
}
