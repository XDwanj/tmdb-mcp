package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/XDwanj/tmdb-mcp/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// testLogger creates a test logger
func testLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

// testConfig creates a test configuration
func testConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			Mode: "sse",
			SSE: config.SSEConfig{
				Enabled: true,
				Host:    "localhost",
				Port:    8910,
			},
		},
	}
}

// TestNewHTTPServer tests the HTTPServer constructor
func TestNewHTTPServer(t *testing.T) {
	cfg := testConfig()
	logger := testLogger()

	server := NewHTTPServer(cfg, nil, logger)

	require.NotNil(t, server)
	assert.NotNil(t, server.server)
	assert.NotNil(t, server.config)
	assert.NotNil(t, server.logger)
	assert.NotNil(t, server.mux)
	assert.Equal(t, "localhost:8910", server.server.Addr)
	assert.Equal(t, 30*time.Second, server.server.ReadTimeout)    // Story 4.4: Updated for SSE
	assert.Equal(t, time.Duration(0), server.server.WriteTimeout) // Story 4.4: Set to 0 for SSE long-lived connections
	assert.Equal(t, 120*time.Second, server.server.IdleTimeout)
}

// TestHealthHandler tests the health check endpoint
func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name           string
		wantStatus     int
		wantBodyFields map[string]interface{}
	}{
		{
			name:       "returns 200 OK with correct fields",
			wantStatus: http.StatusOK,
			wantBodyFields: map[string]interface{}{
				"status":  "ok",
				"version": "1.0.0",
				"mode":    "sse",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := testConfig()
			logger := testLogger()
			server := NewHTTPServer(cfg, nil, logger)

			req := httptest.NewRequest("GET", "/health", nil)
			rr := httptest.NewRecorder()

			handler := server.healthHandler()
			handler.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, tt.wantStatus, rr.Code)

			// Check Content-Type header
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

			// Check response body
			var response map[string]interface{}
			err := json.NewDecoder(rr.Body).Decode(&response)
			require.NoError(t, err)

			for key, expectedValue := range tt.wantBodyFields {
				assert.Equal(t, expectedValue, response[key], "field %s mismatch", key)
			}

			// Check timestamp field exists and is valid
			assert.Contains(t, response, "timestamp")
			timestamp, ok := response["timestamp"].(string)
			assert.True(t, ok, "timestamp should be a string")
			_, err = time.Parse(time.RFC3339, timestamp)
			assert.NoError(t, err, "timestamp should be in RFC3339 format")
		})
	}
}

// TestHealthHandler_ResponseFormat tests the complete JSON response format
func TestHealthHandler_ResponseFormat(t *testing.T) {
	cfg := testConfig()
	logger := testLogger()
	server := NewHTTPServer(cfg, nil, logger)

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	handler := server.healthHandler()
	handler.ServeHTTP(rr, req)

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	// Verify all required fields are present
	requiredFields := []string{"status", "version", "mode", "timestamp"}
	for _, field := range requiredFields {
		assert.Contains(t, response, field, "response should contain field: %s", field)
	}

	// Verify field types
	assert.IsType(t, "", response["status"], "status should be string")
	assert.IsType(t, "", response["version"], "version should be string")
	assert.IsType(t, "", response["mode"], "mode should be string")
	assert.IsType(t, "", response["timestamp"], "timestamp should be string")
}

// TestLoggingMiddleware tests the logging middleware
func TestLoggingMiddleware(t *testing.T) {
	logger := testLogger()

	// Create a simple handler that returns 200 OK
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Wrap with logging middleware
	wrapped := LoggingMiddleware(logger)(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

// TestLoggingMiddleware_StatusCodeCapture tests that the middleware correctly captures status codes
func TestLoggingMiddleware_StatusCodeCapture(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"200 OK", http.StatusOK},
		{"404 Not Found", http.StatusNotFound},
		{"500 Internal Server Error", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := testLogger()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			})

			wrapped := LoggingMiddleware(logger)(handler)

			req := httptest.NewRequest("GET", "/test", nil)
			rr := httptest.NewRecorder()

			wrapped.ServeHTTP(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code)
		})
	}
}

// TestResponseWriter_DefaultStatusCode tests that responseWriter defaults to 200 OK
func TestResponseWriter_DefaultStatusCode(t *testing.T) {
	logger := testLogger()

	// Handler that doesn't call WriteHeader explicitly
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	wrapped := LoggingMiddleware(logger)(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestRecoveryMiddleware tests that the recovery middleware catches panics
func TestRecoveryMiddleware(t *testing.T) {
	logger := testLogger()

	// Create a handler that panics
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	// Wrap with recovery middleware
	wrapped := RecoveryMiddleware(logger)(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Should not panic
	require.NotPanics(t, func() {
		wrapped.ServeHTTP(rr, req)
	})

	// Should return 500 error
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Internal server error")
}

// TestAuthMiddleware tests the Bearer Token authentication middleware
func TestAuthMiddleware(t *testing.T) {
	expectedToken := "abcd1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	logger := testLogger()

	tests := []struct {
		name           string
		authHeader     string
		wantStatus     int
		wantBody       string
		shouldCallNext bool
	}{
		{
			name:           "valid bearer token",
			authHeader:     "Bearer " + expectedToken,
			wantStatus:     http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			wantStatus:     http.StatusUnauthorized,
			wantBody:       `"error":"unauthorized"`,
			shouldCallNext: false,
		},
		{
			name:           "invalid format - Token instead of Bearer",
			authHeader:     "Token " + expectedToken,
			wantStatus:     http.StatusUnauthorized,
			wantBody:       `"error":"unauthorized"`,
			shouldCallNext: false,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer wrongtoken123456789012345678901234567890123456789012345678",
			wantStatus:     http.StatusUnauthorized,
			wantBody:       `"error":"unauthorized"`,
			shouldCallNext: false,
		},
		{
			name:           "empty token",
			authHeader:     "Bearer ",
			wantStatus:     http.StatusUnauthorized,
			wantBody:       `"error":"unauthorized"`,
			shouldCallNext: false,
		},
		{
			name:           "bearer lowercase",
			authHeader:     "bearer " + expectedToken,
			wantStatus:     http.StatusUnauthorized,
			wantBody:       `"error":"unauthorized"`,
			shouldCallNext: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler that tracks if it was called
			nextCalled := false
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			})

			// Wrap with AuthMiddleware
			handler := AuthMiddleware(expectedToken, logger)(nextHandler)

			// Create test request
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()

			// Execute
			handler.ServeHTTP(rr, req)

			// Verify status code
			assert.Equal(t, tt.wantStatus, rr.Code, "status code mismatch")

			// Verify next handler was called or not
			assert.Equal(t, tt.shouldCallNext, nextCalled, "next handler call mismatch")

			// Verify error response body if expected
			if tt.wantBody != "" {
				assert.Contains(t, rr.Body.String(), tt.wantBody, "response body should contain error")

				// Verify JSON structure
				var response map[string]string
				err := json.NewDecoder(rr.Body).Decode(&response)
				require.NoError(t, err, "response should be valid JSON")
				assert.Equal(t, "unauthorized", response["error"])
				assert.NotEmpty(t, response["message"])
			}

			// Verify Content-Type for error responses
			if tt.wantStatus == http.StatusUnauthorized {
				assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
			}
		})
	}
}

// TestCompareTokens tests the constant-time token comparison function
func TestCompareTokens(t *testing.T) {
	tests := []struct {
		name     string
		provided string
		expected string
		want     bool
	}{
		{
			name:     "identical tokens",
			provided: "token123",
			expected: "token123",
			want:     true,
		},
		{
			name:     "different tokens",
			provided: "token123",
			expected: "token456",
			want:     false,
		},
		{
			name:     "different length tokens",
			provided: "token123",
			expected: "token12345",
			want:     false,
		},
		{
			name:     "empty tokens",
			provided: "",
			expected: "",
			want:     true,
		},
		{
			name:     "one empty token",
			provided: "token123",
			expected: "",
			want:     false,
		},
		{
			name:     "case sensitive",
			provided: "Token123",
			expected: "token123",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareTokens(tt.provided, tt.expected)
			assert.Equal(t, tt.want, result)
		})
	}
}

// TestAuthMiddleware_ContentType tests that error responses have correct Content-Type
func TestAuthMiddleware_ContentType(t *testing.T) {
	logger := testLogger()
	expectedToken := "test-token"

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := AuthMiddleware(expectedToken, logger)(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	// No Authorization header - should return 401
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

// TestConnectionTrackingMiddleware tests the connection tracking middleware (Story 4.4)
func TestConnectionTrackingMiddleware(t *testing.T) {
	cfg := testConfig()
	cfg.Server.SSE.Token = "test-token"
	logger := testLogger()

	// 创建带有 MCP server 的 HTTPServer
	server := NewHTTPServer(cfg, nil, logger)

	// Verify initial counter is 0
	initialCount := server.activeConnections.Load()
	assert.Equal(t, int32(0), initialCount, "Initial active connections should be 0")

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify counter was incremented
		count := server.activeConnections.Load()
		assert.Equal(t, int32(1), count, "Active connections should be 1 during request")
		w.WriteHeader(http.StatusOK)
	})

	// Wrap with ConnectionTrackingMiddleware
	tracked := ConnectionTrackingMiddleware(server)(testHandler)

	// Execute request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	tracked.ServeHTTP(rr, req)

	// Verify response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify counter decremented after request
	finalCount := server.activeConnections.Load()
	assert.Equal(t, int32(0), finalCount, "Active connections should be 0 after request completes")
}

// TestConnectionTrackingMiddleware_Concurrent tests concurrent connection tracking (Story 4.4)
func TestConnectionTrackingMiddleware_Concurrent(t *testing.T) {
	cfg := testConfig()
	logger := testLogger()
	server := NewHTTPServer(cfg, nil, logger)

	// Handler that simulates work
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	tracked := ConnectionTrackingMiddleware(server)(testHandler)

	// Run 5 concurrent requests
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func() {
			req := httptest.NewRequest("GET", "/test", nil)
			rr := httptest.NewRecorder()
			tracked.ServeHTTP(rr, req)
			done <- true
		}()
	}

	// Wait for all to complete
	for i := 0; i < 5; i++ {
		<-done
	}

	// Verify final count is 0
	finalCount := server.activeConnections.Load()
	assert.Equal(t, int32(0), finalCount, "All connections should be closed")
}
