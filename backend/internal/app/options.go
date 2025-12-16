package app

import (
	"log/slog"
	"os"

	instanceapi "phasor/backend/internal/instance"
)

// Options configures the application setup.
type Options struct {
	Version     string
	GetHostname instanceapi.HostnameFunc
	Environment string
	Logger      *slog.Logger
}

// Option is a functional option for configuring Options.
type Option func(*Options)

// SystemHostname returns the system hostname or "unknown" if it cannot be determined.
func SystemHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}

	return hostname
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Version:     "unknown",
		GetHostname: SystemHostname,
		Environment: "production",
		Logger:      slog.New(slog.DiscardHandler),
	}
}

// WithVersion sets the application version.
func WithVersion(v string) Option {
	return func(o *Options) {
		o.Version = v
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

// WithHostnameFunc sets the hostname function.
func WithHostnameFunc(fn instanceapi.HostnameFunc) Option {
	return func(o *Options) {
		o.GetHostname = fn
	}
}
