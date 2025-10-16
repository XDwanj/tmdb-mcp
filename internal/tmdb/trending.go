package tmdb

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// GetTrending gets trending movies, TV shows, or people for a specific time window
func (c *Client) GetTrending(ctx context.Context, mediaType, timeWindow string, page int) (*TrendingResponse, error) {
	// 记录请求开始时间
	startTime := time.Now()

	// 验证 mediaType 参数
	if mediaType != "movie" && mediaType != "tv" && mediaType != "person" {
		return nil, fmt.Errorf("invalid media_type: %s, must be movie, tv, or person", mediaType)
	}

	// 验证 timeWindow 参数
	if timeWindow != "day" && timeWindow != "week" {
		return nil, fmt.Errorf("invalid time_window: %s, must be day or week", timeWindow)
	}

	// 设置默认页码
	if page == 0 {
		page = 1
	}

	// 构建端点路径
	endpoint := fmt.Sprintf("/trending/%s/%s", mediaType, timeWindow)

	c.logger.Debug("Starting TMDB API request",
		zap.String("endpoint", endpoint),
		zap.String("media_type", mediaType),
		zap.String("time_window", timeWindow),
		zap.Int("page", page),
	)

	// Rate limiting is handled by OnBeforeRequest middleware
	// 调用 TMDB API /trending/{media_type}/{time_window} 端点
	var trendingResp TrendingResponse
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetQueryParam("page", fmt.Sprintf("%d", page)).
		SetResult(&trendingResp).
		Get(endpoint)
	responseTime := time.Since(startTime)

	if err != nil {
		c.logger.Error("GetTrending failed",
			zap.String("endpoint", endpoint),
			zap.String("media_type", mediaType),
			zap.String("time_window", timeWindow),
			zap.String("error_type", ErrorTypeNetwork),
			zap.Duration("response_time", responseTime),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get trending failed: %w", err)
	}

	// 处理 HTTP 错误
	if resp.IsError() {
		statusCode := resp.StatusCode()

		// 404 返回空结果，不返回错误
		if statusCode == 404 {
			c.logger.Info("GetTrending returned no results",
				zap.String("endpoint", endpoint),
				zap.String("media_type", mediaType),
				zap.String("time_window", timeWindow),
				zap.Int("status_code", statusCode),
				zap.Duration("response_time", responseTime),
			)
			return &TrendingResponse{
				Page:         page,
				Results:      []TrendingResult{},
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

		c.logger.Error("GetTrending API error",
			zap.String("endpoint", endpoint),
			zap.String("media_type", mediaType),
			zap.String("time_window", timeWindow),
			zap.String("error_type", errorType),
			zap.Int("status_code", statusCode),
			zap.Duration("response_time", responseTime),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get trending API error: %w", err)
	}

	c.logger.Info("GetTrending completed successfully",
		zap.String("endpoint", endpoint),
		zap.String("media_type", mediaType),
		zap.String("time_window", timeWindow),
		zap.Int("status_code", resp.StatusCode()),
		zap.Duration("response_time", responseTime),
		zap.Int("result_count", len(trendingResp.Results)),
		zap.Int("total_results", trendingResp.TotalResults),
	)

	return &trendingResp, nil
}
