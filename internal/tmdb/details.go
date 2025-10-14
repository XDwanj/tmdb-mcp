package tmdb

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// GetMovieDetails gets detailed information about a movie using its TMDB ID
func (c *Client) GetMovieDetails(ctx context.Context, id int) (*MovieDetails, error) {
	// 验证参数
	if id <= 0 {
		return nil, fmt.Errorf("invalid movie ID: %d", id)
	}

	c.logger.Debug("Fetching movie details",
		zap.Int("id", id),
	)

	// 等待速率限制
	if err := c.rateLimiter.Wait(ctx); err != nil {
		c.logger.Error("rate limit wait failed", zap.Error(err))
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	// 调用 TMDB API /movie/{id} 端点
	var details MovieDetails
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetQueryParam("append_to_response", "credits,videos").
		SetResult(&details).
		Get(fmt.Sprintf("/movie/%d", id))

	if err != nil {
		c.logger.Error("Get movie details failed", zap.Error(err), zap.Int("id", id))
		return nil, fmt.Errorf("get movie details failed: %w", err)
	}

	// 处理 HTTP 错误
	if resp.IsError() {
		// 404 返回 nil, nil（资源不存在不算错误）
		if resp.StatusCode() == 404 {
			c.logger.Info("Movie not found", zap.Int("id", id))
			return nil, nil
		}

		// 其他错误使用 handleError 处理
		err := handleError(resp)
		c.logger.Error("Get movie details API error", zap.Error(err), zap.Int("id", id))
		return nil, fmt.Errorf("get movie details API error: %w", err)
	}

	c.logger.Info("Movie details fetched successfully",
		zap.Int("id", id),
		zap.String("title", details.Title),
	)

	return &details, nil
}

// GetTVDetails gets detailed information about a TV show using its TMDB ID
func (c *Client) GetTVDetails(ctx context.Context, id int) (*TVDetails, error) {
	// 验证参数
	if id <= 0 {
		return nil, fmt.Errorf("invalid TV ID: %d", id)
	}

	c.logger.Debug("Fetching TV details",
		zap.Int("id", id),
	)

	// 等待速率限制
	if err := c.rateLimiter.Wait(ctx); err != nil {
		c.logger.Error("rate limit wait failed", zap.Error(err))
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	// 调用 TMDB API /tv/{id} 端点
	var details TVDetails
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetQueryParam("append_to_response", "credits,videos").
		SetResult(&details).
		Get(fmt.Sprintf("/tv/%d", id))

	if err != nil {
		c.logger.Error("Get TV details failed", zap.Error(err), zap.Int("id", id))
		return nil, fmt.Errorf("get TV details failed: %w", err)
	}

	// 处理 HTTP 错误
	if resp.IsError() {
		// 404 返回 nil, nil（资源不存在不算错误）
		if resp.StatusCode() == 404 {
			c.logger.Info("TV show not found", zap.Int("id", id))
			return nil, nil
		}

		// 其他错误使用 handleError 处理
		err := handleError(resp)
		c.logger.Error("Get TV details API error", zap.Error(err), zap.Int("id", id))
		return nil, fmt.Errorf("get TV details API error: %w", err)
	}

	c.logger.Info("TV details fetched successfully",
		zap.Int("id", id),
		zap.String("name", details.Name),
	)

	return &details, nil
}

// GetPersonDetails gets detailed information about a person using their TMDB ID
func (c *Client) GetPersonDetails(ctx context.Context, id int) (*PersonDetails, error) {
	// 验证参数
	if id <= 0 {
		return nil, fmt.Errorf("invalid person ID: %d", id)
	}

	c.logger.Debug("Fetching person details",
		zap.Int("id", id),
	)

	// 等待速率限制
	if err := c.rateLimiter.Wait(ctx); err != nil {
		c.logger.Error("rate limit wait failed", zap.Error(err))
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	// 调用 TMDB API /person/{id} 端点
	var details PersonDetails
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetQueryParam("append_to_response", "combined_credits").
		SetResult(&details).
		Get(fmt.Sprintf("/person/%d", id))

	if err != nil {
		c.logger.Error("Get person details failed", zap.Error(err), zap.Int("id", id))
		return nil, fmt.Errorf("get person details failed: %w", err)
	}

	// 处理 HTTP 错误
	if resp.IsError() {
		// 404 返回 nil, nil（资源不存在不算错误）
		if resp.StatusCode() == 404 {
			c.logger.Info("Person not found", zap.Int("id", id))
			return nil, nil
		}

		// 其他错误使用 handleError 处理
		err := handleError(resp)
		c.logger.Error("Get person details API error", zap.Error(err), zap.Int("id", id))
		return nil, fmt.Errorf("get person details API error: %w", err)
	}

	c.logger.Info("Person details fetched successfully",
		zap.Int("id", id),
		zap.String("name", details.Name),
	)

	return &details, nil
}
