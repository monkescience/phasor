package integration_test

import (
	"context"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/monkescience/testastic"
)

// templatesPath returns the path to test templates directory.
func templatesPath() string {
	//nolint:dogsled // runtime.Caller returns 4 values, we only need filename.
	_, filename, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(filename), "..", "testdata", "templates")
}

// testdataPath returns the path to a testdata file for the given test case.
func testdataPath(testcase, filename string) string {
	//nolint:dogsled // runtime.Caller returns 4 values, we only need filename.
	_, callerFile, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(callerFile), "..", "testdata", testcase, filename)
}

// httpGet performs an HTTP GET request with context.
func httpGet(t *testing.T, url string) *http.Response {
	t.Helper()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	testastic.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	testastic.NoError(t, err)

	return resp
}
