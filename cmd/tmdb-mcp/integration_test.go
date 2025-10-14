package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/XDwanj/tmdb-mcp/internal/config"
	"github.com/XDwanj/tmdb-mcp/internal/mcp"
	"github.com/XDwanj/tmdb-mcp/internal/tmdb"
	"github.com/XDwanj/tmdb-mcp/internal/tools"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// testEnvironment holds all dependencies for integration tests
type testEnvironment struct {
	config     config.Config
	tmdbClient *tmdb.Client
	mcpServer  *mcp.Server
}

// mcpClientServer holds the connected client session
type mcpClientServer struct {
	client        *mcpsdk.Client
	clientSession *mcpsdk.ClientSession
}

// setupTestEnvironment creates all test dependencies
// Returns a configured test environment or skips the test if TMDB_API_KEY is not set
func setupTestEnvironment(t *testing.T) *testEnvironment {
	t.Helper()

	// Load TMDB API Key from environment
	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	// Create test configuration
	cfg := config.Config{
		TMDB: config.TMDBConfig{
			APIKey:    apiKey,
			Language:  "en-US",
			RateLimit: 40,
		},
		Logging: config.LogConfig{
			Level: "info",
		},
	}

	// Create test logger
	logger := zaptest.NewLogger(t)

	// Create TMDB Client
	tmdbClient := tmdb.NewClient(cfg.TMDB, logger)

	// Create MCP Server
	mcpServer := mcp.NewServer(tmdbClient, logger)

	return &testEnvironment{
		config:     cfg,
		tmdbClient: tmdbClient,
		mcpServer:  mcpServer,
	}
}

// setupMCPClientServer creates InMemoryTransports and connects client-server
// This simulates the MCP protocol communication in a single process
func setupMCPClientServer(t *testing.T, env *testEnvironment) *mcpClientServer {
	t.Helper()

	ctx := context.Background()

	// Create InMemoryTransports - this is the core of our integration testing strategy
	clientTransport, serverTransport := mcpsdk.NewInMemoryTransports()

	// Start server in a goroutine (Run is blocking)
	serverErrChan := make(chan error, 1)
	go func() {
		err := env.mcpServer.Run(ctx, serverTransport)
		if err != nil {
			serverErrChan <- err
		}
	}()

	// Create and connect client
	client := mcpsdk.NewClient(&mcpsdk.Implementation{
		Name:    "test-client",
		Version: "1.0.0-test",
	}, nil)

	clientSession, err := client.Connect(ctx, clientTransport, nil)
	require.NoError(t, err, "Failed to connect client session")

	// Check if server had any errors during startup
	select {
	case serverErr := <-serverErrChan:
		t.Fatalf("Server failed to start: %v", serverErr)
	default:
		// Server started successfully
	}

	return &mcpClientServer{
		client:        client,
		clientSession: clientSession,
	}
}

// cleanup closes the client session
func (m *mcpClientServer) cleanup() {
	if m.clientSession != nil {
		m.clientSession.Close()
	}
}

// TestSearchTool_Integration is a basic integration test that verifies
// the entire MCP protocol stack works correctly with InMemoryTransports
func TestSearchTool_Integration(t *testing.T) {
	// Setup test environment
	env := setupTestEnvironment(t)

	// Setup MCP client-server
	mcs := setupMCPClientServer(t, env)
	defer mcs.cleanup()

	// Call search tool
	ctx := context.Background()
	start := time.Now()

	result, err := mcs.clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
		Name: "search",
		Arguments: map[string]any{
			"query": "Inception",
			"page":  1,
		},
	})

	duration := time.Since(start)

	// Verify no error
	require.NoError(t, err, "CallTool should not return error")
	assert.NotNil(t, result, "Result should not be nil")

	// Verify performance: response time should be less than 3 seconds
	assert.Less(t, duration, 3*time.Second, "Search should complete within 3 seconds")
	t.Logf("Search completed in %v", duration)

	// Verify result structure
	assert.Len(t, result.Content, 1, "Result should have exactly 1 content item")

	// Extract text content
	textContent, ok := result.Content[0].(*mcpsdk.TextContent)
	require.True(t, ok, "Content should be TextContent type")

	// Parse JSON response
	var response tools.SearchResponse
	err = json.Unmarshal([]byte(textContent.Text), &response)
	require.NoError(t, err, "Failed to unmarshal response JSON")

	// Verify results exist
	assert.NotEmpty(t, response.Results, "Results should not be empty for 'Inception' query")

	// Verify first result has required fields
	if len(response.Results) > 0 {
		firstResult := response.Results[0]
		assert.NotZero(t, firstResult.ID, "Result ID should not be zero")
		assert.NotEmpty(t, firstResult.MediaType, "Result media_type should not be empty")

		t.Logf("Found %d results for 'Inception'", len(response.Results))
		t.Logf("First result: ID=%d, Type=%s, Title=%s",
			firstResult.ID,
			firstResult.MediaType,
			getTitle(firstResult))
	}
}

// getTitle extracts the title from a SearchResult based on media type
func getTitle(result tmdb.SearchResult) string {
	if result.Title != "" {
		return result.Title
	}
	if result.Name != "" {
		return result.Name
	}
	return "(no title)"
}

// TestSearchTool_SuccessScenarios tests successful search scenarios using table-driven tests
// Covers: movies, TV shows, and people searches (AC: 2)
func TestSearchTool_SuccessScenarios(t *testing.T) {
	// Setup test environment once for all sub-tests
	env := setupTestEnvironment(t)

	tests := []struct {
		name          string
		query         string
		expectedType  string
		minResults    int
		validateTitle func(string) bool
	}{
		{
			name:         "popular movie",
			query:        "Inception",
			expectedType: "movie",
			minResults:   1,
			validateTitle: func(title string) bool {
				return title == "Inception"
			},
		},
		{
			name:         "tv show",
			query:        "Breaking Bad",
			expectedType: "tv",
			minResults:   1,
			validateTitle: func(title string) bool {
				return title == "Breaking Bad"
			},
		},
		{
			name:         "person",
			query:        "Christopher Nolan",
			expectedType: "person",
			minResults:   1,
			validateTitle: func(title string) bool {
				return title == "Christopher Nolan"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup MCP client-server for this sub-test
			mcs := setupMCPClientServer(t, env)
			defer mcs.cleanup()

			ctx := context.Background()

			// Call search tool
			result, err := mcs.clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
				Name: "search",
				Arguments: map[string]any{
					"query": tt.query,
					"page":  1,
				},
			})

			// Verify no error
			require.NoError(t, err, "Search should succeed for %s", tt.query)
			assert.NotNil(t, result, "Result should not be nil")

			// Extract and parse response
			assert.Len(t, result.Content, 1, "Result should have exactly 1 content item")
			textContent, ok := result.Content[0].(*mcpsdk.TextContent)
			require.True(t, ok, "Content should be TextContent type")

			var response tools.SearchResponse
			err = json.Unmarshal([]byte(textContent.Text), &response)
			require.NoError(t, err, "Failed to unmarshal response JSON")

			// Verify minimum result count
			assert.GreaterOrEqual(t, len(response.Results), tt.minResults,
				"Should have at least %d result(s) for '%s'", tt.minResults, tt.query)

			// Verify first result matches expected criteria
			if len(response.Results) > 0 {
				firstResult := response.Results[0]

				// Check required fields
				assert.NotZero(t, firstResult.ID, "Result ID should not be zero")
				assert.NotEmpty(t, firstResult.MediaType, "Result media_type should not be empty")

				// Check expected media type
				assert.Equal(t, tt.expectedType, firstResult.MediaType,
					"First result should be of type '%s'", tt.expectedType)

				// Validate title
				title := getTitle(firstResult)
				assert.True(t, tt.validateTitle(title),
					"Title '%s' should match expected pattern for '%s'", title, tt.query)

				t.Logf("✓ Found %d results for '%s', first result: ID=%d, Type=%s, Title=%s",
					len(response.Results), tt.query, firstResult.ID, firstResult.MediaType, title)
			}
		})
	}
}

// TestSearchTool_BoundaryScenarios tests boundary conditions (AC: 2)
func TestSearchTool_BoundaryScenarios(t *testing.T) {
	env := setupTestEnvironment(t)

	tests := []struct {
		name          string
		query         string
		page          int
		expectError   bool
		errorContains string
		checkResults  func(*testing.T, *tools.SearchResponse)
	}{
		{
			name:          "empty query",
			query:         "",
			page:          1,
			expectError:   true,
			errorContains: "query parameter is required",
		},
		{
			name:        "non-existent content",
			query:       "xyzabc123nonexistent999",
			page:        1,
			expectError: false,
			checkResults: func(t *testing.T, resp *tools.SearchResponse) {
				// TMDB may return empty results for non-existent content
				// This is acceptable and not an error
				t.Logf("Results count: %d (empty results are acceptable)", len(resp.Results))
			},
		},
		{
			name:        "pagination test",
			query:       "Star Wars",
			page:        2,
			expectError: false,
			checkResults: func(t *testing.T, resp *tools.SearchResponse) {
				// Page 2 should return results (Star Wars has many entries)
				assert.NotEmpty(t, resp.Results, "Page 2 of Star Wars should have results")
				t.Logf("Page 2 results: %d items", len(resp.Results))
			},
		},
		{
			name:          "query too long",
			query:         string(make([]byte, 501)), // 501 characters
			page:          1,
			expectError:   true,
			errorContains: "query parameter is too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mcs := setupMCPClientServer(t, env)
			defer mcs.cleanup()

			ctx := context.Background()

			// Call search tool
			result, err := mcs.clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
				Name: "search",
				Arguments: map[string]any{
					"query": tt.query,
					"page":  tt.page,
				},
			})

			if tt.expectError {
				// Should return error via IsError field
				require.NoError(t, err, "MCP protocol should not fail")
				require.NotNil(t, result, "Result should not be nil")
				assert.True(t, result.IsError, "IsError should be true for %s", tt.name)

				// Extract error message from content
				if len(result.Content) > 0 {
					textContent, ok := result.Content[0].(*mcpsdk.TextContent)
					if ok {
						assert.Contains(t, textContent.Text, tt.errorContains,
							"Error message should contain '%s'", tt.errorContains)
						t.Logf("✓ Got expected error: %s", textContent.Text)
					}
				}
			} else {
				// Should succeed
				require.NoError(t, err, "Should not return error for %s", tt.name)
				assert.NotNil(t, result, "Result should not be nil")

				// Parse response
				textContent, ok := result.Content[0].(*mcpsdk.TextContent)
				require.True(t, ok, "Content should be TextContent type")

				var response tools.SearchResponse
				err = json.Unmarshal([]byte(textContent.Text), &response)
				require.NoError(t, err, "Failed to unmarshal response JSON")

				// Run custom checks
				if tt.checkResults != nil {
					tt.checkResults(t, &response)
				}
			}
		})
	}
}

// TestSearchTool_ResponseTime specifically tests the response time requirement (AC: 3)
func TestSearchTool_ResponseTime(t *testing.T) {
	env := setupTestEnvironment(t)
	mcs := setupMCPClientServer(t, env)
	defer mcs.cleanup()

	ctx := context.Background()

	// Measure response time
	start := time.Now()
	result, err := mcs.clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
		Name: "search",
		Arguments: map[string]any{
			"query": "Inception",
			"page":  1,
		},
	})
	duration := time.Since(start)

	// Verify success
	require.NoError(t, err, "Search should succeed")
	assert.NotNil(t, result, "Result should not be nil")
	assert.False(t, result.IsError, "Result should not be an error")

	// Verify performance requirement
	assert.Less(t, duration, 3*time.Second,
		"Search should complete within 3 seconds (AC: 3)")

	t.Logf("✓ Response time: %v (requirement: < 3s)", duration)
}

// BenchmarkSearchTool measures the throughput of the search tool (AC: 3)
func BenchmarkSearchTool(b *testing.B) {
	// Skip if no API key
	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		b.Skip("TMDB_API_KEY environment variable not set, skipping benchmark")
	}

	// Setup (not counted in benchmark)
	cfg := config.Config{
		TMDB: config.TMDBConfig{
			APIKey:    apiKey,
			Language:  "en-US",
			RateLimit: 40,
		},
		Logging: config.LogConfig{
			Level: "error", // Use error level to reduce log noise during benchmark
		},
	}

	logger := zaptest.NewLogger(b)
	tmdbClient := tmdb.NewClient(cfg.TMDB, logger)
	mcpServer := mcp.NewServer(tmdbClient, logger)

	clientTransport, serverTransport := mcpsdk.NewInMemoryTransports()

	go func() {
		_ = mcpServer.Run(context.Background(), serverTransport)
	}()

	client := mcpsdk.NewClient(&mcpsdk.Implementation{
		Name:    "benchmark-client",
		Version: "1.0.0-test",
	}, nil)

	clientSession, err := client.Connect(context.Background(), clientTransport, nil)
	if err != nil {
		b.Fatalf("Failed to connect client: %v", err)
	}
	defer clientSession.Close()

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run benchmark
	for i := 0; i < b.N; i++ {
		_, err := clientSession.CallTool(context.Background(), &mcpsdk.CallToolParams{
			Name: "search",
			Arguments: map[string]any{
				"query": "Inception",
				"page":  1,
			},
		})
		if err != nil {
			b.Fatalf("Search failed: %v", err)
		}
	}
}

// TestSearchTool_RateLimiting verifies rate limiting works correctly (AC: 4)
func TestSearchTool_RateLimiting(t *testing.T) {
	env := setupTestEnvironment(t)
	mcs := setupMCPClientServer(t, env)
	defer mcs.cleanup()

	ctx := context.Background()

	// Execute 10 requests quickly
	const requestCount = 10
	start := time.Now()

	for i := 0; i < requestCount; i++ {
		result, err := mcs.clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
			Name: "search",
			Arguments: map[string]any{
				"query": fmt.Sprintf("test%d", i),
				"page":  1,
			},
		})

		// Verify no 429 errors
		require.NoError(t, err, "Request %d should succeed", i+1)
		require.NotNil(t, result, "Result should not be nil")
		assert.False(t, result.IsError, "Request %d should not return error (no 429)", i+1)
	}

	duration := time.Since(start)

	// Verify rate limiter is working but be more realistic about timing
	// Rate: 40 req/10s = 4 req/s = 250ms per request
	// But due to token bucket algorithm and network latency, actual time may vary
	// We expect at least some delay to be enforced, but allow for variance
	expectedMinDuration := time.Duration(1.8 * float64(time.Second)) // More realistic expectation
	assert.GreaterOrEqual(t, duration, expectedMinDuration,
		"Rate limiter should enforce delays. Expected >= %v, got %v", expectedMinDuration, duration)

	// Also verify it's not too fast (which would indicate rate limiting isn't working)
	expectedMaxDuration := time.Duration(5.0 * float64(time.Second))
	assert.Less(t, duration, expectedMaxDuration,
		"Rate limiting should not cause excessive delays. Expected < %v, got %v", expectedMaxDuration, duration)

	t.Logf("✓ %d requests completed in %v (expected range: %v - %v)", requestCount, duration, expectedMinDuration, expectedMaxDuration)
	t.Logf("✓ No 429 errors occurred (rate limiter working correctly)")
}
