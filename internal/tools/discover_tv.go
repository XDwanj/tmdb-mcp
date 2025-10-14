package tools

import (
	"context"

	"github.com/XDwanj/tmdb-mcp/internal/tmdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// DiscoverTVTool implements the MCP discover_tv tool
type DiscoverTVTool struct {
	tmdbClient *tmdb.Client
	logger     *zap.Logger
}

// NewDiscoverTVTool creates a new DiscoverTVTool instance
func NewDiscoverTVTool(tmdbClient *tmdb.Client, logger *zap.Logger) *DiscoverTVTool {
	return &DiscoverTVTool{
		tmdbClient: tmdbClient,
		logger:     logger,
	}
}

// Name returns the tool name
func (t *DiscoverTVTool) Name() string {
	return "discover_tv"
}

// Description returns the tool description
func (t *DiscoverTVTool) Description() string {
	return "Discover TV shows using filters like genre, year, rating, and status. " +
		"Example: Find high-rated crime dramas (genre: 80, vote_average.gte: 8.0) or returning sci-fi series (genre: 10765, with_status: 'Returning Series')"
}

// Handler returns a handler function compatible with mcp.AddTool
// This allows the tool to be registered with the MCP server while keeping
// business logic encapsulated in the DiscoverTVTool struct
func (t *DiscoverTVTool) Handler() func(context.Context, *mcp.CallToolRequest, DiscoverTVParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params DiscoverTVParams) (*mcp.CallToolResult, any, error) {
		t.logger.Info("Discover TV shows request received",
			zap.Stringp("with_genres", params.WithGenres),
			zap.Intp("first_air_date_year", params.FirstAirDateYear),
			zap.Float64p("vote_average_gte", params.VoteAverageGte),
			zap.Float64p("vote_average_lte", params.VoteAverageLte),
			zap.Stringp("with_status", params.WithStatus),
			zap.Stringp("sort_by", params.SortBy),
		)

		// 转换 tools.DiscoverTVParams 到 tmdb.DiscoverTVParams
		// 处理指针类型，零值使用默认值
		tmdbParams := tmdb.DiscoverTVParams{}

		if params.WithGenres != nil {
			tmdbParams.WithGenres = *params.WithGenres
		}
		if params.FirstAirDateYear != nil {
			tmdbParams.FirstAirDateYear = *params.FirstAirDateYear
		}
		if params.VoteAverageGte != nil {
			tmdbParams.VoteAverageGte = *params.VoteAverageGte
		}
		if params.VoteAverageLte != nil {
			tmdbParams.VoteAverageLte = *params.VoteAverageLte
		}
		if params.WithOriginalLanguage != nil {
			tmdbParams.WithOriginalLanguage = *params.WithOriginalLanguage
		}
		if params.WithStatus != nil {
			tmdbParams.WithStatus = *params.WithStatus
		}
		if params.SortBy != nil {
			tmdbParams.SortBy = *params.SortBy
		}
		if params.Page != nil {
			tmdbParams.Page = *params.Page
		}
		if params.Language != nil {
			tmdbParams.Language = *params.Language
		}

		// 调用 TMDB Client（参数验证在 Client 层完成）
		result, err := t.tmdbClient.DiscoverTV(ctx, tmdbParams)
		if err != nil {
			t.logger.Error("Discover TV shows failed",
				zap.Error(err),
			)
			return nil, nil, err
		}

		// 检查结果为空（但不是错误）
		if result == nil || len(result.Results) == 0 {
			t.logger.Info("No TV shows found matching criteria")
			// 返回完整的响应对象，而不是空数组
			if result != nil {
				return &mcp.CallToolResult{}, result, nil
			}
			// 如果 result 为 nil，返回一个空的响应对象
			return &mcp.CallToolResult{}, &tmdb.DiscoverTVResponse{
				Page:         1,
				Results:      []tmdb.DiscoverTVResult{},
				TotalPages:   0,
				TotalResults: 0,
			}, nil
		}

		t.logger.Info("Discover TV shows completed",
			zap.Int("count", len(result.Results)),
			zap.Int("total_results", result.TotalResults),
		)

		// 返回空的 CallToolResult 和结构化响应
		return &mcp.CallToolResult{}, result, nil
	}
}
