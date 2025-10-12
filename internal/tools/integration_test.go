//go:build integration
// +build integration

package tools

import (
	"context"
	"os"
	"testing"

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
	results, err := client.Search(ctx, "Inception", 1)

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
	results, err := client.Search(ctx, "Christopher Nolan", 1)

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
	results, err := client.Search(ctx, "xyzabc123nonexistent9999", 1)

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
	page1, err := client.Search(ctx, "Star Wars", 1)
	assert.NoError(t, err)
	assert.NotNil(t, page1)
	assert.Equal(t, 1, page1.Page)
	assert.Greater(t, len(page1.Results), 0)

	// Get page 2
	page2, err := client.Search(ctx, "Star Wars", 2)
	assert.NoError(t, err)
	assert.NotNil(t, page2)
	assert.Equal(t, 2, page2.Page)

	// Verify pages have different results
	if len(page2.Results) > 0 {
		assert.NotEqual(t, page1.Results[0].ID, page2.Results[0].ID, "Different pages should have different results")
	}
}
