package server

import (
	"log/slog"
)

// TestServerOptions configures the test server.
type TestServerOptions struct {
	BackendURL    string
	TileColors    []string
	TemplatesPath string
	Logger        *slog.Logger
}

// Option is a functional option for configuring TestServerOptions.
type Option func(*TestServerOptions)

// DefaultOptions returns sensible defaults for testing.
func DefaultOptions() TestServerOptions {
	return TestServerOptions{
		BackendURL: "http://localhost:8080/instance/info",
		TileColors: []string{"#667eea", "#f093fb", "#4facfe"},
		Logger:     slog.New(slog.DiscardHandler),
	}
}

// WithBackendURL sets the backend instance API URL.
func WithBackendURL(url string) Option {
	return func(o *TestServerOptions) {
		o.BackendURL = url
	}
}

// WithTileColors sets the tile colors.
func WithTileColors(colors []string) Option {
	return func(o *TestServerOptions) {
		o.TileColors = colors
	}
}

// WithTemplatesPath sets the path to templates directory.
func WithTemplatesPath(path string) Option {
	return func(o *TestServerOptions) {
		o.TemplatesPath = path
	}
}

// WithLogger sets a custom logger.
func WithLogger(l *slog.Logger) Option {
	return func(o *TestServerOptions) {
		o.Logger = l
	}
}
