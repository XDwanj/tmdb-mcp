package tmdb

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// Error type constants for categorizing TMDB API errors
const (
	ErrorTypeAuth      = "authentication"
	ErrorTypeNotFound  = "not_found"
	ErrorTypeRateLimit = "rate_limit"
	ErrorTypeServer    = "server_error"
	ErrorTypeNetwork   = "network_error"
	ErrorTypeParsing   = "parsing_error"
	ErrorTypeUnknown   = "unknown"
)

// TMDBError represents a TMDB API error response
type TMDBError struct {
	StatusCode     int    `json:"status_code"`    // TMDB API 返回的状态码
	StatusMessage  string `json:"status_message"` // TMDB API 返回的错误消息
	HTTPStatusCode int    `json:"-"`              // HTTP 状态码（可能不同于 TMDB 状态码）
	ErrorType      string `json:"-"`              // 错误类型分类
	Retryable      bool   `json:"-"`              // 是否可重试
}

// Error implements the error interface
func (e *TMDBError) Error() string {
	if e.ErrorType != "" && e.ErrorType != ErrorTypeUnknown {
		return fmt.Sprintf("TMDB API Error [%s] %d: %s", e.ErrorType, e.StatusCode, e.StatusMessage)
	}
	return fmt.Sprintf("TMDB API Error %d: %s", e.StatusCode, e.StatusMessage)
}

// handleError parses and handles TMDB API error responses
func handleError(resp *resty.Response) error {
	statusCode := resp.StatusCode()

	// 处理特定错误码
	switch statusCode {
	case 401:
		return &TMDBError{
			StatusCode:     401,
			StatusMessage:  "Invalid or missing TMDB API Key",
			HTTPStatusCode: 401,
			ErrorType:      ErrorTypeAuth,
			Retryable:      false,
		}

	case 404:
		return &TMDBError{
			StatusCode:     404,
			StatusMessage:  "Resource not found",
			HTTPStatusCode: 404,
			ErrorType:      ErrorTypeNotFound,
			Retryable:      false,
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
			StatusCode:     429,
			StatusMessage:  message,
			HTTPStatusCode: 429,
			ErrorType:      ErrorTypeRateLimit,
			Retryable:      true,
		}

	case 500, 502, 503:
		// 服务器错误，可以重试
		return &TMDBError{
			StatusCode:     statusCode,
			StatusMessage:  "TMDB API server error",
			HTTPStatusCode: statusCode,
			ErrorType:      ErrorTypeServer,
			Retryable:      true,
		}
	}

	// 尝试解析 TMDB API 错误响应
	var tmdbErr TMDBError
	if err := json.Unmarshal(resp.Body(), &tmdbErr); err == nil && tmdbErr.StatusMessage != "" {
		tmdbErr.HTTPStatusCode = statusCode
		tmdbErr.StatusCode = statusCode
		tmdbErr.ErrorType = ErrorTypeUnknown
		tmdbErr.Retryable = false
		return &tmdbErr
	}

	// JSON 解析失败或响应体为空
	if len(resp.Body()) > 0 {
		// 有响应体但解析失败
		return &TMDBError{
			StatusCode:     statusCode,
			StatusMessage:  "Failed to parse TMDB API response",
			HTTPStatusCode: statusCode,
			ErrorType:      ErrorTypeParsing,
			Retryable:      false,
		}
	}

	// 默认错误消息
	return &TMDBError{
		StatusCode:     statusCode,
		StatusMessage:  fmt.Sprintf("HTTP %d: %s", statusCode, resp.Status()),
		HTTPStatusCode: statusCode,
		ErrorType:      ErrorTypeUnknown,
		Retryable:      false,
	}
}
