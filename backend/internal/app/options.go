package app

import (
	"log/slog"
)

// Options configures the application setup.
type Options struct {
	Version     string
	Environment string
	Logger      *slog.Logger
}

// Option is a functional option for configuring Options.
type Option func(*Options)

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Version:     "unknown",
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
