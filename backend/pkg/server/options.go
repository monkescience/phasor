package server

import (
	"log/slog"

	instanceapi "phasor/backend/internal/instance"
)

// TestServerOptions configures the test server.
type TestServerOptions struct {
	Version     string
	GetHostname instanceapi.HostnameFunc
	Logger      *slog.Logger
}

// Option is a functional option for configuring TestServerOptions.
type Option func(*TestServerOptions)

// DefaultOptions returns sensible defaults for testing.
func DefaultOptions() TestServerOptions {
	return TestServerOptions{
		Version:     "test-version",
		GetHostname: func() string { return "test-host" },
		Logger:      slog.New(slog.DiscardHandler),
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

// WithHostname sets a static hostname for the test server.
func WithHostname(hostname string) Option {
	return func(o *TestServerOptions) {
		o.GetHostname = func() string { return hostname }
	}
}

// WithHostnameFunc sets the hostname function for the test server.
func WithHostnameFunc(fn instanceapi.HostnameFunc) Option {
	return func(o *TestServerOptions) {
		o.GetHostname = fn
	}
}
