package server

import (
	"log/slog"
)

// TestServerOptions configures the test server.
type TestServerOptions struct {
	Version string
	Logger  *slog.Logger
}

// Option is a functional option for configuring TestServerOptions.
type Option func(*TestServerOptions)

// DefaultOptions returns sensible defaults for testing.
func DefaultOptions() TestServerOptions {
	return TestServerOptions{
		Version: "test-version",
		Logger:  slog.New(slog.DiscardHandler),
	}
}

// WithVersion sets the application version.
func WithVersion(v string) Option {
	return func(o *TestServerOptions) {
		o.Version = v
	}
}

// WithLogger sets a custom logger.
func WithLogger(l *slog.Logger) Option {
	return func(o *TestServerOptions) {
		o.Logger = l
	}
}
