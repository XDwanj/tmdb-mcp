package tmdb

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/XDwanj/tmdb-mcp/internal/config"
)

// TestHandleError tests the handleError function
func TestHandleError(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		retryAfter     string
		expectedErr    string
		expectedStatus int
	}{
		{
			name:           "401_unauthorized",
			statusCode:     401,
			responseBody:   `{"status_code":7,"status_message":"Invalid API key"}`,
			expectedErr:    "Invalid or missing TMDB API Key",
			expectedStatus: 401,
		},
		{
			name:           "404_not_found",
			statusCode:     404,
			responseBody:   `{"status_code":34,"status_message":"Not found"}`,
			expectedErr:    "Resource not found",
			expectedStatus: 404,
		},
		{
			name:           "429_rate_limit_with_retry",
			statusCode:     429,
			responseBody:   `{"status_code":25,"status_message":"Rate limit"}`,
			retryAfter:     "60",
			expectedErr:    "retry after 60 seconds",
			expectedStatus: 429,
		},
		{
			name:           "429_rate_limit_without_retry",
			statusCode:     429,
			responseBody:   `{"status_code":25,"status_message":"Rate limit"}`,
			expectedErr:    "Rate limit exceeded",
			expectedStatus: 429,
		},
		{
			name:           "500_server_error",
			statusCode:     500,
			responseBody:   `{"status_code":11,"status_message":"Internal error"}`,
			expectedErr:    "TMDB API server error",
			expectedStatus: 500,
		},
		{
			name:           "generic_error",
			statusCode:     503,
			responseBody:   `Service Unavailable`,
			expectedErr:    "TMDB API server error",
			expectedStatus: 503,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试服务器
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.retryAfter != "" {
					w.Header().Set("Retry-After", tt.retryAfter)
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// 创建客户端并发送请求
			logger := zap.NewNop()
			cfg := config.TMDBConfig{
				APIKey:    "test_api_key",
				Language:  "en-US",
				RateLimit: 40,
			}
			client := NewClient(cfg, logger)
			client.httpClient.SetBaseURL(server.URL)

			resp, err := client.httpClient.R().Get("/test")
			require.NoError(t, err)

			// 测试 handleError
			handledErr := handleError(resp)
			require.Error(t, handledErr)

			tmdbErr, ok := handledErr.(*TMDBError)
			require.True(t, ok, "Error should be of type *TMDBError")
			assert.Equal(t, tt.expectedStatus, tmdbErr.StatusCode)
			assert.Contains(t, tmdbErr.Error(), tt.expectedErr)
		})
	}
}

// TestTMDBError_Error tests the TMDBError Error method
func TestTMDBError_Error(t *testing.T) {
	err := &TMDBError{
		StatusCode:    401,
		StatusMessage: "Invalid API key",
	}

	expected := "TMDB API Error 401: Invalid API key"
	assert.Equal(t, expected, err.Error())
}
