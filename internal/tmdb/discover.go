package tmdb

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// DiscoverMoviesParams represents parameters for discovering movies
type DiscoverMoviesParams struct {
	WithGenres           string
	PrimaryReleaseYear   int
	VoteAverageGte       float64
	VoteAverageLte       float64
	WithOriginalLanguage string
	SortBy               string
	Page                 int
	Language             string
}

// DiscoverTVParams represents parameters for discovering TV shows
type DiscoverTVParams struct {
	WithGenres           string
	FirstAirDateYear     int
	VoteAverageGte       float64
	VoteAverageLte       float64
	WithOriginalLanguage string
	WithStatus           string
	SortBy               string
	Page                 int
	Language             string
}

// DiscoverMovies discovers movies using various filters
func (c *Client) DiscoverMovies(ctx context.Context, params DiscoverMoviesParams) (*DiscoverMoviesResponse, error) {
	// 参数验证
	if params.VoteAverageGte < 0 || params.VoteAverageGte > 10 {
		return nil, fmt.Errorf("vote_average.gte must be between 0 and 10")
	}
	if params.VoteAverageLte < 0 || params.VoteAverageLte > 10 {
		return nil, fmt.Errorf("vote_average.lte must be between 0 and 10")
	}

	// 设置默认值
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.SortBy == "" {
		params.SortBy = "popularity.desc"
	}

	c.logger.Debug("Discovering movies",
		zap.String("with_genres", params.WithGenres),
		zap.Int("primary_release_year", params.PrimaryReleaseYear),
		zap.Float64("vote_average_gte", params.VoteAverageGte),
		zap.Float64("vote_average_lte", params.VoteAverageLte),
		zap.String("sort_by", params.SortBy),
		zap.Int("page", params.Page),
	)

	// 等待速率限制
	if err := c.rateLimiter.Wait(ctx); err != nil {
		c.logger.Error("rate limit wait failed", zap.Error(err))
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	// 构建请求
	req := c.httpClient.R().
		SetContext(ctx).
		SetQueryParam("sort_by", params.SortBy).
		SetQueryParam("page", fmt.Sprintf("%d", params.Page))

	// 添加可选参数（仅当非空时）
	if params.WithGenres != "" {
		req.SetQueryParam("with_genres", params.WithGenres)
	}
	if params.PrimaryReleaseYear > 0 {
		req.SetQueryParam("primary_release_year", fmt.Sprintf("%d", params.PrimaryReleaseYear))
	}
	if params.VoteAverageGte > 0 {
		req.SetQueryParam("vote_average.gte", fmt.Sprintf("%.1f", params.VoteAverageGte))
	}
	if params.VoteAverageLte > 0 {
		req.SetQueryParam("vote_average.lte", fmt.Sprintf("%.1f", params.VoteAverageLte))
	}
	if params.WithOriginalLanguage != "" {
		req.SetQueryParam("with_original_language", params.WithOriginalLanguage)
	}
	// 如果指定了 language 参数，添加到请求中（会覆盖 OnBeforeRequest 中的默认值）
	if params.Language != "" {
		req.SetQueryParam("language", params.Language)
	}

	// 调用 TMDB API /discover/movie 端点
	var response DiscoverMoviesResponse
	resp, err := req.SetResult(&response).Get("/discover/movie")

	if err != nil {
		c.logger.Error("Discover movies failed", zap.Error(err))
		return nil, fmt.Errorf("discover movies failed: %w", err)
	}

	// 处理 HTTP 错误
	if resp.IsError() {
		// 404 返回空结果
		if resp.StatusCode() == 404 {
			c.logger.Info("Discover movies returned no results")
			return &DiscoverMoviesResponse{
				Page:         params.Page,
				Results:      []DiscoverMovieResult{},
				TotalPages:   0,
				TotalResults: 0,
			}, nil
		}

		// 其他错误使用 handleError 处理
		err := handleError(resp)
		c.logger.Error("Discover movies API error", zap.Error(err))
		return nil, fmt.Errorf("discover movies API error: %w", err)
	}

	c.logger.Info("Movies discovered successfully",
		zap.Int("count", len(response.Results)),
		zap.Int("total_results", response.TotalResults),
	)

	return &response, nil
}

// DiscoverTV discovers TV shows using various filters
func (c *Client) DiscoverTV(ctx context.Context, params DiscoverTVParams) (*DiscoverTVResponse, error) {
	// 参数验证
	if params.VoteAverageGte < 0 || params.VoteAverageGte > 10 {
		return nil, fmt.Errorf("vote_average.gte must be between 0 and 10")
	}
	if params.VoteAverageLte < 0 || params.VoteAverageLte > 10 {
		return nil, fmt.Errorf("vote_average.lte must be between 0 and 10")
	}

	// 设置默认值
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.SortBy == "" {
		params.SortBy = "popularity.desc"
	}

	c.logger.Debug("Discovering TV shows",
		zap.String("with_genres", params.WithGenres),
		zap.Int("first_air_date_year", params.FirstAirDateYear),
		zap.Float64("vote_average_gte", params.VoteAverageGte),
		zap.Float64("vote_average_lte", params.VoteAverageLte),
		zap.String("with_status", params.WithStatus),
		zap.String("sort_by", params.SortBy),
		zap.Int("page", params.Page),
	)

	// 等待速率限制
	if err := c.rateLimiter.Wait(ctx); err != nil {
		c.logger.Error("rate limit wait failed", zap.Error(err))
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	// 构建请求
	req := c.httpClient.R().
		SetContext(ctx).
		SetQueryParam("sort_by", params.SortBy).
		SetQueryParam("page", fmt.Sprintf("%d", params.Page))

	// 添加可选参数（仅当非空时）
	if params.WithGenres != "" {
		req.SetQueryParam("with_genres", params.WithGenres)
	}
	if params.FirstAirDateYear > 0 {
		req.SetQueryParam("first_air_date_year", fmt.Sprintf("%d", params.FirstAirDateYear))
	}
	if params.VoteAverageGte > 0 {
		req.SetQueryParam("vote_average.gte", fmt.Sprintf("%.1f", params.VoteAverageGte))
	}
	if params.VoteAverageLte > 0 {
		req.SetQueryParam("vote_average.lte", fmt.Sprintf("%.1f", params.VoteAverageLte))
	}
	if params.WithOriginalLanguage != "" {
		req.SetQueryParam("with_original_language", params.WithOriginalLanguage)
	}
	if params.WithStatus != "" {
		req.SetQueryParam("with_status", params.WithStatus)
	}
	// 如果指定了 language 参数，添加到请求中（会覆盖 OnBeforeRequest 中的默认值）
	if params.Language != "" {
		req.SetQueryParam("language", params.Language)
	}

	// 调用 TMDB API /discover/tv 端点
	var response DiscoverTVResponse
	resp, err := req.SetResult(&response).Get("/discover/tv")

	if err != nil {
		c.logger.Error("Discover TV shows failed", zap.Error(err))
		return nil, fmt.Errorf("discover TV shows failed: %w", err)
	}

	// 处理 HTTP 错误
	if resp.IsError() {
		// 404 返回空结果
		if resp.StatusCode() == 404 {
			c.logger.Info("Discover TV shows returned no results")
			return &DiscoverTVResponse{
				Page:         params.Page,
				Results:      []DiscoverTVResult{},
				TotalPages:   0,
				TotalResults: 0,
			}, nil
		}

		// 其他错误使用 handleError 处理
		err := handleError(resp)
		c.logger.Error("Discover TV shows API error", zap.Error(err))
		return nil, fmt.Errorf("discover TV shows API error: %w", err)
	}

	c.logger.Info("TV shows discovered successfully",
		zap.Int("count", len(response.Results)),
		zap.Int("total_results", response.TotalResults),
	)

	return &response, nil
}
