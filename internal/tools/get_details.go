package tools

import (
	"context"
	"fmt"

	"github.com/XDwanj/tmdb-mcp/internal/tmdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// GetDetailsTool implements the MCP get_details tool
type GetDetailsTool struct {
	tmdbClient *tmdb.Client
	logger     *zap.Logger
}

// NewGetDetailsTool creates a new GetDetailsTool instance
func NewGetDetailsTool(tmdbClient *tmdb.Client, logger *zap.Logger) *GetDetailsTool {
	return &GetDetailsTool{
		tmdbClient: tmdbClient,
		logger:     logger,
	}
}

// Name returns the tool name
func (t *GetDetailsTool) Name() string {
	return "get_details"
}

// Description returns the tool description
func (t *GetDetailsTool) Description() string {
	return "Get detailed information about a movie, TV show, or person using their TMDB ID"
}

// Handler returns a handler function compatible with mcp.AddTool
// This allows the tool to be registered with the MCP server while keeping
// business logic encapsulated in the GetDetailsTool struct
func (t *GetDetailsTool) Handler() func(context.Context, *mcp.CallToolRequest, GetDetailsParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params GetDetailsParams) (*mcp.CallToolResult, any, error) {
		// 验证 media_type 参数
		validMediaTypes := map[string]bool{
			"movie":  true,
			"tv":     true,
			"person": true,
		}
		if !validMediaTypes[params.MediaType] {
			return nil, nil, fmt.Errorf("invalid media_type: must be 'movie', 'tv', or 'person'")
		}

		t.logger.Info("Get details request received",
			zap.String("media_type", params.MediaType),
			zap.Int("id", params.ID),
		)

		// 根据 media_type 调用相应的 TMDB Client 方法
		switch params.MediaType {
		case "movie":
			movieDetails, err := t.tmdbClient.GetMovieDetails(ctx, params.ID)
			if err != nil {
				t.logger.Error("Get details failed",
					zap.Error(err),
					zap.String("media_type", params.MediaType),
					zap.Int("id", params.ID),
				)
				return nil, nil, err
			}
			// 检查资源是否存在（404 情况）
			if movieDetails == nil {
				t.logger.Warn("Resource not found",
					zap.String("media_type", params.MediaType),
					zap.Int("id", params.ID),
				)
				return nil, nil, fmt.Errorf("%s with ID %d not found", params.MediaType, params.ID)
			}
			t.logger.Info("Get details completed",
				zap.String("media_type", params.MediaType),
				zap.Int("id", params.ID),
			)
			return &mcp.CallToolResult{}, movieDetails, nil

		case "tv":
			tvDetails, err := t.tmdbClient.GetTVDetails(ctx, params.ID)
			if err != nil {
				t.logger.Error("Get details failed",
					zap.Error(err),
					zap.String("media_type", params.MediaType),
					zap.Int("id", params.ID),
				)
				return nil, nil, err
			}
			// 检查资源是否存在（404 情况）
			if tvDetails == nil {
				t.logger.Warn("Resource not found",
					zap.String("media_type", params.MediaType),
					zap.Int("id", params.ID),
				)
				return nil, nil, fmt.Errorf("%s with ID %d not found", params.MediaType, params.ID)
			}
			t.logger.Info("Get details completed",
				zap.String("media_type", params.MediaType),
				zap.Int("id", params.ID),
			)
			return &mcp.CallToolResult{}, tvDetails, nil

		case "person":
			personDetails, err := t.tmdbClient.GetPersonDetails(ctx, params.ID)
			if err != nil {
				t.logger.Error("Get details failed",
					zap.Error(err),
					zap.String("media_type", params.MediaType),
					zap.Int("id", params.ID),
				)
				return nil, nil, err
			}
			// 检查资源是否存在（404 情况）
			if personDetails == nil {
				t.logger.Warn("Resource not found",
					zap.String("media_type", params.MediaType),
					zap.Int("id", params.ID),
				)
				return nil, nil, fmt.Errorf("%s with ID %d not found", params.MediaType, params.ID)
			}
			t.logger.Info("Get details completed",
				zap.String("media_type", params.MediaType),
				zap.Int("id", params.ID),
			)
			return &mcp.CallToolResult{}, personDetails, nil
		}

		// 不应该到达这里（已经验证了 media_type）
		return nil, nil, fmt.Errorf("invalid media_type: %s", params.MediaType)
	}
}
