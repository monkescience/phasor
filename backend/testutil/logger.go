package testutil

import (
	"context"
	"log/slog"
	"testing"
)

type testLogHandler struct {
	t     *testing.T
	level slog.Level
}

func (h *testLogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *testLogHandler) Handle(_ context.Context, r slog.Record) error {
	h.t.Log(r.Message)

	return nil
}

func (h *testLogHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }
func (h *testLogHandler) WithGroup(_ string) slog.Handler      { return h }

// NewTestLogger creates a logger that writes to t.Log().
// Logs only appear when a test fails or when running with -v flag.
func NewTestLogger(t *testing.T) *slog.Logger {
	t.Helper()

	return slog.New(&testLogHandler{t: t, level: slog.LevelDebug})
}
