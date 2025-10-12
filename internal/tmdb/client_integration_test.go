package tmdb

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/XDwanj/tmdb-mcp/internal/config"
)

// TestClient_RateLimiter_Integration tests rate limiter integration with TMDB Client
func TestClient_RateLimiter_Integration(t *testing.T) {
	// Create mock server
	requestCount := 0
	var requestMutex sync.Mutex
	var requestTimes []time.Time

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestMutex.Lock()
		requestCount++
		requestTimes = append(requestTimes, time.Now())
		requestMutex.Unlock()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"images":{"base_url":"http://image.tmdb.org/t/p/"}}`))
	}))
	defer server.Close()

	// Create client with rate limit of 40 req/10s
	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    "test_api_key",
		Language:  "en-US",
		RateLimit: 40,
	}
	client := NewClient(cfg, logger)
	client.httpClient.SetBaseURL(server.URL)

	ctx := context.Background()
	start := time.Now()

	// Send 50 requests sequentially
	for i := 0; i < 50; i++ {
		err := client.Ping(ctx)
		require.NoError(t, err, "Request %d should succeed", i+1)
	}

	totalElapsed := time.Since(start)

	// Verify all 50 requests completed
	assert.Equal(t, 50, requestCount, "All 50 requests should have been sent")

	// Verify timing:
	// First 40 requests use burst capacity (should be fast)
	// Requests 41-50 need to wait for tokens (250ms per request = ~2.5s for 10 requests)
	// Total time should be > 2 seconds (for the delayed requests)
	assert.Greater(t, totalElapsed, 2*time.Second, "Total time should be > 2s due to rate limiting")
	assert.Less(t, totalElapsed, 5*time.Second, "Total time should be reasonable (< 5s)")

	// Verify first 40 requests completed quickly
	requestMutex.Lock()
	if len(requestTimes) >= 40 {
		first40Duration := requestTimes[39].Sub(requestTimes[0])
		assert.Less(t, first40Duration, 1*time.Second, "First 40 requests should complete in < 1s using burst capacity")
	}
	requestMutex.Unlock()

	// Verify no 429 errors were triggered (all requests succeeded)
	assert.Equal(t, 50, requestCount, "No requests should have failed with 429")
}

// TestClient_RateLimiter_Concurrent tests rate limiter thread safety with concurrent requests
func TestClient_RateLimiter_Concurrent(t *testing.T) {
	// Create mock server
	var requestMutex sync.Mutex
	requestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestMutex.Lock()
		requestCount++
		requestMutex.Unlock()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"images":{"base_url":"http://image.tmdb.org/t/p/"}}`))
	}))
	defer server.Close()

	// Create client with rate limit of 10 req/10s (faster test)
	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    "test_api_key",
		Language:  "en-US",
		RateLimit: 10,
	}
	client := NewClient(cfg, logger)
	client.httpClient.SetBaseURL(server.URL)

	ctx := context.Background()
	start := time.Now()

	// Launch 20 concurrent goroutines
	var wg sync.WaitGroup
	errorsChan := make(chan error, 20)

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(reqNum int) {
			defer wg.Done()
			err := client.Ping(ctx)
			if err != nil {
				errorsChan <- err
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errorsChan)

	elapsed := time.Since(start)

	// Verify all requests succeeded
	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}
	assert.Empty(t, errors, "All concurrent requests should succeed")
	assert.Equal(t, 20, requestCount, "All 20 concurrent requests should complete")

	// Verify rate limiting still works with concurrency
	// 10 burst + 10 delayed = at least 9 seconds (10 requests * 1s interval)
	assert.Greater(t, elapsed, 9*time.Second, "Concurrent requests should still be rate limited")
	assert.Less(t, elapsed, 12*time.Second, "Should complete within reasonable time")
}

// TestClient_RateLimiter_ContextCancellation tests that rate limiter respects context cancellation
func TestClient_RateLimiter_ContextCancellation(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	// Create client with low rate limit
	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    "test_api_key",
		Language:  "en-US",
		RateLimit: 5, // 5 req/10s = 1 req/2s
	}
	client := NewClient(cfg, logger)
	client.httpClient.SetBaseURL(server.URL)

	ctx := context.Background()

	// Exhaust burst capacity
	for i := 0; i < 5; i++ {
		err := client.Ping(ctx)
		require.NoError(t, err)
	}

	// Create context with short timeout
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// This request should timeout waiting for rate limiter
	err := client.Ping(ctxTimeout)
	require.Error(t, err, "Request should fail due to context timeout")
	assert.Contains(t, err.Error(), "rate limit wait failed", "Error should mention rate limit wait failure")
}
