package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/XDwanj/tmdb-mcp/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getFreePort returns an available port on the system
func getFreePort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err, "Failed to find free port")
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

// TestHTTPServer_Integration tests the full HTTP server lifecycle
func TestHTTPServer_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Get a free port for testing
	port := getFreePort(t)

	cfg := &config.Config{
		Server: config.ServerConfig{
			Mode: "sse",
			SSE: config.SSEConfig{
				Enabled: true,
				Host:    "127.0.0.1",
				Port:    port,
			},
		},
	}

	logger := testLogger()
	server := NewHTTPServer(cfg, nil, logger)

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := server.Start(); err != nil {
			errChan <- err
		}
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Ensure server is shut down after test
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	// Test health endpoint
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/health", port))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "1.0.0", response["version"])
	assert.Equal(t, "sse", response["mode"])
	assert.Contains(t, response, "timestamp")

	// Check for any start errors
	select {
	case err := <-errChan:
		t.Fatalf("Server failed to start: %v", err)
	default:
		// No error, continue
	}
}

// TestHTTPServer_GracefulShutdown tests graceful shutdown
func TestHTTPServer_GracefulShutdown(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Get a free port for testing
	port := getFreePort(t)

	cfg := &config.Config{
		Server: config.ServerConfig{
			Mode: "sse",
			SSE: config.SSEConfig{
				Enabled: true,
				Host:    "127.0.0.1",
				Port:    port,
			},
		},
	}

	logger := testLogger()
	server := NewHTTPServer(cfg, nil, logger)

	// Start server
	go func() {
		server.Start()
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Make a request
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/health", port))
	require.NoError(t, err)
	resp.Body.Close()

	// Gracefully shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	assert.NoError(t, err)

	// Verify server is shut down by trying to connect again
	time.Sleep(100 * time.Millisecond)
	_, err = http.Get(fmt.Sprintf("http://127.0.0.1:%d/health", port))
	assert.Error(t, err, "Server should be shut down")
}

// TestHTTPServer_ConcurrentRequests tests concurrent health check requests
func TestHTTPServer_ConcurrentRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Get a free port for testing
	port := getFreePort(t)

	cfg := &config.Config{
		Server: config.ServerConfig{
			Mode: "sse",
			SSE: config.SSEConfig{
				Enabled: true,
				Host:    "127.0.0.1",
				Port:    port,
			},
		},
	}

	logger := testLogger()
	server := NewHTTPServer(cfg, nil, logger)

	// Start server
	go func() {
		server.Start()
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	// Make concurrent requests
	concurrency := 10
	var wg sync.WaitGroup
	wg.Add(concurrency)

	errors := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()

			resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/health", port))
			if err != nil {
				errors <- fmt.Errorf("request %d failed: %w", id, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				errors <- fmt.Errorf("request %d got status %d", id, resp.StatusCode)
				return
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		t.Error(err)
	}
}

// TestHTTPServer_ShutdownTimeout tests shutdown with context timeout
func TestHTTPServer_ShutdownTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Get a free port for testing
	port := getFreePort(t)

	cfg := &config.Config{
		Server: config.ServerConfig{
			Mode: "sse",
			SSE: config.SSEConfig{
				Enabled: true,
				Host:    "127.0.0.1",
				Port:    port,
			},
		},
	}

	logger := testLogger()
	server := NewHTTPServer(cfg, nil, logger)

	// Start server
	go func() {
		server.Start()
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Shutdown with a very short timeout to test timeout handling
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// This should complete quickly (either success or timeout error)
	err := server.Shutdown(ctx)
	// We don't assert the error because it could succeed or timeout depending on timing
	t.Logf("Shutdown result: %v", err)
}
