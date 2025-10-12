package tmdb

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// TMDBError represents a TMDB API error response
type TMDBError struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

// Error implements the error interface
func (e *TMDBError) Error() string {
	return fmt.Sprintf("TMDB API Error %d: %s", e.StatusCode, e.StatusMessage)
}

// handleError parses and handles TMDB API error responses
func handleError(resp *resty.Response) error {
	statusCode := resp.StatusCode()

	// 处理特定错误码
	switch statusCode {
	case 401:
		return &TMDBError{
			StatusCode:    401,
			StatusMessage: "Invalid or missing TMDB API Key",
		}

	case 404:
		return &TMDBError{
			StatusCode:    404,
			StatusMessage: "Resource not found",
		}

	case 429:
		// 解析 Retry-After header
		retryAfter := resp.Header().Get("Retry-After")
		message := "Rate limit exceeded"
		if retryAfter != "" {
			if seconds, err := strconv.Atoi(retryAfter); err == nil {
				message = fmt.Sprintf("Rate limit exceeded, retry after %d seconds", seconds)
			} else {
				message = fmt.Sprintf("Rate limit exceeded, retry after %s", retryAfter)
			}
		}
		return &TMDBError{
			StatusCode:    429,
			StatusMessage: message,
		}
	}

	// 尝试解析 TMDB API 错误响应
	var tmdbErr TMDBError
	if err := json.Unmarshal(resp.Body(), &tmdbErr); err == nil && tmdbErr.StatusMessage != "" {
		tmdbErr.StatusCode = statusCode
		return &tmdbErr
	}

	// 默认错误消息
	return &TMDBError{
		StatusCode:    statusCode,
		StatusMessage: fmt.Sprintf("HTTP %d: %s", statusCode, resp.Status()),
	}
}
