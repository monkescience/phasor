package frontend

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"slices"
	"strconv"
	"time"
)

const (
	defaultTileCount         = 3
	maxTileCount             = 20
	httpClientTimeout        = 5 * time.Second
	defaultFallbackColor     = "#667eea"
	httpRequestTimeout       = 3 * time.Second
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

// colorTracker assigns unique colors to keys within a single request.
type colorTracker struct {
	colors   []string
	assigned map[string]int
	nextIdx  int
}

// newColorTracker creates a new color tracker with the given color palette.
func newColorTracker(colors []string) *colorTracker {
	return &colorTracker{
		colors:   colors,
		assigned: make(map[string]int),
		nextIdx:  0,
	}
}

// getColor returns a unique color for the given key.
// Same key always returns the same color within a request.
// Different keys get different colors until the palette is exhausted.
func (ct *colorTracker) getColor(key string) string {
	if len(ct.colors) == 0 {
		return defaultFallbackColor
	}

	if idx, ok := ct.assigned[key]; ok {
		return ct.colors[idx]
	}

	idx := ct.nextIdx % len(ct.colors)
	ct.assigned[key] = idx
	ct.nextIdx++

	return ct.colors[idx]
}

// IndexData contains data for rendering the index page.
type IndexData struct {
	Count int
}

// errorInstanceInfo returns an InstanceInfoResponse for error cases.
func errorInstanceInfo() InstanceInfoResponse {
	return InstanceInfoResponse{
		Version:   "error",
		Hostname:  "failed to fetch",
		Uptime:    "N/A",
		GoVersion: "N/A",
		Timestamp: time.Now(),
	}
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
		instanceClient: &http.Client{
			Timeout: httpClientTimeout,
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
		parsedCount, err := strconv.Atoi(countStr)
		if err == nil && parsedCount > 0 && parsedCount <= maxTileCount {
			count = parsedCount
		}
	}

	colorTracker := newColorTracker(h.tileColors)

	instances := make([]InstanceTileData, count)
	for i := range count {
		info, err := h.fetchInstanceInfo(req.Context())
		if err != nil {
			info = errorInstanceInfo()
		}

		tileColor := colorTracker.getColor(info.Hostname + "|" + info.Version)
		instances[i] = InstanceTileData{
			Index:         i + 1,
			Info:          info,
			Color:         tileColor,
			HostnameColor: tileColor,
		}
	}

	// Sort by Hostname (descending), then Version (descending)
	slices.SortFunc(instances, func(a, b InstanceTileData) int {
		if result := cmp.Compare(b.Info.Hostname, a.Info.Hostname); result != 0 {
			return result
		}

		return cmp.Compare(b.Info.Version, a.Info.Version)
	})

	for i := range instances {
		instances[i].Index = i + 1
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

	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
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
