package lib

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTestValidator(transport http.RoundTripper, retryProxy bool, t *testing.T) *GithubTokenValidator {
	t.Helper()
	dir := t.TempDir()
	all, err := NewResultWriter(filepath.Join(dir, "tokens.txt"))
	if err != nil {
		t.Fatal(err)
	}
	valid, err := NewResultWriter(filepath.Join(dir, "valid-tokens.txt"))
	if err != nil {
		_ = all.Close()
		t.Fatal(err)
	}
	validator := &GithubTokenValidator{
		all:        all,
		valid:      valid,
		client:     &http.Client{Transport: transport},
		retryProxy: retryProxy,
	}
	t.Cleanup(func() { _ = validator.Close() })
	return validator
}

func TestGithubTokenValidatorProxyRetry(t *testing.T) {
	tests := []struct {
		name       string
		retryProxy bool
		attempts   int
	}{
		{name: "without proxy", attempts: 1},
		{name: "with proxy", retryProxy: true, attempts: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attempts := 0
			transport := roundTripFunc(func(*http.Request) (*http.Response, error) {
				attempts++
				return nil, errors.New("proxy unavailable")
			})
			validator := newTestValidator(transport, tt.retryProxy, t)

			if err := validator.Validate(context.Background(), "http://target", "synthetic-token"); err != nil {
				t.Fatal(err)
			}
			if attempts != tt.attempts {
				t.Fatalf("got %d attempts, want %d", attempts, tt.attempts)
			}
		})
	}
}

func TestGithubTokenValidatorDoesNotRetryHTTPResponses(t *testing.T) {
	attempts := 0
	transport := roundTripFunc(func(*http.Request) (*http.Response, error) {
		attempts++
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Header:     make(http.Header),
			Body:       http.NoBody,
		}, nil
	})
	validator := newTestValidator(transport, true, t)

	if err := validator.Validate(context.Background(), "http://target", "synthetic-token"); err != nil {
		t.Fatal(err)
	}
	if attempts != 1 {
		t.Fatalf("got %d attempts, want 1", attempts)
	}
}

func TestResultWriterPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix permissions are not supported on Windows")
	}

	dir := filepath.Join(t.TempDir(), "results")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "tokens.txt")
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}

	writer, err := NewResultWriter(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = writer.Close() })

	dirInfo, err := os.Stat(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got := dirInfo.Mode().Perm(); got != 0o700 {
		t.Fatalf("directory mode is %04o, want 0700", got)
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := fileInfo.Mode().Perm(); got != 0o600 {
		t.Fatalf("file mode is %04o, want 0600", got)
	}
}

func TestScannerRunReturnsFlushError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("[core]\n"))
	}))
	defer server.Close()

	writer, err := NewResultWriter(filepath.Join(t.TempDir(), "results", "git-config.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if err := writer.file.Close(); err != nil {
		t.Fatal(err)
	}
	scanner := &Scanner{
		urls:   []string{server.URL},
		client: server.Client(),
		writer: writer,
	}

	if err := scanner.Run(1, nil); err == nil {
		t.Fatal("Run returned nil after the result flush failed")
	}
}
