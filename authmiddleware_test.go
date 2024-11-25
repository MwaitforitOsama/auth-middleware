package authmiddleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	authmiddleware "github.com/MwaitforitOsama/auth-middleware" // Replace with your plugin import path
)

func TestAuthMiddleware(t *testing.T) {
	cfg := authmiddleware.CreateConfig()
	cfg.Headers["X-Host"] = "[[.Host]]"
	cfg.Headers["X-Method"] = "[[.Method]]"
	cfg.Headers["X-URL"] = "[[.URL]]"
	cfg.Headers["X-Token"] = "[[.Token]]" // Assuming the token is part of the header
	cfg.Headers["X-Demo"] = "test"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// You could simulate your middleware logic here
	})

	// Create the handler with your plugin
	handler, err := authmiddleware.New(ctx, next, cfg, "auth-middleware")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	// Prepare a test request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Call your middleware handler
	handler.ServeHTTP(recorder, req)

	// Assert that headers are set correctly
	assertHeader(t, req, "X-Host", "localhost")
	assertHeader(t, req, "X-Method", "GET")
	assertHeader(t, req, "X-URL", "http://localhost")
	assertHeader(t, req, "X-Token", "[[.Token]]") // Adjust the token value based on your logic
	assertHeader(t, req, "X-Demo", "test")
}

// Helper function to assert header values
func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value for %s: got %s, want %s", key, req.Header.Get(key), expected)
	}
}
