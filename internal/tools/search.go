package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/XDwanj/tmdb-mcp/internal/tmdb"
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

// Call executes the search tool
func (t *SearchTool) Call(ctx context.Context, params json.RawMessage) (interface{}, error) {
	// 解析参数
	var searchParams SearchParams
	if err := json.Unmarshal(params, &searchParams); err != nil {
		return nil, fmt.Errorf("failed to parse parameters: %w", err)
	}

	// 验证 Query 参数
	if searchParams.Query == "" {
		return nil, errors.New("query parameter is required")
	}

	// 设置默认 Page
	if searchParams.Page == 0 {
		searchParams.Page = 1
	}

	// 调用 TMDB Client
	results, err := t.tmdbClient.Search(ctx, searchParams.Query, searchParams.Page)
	if err != nil {
		t.logger.Error("Search failed",
			zap.Error(err),
			zap.String("query", searchParams.Query),
		)
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// 记录成功日志
	t.logger.Info("Search completed",
		zap.String("query", searchParams.Query),
		zap.Int("results", len(results.Results)),
	)

	return results.Results, nil
}
