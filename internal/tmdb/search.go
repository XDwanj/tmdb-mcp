package tmdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

const (
	// maxQueryLength is the maximum allowed length for search queries
	maxQueryLength = 500
)

// Search searches for movies, TV shows, and people using a query string
func (c *Client) Search(ctx context.Context, query string, page int, language *string) (*SearchResponse, error) {
	// 记录请求开始时间
	startTime := time.Now()
	endpoint := "/search/multi"

	// 验证 query 参数
	if query == "" {
		return nil, errors.New("query parameter is required")
	}

	// 验证 query 长度
	if len(query) > maxQueryLength {
		return nil, fmt.Errorf("query parameter is too long: maximum length is %d characters", maxQueryLength)
	}

	// 设置默认页码
	if page == 0 {
		page = 1
	}

	c.logger.Debug("Starting TMDB API request",
		zap.String("endpoint", endpoint),
		zap.String("query", query),
		zap.Int("page", page),
	)

	// Wait for rate limit
	if err := c.rateLimiter.Wait(ctx); err != nil {
		c.logger.Error("rate limit wait failed", zap.Error(err))
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	// 调用 TMDB API /search/multi 端点
	var searchResp SearchResponse
	req := c.httpClient.R().
		SetContext(ctx).
		SetQueryParam("query", query).
		SetQueryParam("page", fmt.Sprintf("%d", page)).
		SetResult(&searchResp)

	// 如果指定了 language 参数，添加到请求中（会覆盖 OnBeforeRequest 中的默认值）
	if language != nil && *language != "" {
		req.SetQueryParam("language", *language)
	}

	resp, err := req.Get(endpoint)
	responseTime := time.Since(startTime)

	if err != nil {
		c.logger.Error("Search failed",
			zap.String("endpoint", endpoint),
			zap.String("query", query),
			zap.String("error_type", ErrorTypeNetwork),
			zap.Duration("response_time", responseTime),
			zap.Error(err),
		)
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// 处理 HTTP 错误
	if resp.IsError() {
		statusCode := resp.StatusCode()

		// 404 返回空结果，不返回错误
		if statusCode == 404 {
			c.logger.Info("Search returned no results",
				zap.String("endpoint", endpoint),
				zap.String("query", query),
				zap.Int("status_code", statusCode),
				zap.Duration("response_time", responseTime),
			)
			return &SearchResponse{
				Page:         page,
				Results:      []SearchResult{},
				TotalPages:   0,
				TotalResults: 0,
			}, nil
		}

		// 其他错误使用 handleError 处理
		err := handleError(resp)
		errorType := ErrorTypeUnknown
		if tmdbErr, ok := err.(*TMDBError); ok {
			errorType = tmdbErr.ErrorType
		}

		c.logger.Error("Search API error",
			zap.String("endpoint", endpoint),
			zap.String("query", query),
			zap.String("error_type", errorType),
			zap.Int("status_code", statusCode),
			zap.Duration("response_time", responseTime),
			zap.Error(err),
		)
		return nil, fmt.Errorf("search API error: %w", err)
	}

	c.logger.Info("Search completed successfully",
		zap.String("endpoint", endpoint),
		zap.String("query", query),
		zap.Int("status_code", resp.StatusCode()),
		zap.Duration("response_time", responseTime),
		zap.Int("result_count", len(searchResp.Results)),
		zap.Int("total_results", searchResp.TotalResults),
	)

	return &searchResp, nil
}
