// Package fixtures provides shared test utilities for integration tests.
package fixtures

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/monkescience/testastic"
)

// ReadBody reads and returns the response body as a string.
func ReadBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	buf := new(strings.Builder)
	_, err := io.Copy(buf, resp.Body)
	testastic.NoError(t, err)
	return buf.String()
}
