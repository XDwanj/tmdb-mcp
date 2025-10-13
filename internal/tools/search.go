package tools

import (
	"context"

	"github.com/XDwanj/tmdb-mcp/internal/tmdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// SearchTool implements the MCP search tool
type SearchTool struct {
	tmdbClient *tmdb.Client
	logger     *zap.Logger
}

// NewSearchTool creates a new SearchTool instance
func NewSearchTool(tmdbClient *tmdb.Client, logger *zap.Logger) *SearchTool {
	return &SearchTool{
		tmdbClient: tmdbClient,
		logger:     logger,
	}
}

// Name returns the tool name
func (t *SearchTool) Name() string {
	return "search"
}

// Description returns the tool description
func (t *SearchTool) Description() string {
	return "Search for movies, TV shows, and people on TMDB using a query string"
}

// Handler returns a handler function compatible with mcp.AddTool
// This allows the tool to be registered with the MCP server while keeping
// business logic encapsulated in the SearchTool struct
func (t *SearchTool) Handler() func(context.Context, *mcp.CallToolRequest, SearchParams) (*mcp.CallToolResult, SearchResponse, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params SearchParams) (*mcp.CallToolResult, SearchResponse, error) {
		// Set default page
		if params.Page == 0 {
			params.Page = 1
		}

		t.logger.Info("Search request received",
			zap.String("query", params.Query),
			zap.Int("page", params.Page),
		)

		// Call TMDB Client (validation is done in the client layer)
		results, err := t.tmdbClient.Search(ctx, params.Query, params.Page)
		if err != nil {
			t.logger.Error("Search failed",
				zap.Error(err),
				zap.String("query", params.Query),
			)
			return nil, SearchResponse{}, err
		}

		t.logger.Info("Search completed",
			zap.String("query", params.Query),
			zap.Int("results", len(results.Results)),
		)

		// Return empty result metadata and structured response
		return &mcp.CallToolResult{}, SearchResponse{Results: results.Results}, nil
	}
}
