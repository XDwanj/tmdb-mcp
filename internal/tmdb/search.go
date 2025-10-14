package tmdb

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

const (
	// maxQueryLength is the maximum allowed length for search queries
	maxQueryLength = 500
)

// Search searches for movies, TV shows, and people using a query string
func (c *Client) Search(ctx context.Context, query string, page int, language *string) (*SearchResponse, error) {
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

	c.logger.Info("Searching TMDB",
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

	resp, err := req.Get("/search/multi")

	if err != nil {
		c.logger.Error("Search failed", zap.Error(err), zap.String("query", query))
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// 处理 HTTP 错误
	if resp.IsError() {
		// 404 返回空结果，不返回错误
		if resp.StatusCode() == 404 {
			c.logger.Info("Search returned no results",
				zap.String("query", query),
				zap.Int("status_code", 404),
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
		c.logger.Error("Search API error", zap.Error(err), zap.String("query", query))
		return nil, fmt.Errorf("search API error: %w", err)
	}

	c.logger.Info("Search completed",
		zap.String("query", query),
		zap.Int("results", len(searchResp.Results)),
		zap.Int("total_results", searchResp.TotalResults),
	)

	return &searchResp, nil
}
