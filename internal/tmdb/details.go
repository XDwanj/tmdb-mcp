package tmdb

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// GetMovieDetails gets detailed information about a movie using its TMDB ID
func (c *Client) GetMovieDetails(ctx context.Context, id int, language *string) (*MovieDetails, error) {
	// 记录请求开始时间
	startTime := time.Now()
	endpoint := fmt.Sprintf("/movie/%d", id)

	// 验证参数
	if id <= 0 {
		return nil, fmt.Errorf("invalid movie ID: %d", id)
	}

	c.logger.Debug("Starting TMDB API request",
		zap.String("endpoint", endpoint),
		zap.Int("id", id),
	)

	// Rate limiting is handled by OnBeforeRequest middleware
	// 调用 TMDB API /movie/{id} 端点
	var details MovieDetails
	req := c.httpClient.R().
		SetContext(ctx).
		SetQueryParam("append_to_response", "credits,videos").
		SetResult(&details)

	// 如果指定了 language 参数，添加到请求中（会覆盖 OnBeforeRequest 中的默认值）
	if language != nil && *language != "" {
		req.SetQueryParam("language", *language)
	}

	resp, err := req.Get(endpoint)
	responseTime := time.Since(startTime)

	if err != nil {
		c.logger.Error("Get movie details failed",
			zap.String("endpoint", endpoint),
			zap.Int("id", id),
			zap.String("error_type", ErrorTypeNetwork),
			zap.Duration("response_time", responseTime),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get movie details failed: %w", err)
	}

	// 处理 HTTP 错误
	if resp.IsError() {
		statusCode := resp.StatusCode()

		// 404 返回 nil, nil（资源不存在不算错误）
		if statusCode == 404 {
			c.logger.Info("Movie not found",
				zap.String("endpoint", endpoint),
				zap.Int("id", id),
				zap.Int("status_code", statusCode),
				zap.Duration("response_time", responseTime),
			)
			return nil, nil
		}

		// 其他错误使用 handleError 处理
		err := handleError(resp)
		errorType := ErrorTypeUnknown
		if tmdbErr, ok := err.(*TMDBError); ok {
			errorType = tmdbErr.ErrorType
		}

		c.logger.Error("Get movie details API error",
			zap.String("endpoint", endpoint),
			zap.Int("id", id),
			zap.String("error_type", errorType),
			zap.Int("status_code", statusCode),
			zap.Duration("response_time", responseTime),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get movie details API error: %w", err)
	}

	c.logger.Info("Movie details fetched successfully",
		zap.String("endpoint", endpoint),
		zap.Int("id", id),
		zap.String("title", details.Title),
		zap.Int("status_code", resp.StatusCode()),
		zap.Duration("response_time", responseTime),
	)

	return &details, nil
}

// GetTVDetails gets detailed information about a TV show using its TMDB ID
func (c *Client) GetTVDetails(ctx context.Context, id int, language *string) (*TVDetails, error) {
	// 记录请求开始时间
	startTime := time.Now()
	endpoint := fmt.Sprintf("/tv/%d", id)

	// 验证参数
	if id <= 0 {
		return nil, fmt.Errorf("invalid TV ID: %d", id)
	}

	c.logger.Debug("Starting TMDB API request",
		zap.String("endpoint", endpoint),
		zap.Int("id", id),
	)

	// Rate limiting is handled by OnBeforeRequest middleware
	// 调用 TMDB API /tv/{id} 端点
	var details TVDetails
	req := c.httpClient.R().
		SetContext(ctx).
		SetQueryParam("append_to_response", "credits,videos").
		SetResult(&details)

	// 如果指定了 language 参数，添加到请求中（会覆盖 OnBeforeRequest 中的默认值）
	if language != nil && *language != "" {
		req.SetQueryParam("language", *language)
	}

	resp, err := req.Get(endpoint)
	responseTime := time.Since(startTime)

	if err != nil {
		c.logger.Error("Get TV details failed",
			zap.String("endpoint", endpoint),
			zap.Int("id", id),
			zap.String("error_type", ErrorTypeNetwork),
			zap.Duration("response_time", responseTime),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get TV details failed: %w", err)
	}

	// 处理 HTTP 错误
	if resp.IsError() {
		statusCode := resp.StatusCode()

		// 404 返回 nil, nil（资源不存在不算错误）
		if statusCode == 404 {
			c.logger.Info("TV show not found",
				zap.String("endpoint", endpoint),
				zap.Int("id", id),
				zap.Int("status_code", statusCode),
				zap.Duration("response_time", responseTime),
			)
			return nil, nil
		}

		// 其他错误使用 handleError 处理
		err := handleError(resp)
		errorType := ErrorTypeUnknown
		if tmdbErr, ok := err.(*TMDBError); ok {
			errorType = tmdbErr.ErrorType
		}

		c.logger.Error("Get TV details API error",
			zap.String("endpoint", endpoint),
			zap.Int("id", id),
			zap.String("error_type", errorType),
			zap.Int("status_code", statusCode),
			zap.Duration("response_time", responseTime),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get TV details API error: %w", err)
	}

	c.logger.Info("TV details fetched successfully",
		zap.String("endpoint", endpoint),
		zap.Int("id", id),
		zap.String("name", details.Name),
		zap.Int("status_code", resp.StatusCode()),
		zap.Duration("response_time", responseTime),
	)

	return &details, nil
}

// GetPersonDetails gets detailed information about a person using their TMDB ID
func (c *Client) GetPersonDetails(ctx context.Context, id int, language *string) (*PersonDetails, error) {
	// 记录请求开始时间
	startTime := time.Now()
	endpoint := fmt.Sprintf("/person/%d", id)

	// 验证参数
	if id <= 0 {
		return nil, fmt.Errorf("invalid person ID: %d", id)
	}

	c.logger.Debug("Starting TMDB API request",
		zap.String("endpoint", endpoint),
		zap.Int("id", id),
	)

	// Rate limiting is handled by OnBeforeRequest middleware
	// 调用 TMDB API /person/{id} 端点
	var details PersonDetails
	req := c.httpClient.R().
		SetContext(ctx).
		SetQueryParam("append_to_response", "combined_credits").
		SetResult(&details)

	// 如果指定了 language 参数，添加到请求中（会覆盖 OnBeforeRequest 中的默认值）
	if language != nil && *language != "" {
		req.SetQueryParam("language", *language)
	}

	resp, err := req.Get(endpoint)
	responseTime := time.Since(startTime)

	if err != nil {
		c.logger.Error("Get person details failed",
			zap.String("endpoint", endpoint),
			zap.Int("id", id),
			zap.String("error_type", ErrorTypeNetwork),
			zap.Duration("response_time", responseTime),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get person details failed: %w", err)
	}

	// 处理 HTTP 错误
	if resp.IsError() {
		statusCode := resp.StatusCode()

		// 404 返回 nil, nil（资源不存在不算错误）
		if statusCode == 404 {
			c.logger.Info("Person not found",
				zap.String("endpoint", endpoint),
				zap.Int("id", id),
				zap.Int("status_code", statusCode),
				zap.Duration("response_time", responseTime),
			)
			return nil, nil
		}

		// 其他错误使用 handleError 处理
		err := handleError(resp)
		errorType := ErrorTypeUnknown
		if tmdbErr, ok := err.(*TMDBError); ok {
			errorType = tmdbErr.ErrorType
		}

		c.logger.Error("Get person details API error",
			zap.String("endpoint", endpoint),
			zap.Int("id", id),
			zap.String("error_type", errorType),
			zap.Int("status_code", statusCode),
			zap.Duration("response_time", responseTime),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get person details API error: %w", err)
	}

	c.logger.Info("Person details fetched successfully",
		zap.String("endpoint", endpoint),
		zap.Int("id", id),
		zap.String("name", details.Name),
		zap.Int("status_code", resp.StatusCode()),
		zap.Duration("response_time", responseTime),
	)

	return &details, nil
}
