package tools

import (
	"context"

	"github.com/XDwanj/tmdb-mcp/internal/tmdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// GetTrendingTool implements the MCP get_trending tool
type GetTrendingTool struct {
	tmdbClient *tmdb.Client
	logger     *zap.Logger
}

// NewGetTrendingTool creates a new GetTrendingTool instance
func NewGetTrendingTool(tmdbClient *tmdb.Client, logger *zap.Logger) *GetTrendingTool {
	return &GetTrendingTool{
		tmdbClient: tmdbClient,
		logger:     logger,
	}
}

// Name returns the tool name
func (t *GetTrendingTool) Name() string {
	return "get_trending"
}

// Description returns the tool description
func (t *GetTrendingTool) Description() string {
	return `Get trending movies, TV shows, or people for a specific time window (day or week).

Examples:
- Get today's trending movies: media_type=movie, time_window=day
- Get this week's trending TV shows: media_type=tv, time_window=week
- Get today's trending people: media_type=person, time_window=day

Parameters:
- media_type: Type of media (movie/tv/person)
- time_window: Time window for trending items (day/week)
- page: Page number (optional, default: 1)
- language: ISO 639-1 language code (optional, uses config default if not specified)`
}

// Handler returns a handler function compatible with mcp.AddTool
// This allows the tool to be registered with the MCP server while keeping
// business logic encapsulated in the GetTrendingTool struct
func (t *GetTrendingTool) Handler() func(context.Context, *mcp.CallToolRequest, GetTrendingParams) (*mcp.CallToolResult, GetTrendingResponse, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params GetTrendingParams) (*mcp.CallToolResult, GetTrendingResponse, error) {
		// Set default page
		page := 1
		if params.Page != nil {
			page = *params.Page
		}

		// Call TMDB Client (validation is done in the client layer)
		results, err := t.tmdbClient.GetTrending(ctx, params.MediaType, params.TimeWindow, page)
		if err != nil {
			return nil, GetTrendingResponse{}, convertTMDBError(err, "content")
		}

		// Return empty result metadata and structured response
		return &mcp.CallToolResult{}, GetTrendingResponse{Results: results.Results}, nil
	}
}
