package tmdb

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestClient_GetTrending_MovieDay tests getting today's trending movies
func TestClient_GetTrending_MovieDay(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/trending/movie/day", r.URL.Path)
		assert.Equal(t, "1", r.URL.Query().Get("page"))
		assert.Equal(t, "test-api-key", r.URL.Query().Get("api_key"))

		response := TrendingResponse{
			Page: 1,
			Results: []TrendingResult{
				{
					ID:          27205,
					MediaType:   "movie",
					Title:       "Inception",
					ReleaseDate: "2010-07-16",
					VoteAverage: 8.4,
					Overview:    "A skilled thief...",
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
	result, err := client.GetTrending(ctx, "movie", "day", 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Results))
	assert.Equal(t, "movie", result.Results[0].MediaType)
	assert.Equal(t, "Inception", result.Results[0].Title)
	assert.Equal(t, 8.4, result.Results[0].VoteAverage)
}

// TestClient_GetTrending_TVWeek tests getting this week's trending TV shows
func TestClient_GetTrending_TVWeek(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/trending/tv/week", r.URL.Path)

		response := TrendingResponse{
			Page: 1,
			Results: []TrendingResult{
				{
					ID:           1396,
					MediaType:    "tv",
					Name:         "Breaking Bad",
					FirstAirDate: "2008-01-20",
					VoteAverage:  9.3,
					Overview:     "A chemistry teacher...",
					Popularity:   5678.9,
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
	result, err := client.GetTrending(ctx, "tv", "week", 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Results))
	assert.Equal(t, "tv", result.Results[0].MediaType)
	assert.Equal(t, "Breaking Bad", result.Results[0].Name)
	assert.Equal(t, 9.3, result.Results[0].VoteAverage)
}

// TestClient_GetTrending_PersonDay tests getting today's trending people
func TestClient_GetTrending_PersonDay(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/trending/person/day", r.URL.Path)

		response := TrendingResponse{
			Page: 1,
			Results: []TrendingResult{
				{
					ID:                 525,
					MediaType:          "person",
					Name:               "Christopher Nolan",
					KnownForDepartment: "Directing",
					Popularity:         9876.5,
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
	result, err := client.GetTrending(ctx, "person", "day", 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Results))
	assert.Equal(t, "person", result.Results[0].MediaType)
	assert.Equal(t, "Christopher Nolan", result.Results[0].Name)
	assert.Equal(t, "Directing", result.Results[0].KnownForDepartment)
}

// TestClient_GetTrending_InvalidMediaType tests with invalid media_type
func TestClient_GetTrending_InvalidMediaType(t *testing.T) {
	client := createTestClient(t, "http://localhost", "test-api-key")

	ctx := context.Background()
	result, err := client.GetTrending(ctx, "invalid", "day", 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid media_type")
	assert.Contains(t, err.Error(), "must be movie, tv, or person")
}

// TestClient_GetTrending_InvalidTimeWindow tests with invalid time_window
func TestClient_GetTrending_InvalidTimeWindow(t *testing.T) {
	client := createTestClient(t, "http://localhost", "test-api-key")

	ctx := context.Background()
	result, err := client.GetTrending(ctx, "movie", "month", 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid time_window")
	assert.Contains(t, err.Error(), "must be day or week")
}

// TestClient_GetTrending_DefaultPage tests with page=0 (should default to 1)
func TestClient_GetTrending_DefaultPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that page is set to 1 when 0 is provided
		assert.Equal(t, "1", r.URL.Query().Get("page"))

		response := TrendingResponse{
			Page:         1,
			Results:      []TrendingResult{},
			TotalPages:   0,
			TotalResults: 0,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL, "test-api-key")

	ctx := context.Background()
	result, err := client.GetTrending(ctx, "movie", "day", 0)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Page)
}

// TestClient_GetTrending_Page2 tests with page=2
func TestClient_GetTrending_Page2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "2", r.URL.Query().Get("page"))

		response := TrendingResponse{
			Page:         2,
			Results:      []TrendingResult{},
			TotalPages:   5,
			TotalResults: 100,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL, "test-api-key")

	ctx := context.Background()
	result, err := client.GetTrending(ctx, "movie", "day", 2)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Page)
}

// TestClient_GetTrending_404NotFound tests 404 response
func TestClient_GetTrending_404NotFound(t *testing.T) {
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
	result, err := client.GetTrending(ctx, "movie", "day", 1)

	// 404 should return empty results, not error
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Results))
	assert.Equal(t, 0, result.TotalResults)
}

// TestClient_GetTrending_401Unauthorized tests 401 response
func TestClient_GetTrending_401Unauthorized(t *testing.T) {
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
	result, err := client.GetTrending(ctx, "movie", "day", 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Invalid or missing TMDB API Key")
}

// TestClient_GetTrending_429RateLimit tests 429 response
func TestClient_GetTrending_429RateLimit(t *testing.T) {
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
	result, err := client.GetTrending(ctx, "movie", "day", 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Rate limit exceeded")
}

// TestClient_GetTrending_500ServerError tests 500 response
func TestClient_GetTrending_500ServerError(t *testing.T) {
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
	result, err := client.GetTrending(ctx, "movie", "day", 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "TMDB API server error")
}
