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
	assert.Equal(t, 15*time.Second, server.server.ReadTimeout)
	assert.Equal(t, 15*time.Second, server.server.WriteTimeout)
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
