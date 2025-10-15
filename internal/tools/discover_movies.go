package tools

import (
	"context"

	"github.com/XDwanj/tmdb-mcp/internal/tmdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// DiscoverMoviesTool implements the MCP discover_movies tool
type DiscoverMoviesTool struct {
	tmdbClient *tmdb.Client
	logger     *zap.Logger
}

// NewDiscoverMoviesTool creates a new DiscoverMoviesTool instance
func NewDiscoverMoviesTool(tmdbClient *tmdb.Client, logger *zap.Logger) *DiscoverMoviesTool {
	return &DiscoverMoviesTool{
		tmdbClient: tmdbClient,
		logger:     logger,
	}
}

// Name returns the tool name
func (t *DiscoverMoviesTool) Name() string {
	return "discover_movies"
}

// Description returns the tool description
func (t *DiscoverMoviesTool) Description() string {
	return "Discover movies using filters like genre, year, rating, and language. " +
		"Example: Find science fiction movies (genre: 878) released after 2020 with rating ≥ 8.0"
}

// Handler returns a handler function compatible with mcp.AddTool
// This allows the tool to be registered with the MCP server while keeping
// business logic encapsulated in the DiscoverMoviesTool struct
func (t *DiscoverMoviesTool) Handler() func(context.Context, *mcp.CallToolRequest, DiscoverMoviesParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params DiscoverMoviesParams) (*mcp.CallToolResult, any, error) {
		t.logger.Info("Discover movies request received",
			zap.Stringp("with_genres", params.WithGenres),
			zap.Intp("primary_release_year", params.PrimaryReleaseYear),
			zap.Float64p("vote_average_gte", params.VoteAverageGte),
			zap.Float64p("vote_average_lte", params.VoteAverageLte),
			zap.Stringp("sort_by", params.SortBy),
		)

		// 转换 tools.DiscoverMoviesParams 到 tmdb.DiscoverMoviesParams
		// 处理指针类型，零值使用默认值
		tmdbParams := tmdb.DiscoverMoviesParams{}

		if params.WithGenres != nil {
			tmdbParams.WithGenres = *params.WithGenres
		}
		if params.PrimaryReleaseYear != nil {
			tmdbParams.PrimaryReleaseYear = *params.PrimaryReleaseYear
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
		result, err := t.tmdbClient.DiscoverMovies(ctx, tmdbParams)
		if err != nil {
			t.logger.Error("Discover movies failed",
				zap.Error(err),
			)
			return nil, nil, convertTMDBError(err, "movies")
		}

		// 检查结果为空（但不是错误）
		if result == nil || len(result.Results) == 0 {
			t.logger.Info("No movies found matching criteria")
			// 返回完整的响应对象，而不是空数组
			if result != nil {
				return &mcp.CallToolResult{}, result, nil
			}
			// 如果 result 为 nil，返回一个空的响应对象
			return &mcp.CallToolResult{}, &tmdb.DiscoverMoviesResponse{
				Page:         1,
				Results:      []tmdb.DiscoverMovieResult{},
				TotalPages:   0,
				TotalResults: 0,
			}, nil
		}

		t.logger.Info("Discover movies completed",
			zap.Int("count", len(result.Results)),
			zap.Int("total_results", result.TotalResults),
		)

		// 返回空的 CallToolResult 和结构化响应
		return &mcp.CallToolResult{}, result, nil
	}
}
