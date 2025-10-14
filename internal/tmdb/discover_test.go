package tmdb

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestClient_DiscoverMovies_ValidParams tests discover movies with valid parameters
func TestClient_DiscoverMovies_ValidParams(t *testing.T) {
	// Mock TMDB API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/discover/movie", r.URL.Path)
		assert.Equal(t, "878", r.URL.Query().Get("with_genres"))
		assert.Equal(t, "2020", r.URL.Query().Get("primary_release_year"))
		assert.Equal(t, "8.0", r.URL.Query().Get("vote_average.gte"))
		assert.Equal(t, "popularity.desc", r.URL.Query().Get("sort_by"))
		assert.Equal(t, "1", r.URL.Query().Get("page"))

		// Return mock response
		response := DiscoverMoviesResponse{
			Page: 1,
			Results: []DiscoverMovieResult{
				{
					ID:          550,
					Title:       "Fight Club",
					ReleaseDate: "1999-10-15",
					VoteAverage: 8.4,
					Overview:    "A ticking-time-bomb insomniac...",
					GenreIDs:    []int{18},
					Popularity:  63.869,
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

	// Call DiscoverMovies
	ctx := context.Background()
	params := DiscoverMoviesParams{
		WithGenres:         "878",
		PrimaryReleaseYear: 2020,
		VoteAverageGte:     8.0,
		SortBy:             "popularity.desc",
		Page:               1,
	}
	result, err := client.DiscoverMovies(ctx, params)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 1, len(result.Results))
	assert.Equal(t, 550, result.Results[0].ID)
	assert.Equal(t, "Fight Club", result.Results[0].Title)
	assert.Equal(t, 8.4, result.Results[0].VoteAverage)
}

// TestClient_DiscoverMovies_DefaultValues tests discover movies with default values
func TestClient_DiscoverMovies_DefaultValues(t *testing.T) {
	// Mock TMDB API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify default sort_by is set
		assert.Equal(t, "popularity.desc", r.URL.Query().Get("sort_by"))
		// Verify default page is set
		assert.Equal(t, "1", r.URL.Query().Get("page"))

		response := DiscoverMoviesResponse{
			Page:         1,
			Results:      []DiscoverMovieResult{},
			TotalPages:   0,
			TotalResults: 0,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL, "test-api-key")

	ctx := context.Background()
	// Call with empty parameters
	params := DiscoverMoviesParams{}
	result, err := client.DiscoverMovies(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Page)
}

// TestClient_DiscoverMovies_VoteAverageGteInvalid tests with invalid vote_average.gte
func TestClient_DiscoverMovies_VoteAverageGteInvalid(t *testing.T) {
	client := createTestClient(t, "http://dummy.url", "test-api-key")

	ctx := context.Background()

	// Test vote_average.gte > 10
	params := DiscoverMoviesParams{
		VoteAverageGte: 11.0,
	}
	result, err := client.DiscoverMovies(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "vote_average.gte must be between 0 and 10")

	// Test vote_average.gte < 0
	params = DiscoverMoviesParams{
		VoteAverageGte: -1.0,
	}
	result, err = client.DiscoverMovies(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "vote_average.gte must be between 0 and 10")
}

// TestClient_DiscoverMovies_VoteAverageLteInvalid tests with invalid vote_average.lte
func TestClient_DiscoverMovies_VoteAverageLteInvalid(t *testing.T) {
	client := createTestClient(t, "http://dummy.url", "test-api-key")

	ctx := context.Background()

	// Test vote_average.lte > 10
	params := DiscoverMoviesParams{
		VoteAverageLte: 11.0,
	}
	result, err := client.DiscoverMovies(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "vote_average.lte must be between 0 and 10")

	// Test vote_average.lte < 0
	params = DiscoverMoviesParams{
		VoteAverageLte: -1.0,
	}
	result, err = client.DiscoverMovies(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "vote_average.lte must be between 0 and 10")
}

// TestClient_DiscoverMovies_404NotFound tests discover movies with no results (404)
func TestClient_DiscoverMovies_404NotFound(t *testing.T) {
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
	params := DiscoverMoviesParams{
		WithGenres: "999999", // Non-existent genre
	}
	result, err := client.DiscoverMovies(ctx, params)

	// 404 should return empty results, not error
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Results))
	assert.Equal(t, 0, result.TotalResults)
}

// TestClient_DiscoverMovies_401Unauthorized tests discover movies with invalid API key
func TestClient_DiscoverMovies_401Unauthorized(t *testing.T) {
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
	params := DiscoverMoviesParams{}
	result, err := client.DiscoverMovies(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Invalid or missing TMDB API Key")
}

// TestClient_DiscoverMovies_ParameterMapping tests that parameters are correctly mapped to query params
func TestClient_DiscoverMovies_ParameterMapping(t *testing.T) {
	// Mock TMDB API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify all parameters are correctly mapped
		assert.Equal(t, "/discover/movie", r.URL.Path)
		assert.Equal(t, "28,12", r.URL.Query().Get("with_genres"))
		assert.Equal(t, "2023", r.URL.Query().Get("primary_release_year"))
		assert.Equal(t, "7.5", r.URL.Query().Get("vote_average.gte"))
		assert.Equal(t, "9.5", r.URL.Query().Get("vote_average.lte"))
		assert.Equal(t, "en", r.URL.Query().Get("with_original_language"))
		assert.Equal(t, "vote_average.desc", r.URL.Query().Get("sort_by"))
		assert.Equal(t, "2", r.URL.Query().Get("page"))

		response := DiscoverMoviesResponse{
			Page:         2,
			Results:      []DiscoverMovieResult{},
			TotalPages:   5,
			TotalResults: 100,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL, "test-api-key")

	ctx := context.Background()
	params := DiscoverMoviesParams{
		WithGenres:           "28,12",
		PrimaryReleaseYear:   2023,
		VoteAverageGte:       7.5,
		VoteAverageLte:       9.5,
		WithOriginalLanguage: "en",
		SortBy:               "vote_average.desc",
		Page:                 2,
	}
	result, err := client.DiscoverMovies(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Page)
}

// TestClient_DiscoverMovies_EmptyResults tests discover movies with valid params but no results
func TestClient_DiscoverMovies_EmptyResults(t *testing.T) {
	// Mock TMDB API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := DiscoverMoviesResponse{
			Page:         1,
			Results:      []DiscoverMovieResult{},
			TotalPages:   0,
			TotalResults: 0,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL, "test-api-key")

	ctx := context.Background()
	params := DiscoverMoviesParams{
		WithGenres:     "878",
		VoteAverageGte: 9.9, // Very high rating
	}
	result, err := client.DiscoverMovies(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Results))
}
