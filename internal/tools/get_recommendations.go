package tools

import (
	"context"
	"fmt"

	"github.com/XDwanj/tmdb-mcp/internal/tmdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// GetRecommendationsTool implements the MCP get_recommendations tool
type GetRecommendationsTool struct {
	tmdbClient *tmdb.Client
	logger     *zap.Logger
}

// NewGetRecommendationsTool creates a new GetRecommendationsTool instance
func NewGetRecommendationsTool(tmdbClient *tmdb.Client, logger *zap.Logger) *GetRecommendationsTool {
	return &GetRecommendationsTool{
		tmdbClient: tmdbClient,
		logger:     logger,
	}
}

// Name returns the tool name
func (t *GetRecommendationsTool) Name() string {
	return "get_recommendations"
}

// Description returns the tool description
func (t *GetRecommendationsTool) Description() string {
	return `Get movie or TV show recommendations based on a specific title you like.

Examples:
- Get movie recommendations based on Inception (ID: 27205): media_type=movie, id=27205
- Get TV show recommendations based on Breaking Bad (ID: 1396): media_type=tv, id=1396

Parameters:
- media_type: Type of media to get recommendations for (movie/tv)
- id: TMDB ID of the movie or TV show
- page: Page number (optional, default: 1)
- language: ISO 639-1 language code (optional, uses config default if not specified)`
}

// Handler returns a handler function compatible with mcp.AddTool
// This allows the tool to be registered with the MCP server while keeping
// business logic encapsulated in the GetRecommendationsTool struct
func (t *GetRecommendationsTool) Handler() func(context.Context, *mcp.CallToolRequest, GetRecommendationsParams) (*mcp.CallToolResult, GetRecommendationsResponse, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params GetRecommendationsParams) (*mcp.CallToolResult, GetRecommendationsResponse, error) {
		// Set default page
		page := 1
		if params.Page != nil {
			page = *params.Page
		}

		t.logger.Info("GetRecommendations request received",
			zap.String("media_type", params.MediaType),
			zap.Int("id", params.ID),
			zap.Int("page", page),
		)

		// Call appropriate TMDB Client method based on media type
		var results *tmdb.RecommendationsResponse
		var err error

		switch params.MediaType {
		case "movie":
			results, err = t.tmdbClient.GetMovieRecommendations(ctx, params.ID, page)
		case "tv":
			results, err = t.tmdbClient.GetTVRecommendations(ctx, params.ID, page)
		default:
			err = fmt.Errorf("invalid media_type: %s, must be movie or tv", params.MediaType)
		}

		if err != nil {
			t.logger.Error("GetRecommendations failed",
				zap.Error(err),
				zap.String("media_type", params.MediaType),
				zap.Int("id", params.ID),
			)
			return nil, GetRecommendationsResponse{}, convertTMDBError(err, "content")
		}

		t.logger.Info("GetRecommendations completed",
			zap.String("media_type", params.MediaType),
			zap.Int("id", params.ID),
			zap.Int("results", len(results.Results)),
		)

		// Return empty result metadata and structured response
		return &mcp.CallToolResult{}, GetRecommendationsResponse{Results: results.Results}, nil
	}
}
