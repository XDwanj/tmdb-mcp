package tmdb

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XDwanj/tmdb-mcp/internal/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestClient_Search_ValidQuery tests search with a valid query
func TestClient_Search_ValidQuery(t *testing.T) {
	// Mock TMDB API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/search/multi", r.URL.Path)
		assert.Equal(t, "Inception", r.URL.Query().Get("query"))
		assert.Equal(t, "1", r.URL.Query().Get("page"))
		assert.Equal(t, "test-api-key", r.URL.Query().Get("api_key"))
		assert.Equal(t, "en-US", r.URL.Query().Get("language"))

		// Return mock response
		response := SearchResponse{
			Page: 1,
			Results: []SearchResult{
				{
					ID:          27205,
					MediaType:   "movie",
					Title:       "Inception",
					ReleaseDate: "2010-07-16",
					VoteAverage: 8.4,
					Overview:    "Cobb, a skilled thief...",
				},
			},
			TotalPages:   1,
			TotalResults: 1,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with mock server
	client := createTestClient(t, server.URL, "test-api-key")

	// Call Search
	ctx := context.Background()
	result, err := client.Search(ctx, "Inception", 1, nil)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 1, len(result.Results))
	assert.Equal(t, 27205, result.Results[0].ID)
	assert.Equal(t, "movie", result.Results[0].MediaType)
	assert.Equal(t, "Inception", result.Results[0].Title)
	assert.Equal(t, 8.4, result.Results[0].VoteAverage)
}

// TestClient_Search_EmptyQuery tests search with empty query
func TestClient_Search_EmptyQuery(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    "test-api-key",
		Language:  "en-US",
		RateLimit: 40,
	}
	client := NewClient(cfg, logger)

	ctx := context.Background()
	result, err := client.Search(ctx, "", 1, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "query parameter is required")
}

// TestClient_Search_QueryTooLong tests search with query exceeding max length
func TestClient_Search_QueryTooLong(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    "test-api-key",
		Language:  "en-US",
		RateLimit: 40,
	}
	client := NewClient(cfg, logger)

	// Create a query longer than 500 characters
	longQuery := string(make([]byte, 501))
	ctx := context.Background()
	result, err := client.Search(ctx, longQuery, 1, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "query parameter is too long")
	assert.Contains(t, err.Error(), "500 characters")
}

// TestClient_Search_404NotFound tests search with no results (404)
func TestClient_Search_404NotFound(t *testing.T) {
	// Mock TMDB API server returning 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"status_code":    34,
			"status_message": "The resource you requested could not be found.",
		})
	}))
	defer server.Close()

	client := createTestClient(t, server.URL, "test-api-key")

	ctx := context.Background()
	result, err := client.Search(ctx, "xyzabc123nonexistent", 1, nil)

	// 404 should return empty results, not error
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Results))
	assert.Equal(t, 0, result.TotalResults)
}

// TestClient_Search_401Unauthorized tests search with invalid API key
func TestClient_Search_401Unauthorized(t *testing.T) {
	// Mock TMDB API server returning 401
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"status_code":    7,
			"status_message": "Invalid API key: You must be granted a valid key.",
		})
	}))
	defer server.Close()

	client := createTestClient(t, server.URL, "invalid-key")

	ctx := context.Background()
	result, err := client.Search(ctx, "Inception", 1, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Invalid or missing TMDB API Key")
}

// TestClient_Search_DefaultPage tests search with default page (0)
func TestClient_Search_DefaultPage(t *testing.T) {
	// Mock TMDB API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that page is set to 1 when 0 is provided
		assert.Equal(t, "1", r.URL.Query().Get("page"))

		response := SearchResponse{
			Page:         1,
			Results:      []SearchResult{},
			TotalPages:   0,
			TotalResults: 0,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL, "test-api-key")

	ctx := context.Background()
	result, err := client.Search(ctx, "test", 0, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Page)
}

// TestClient_Search_MultipleMediaTypes tests search with mixed results
func TestClient_Search_MultipleMediaTypes(t *testing.T) {
	// Mock TMDB API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := SearchResponse{
			Page: 1,
			Results: []SearchResult{
				{
					ID:          27205,
					MediaType:   "movie",
					Title:       "Inception",
					ReleaseDate: "2010-07-16",
					VoteAverage: 8.4,
				},
				{
					ID:        525,
					MediaType: "person",
					Name:      "Christopher Nolan",
					Overview:  "Director and screenwriter",
				},
				{
					ID:           1396,
					MediaType:    "tv",
					Name:         "Breaking Bad",
					FirstAirDate: "2008-01-20",
					VoteAverage:  9.3,
				},
			},
			TotalPages:   1,
			TotalResults: 3,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL, "test-api-key")

	ctx := context.Background()
	result, err := client.Search(ctx, "test", 1, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, len(result.Results))
	assert.Equal(t, "movie", result.Results[0].MediaType)
	assert.Equal(t, "person", result.Results[1].MediaType)
	assert.Equal(t, "tv", result.Results[2].MediaType)
}

// createTestClient creates a test client with custom base URL
func createTestClient(t *testing.T, baseURL, apiKey string) *Client {
	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := NewClient(cfg, logger)

	// Override base URL for testing
	client.httpClient.SetBaseURL(baseURL)

	return client
}
