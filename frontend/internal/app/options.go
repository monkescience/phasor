package app

import (
	"log/slog"
)

// Options configures the application setup.
type Options struct {
	BackendURL    string
	TileColors    []string
	TemplatesPath string
	Environment   string
	Logger        *slog.Logger
}

// Option is a functional option for configuring Options.
type Option func(*Options)

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		BackendURL:  "http://localhost:8080/instance/info",
		TileColors:  []string{"#667eea", "#f093fb", "#4facfe"},
		Environment: "production",
		Logger:      slog.New(slog.DiscardHandler),
	}
}

// WithBackendURL sets the backend instance API URL.
func WithBackendURL(url string) Option {
	return func(o *Options) {
		o.BackendURL = url
	}
}

// WithTileColors sets the tile colors.
func WithTileColors(colors []string) Option {
	return func(o *Options) {
		o.TileColors = colors
	}
}

// WithTemplatesPath sets the path to templates directory.
func WithTemplatesPath(path string) Option {
	return func(o *Options) {
		o.TemplatesPath = path
	}
}

// WithEnvironment sets the environment name.
func WithEnvironment(env string) Option {
	return func(o *Options) {
		o.Environment = env
	}
}

// WithLogger sets the logger.
func WithLogger(l *slog.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}
