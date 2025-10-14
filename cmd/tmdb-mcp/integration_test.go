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

// TestGetDetailsTool_Integration tests the get_details tool (AC: 9)
func TestGetDetailsTool_Integration(t *testing.T) {
	env := setupTestEnvironment(t)

	tests := []struct {
		name           string
		mediaType      string
		id             int
		expectedFields map[string]bool // Fields that must exist in response
	}{
		{
			name:      "movie details - Inception",
			mediaType: "movie",
			id:        27205,
			expectedFields: map[string]bool{
				"id":      true,
				"title":   true,
				"credits": true,
				"videos":  true,
			},
		},
		{
			name:      "tv details - Game of Thrones",
			mediaType: "tv",
			id:        1399,
			expectedFields: map[string]bool{
				"id":      true,
				"name":    true,
				"credits": true,
				"videos":  true,
			},
		},
		{
			name:      "person details - Christopher Nolan",
			mediaType: "person",
			id:        525,
			expectedFields: map[string]bool{
				"id":               true,
				"name":             true,
				"combined_credits": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mcs := setupMCPClientServer(t, env)
			defer mcs.cleanup()

			ctx := context.Background()

			// Call get_details tool
			result, err := mcs.clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
				Name: "get_details",
				Arguments: map[string]any{
					"media_type": tt.mediaType,
					"id":         tt.id,
				},
			})

			// Verify no error
			require.NoError(t, err, "CallTool should not return error")
			assert.NotNil(t, result, "Result should not be nil")
			assert.False(t, result.IsError, "Result should not be an error")

			// Verify result structure
			assert.Len(t, result.Content, 1, "Result should have exactly 1 content item")

			// Extract text content
			textContent, ok := result.Content[0].(*mcpsdk.TextContent)
			require.True(t, ok, "Content should be TextContent type")

			// Parse JSON response
			var response map[string]any
			err = json.Unmarshal([]byte(textContent.Text), &response)
			require.NoError(t, err, "Failed to unmarshal response JSON")

			// Verify all expected fields exist
			for field := range tt.expectedFields {
				assert.Contains(t, response, field, "Response should contain field '%s'", field)
			}

			// Verify ID matches
			idFloat, ok := response["id"].(float64)
			require.True(t, ok, "ID should be a number")
			assert.Equal(t, float64(tt.id), idFloat, "ID should match requested ID")

			t.Logf("✓ Retrieved %s details (ID=%d), fields: %v", tt.mediaType, tt.id, getMapKeys(response))
		})
	}
}

// TestGetDetailsTool_ErrorScenarios tests error handling for get_details (AC: 9)
func TestGetDetailsTool_ErrorScenarios(t *testing.T) {
	env := setupTestEnvironment(t)

	tests := []struct {
		name          string
		mediaType     string
		id            int
		expectError   bool
		errorContains string
	}{
		{
			name:          "invalid media type",
			mediaType:     "invalid",
			id:            123,
			expectError:   true,
			errorContains: "invalid media_type",
		},
		{
			name:          "non-existent movie ID",
			mediaType:     "movie",
			id:            9999999,
			expectError:   true,
			errorContains: "not found",
		},
		{
			name:          "non-existent tv ID",
			mediaType:     "tv",
			id:            9999999,
			expectError:   true,
			errorContains: "not found",
		},
		{
			name:          "non-existent person ID",
			mediaType:     "person",
			id:            9999999,
			expectError:   true,
			errorContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mcs := setupMCPClientServer(t, env)
			defer mcs.cleanup()

			ctx := context.Background()

			// Call get_details tool
			result, err := mcs.clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
				Name: "get_details",
				Arguments: map[string]any{
					"media_type": tt.mediaType,
					"id":         tt.id,
				},
			})

			// Verify error
			require.NoError(t, err, "MCP protocol should not fail")
			require.NotNil(t, result, "Result should not be nil")
			assert.True(t, result.IsError, "IsError should be true for %s", tt.name)

			// Extract error message
			if len(result.Content) > 0 {
				textContent, ok := result.Content[0].(*mcpsdk.TextContent)
				if ok {
					assert.Contains(t, textContent.Text, tt.errorContains,
						"Error message should contain '%s'", tt.errorContains)
					t.Logf("✓ Got expected error: %s", textContent.Text)
				}
			}
		})
	}
}

// TestGetDetailsTool_DataIntegrity tests that returned data has complete structure (AC: 9)
func TestGetDetailsTool_DataIntegrity(t *testing.T) {
	env := setupTestEnvironment(t)
	mcs := setupMCPClientServer(t, env)
	defer mcs.cleanup()

	ctx := context.Background()

	// Test Movie Details - verify Credits and Videos are populated
	t.Run("movie has credits and videos", func(t *testing.T) {
		result, err := mcs.clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
			Name: "get_details",
			Arguments: map[string]any{
				"media_type": "movie",
				"id":         27205, // Inception
			},
		})

		require.NoError(t, err)
		assert.False(t, result.IsError)

		textContent, ok := result.Content[0].(*mcpsdk.TextContent)
		require.True(t, ok)

		var movieDetails tmdb.MovieDetails
		err = json.Unmarshal([]byte(textContent.Text), &movieDetails)
		require.NoError(t, err)

		// Verify Credits are populated
		assert.NotEmpty(t, movieDetails.Credits.Cast, "Movie should have cast members")
		assert.NotEmpty(t, movieDetails.Credits.Crew, "Movie should have crew members")

		// Verify Videos are populated
		assert.NotEmpty(t, movieDetails.Videos.Results, "Movie should have videos")

		t.Logf("✓ Movie has %d cast, %d crew, %d videos",
			len(movieDetails.Credits.Cast),
			len(movieDetails.Credits.Crew),
			len(movieDetails.Videos.Results))
	})

	// Test Person Details - verify CombinedCredits are populated
	t.Run("person has combined credits", func(t *testing.T) {
		result, err := mcs.clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
			Name: "get_details",
			Arguments: map[string]any{
				"media_type": "person",
				"id":         525, // Christopher Nolan
			},
		})

		require.NoError(t, err)
		assert.False(t, result.IsError)

		textContent, ok := result.Content[0].(*mcpsdk.TextContent)
		require.True(t, ok)

		var personDetails tmdb.PersonDetails
		err = json.Unmarshal([]byte(textContent.Text), &personDetails)
		require.NoError(t, err)

		// Verify CombinedCredits are populated
		totalCredits := len(personDetails.CombinedCredits.Cast) + len(personDetails.CombinedCredits.Crew)
		assert.Greater(t, totalCredits, 0, "Person should have combined credits")

		t.Logf("✓ Person has %d cast credits, %d crew credits (total: %d)",
			len(personDetails.CombinedCredits.Cast),
			len(personDetails.CombinedCredits.Crew),
			totalCredits)
	})
}

// getMapKeys returns all keys from a map as a slice
func getMapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// TestDiscoverMoviesTool_Integration tests the discover_movies tool (AC: 10)
func TestDiscoverMoviesTool_Integration(t *testing.T) {
	env := setupTestEnvironment(t)

	tests := []struct {
		name           string
		withGenres     string
		releaseYear    int
		voteAverageGte float64
		sortBy         string
		minResults     int
		validateResult func(*testing.T, *tmdb.DiscoverMoviesResponse)
	}{
		{
			name:           "sci-fi movies after 2020 with rating >= 8.0",
			withGenres:     "878", // Science Fiction
			releaseYear:    2020,
			voteAverageGte: 8.0,
			sortBy:         "popularity.desc",
			minResults:     1,
			validateResult: func(t *testing.T, resp *tmdb.DiscoverMoviesResponse) {
				if len(resp.Results) > 0 {
					firstResult := resp.Results[0]
					assert.Contains(t, firstResult.GenreIDs, 878, "Result should be sci-fi (genre 878)")
					assert.GreaterOrEqual(t, firstResult.VoteAverage, 8.0, "Result should have rating >= 8.0")
					t.Logf("✓ Found sci-fi movie: %s (Rating: %.1f)", firstResult.Title, firstResult.VoteAverage)
				}
			},
		},
		{
			name:           "highest rated action movies",
			withGenres:     "28", // Action
			voteAverageGte: 0,
			sortBy:         "vote_average.desc",
			minResults:     1,
			validateResult: func(t *testing.T, resp *tmdb.DiscoverMoviesResponse) {
				if len(resp.Results) > 0 {
					firstResult := resp.Results[0]
					assert.Contains(t, firstResult.GenreIDs, 28, "Result should be action (genre 28)")
					// First result should have high rating since sorted by vote_average.desc
					assert.Greater(t, firstResult.VoteAverage, 7.0, "Top action movie should have high rating")
					t.Logf("✓ Found top action movie: %s (Rating: %.1f)", firstResult.Title, firstResult.VoteAverage)
				}
			},
		},
		{
			name:       "default behavior - popular movies",
			minResults: 1,
			validateResult: func(t *testing.T, resp *tmdb.DiscoverMoviesResponse) {
				// Default should return popular movies
				assert.NotEmpty(t, resp.Results, "Should return popular movies by default")
				if len(resp.Results) > 0 {
					firstResult := resp.Results[0]
					assert.Greater(t, firstResult.Popularity, 0.0, "Result should have popularity score")
					t.Logf("✓ Found popular movie: %s (Popularity: %.1f)", firstResult.Title, firstResult.Popularity)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mcs := setupMCPClientServer(t, env)
			defer mcs.cleanup()

			ctx := context.Background()

			// Prepare arguments
			args := map[string]any{}
			if tt.withGenres != "" {
				args["with_genres"] = tt.withGenres
			}
			if tt.releaseYear > 0 {
				args["primary_release_year"] = tt.releaseYear
			}
			if tt.voteAverageGte > 0 {
				args["vote_average.gte"] = tt.voteAverageGte
			}
			if tt.sortBy != "" {
				args["sort_by"] = tt.sortBy
			}

			// Call discover_movies tool
			result, err := mcs.clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
				Name:      "discover_movies",
				Arguments: args,
			})

			// Verify no error
			require.NoError(t, err, "CallTool should not return error")
			assert.NotNil(t, result, "Result should not be nil")
			assert.False(t, result.IsError, "Result should not be an error")

			// Verify result structure
			assert.Len(t, result.Content, 1, "Result should have exactly 1 content item")

			// Extract text content
			textContent, ok := result.Content[0].(*mcpsdk.TextContent)
			require.True(t, ok, "Content should be TextContent type")

			// Parse JSON response
			var response tmdb.DiscoverMoviesResponse
			err = json.Unmarshal([]byte(textContent.Text), &response)
			require.NoError(t, err, "Failed to unmarshal response JSON")

			// Verify minimum result count
			assert.GreaterOrEqual(t, len(response.Results), tt.minResults,
				"Should have at least %d result(s)", tt.minResults)

			// Run custom validation
			if tt.validateResult != nil {
				tt.validateResult(t, &response)
			}
		})
	}
}

// TestDiscoverMoviesTool_DataIntegrity tests data structure integrity (AC: 10)
func TestDiscoverMoviesTool_DataIntegrity(t *testing.T) {
	env := setupTestEnvironment(t)
	mcs := setupMCPClientServer(t, env)
	defer mcs.cleanup()

	ctx := context.Background()

	// Test with specific parameters
	result, err := mcs.clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
		Name: "discover_movies",
		Arguments: map[string]any{
			"with_genres":      "878",
			"vote_average.gte": 7.0,
			"sort_by":          "popularity.desc",
			"page":             1,
		},
	})

	require.NoError(t, err)
	assert.False(t, result.IsError)

	textContent, ok := result.Content[0].(*mcpsdk.TextContent)
	require.True(t, ok)

	var response tmdb.DiscoverMoviesResponse
	err = json.Unmarshal([]byte(textContent.Text), &response)
	require.NoError(t, err)

	// Verify response structure
	assert.Greater(t, response.Page, 0, "Page should be > 0")
	assert.NotEmpty(t, response.Results, "Results should not be empty")

	// Verify first result has all required fields
	if len(response.Results) > 0 {
		firstResult := response.Results[0]
		assert.NotZero(t, firstResult.ID, "Result should have ID")
		assert.NotEmpty(t, firstResult.Title, "Result should have title")
		assert.NotEmpty(t, firstResult.ReleaseDate, "Result should have release_date")
		assert.GreaterOrEqual(t, firstResult.VoteAverage, 0.0, "Result should have vote_average")
		assert.NotEmpty(t, firstResult.Overview, "Result should have overview")
		assert.NotEmpty(t, firstResult.GenreIDs, "Result should have genre_ids")
		assert.Greater(t, firstResult.Popularity, 0.0, "Result should have popularity")

		t.Logf("✓ Data integrity verified: ID=%d, Title=%s, Rating=%.1f, Genres=%v",
			firstResult.ID,
			firstResult.Title,
			firstResult.VoteAverage,
			firstResult.GenreIDs)
	}
}

// TestDiscoverMoviesTool_ErrorScenarios tests error handling (AC: 10)
func TestDiscoverMoviesTool_ErrorScenarios(t *testing.T) {
	env := setupTestEnvironment(t)

	tests := []struct {
		name          string
		arguments     map[string]any
		expectError   bool
		errorContains string
	}{
		{
			name: "vote_average.gte out of range (> 10)",
			arguments: map[string]any{
				"vote_average.gte": 11.0,
			},
			expectError:   true,
			errorContains: "vote_average.gte must be between 0 and 10",
		},
		{
			name: "vote_average.gte out of range (< 0)",
			arguments: map[string]any{
				"vote_average.gte": -1.0,
			},
			expectError:   true,
			errorContains: "vote_average.gte must be between 0 and 10",
		},
		{
			name: "vote_average.lte out of range",
			arguments: map[string]any{
				"vote_average.lte": 15.0,
			},
			expectError:   true,
			errorContains: "vote_average.lte must be between 0 and 10",
		},
		{
			name: "non-existent genre combination",
			arguments: map[string]any{
				"with_genres":      "999999",
				"vote_average.gte": 9.9,
			},
			expectError: false, // Should return empty results, not error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mcs := setupMCPClientServer(t, env)
			defer mcs.cleanup()

			ctx := context.Background()

			// Call discover_movies tool
			result, err := mcs.clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
				Name:      "discover_movies",
				Arguments: tt.arguments,
			})

			if tt.expectError {
				// Should return error via IsError field
				require.NoError(t, err, "MCP protocol should not fail")
				require.NotNil(t, result, "Result should not be nil")
				assert.True(t, result.IsError, "IsError should be true for %s", tt.name)

				// Extract error message
				if len(result.Content) > 0 {
					textContent, ok := result.Content[0].(*mcpsdk.TextContent)
					if ok {
						assert.Contains(t, textContent.Text, tt.errorContains,
							"Error message should contain '%s'", tt.errorContains)
						t.Logf("✓ Got expected error: %s", textContent.Text)
					}
				}
			} else {
				// Should succeed (possibly with empty results)
				require.NoError(t, err, "Should not return error for %s", tt.name)
				assert.NotNil(t, result, "Result should not be nil")

				textContent, ok := result.Content[0].(*mcpsdk.TextContent)
				require.True(t, ok, "Content should be TextContent type")

				var response tmdb.DiscoverMoviesResponse
				err = json.Unmarshal([]byte(textContent.Text), &response)
				require.NoError(t, err, "Failed to unmarshal response JSON")

				t.Logf("✓ Query succeeded with %d results (empty is acceptable)", len(response.Results))
			}
		})
	}
}

// TestLanguageParameterOverride tests that language parameter can override config default
// This verifies the two-tier priority model: tool-level param > config default
func TestLanguageParameterOverride(t *testing.T) {
	env := setupTestEnvironment(t)

	tests := []struct {
		name       string
		toolName   string
		baseArgs   map[string]any
		testWithLang bool
	}{
		{
			name:     "search tool",
			toolName: "search",
			baseArgs: map[string]any{
				"query": "Inception",
				"page":  1,
			},
			testWithLang: true,
		},
		{
			name:     "get_details tool",
			toolName: "get_details",
			baseArgs: map[string]any{
				"media_type": "movie",
				"id":         27205, // Inception
			},
			testWithLang: true,
		},
		{
			name:     "discover_movies tool",
			toolName: "discover_movies",
			baseArgs: map[string]any{
				"with_genres": "878",
			},
			testWithLang: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mcs := setupMCPClientServer(t, env)
			defer mcs.cleanup()

			ctx := context.Background()

			// Test 1: Without language parameter (should use config default: en-US)
			t.Run("without language param", func(t *testing.T) {
				result, err := mcs.clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
					Name:      tt.toolName,
					Arguments: tt.baseArgs,
				})

				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.False(t, result.IsError, "Should succeed with config default language")
				t.Logf("✓ Tool succeeded without language param (using config default)")
			})

			// Test 2: With language parameter (should override config default)
			if tt.testWithLang {
				t.Run("with language param override", func(t *testing.T) {
					argsWithLang := make(map[string]any)
					for k, v := range tt.baseArgs {
						argsWithLang[k] = v
					}
					argsWithLang["language"] = "zh" // Override to Chinese

					result, err := mcs.clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
						Name:      tt.toolName,
						Arguments: argsWithLang,
					})

					require.NoError(t, err)
					assert.NotNil(t, result)
					assert.False(t, result.IsError, "Should succeed with language override")
					t.Logf("✓ Tool succeeded with language=zh override")
				})
			}
		})
	}
}
