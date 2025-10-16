package tmdb

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// GetMovieRecommendations gets movie recommendations based on a movie ID
func (c *Client) GetMovieRecommendations(ctx context.Context, id int, page int) (*RecommendationsResponse, error) {
	// 验证 ID 参数
	if id <= 0 {
		return nil, fmt.Errorf("invalid movie ID: %d, must be greater than 0", id)
	}

	// 设置默认页码
	if page == 0 {
		page = 1
	}

	// 构建端点路径
	endpoint := fmt.Sprintf("/movie/%d/recommendations", id)

	return c.getRecommendations(ctx, endpoint, "movie", id, page)
}

// GetTVRecommendations gets TV show recommendations based on a TV show ID
func (c *Client) GetTVRecommendations(ctx context.Context, id int, page int) (*RecommendationsResponse, error) {
	// 验证 ID 参数
	if id <= 0 {
		return nil, fmt.Errorf("invalid TV show ID: %d, must be greater than 0", id)
	}

	// 设置默认页码
	if page == 0 {
		page = 1
	}

	// 构建端点路径
	endpoint := fmt.Sprintf("/tv/%d/recommendations", id)

	return c.getRecommendations(ctx, endpoint, "tv", id, page)
}

// getRecommendations is a shared helper method for getting recommendations
func (c *Client) getRecommendations(ctx context.Context, endpoint, mediaType string, id, page int) (*RecommendationsResponse, error) {
	// Rate limiting is handled by OnBeforeRequest middleware
	// 调用 TMDB API /movie/{id}/recommendations 或 /tv/{id}/recommendations 端点
	var recommendationsResp RecommendationsResponse
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetQueryParam("page", fmt.Sprintf("%d", page)).
		SetResult(&recommendationsResp).
		Get(endpoint)

	if err != nil {
		return nil, fmt.Errorf("get recommendations failed: %w", err)
	}

	// 处理 HTTP 错误
	if resp.IsError() {
		statusCode := resp.StatusCode()

		// 404 返回空结果，不返回错误
		if statusCode == 404 {
			c.logger.Info("GetRecommendations returned no results",
				zap.String("endpoint", endpoint),
				zap.String("media_type", mediaType),
				zap.Int("id", id),
				zap.Int("status_code", statusCode),
			)
			return &RecommendationsResponse{
				Page:         page,
				Results:      []RecommendationResult{},
				TotalPages:   0,
				TotalResults: 0,
			}, nil
		}

		// 其他错误使用 handleError 处理
		err := handleError(resp)
		return nil, fmt.Errorf("get recommendations API error: %w", err)
	}

	return &recommendationsResp, nil
}
