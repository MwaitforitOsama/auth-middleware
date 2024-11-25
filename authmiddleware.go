package authmiddleware

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Config struct includes a map for headers and the AuthAPIURL.
type Config struct {
	AuthAPIURL string            `json:"authApiUrl,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"` // Add this field
}

// CreateConfig initializes the middleware configuration.
func CreateConfig() *Config {
	return &Config{
		Headers: make(map[string]string), // Initialize the Headers map
	}
}

type AuthMiddleware struct {
	next       http.Handler
	authAPIURL string
	headers    map[string]string
}

// New creates a new instance of the middleware.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	fmt.Println("Initializing authMiddleware plugin...")

	if config.AuthAPIURL == "" {
		return nil, fmt.Errorf("authApiUrl is required")
	}

	if config.Headers == nil {
		config.Headers = make(map[string]string) // Ensure Headers map is initialized
	}

	fmt.Println("IF done")

	return &AuthMiddleware{
		next:       next,
		authAPIURL: config.AuthAPIURL,
		headers:    config.Headers, // Set the headers here
	}, nil
}

// ServeHTTP processes the HTTP request.
func (a *AuthMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Skip non-HTTP/HTTPS or WebSocket protocols
	if !isHTTPRequest(req) {
		a.next.ServeHTTP(rw, req)
		return
	}

	// Forward request to /api/auth
	authReq, err := http.NewRequest(http.MethodPost, a.authAPIURL, nil)
	if err != nil {
		// Log error and terminate request
		fmt.Println("Error creating auth request:", err)
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Copy headers from the original request and add the middleware headers
	for key, value := range a.headers {
		req.Header.Add(key, value)
	}
	authReq.Header = req.Header

	// Call the auth API
	resp, err := http.DefaultClient.Do(authReq)
	if err != nil {
		// Log error and terminate request
		fmt.Println("Error calling auth API:", err)
		http.Error(rw, "Forbidden", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	// Check if the auth API returned a 200 response
	if resp.StatusCode != http.StatusOK {
		// Log response and terminate request
		fmt.Printf("Auth API returned status: %d\n", resp.StatusCode)
		http.Error(rw, "Forbidden", http.StatusForbidden)
		return
	}

	// Discard the response body since it's not needed
	io.Copy(io.Discard, resp.Body)

	// Forward the original request to the next handler (Traefik)
	a.next.ServeHTTP(rw, req)
}

// isHTTPRequest checks if the request is HTTP/HTTPS or WebSocket.
func isHTTPRequest(req *http.Request) bool {
	return strings.HasPrefix(req.Proto, "HTTP/") || strings.HasPrefix(req.Header.Get("Upgrade"), "websocket")
}
