package tmdb

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestClient_GetMovieRecommendations_Success tests getting movie recommendations successfully
func TestClient_GetMovieRecommendations_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/movie/27205/recommendations", r.URL.Path)
		assert.Equal(t, "1", r.URL.Query().Get("page"))
		assert.Equal(t, "test-api-key", r.URL.Query().Get("api_key"))

		response := RecommendationsResponse{
			Page: 1,
			Results: []RecommendationResult{
				{
					ID:          155,
					Title:       "The Dark Knight",
					ReleaseDate: "2008-07-16",
					VoteAverage: 8.5,
					Overview:    "When the menace known as the Joker emerges...",
					Popularity:  1234.5,
				},
			},
			TotalPages:   1,
			TotalResults: 1,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL, "test-api-key")

	ctx := context.Background()
	result, err := client.GetMovieRecommendations(ctx, 27205, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Results))
	assert.Equal(t, "The Dark Knight", result.Results[0].Title)
	assert.Equal(t, 8.5, result.Results[0].VoteAverage)
}

// TestClient_GetTVRecommendations_Success tests getting TV recommendations successfully
func TestClient_GetTVRecommendations_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/tv/1396/recommendations", r.URL.Path)
		assert.Equal(t, "1", r.URL.Query().Get("page"))

		response := RecommendationsResponse{
			Page: 1,
			Results: []RecommendationResult{
				{
					ID:           60059,
					Name:         "Better Call Saul",
					FirstAirDate: "2015-02-08",
					VoteAverage:  8.7,
					Overview:     "Six years before Saul Goodman meets Walter White...",
					Popularity:   987.6,
				},
			},
			TotalPages:   1,
			TotalResults: 1,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL, "test-api-key")

	ctx := context.Background()
	result, err := client.GetTVRecommendations(ctx, 1396, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Results))
	assert.Equal(t, "Better Call Saul", result.Results[0].Name)
	assert.Equal(t, 8.7, result.Results[0].VoteAverage)
}

// TestClient_GetMovieRecommendations_InvalidID tests with invalid movie ID
func TestClient_GetMovieRecommendations_InvalidID(t *testing.T) {
	client := createTestClient(t, "http://localhost", "test-api-key")

	ctx := context.Background()
	result, err := client.GetMovieRecommendations(ctx, 0, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid movie ID")
	assert.Contains(t, err.Error(), "must be greater than 0")
}

// TestClient_GetTVRecommendations_InvalidID tests with invalid TV show ID
func TestClient_GetTVRecommendations_InvalidID(t *testing.T) {
	client := createTestClient(t, "http://localhost", "test-api-key")

	ctx := context.Background()
	result, err := client.GetTVRecommendations(ctx, -1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid TV show ID")
	assert.Contains(t, err.Error(), "must be greater than 0")
}

// TestClient_GetMovieRecommendations_DefaultPage tests with page=0 (should default to 1)
func TestClient_GetMovieRecommendations_DefaultPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that page is set to 1 when 0 is provided
		assert.Equal(t, "1", r.URL.Query().Get("page"))

		response := RecommendationsResponse{
			Page:         1,
			Results:      []RecommendationResult{},
			TotalPages:   0,
			TotalResults: 0,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL, "test-api-key")

	ctx := context.Background()
	result, err := client.GetMovieRecommendations(ctx, 27205, 0)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Page)
}

// TestClient_GetMovieRecommendations_Page2 tests with page=2
func TestClient_GetMovieRecommendations_Page2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "2", r.URL.Query().Get("page"))

		response := RecommendationsResponse{
			Page:         2,
			Results:      []RecommendationResult{},
			TotalPages:   5,
			TotalResults: 100,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL, "test-api-key")

	ctx := context.Background()
	result, err := client.GetMovieRecommendations(ctx, 27205, 2)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Page)
}

// TestClient_GetMovieRecommendations_404NotFound tests 404 response
func TestClient_GetMovieRecommendations_404NotFound(t *testing.T) {
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
	result, err := client.GetMovieRecommendations(ctx, 999999, 1)

	// 404 should return empty results, not error
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Results))
	assert.Equal(t, 0, result.TotalResults)
}

// TestClient_GetTVRecommendations_404NotFound tests 404 response for TV
func TestClient_GetTVRecommendations_404NotFound(t *testing.T) {
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
	result, err := client.GetTVRecommendations(ctx, 999999, 1)

	// 404 should return empty results, not error
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Results))
	assert.Equal(t, 0, result.TotalResults)
}

// TestClient_GetMovieRecommendations_401Unauthorized tests 401 response
func TestClient_GetMovieRecommendations_401Unauthorized(t *testing.T) {
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
	result, err := client.GetMovieRecommendations(ctx, 27205, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Invalid or missing TMDB API Key")
}

// TestClient_GetMovieRecommendations_429RateLimit tests 429 response
func TestClient_GetMovieRecommendations_429RateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]any{
			"status_code":    25,
			"status_message": "Your request count (#) is over the allowed limit of (40).",
		})
	}))
	defer server.Close()

	client := createTestClient(t, server.URL, "test-api-key")

	ctx := context.Background()
	result, err := client.GetMovieRecommendations(ctx, 27205, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Rate limit exceeded")
}

// TestClient_GetMovieRecommendations_500ServerError tests 500 response
func TestClient_GetMovieRecommendations_500ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"status_code":    11,
			"status_message": "Internal error: Something went wrong, contact TMDb.",
		})
	}))
	defer server.Close()

	client := createTestClient(t, server.URL, "test-api-key")

	ctx := context.Background()
	result, err := client.GetMovieRecommendations(ctx, 27205, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "TMDB API server error")
}
