package test

import (
	"testing"
)

// Integration tests are now located within their respective modules:
// - Backend handler tests: backend/internal/instance/handler_test.go
// - Backend config tests: backend/internal/config/config_test.go
// - Frontend handler tests: frontend/internal/frontend/handler_test.go
// - Frontend config tests: frontend/internal/config/config_test.go
//
// This test module can be used for E2E tests that start actual services.
// For coverage, run tests from each module:
//   cd backend && go test -cover ./...
//   cd frontend && go test -cover ./...

func TestPlaceholder(t *testing.T) {
	t.Skip("Tests are in backend and frontend modules")
}
