package tmdb

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// GetMovieDetails gets detailed information about a movie using its TMDB ID
func (c *Client) GetMovieDetails(ctx context.Context, id int, language *string) (*MovieDetails, error) {
	endpoint := fmt.Sprintf("/movie/%d", id)

	// 验证参数
	if id <= 0 {
		return nil, fmt.Errorf("invalid movie ID: %d", id)
	}

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

	if err != nil {
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
			)
			return nil, nil
		}

		// 其他错误使用 handleError 处理
		err := handleError(resp)
		return nil, fmt.Errorf("get movie details API error: %w", err)
	}

	return &details, nil
}

// GetTVDetails gets detailed information about a TV show using its TMDB ID
func (c *Client) GetTVDetails(ctx context.Context, id int, language *string) (*TVDetails, error) {
	endpoint := fmt.Sprintf("/tv/%d", id)

	// 验证参数
	if id <= 0 {
		return nil, fmt.Errorf("invalid TV ID: %d", id)
	}

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

	if err != nil {
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
			)
			return nil, nil
		}

		// 其他错误使用 handleError 处理
		err := handleError(resp)
		return nil, fmt.Errorf("get TV details API error: %w", err)
	}

	return &details, nil
}

// GetPersonDetails gets detailed information about a person using their TMDB ID
func (c *Client) GetPersonDetails(ctx context.Context, id int, language *string) (*PersonDetails, error) {
	endpoint := fmt.Sprintf("/person/%d", id)

	// 验证参数
	if id <= 0 {
		return nil, fmt.Errorf("invalid person ID: %d", id)
	}

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

	if err != nil {
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
			)
			return nil, nil
		}

		// 其他错误使用 handleError 处理
		err := handleError(resp)
		return nil, fmt.Errorf("get person details API error: %w", err)
	}

	return &details, nil
}
