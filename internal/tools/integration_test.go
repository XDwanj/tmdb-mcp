//go:build integration
// +build integration

package tools

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/XDwanj/tmdb-mcp/internal/config"
	"github.com/XDwanj/tmdb-mcp/internal/tmdb"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestSearchIntegration_Inception tests searching for "Inception" movie
func TestSearchIntegration_Inception(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Get TMDB API Key from environment
	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	// Create TMDB client
	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	// Search for "Inception"
	ctx := context.Background()
	results, err := client.Search(ctx, "Inception", 1, nil)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Greater(t, len(results.Results), 0, "Should return at least one result")

	// Verify the first result is likely the Inception movie
	found := false
	for _, result := range results.Results {
		if result.MediaType == "movie" && result.Title == "Inception" {
			found = true
			assert.Equal(t, 27205, result.ID, "Inception movie ID should be 27205")
			assert.Contains(t, result.Overview, "Cobb")
			break
		}
	}
	assert.True(t, found, "Should find Inception movie in results")
}

// TestSearchIntegration_ChristopherNolan tests searching for "Christopher Nolan" person
func TestSearchIntegration_ChristopherNolan(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Get TMDB API Key from environment
	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	// Create TMDB client
	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	// Search for "Christopher Nolan"
	ctx := context.Background()
	results, err := client.Search(ctx, "Christopher Nolan", 1, nil)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Greater(t, len(results.Results), 0, "Should return at least one result")

	// Verify we get person results
	found := false
	for _, result := range results.Results {
		if result.MediaType == "person" && result.Name == "Christopher Nolan" {
			found = true
			assert.Equal(t, 525, result.ID, "Christopher Nolan person ID should be 525")
			break
		}
	}
	assert.True(t, found, "Should find Christopher Nolan person in results")
}

// TestSearchIntegration_BreakingBad tests searching for "Breaking Bad" TV show
func TestSearchIntegration_BreakingBad(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Get TMDB API Key from environment
	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	// Create TMDB client
	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	// Search for "Breaking Bad"
	ctx := context.Background()
	results, err := client.Search(ctx, "Breaking Bad", 1, nil)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Greater(t, len(results.Results), 0, "Should return at least one result")

	// Verify the first result is likely the Breaking Bad TV show
	found := false
	for _, result := range results.Results {
		if result.MediaType == "tv" && result.Name == "Breaking Bad" {
			found = true
			assert.Equal(t, 1396, result.ID, "Breaking Bad TV show ID should be 1396")
			break
		}
	}
	assert.True(t, found, "Should find Breaking Bad TV show in results")
}

// TestSearchIntegration_NonExistentContent tests searching for non-existent content
func TestSearchIntegration_NonExistentContent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Get TMDB API Key from environment
	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	// Create TMDB client
	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	// Search for non-existent content
	ctx := context.Background()
	results, err := client.Search(ctx, "xyzabc123nonexistent9999", 1, nil)

	// Assertions
	assert.NoError(t, err, "Should not return error for no results")
	assert.NotNil(t, results)
	assert.Equal(t, 0, len(results.Results), "Should return empty results")
	assert.Equal(t, 0, results.TotalResults)
}

// TestSearchIntegration_Pagination tests search pagination
func TestSearchIntegration_Pagination(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Get TMDB API Key from environment
	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	// Create TMDB client
	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	// Search for "Star Wars" which should have many results
	ctx := context.Background()
	page1, err := client.Search(ctx, "Star Wars", 1, nil)
	assert.NoError(t, err)
	assert.NotNil(t, page1)
	assert.Equal(t, 1, page1.Page)
	assert.Greater(t, len(page1.Results), 0)

	// Get page 2
	page2, err := client.Search(ctx, "Star Wars", 2, nil)
	assert.NoError(t, err)
	assert.NotNil(t, page2)
	assert.Equal(t, 2, page2.Page)

	// Verify pages have different results
	if len(page2.Results) > 0 {
		assert.NotEqual(t, page1.Results[0].ID, page2.Results[0].ID, "Different pages should have different results")
	}
}

// TestGetRecommendationsIntegration_InceptionMovie tests getting movie recommendations based on Inception
func TestGetRecommendationsIntegration_InceptionMovie(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Get TMDB API Key from environment
	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	// Create TMDB client
	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	// Get movie recommendations based on Inception (ID: 27205)
	ctx := context.Background()
	results, err := client.GetMovieRecommendations(ctx, 27205, 1)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Greater(t, len(results.Results), 0, "Inception should have movie recommendations")

	// Verify results have expected fields
	for _, result := range results.Results {
		assert.NotZero(t, result.ID, "Result should have ID")
		assert.NotEmpty(t, result.Title, "Movie result should have title")
		assert.NotZero(t, result.VoteAverage, "Result should have vote average")
	}
}

// TestGetRecommendationsIntegration_BreakingBadTV tests getting TV recommendations based on Breaking Bad
func TestGetRecommendationsIntegration_BreakingBadTV(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Get TMDB API Key from environment
	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	// Create TMDB client
	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	// Get TV recommendations based on Breaking Bad (ID: 1396)
	ctx := context.Background()
	results, err := client.GetTVRecommendations(ctx, 1396, 1)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Greater(t, len(results.Results), 0, "Breaking Bad should have TV recommendations")

	// Verify results have expected fields
	for _, result := range results.Results {
		assert.NotZero(t, result.ID, "Result should have ID")
		assert.NotEmpty(t, result.Name, "TV result should have name")
		assert.NotZero(t, result.VoteAverage, "Result should have vote average")
	}
}

// TestGetRecommendationsIntegration_NoRecommendations tests getting recommendations for content with no results
// Note: TMDB may return recommendations even for very high IDs, so this test validates
// that the API call succeeds without error, rather than asserting empty results
func TestGetRecommendationsIntegration_NoRecommendations(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Get TMDB API Key from environment
	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	// Create TMDB client
	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	// Get recommendations for a very high movie ID (ID: 999999)
	// Note: TMDB may still return recommendations via collaborative filtering
	ctx := context.Background()
	results, err := client.GetMovieRecommendations(ctx, 999999, 1)

	// Assertions - should not return error, results may or may not be empty
	assert.NoError(t, err, "Should not return error even if ID doesn't exist")
	assert.NotNil(t, results, "Should return non-nil response")
	// We don't assert empty results because TMDB may return recommendations anyway
}

// ==================== get_details Tool Tests ====================

// TestGetDetailsIntegration_InceptionMovie tests getting details for Inception movie
func TestGetDetailsIntegration_InceptionMovie(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	ctx := context.Background()
	details, err := client.GetMovieDetails(ctx, 27205, nil) // Inception

	assert.NoError(t, err, "Should successfully get movie details")
	if err != nil {
		return // Skip remaining assertions if API call failed
	}
	assert.NotNil(t, details, "Details should not be nil")
	assert.Equal(t, 27205, details.ID, "Movie ID should match")
	assert.Equal(t, "Inception", details.Title, "Movie title should be Inception")
	assert.NotEmpty(t, details.Overview, "Overview should not be empty")
	assert.NotZero(t, details.Runtime, "Runtime should not be zero")
	assert.NotNil(t, details.Credits, "Credits should be included")
	assert.NotNil(t, details.Videos, "Videos should be included")
}

// TestGetDetailsIntegration_BreakingBadTV tests getting details for Breaking Bad TV show
func TestGetDetailsIntegration_BreakingBadTV(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	ctx := context.Background()
	details, err := client.GetTVDetails(ctx, 1396, nil) // Breaking Bad

	assert.NoError(t, err, "Should successfully get TV show details")
	if err != nil {
		return // Skip remaining assertions if API call failed
	}
	assert.NotNil(t, details, "Details should not be nil")
	assert.Equal(t, 1396, details.ID, "TV show ID should match")
	assert.Equal(t, "Breaking Bad", details.Name, "TV show name should be Breaking Bad")
	assert.NotEmpty(t, details.Overview, "Overview should not be empty")
	assert.NotZero(t, details.NumberOfSeasons, "Number of seasons should not be zero")
	assert.NotNil(t, details.Credits, "Credits should be included")
	assert.NotNil(t, details.Videos, "Videos should be included")
}

// TestGetDetailsIntegration_ChristopherNolanPerson tests getting details for Christopher Nolan person
func TestGetDetailsIntegration_ChristopherNolanPerson(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	ctx := context.Background()
	details, err := client.GetPersonDetails(ctx, 525, nil) // Christopher Nolan

	assert.NoError(t, err, "Should successfully get person details")
	if err != nil {
		return // Skip remaining assertions if API call failed
	}
	assert.NotNil(t, details, "Details should not be nil")
	assert.Equal(t, 525, details.ID, "Person ID should match")
	assert.Equal(t, "Christopher Nolan", details.Name, "Person name should be Christopher Nolan")
	assert.NotEmpty(t, details.Biography, "Biography should not be empty")
	assert.NotEmpty(t, details.KnownForDepartment, "Known for department should not be empty")
	assert.NotNil(t, details.CombinedCredits, "Combined credits should be included")
}

// ==================== discover_movies Tool Tests ====================

// TestDiscoverMoviesIntegration_HighRatedSciFi tests discovering high-rated sci-fi movies
func TestDiscoverMoviesIntegration_HighRatedSciFi(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	ctx := context.Background()
	params := tmdb.DiscoverMoviesParams{
		WithGenres:     "878", // Science Fiction
		VoteAverageGte: 8.0,
		SortBy:         "vote_average.desc",
	}
	results, err := client.DiscoverMovies(ctx, params)

	assert.NoError(t, err, "Should successfully discover movies")
	if err != nil {
		return // Skip remaining assertions if API call failed
	}
	assert.NotNil(t, results, "Results should not be nil")
	assert.Greater(t, len(results.Results), 0, "Should return at least one movie")

	// Verify genre and rating constraints
	for _, movie := range results.Results {
		assert.NotZero(t, movie.ID, "Movie should have ID")
		assert.NotEmpty(t, movie.Title, "Movie should have title")
		assert.GreaterOrEqual(t, movie.VoteAverage, 8.0, "Movie vote average should be >= 8.0")
		// Genre ID 878 (Science Fiction) should be in genre_ids
		hasSciFiGenre := false
		for _, genreID := range movie.GenreIDs {
			if genreID == 878 {
				hasSciFiGenre = true
				break
			}
		}
		assert.True(t, hasSciFiGenre, "Movie should have Science Fiction genre")
	}
}

// TestDiscoverMoviesIntegration_RecentActionMovies tests discovering recent action movies
func TestDiscoverMoviesIntegration_RecentActionMovies(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	ctx := context.Background()
	params := tmdb.DiscoverMoviesParams{
		WithGenres:         "28", // Action
		PrimaryReleaseYear: 2020,
		SortBy:             "popularity.desc",
	}
	results, err := client.DiscoverMovies(ctx, params)

	assert.NoError(t, err, "Should successfully discover movies")
	if err != nil {
		return // Skip remaining assertions if API call failed
	}
	assert.NotNil(t, results, "Results should not be nil")
	assert.Greater(t, len(results.Results), 0, "Should return at least one movie")

	// Verify results
	for _, movie := range results.Results {
		assert.NotZero(t, movie.ID, "Movie should have ID")
		assert.NotEmpty(t, movie.Title, "Movie should have title")
		// Verify release year is 2020
		if movie.ReleaseDate != "" && len(movie.ReleaseDate) >= 4 {
			assert.Equal(t, "2020", movie.ReleaseDate[:4], "Movie should be released in 2020")
		}
	}
}

// TestDiscoverMoviesIntegration_DefaultBehavior tests discover movies with no filters
func TestDiscoverMoviesIntegration_DefaultBehavior(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	ctx := context.Background()
	params := tmdb.DiscoverMoviesParams{} // Empty params
	results, err := client.DiscoverMovies(ctx, params)

	assert.NoError(t, err, "Should successfully discover movies with default params")
	if err != nil {
		return // Skip remaining assertions if API call failed
	}
	assert.NotNil(t, results, "Results should not be nil")
	assert.Greater(t, len(results.Results), 0, "Should return popular movies by default")

	// Verify basic structure
	for _, movie := range results.Results {
		assert.NotZero(t, movie.ID, "Movie should have ID")
		assert.NotEmpty(t, movie.Title, "Movie should have title")
	}
}

// ==================== discover_tv Tool Tests ====================

// TestDiscoverTVIntegration_HighRatedCrimeDrama tests discovering high-rated crime dramas
func TestDiscoverTVIntegration_HighRatedCrimeDrama(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	ctx := context.Background()
	params := tmdb.DiscoverTVParams{
		WithGenres:     "80", // Crime
		VoteAverageGte: 8.0,
		SortBy:         "vote_average.desc",
	}
	results, err := client.DiscoverTV(ctx, params)

	assert.NoError(t, err, "Should successfully discover TV shows")
	if err != nil {
		return // Skip remaining assertions if API call failed
	}
	assert.NotNil(t, results, "Results should not be nil")
	assert.Greater(t, len(results.Results), 0, "Should return at least one TV show")

	// Verify genre and rating constraints
	for _, show := range results.Results {
		assert.NotZero(t, show.ID, "TV show should have ID")
		assert.NotEmpty(t, show.Name, "TV show should have name")
		assert.GreaterOrEqual(t, show.VoteAverage, 8.0, "TV show vote average should be >= 8.0")
	}
}

// TestDiscoverTVIntegration_OngoingSciFiSeries tests discovering ongoing sci-fi series
func TestDiscoverTVIntegration_OngoingSciFiSeries(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	ctx := context.Background()
	params := tmdb.DiscoverTVParams{
		WithGenres: "10765", // Sci-Fi & Fantasy
		SortBy:     "popularity.desc",
	}
	results, err := client.DiscoverTV(ctx, params)

	assert.NoError(t, err, "Should successfully discover TV shows")
	if err != nil {
		return // Skip remaining assertions if API call failed
	}
	assert.NotNil(t, results, "Results should not be nil")
	assert.Greater(t, len(results.Results), 0, "Should return at least one TV show")

	// Verify results
	for _, show := range results.Results {
		assert.NotZero(t, show.ID, "TV show should have ID")
		assert.NotEmpty(t, show.Name, "TV show should have name")
	}
}

// TestDiscoverTVIntegration_DefaultBehavior tests discover TV with no filters
func TestDiscoverTVIntegration_DefaultBehavior(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	ctx := context.Background()
	params := tmdb.DiscoverTVParams{} // Empty params
	results, err := client.DiscoverTV(ctx, params)

	assert.NoError(t, err, "Should successfully discover TV shows with default params")
	if err != nil {
		return // Skip remaining assertions if API call failed
	}
	assert.NotNil(t, results, "Results should not be nil")
	assert.Greater(t, len(results.Results), 0, "Should return popular TV shows by default")

	// Verify basic structure
	for _, show := range results.Results {
		assert.NotZero(t, show.ID, "TV show should have ID")
		assert.NotEmpty(t, show.Name, "TV show should have name")
	}
}

// ==================== get_trending Tool Tests ====================

// TestGetTrendingIntegration_TodayMovies tests getting today's trending movies
func TestGetTrendingIntegration_TodayMovies(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	ctx := context.Background()
	results, err := client.GetTrending(ctx, "movie", "day", 1)

	assert.NoError(t, err, "Should successfully get trending movies")
	if err != nil {
		return // Skip remaining assertions if API call failed
	}
	assert.NotNil(t, results, "Results should not be nil")
	assert.Greater(t, len(results.Results), 0, "Should return at least one trending movie")

	// Verify results structure
	for _, movie := range results.Results {
		assert.NotZero(t, movie.ID, "Movie should have ID")
		assert.NotEmpty(t, movie.Title, "Movie should have title")
		assert.NotZero(t, movie.Popularity, "Movie should have popularity score")
	}
}

// TestGetTrendingIntegration_WeeklyTVShows tests getting this week's trending TV shows
func TestGetTrendingIntegration_WeeklyTVShows(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	ctx := context.Background()
	results, err := client.GetTrending(ctx, "tv", "week", 1)

	assert.NoError(t, err, "Should successfully get trending TV shows")
	if err != nil {
		return // Skip remaining assertions if API call failed
	}
	assert.NotNil(t, results, "Results should not be nil")
	assert.Greater(t, len(results.Results), 0, "Should return at least one trending TV show")

	// Verify results structure
	for _, show := range results.Results {
		assert.NotZero(t, show.ID, "TV show should have ID")
		assert.NotEmpty(t, show.Name, "TV show should have name")
		assert.NotZero(t, show.Popularity, "TV show should have popularity score")
	}
}

// TestGetTrendingIntegration_TodayPeople tests getting today's trending people
func TestGetTrendingIntegration_TodayPeople(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)

	ctx := context.Background()
	results, err := client.GetTrending(ctx, "person", "day", 1)

	assert.NoError(t, err, "Should successfully get trending people")
	if err != nil {
		return // Skip remaining assertions if API call failed
	}
	assert.NotNil(t, results, "Results should not be nil")
	assert.Greater(t, len(results.Results), 0, "Should return at least one trending person")

	// Verify results structure
	for _, person := range results.Results {
		assert.NotZero(t, person.ID, "Person should have ID")
		assert.NotEmpty(t, person.Name, "Person should have name")
		assert.NotZero(t, person.Popularity, "Person should have popularity score")
		// KnownForDepartment may be empty for some people, so we don't assert it
	}
}

// ==================== Multi-Tool Combination Tests ====================

// TestMultiToolCombination_SearchToGetDetails tests search → get_details workflow
func TestMultiToolCombination_SearchToGetDetails(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)
	ctx := context.Background()

	// Step 1: Search for "Inception"
	searchResults, err := client.Search(ctx, "Inception", 1, nil)
	assert.NoError(t, err, "Search should succeed")
	if err != nil {
		return
	}
	assert.Greater(t, len(searchResults.Results), 0, "Search should return results")

	// Step 2: Extract first result's ID and media_type
	firstResult := searchResults.Results[0]
	assert.NotZero(t, firstResult.ID, "First result should have ID")
	assert.NotEmpty(t, firstResult.MediaType, "First result should have media_type")

	// Step 3: Call get_details based on media_type
	var detailsErr error
	if firstResult.MediaType == "movie" {
		details, err := client.GetMovieDetails(ctx, firstResult.ID, nil)
		detailsErr = err
		if err == nil {
			assert.NotNil(t, details, "Movie details should not be nil")
			assert.NotNil(t, details.Credits, "Movie details should include credits")
			assert.NotNil(t, details.Videos, "Movie details should include videos")
		}
	} else if firstResult.MediaType == "tv" {
		details, err := client.GetTVDetails(ctx, firstResult.ID, nil)
		detailsErr = err
		if err == nil {
			assert.NotNil(t, details, "TV details should not be nil")
			assert.NotNil(t, details.Credits, "TV details should include credits")
			assert.NotNil(t, details.Videos, "TV details should include videos")
		}
	}

	assert.NoError(t, detailsErr, "Get details should succeed")
}

// TestMultiToolCombination_DiscoverToRecommendations tests discover_movies → get_recommendations workflow
func TestMultiToolCombination_DiscoverToRecommendations(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)
	ctx := context.Background()

	// Step 1: Discover high-rated sci-fi movies
	params := tmdb.DiscoverMoviesParams{
		WithGenres:     "878", // Science Fiction
		VoteAverageGte: 8.0,
		SortBy:         "vote_average.desc",
	}
	discoverResults, err := client.DiscoverMovies(ctx, params)
	assert.NoError(t, err, "Discover movies should succeed")
	if err != nil {
		return
	}
	assert.Greater(t, len(discoverResults.Results), 0, "Discover should return results")

	// Step 2: Extract first result's ID
	firstMovie := discoverResults.Results[0]
	assert.NotZero(t, firstMovie.ID, "First movie should have ID")

	// Step 3: Get recommendations based on that movie
	recommendations, err := client.GetMovieRecommendations(ctx, firstMovie.ID, 1)
	assert.NoError(t, err, "Get recommendations should succeed")
	if err != nil {
		return
	}

	// Verify recommendations list is not empty and movies have good ratings
	assert.Greater(t, len(recommendations.Results), 0, "Should have recommendations")
	for _, rec := range recommendations.Results {
		assert.NotZero(t, rec.ID, "Recommendation should have ID")
		assert.NotEmpty(t, rec.Title, "Recommendation should have title")
		assert.Greater(t, rec.VoteAverage, 0.0, "Recommendation should have vote average")
	}
}

// TestMultiToolCombination_TrendingToGetDetails tests get_trending → get_details workflow
func TestMultiToolCombination_TrendingToGetDetails(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)
	ctx := context.Background()

	// Step 1: Get today's trending movies
	trendingResults, err := client.GetTrending(ctx, "movie", "day", 1)
	assert.NoError(t, err, "Get trending should succeed")
	if err != nil {
		return
	}
	assert.Greater(t, len(trendingResults.Results), 0, "Trending should return results")

	// Step 2: Extract first result's ID
	firstMovie := trendingResults.Results[0]
	assert.NotZero(t, firstMovie.ID, "First trending movie should have ID")

	// Step 3: Get details
	details, err := client.GetMovieDetails(ctx, firstMovie.ID, nil)
	assert.NoError(t, err, "Get details should succeed")
	if err != nil {
		return
	}

	// Verify details are complete
	assert.NotNil(t, details, "Details should not be nil")
	assert.Equal(t, firstMovie.ID, details.ID, "Details ID should match trending movie ID")
	assert.NotEmpty(t, details.Title, "Details should have title")
	assert.NotNil(t, details.Credits, "Details should include cast information")
	assert.NotNil(t, details.Credits.Cast, "Details should include cast array")
}

// ==================== Performance Tests ====================

// TestAllTools_PerformanceTest tests sequential calls to all 6 tools with performance metrics
func TestAllTools_PerformanceTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)
	ctx := context.Background()

	// Define tool calls with descriptions
	tools := []struct {
		name string
		call func() error
	}{
		{
			name: "search",
			call: func() error {
				_, err := client.Search(ctx, "Inception", 1, nil)
				return err
			},
		},
		{
			name: "get_details",
			call: func() error {
				_, err := client.GetMovieDetails(ctx, 27205, nil)
				return err
			},
		},
		{
			name: "discover_movies",
			call: func() error {
				params := tmdb.DiscoverMoviesParams{
					WithGenres:     "878",
					VoteAverageGte: 8.0,
				}
				_, err := client.DiscoverMovies(ctx, params)
				return err
			},
		},
		{
			name: "discover_tv",
			call: func() error {
				params := tmdb.DiscoverTVParams{
					WithGenres:     "80",
					VoteAverageGte: 8.0,
				}
				_, err := client.DiscoverTV(ctx, params)
				return err
			},
		},
		{
			name: "get_trending",
			call: func() error {
				_, err := client.GetTrending(ctx, "movie", "day", 1)
				return err
			},
		},
		{
			name: "get_recommendations",
			call: func() error {
				_, err := client.GetMovieRecommendations(ctx, 27205, 1)
				return err
			},
		},
	}

	// Record initial call count
	initialCallCount := client.GetCallCount()

	// Execute all tools sequentially and record performance
	totalStart := time.Now()
	for _, tool := range tools {
		start := time.Now()
		err := tool.call()
		duration := time.Since(start)

		// Log performance
		t.Logf("%s completed in %v", tool.name, duration)

		// Verify no error
		assert.NoError(t, err, "%s should not return error", tool.name)

		// Verify response time < 3 seconds
		assert.Less(t, duration, 3*time.Second, "%s response time should be < 3s", tool.name)
	}
	totalDuration := time.Since(totalStart)

	// Verify total time < 10 seconds
	assert.Less(t, totalDuration, 10*time.Second, "Total time should be < 10 seconds")
	t.Logf("Total execution time: %v", totalDuration)

	// Verify API call count incremented correctly (6 calls)
	finalCallCount := client.GetCallCount()
	callsMade := finalCallCount - initialCallCount
	assert.Equal(t, uint64(6), callsMade, "Should have made 6 API calls")
	t.Logf("API calls made: %d (initial: %d, final: %d)", callsMade, initialCallCount, finalCallCount)
}

// ==================== Concurrent Tests ====================

// TestAllTools_ConcurrentTest tests concurrent calls to multiple tools
func TestAllTools_ConcurrentTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)
	ctx := context.Background()

	// Define 5 concurrent tool calls
	tools := []struct {
		name string
		call func() error
	}{
		{
			name: "search",
			call: func() error {
				_, err := client.Search(ctx, "Inception", 1, nil)
				return err
			},
		},
		{
			name: "get_details",
			call: func() error {
				_, err := client.GetMovieDetails(ctx, 27205, nil)
				return err
			},
		},
		{
			name: "discover_movies",
			call: func() error {
				params := tmdb.DiscoverMoviesParams{
					WithGenres:     "878",
					VoteAverageGte: 8.0,
				}
				_, err := client.DiscoverMovies(ctx, params)
				return err
			},
		},
		{
			name: "get_trending",
			call: func() error {
				_, err := client.GetTrending(ctx, "movie", "day", 1)
				return err
			},
		},
		{
			name: "get_recommendations",
			call: func() error {
				_, err := client.GetMovieRecommendations(ctx, 27205, 1)
				return err
			},
		},
	}

	// Use sync.WaitGroup to manage goroutines
	var wg sync.WaitGroup
	errors := make(chan error, len(tools))

	// Execute all tools concurrently
	start := time.Now()
	for _, tool := range tools {
		wg.Add(1)
		go func(toolName string, toolCall func() error) {
			defer wg.Done()
			err := toolCall()
			if err != nil {
				errors <- fmt.Errorf("%s failed: %w", toolName, err)
			} else {
				t.Logf("%s completed successfully", toolName)
			}
		}(tool.name, tool.call)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errors)
	duration := time.Since(start)

	// Check for errors
	var errs []error
	for err := range errors {
		errs = append(errs, err)
	}

	// Verify all calls succeeded (no 429 errors or other failures)
	assert.Empty(t, errs, "All concurrent calls should succeed without errors")
	if len(errs) > 0 {
		for _, err := range errs {
			t.Logf("Error: %v", err)
		}
	}

	// Log concurrent execution time
	t.Logf("Concurrent execution completed in %v", duration)

	// Note: Due to rate limiting, concurrent execution may not be much faster than sequential
	// The main goal is to verify that rate limiting works correctly without errors
}

// ==================== Error Scenario Tests ====================

// TestErrorScenario_InvalidParameters tests handling of invalid parameters
func TestErrorScenario_InvalidParameters(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)
	ctx := context.Background()

	// Test 1: Invalid vote_average range for discover_movies
	params := tmdb.DiscoverMoviesParams{
		VoteAverageGte: 11.0, // Invalid: should be 0-10
	}
	_, err := client.DiscoverMovies(ctx, params)
	assert.Error(t, err, "Should return error for invalid vote_average")
	assert.Contains(t, err.Error(), "vote_average.gte must be between 0 and 10")

	// Test 2: Invalid media_type for get_trending
	_, err = client.GetTrending(ctx, "invalid_type", "day", 1)
	assert.Error(t, err, "Should return error for invalid media_type")
	assert.Contains(t, err.Error(), "invalid media_type")

	// Test 3: Invalid time_window for get_trending
	_, err = client.GetTrending(ctx, "movie", "invalid_window", 1)
	assert.Error(t, err, "Should return error for invalid time_window")
	assert.Contains(t, err.Error(), "invalid time_window")

	// Test 4: Invalid ID (ID <= 0) for get_details
	_, err = client.GetMovieDetails(ctx, 0, nil)
	assert.Error(t, err, "Should return error for ID <= 0")
	assert.Contains(t, err.Error(), "invalid movie ID")

	t.Log("All error scenario tests passed - parameters are validated correctly")
}

// TestErrorScenario_NonExistentID tests handling of non-existent IDs
func TestErrorScenario_NonExistentID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("TMDB_API_KEY environment variable not set, skipping integration test")
	}

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    apiKey,
		Language:  "en-US",
		RateLimit: 40,
	}
	client := tmdb.NewClient(cfg, logger)
	ctx := context.Background()

	// Test with a very large ID that doesn't exist
	details, err := client.GetMovieDetails(ctx, 999999999, nil)

	// According to our implementation, 404 returns nil, nil (not an error)
	assert.NoError(t, err, "404 should not return error")
	assert.Nil(t, details, "Non-existent ID should return nil details")

	t.Log("Non-existent ID handled correctly (returns nil without error)")
}
