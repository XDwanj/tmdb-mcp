package tmdb

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestClient_GetMovieDetails_Success tests getting movie details successfully
func TestClient_GetMovieDetails_Success(t *testing.T) {
	// Mock TMDB API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/movie/27205", r.URL.Path)
		assert.Equal(t, "credits,videos", r.URL.Query().Get("append_to_response"))
		assert.Equal(t, "test-api-key", r.URL.Query().Get("api_key"))
		assert.Equal(t, "en-US", r.URL.Query().Get("language"))

		// Return mock response
		response := MovieDetails{
			ID:          27205,
			Title:       "Inception",
			ReleaseDate: "2010-07-16",
			Runtime:     148,
			VoteAverage: 8.4,
			Overview:    "Cobb, a skilled thief...",
			Genres: []Genre{
				{ID: 28, Name: "Action"},
				{ID: 878, Name: "Science Fiction"},
			},
			Credits: Credits{
				Cast: []CastMember{
					{ID: 6193, Name: "Leonardo DiCaprio", Character: "Cobb"},
				},
				Crew: []CrewMember{
					{ID: 525, Name: "Christopher Nolan", Job: "Director", Department: "Directing"},
				},
			},
			Videos: Videos{
				Results: []Video{
					{ID: "1", Key: "YoHD9XEInc0", Name: "Official Trailer", Site: "YouTube", Type: "Trailer"},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with mock server
	client := createTestClient(t, server.URL, "test-api-key")

	// Call GetMovieDetails
	ctx := context.Background()
	result, err := client.GetMovieDetails(ctx, 27205, nil)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 27205, result.ID)
	assert.Equal(t, "Inception", result.Title)
	assert.Equal(t, 148, result.Runtime)
	assert.Equal(t, 8.4, result.VoteAverage)
	assert.Equal(t, 2, len(result.Genres))
	assert.Equal(t, 1, len(result.Credits.Cast))
	assert.Equal(t, "Leonardo DiCaprio", result.Credits.Cast[0].Name)
	assert.Equal(t, 1, len(result.Credits.Crew))
	assert.Equal(t, "Christopher Nolan", result.Credits.Crew[0].Name)
	assert.Equal(t, 1, len(result.Videos.Results))
}

// TestClient_GetTVDetails_Success tests getting TV details successfully
func TestClient_GetTVDetails_Success(t *testing.T) {
	// Mock TMDB API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/tv/1399", r.URL.Path)
		assert.Equal(t, "credits,videos", r.URL.Query().Get("append_to_response"))
		assert.Equal(t, "test-api-key", r.URL.Query().Get("api_key"))
		assert.Equal(t, "en-US", r.URL.Query().Get("language"))

		// Return mock response
		response := TVDetails{
			ID:               1399,
			Name:             "Game of Thrones",
			FirstAirDate:     "2011-04-17",
			LastAirDate:      "2019-05-19",
			NumberOfSeasons:  8,
			NumberOfEpisodes: 73,
			VoteAverage:      8.3,
			Overview:         "Seven noble families fight for control...",
			Genres: []Genre{
				{ID: 10765, Name: "Sci-Fi & Fantasy"},
			},
			Credits: Credits{
				Cast: []CastMember{
					{ID: 239019, Name: "Peter Dinklage", Character: "Tyrion Lannister"},
				},
				Crew: []CrewMember{
					{ID: 1318704, Name: "David Benioff", Job: "Executive Producer", Department: "Production"},
				},
			},
			Videos: Videos{
				Results: []Video{
					{ID: "2", Key: "rlR4PJn8b8I", Name: "Official Trailer", Site: "YouTube", Type: "Trailer"},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with mock server
	client := createTestClient(t, server.URL, "test-api-key")

	// Call GetTVDetails
	ctx := context.Background()
	result, err := client.GetTVDetails(ctx, 1399, nil)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1399, result.ID)
	assert.Equal(t, "Game of Thrones", result.Name)
	assert.Equal(t, 8, result.NumberOfSeasons)
	assert.Equal(t, 73, result.NumberOfEpisodes)
	assert.Equal(t, 8.3, result.VoteAverage)
	assert.Equal(t, 1, len(result.Genres))
	assert.Equal(t, 1, len(result.Credits.Cast))
	assert.Equal(t, "Peter Dinklage", result.Credits.Cast[0].Name)
	assert.Equal(t, 1, len(result.Credits.Crew))
	assert.Equal(t, 1, len(result.Videos.Results))
}

// TestClient_GetPersonDetails_Success tests getting person details successfully
func TestClient_GetPersonDetails_Success(t *testing.T) {
	// Mock TMDB API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/person/525", r.URL.Path)
		assert.Equal(t, "combined_credits", r.URL.Query().Get("append_to_response"))
		assert.Equal(t, "test-api-key", r.URL.Query().Get("api_key"))
		assert.Equal(t, "en-US", r.URL.Query().Get("language"))

		// Return mock response
		response := PersonDetails{
			ID:                 525,
			Name:               "Christopher Nolan",
			Birthday:           "1970-07-30",
			PlaceOfBirth:       "London, England, UK",
			Biography:          "Best known for his cerebral...",
			KnownForDepartment: "Directing",
			CombinedCredits: CombinedCredits{
				Cast: []CombinedCastCredit{},
				Crew: []CombinedCrewCredit{
					{ID: 27205, MediaType: "movie", Title: "Inception", Job: "Director", Department: "Directing", ReleaseDate: "2010-07-16"},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with mock server
	client := createTestClient(t, server.URL, "test-api-key")

	// Call GetPersonDetails
	ctx := context.Background()
	result, err := client.GetPersonDetails(ctx, 525, nil)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 525, result.ID)
	assert.Equal(t, "Christopher Nolan", result.Name)
	assert.Equal(t, "1970-07-30", result.Birthday)
	assert.Equal(t, "Directing", result.KnownForDepartment)
	assert.Equal(t, 0, len(result.CombinedCredits.Cast))
	assert.Equal(t, 1, len(result.CombinedCredits.Crew))
	assert.Equal(t, "Inception", result.CombinedCredits.Crew[0].Title)
}

// TestClient_GetMovieDetails_404NotFound tests getting movie details for non-existent ID
func TestClient_GetMovieDetails_404NotFound(t *testing.T) {
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
	result, err := client.GetMovieDetails(ctx, 9999999, nil)

	// 404 should return nil, nil (not an error)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

// TestClient_GetTVDetails_404NotFound tests getting TV details for non-existent ID
func TestClient_GetTVDetails_404NotFound(t *testing.T) {
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
	result, err := client.GetTVDetails(ctx, 9999999, nil)

	// 404 should return nil, nil (not an error)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

// TestClient_GetPersonDetails_404NotFound tests getting person details for non-existent ID
func TestClient_GetPersonDetails_404NotFound(t *testing.T) {
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
	result, err := client.GetPersonDetails(ctx, 9999999, nil)

	// 404 should return nil, nil (not an error)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

// TestClient_GetMovieDetails_401Unauthorized tests getting movie details with invalid API key
func TestClient_GetMovieDetails_401Unauthorized(t *testing.T) {
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
	result, err := client.GetMovieDetails(ctx, 27205, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Invalid or missing TMDB API Key")
}

// TestClient_GetMovieDetails_InvalidID tests getting movie details with invalid ID
func TestClient_GetMovieDetails_InvalidID(t *testing.T) {
	client := createTestClient(t, "http://example.com", "test-api-key")

	ctx := context.Background()

	// Test with ID = 0
	result, err := client.GetMovieDetails(ctx, 0, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid movie ID")

	// Test with negative ID
	result, err = client.GetMovieDetails(ctx, -1, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid movie ID")
}

// TestClient_GetTVDetails_InvalidID tests getting TV details with invalid ID
func TestClient_GetTVDetails_InvalidID(t *testing.T) {
	client := createTestClient(t, "http://example.com", "test-api-key")

	ctx := context.Background()

	// Test with ID = 0
	result, err := client.GetTVDetails(ctx, 0, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid TV ID")

	// Test with negative ID
	result, err = client.GetTVDetails(ctx, -1, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid TV ID")
}

// TestClient_GetPersonDetails_InvalidID tests getting person details with invalid ID
func TestClient_GetPersonDetails_InvalidID(t *testing.T) {
	client := createTestClient(t, "http://example.com", "test-api-key")

	ctx := context.Background()

	// Test with ID = 0
	result, err := client.GetPersonDetails(ctx, 0, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid person ID")

	// Test with negative ID
	result, err = client.GetPersonDetails(ctx, -1, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid person ID")
}
