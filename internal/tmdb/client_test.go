package tmdb

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/XDwanj/tmdb-mcp/internal/config"
)

// TestNewClient tests the Client constructor
func TestNewClient(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    "test_api_key",
		Language:  "en-US",
		RateLimit: 40,
	}

	client := NewClient(cfg, logger)

	assert.NotNil(t, client)
	assert.NotNil(t, client.httpClient)
	assert.Equal(t, "test_api_key", client.apiKey)
	assert.Equal(t, "en-US", client.language)
	assert.NotNil(t, client.logger)
	assert.NotNil(t, client.rateLimiter)
}

// TestClient_Ping tests the Ping method
func TestClient_Ping(t *testing.T) {
	tests := []struct {
		name           string
		responseStatus int
		responseBody   string
		wantErr        bool
		errMessage     string
	}{
		{
			name:           "success",
			responseStatus: http.StatusOK,
			responseBody:   `{"images":{"base_url":"http://image.tmdb.org/t/p/"}}`,
			wantErr:        false,
		},
		{
			name:           "unauthorized",
			responseStatus: http.StatusUnauthorized,
			responseBody:   `{"status_code":7,"status_message":"Invalid API key"}`,
			wantErr:        true,
			errMessage:     "Invalid or missing TMDB API Key",
		},
		{
			name:           "not_found",
			responseStatus: http.StatusNotFound,
			responseBody:   `{"status_code":34,"status_message":"The resource you requested could not be found"}`,
			wantErr:        true,
			errMessage:     "Resource not found",
		},
		{
			name:           "rate_limit",
			responseStatus: http.StatusTooManyRequests,
			responseBody:   `{"status_code":25,"status_message":"Your request count is over the allowed limit"}`,
			wantErr:        true,
			errMessage:     "Rate limit exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试服务器
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 验证 API Key 查询参数
				apiKey := r.URL.Query().Get("api_key")
				assert.Equal(t, "test_api_key", apiKey, "API key should be added to request")

				// 验证 language 查询参数
				language := r.URL.Query().Get("language")
				assert.Equal(t, "zh-CN", language, "language should be added to request")

				// 验证 User-Agent header
				userAgent := r.Header.Get("User-Agent")
				assert.Contains(t, userAgent, "tmdb-mcp", "User-Agent should contain tmdb-mcp")

				// 返回模拟响应
				w.WriteHeader(tt.responseStatus)
				if tt.responseStatus == http.StatusTooManyRequests {
					w.Header().Set("Retry-After", "30")
				}
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// 创建测试客户端
			logger := zap.NewNop()
			cfg := config.TMDBConfig{
				APIKey:    "test_api_key",
				Language:  "zh-CN",
				RateLimit: 40,
			}
			client := NewClient(cfg, logger)

			// 修改 base URL 指向测试服务器
			client.httpClient.SetBaseURL(server.URL)

			// 测试 Ping
			ctx := context.Background()
			err := client.Ping(ctx)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMessage)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestClient_Ping_Timeout tests Ping method with context timeout
func TestClient_Ping_Timeout(t *testing.T) {
	// 创建慢响应的测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond) // 延迟响应
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    "test_api_key",
		Language:  "en-US",
		RateLimit: 40,
	}
	client := NewClient(cfg, logger)
	client.httpClient.SetBaseURL(server.URL)

	// 使用 100ms 超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := client.Ping(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

// TestClient_Resty_Configuration tests that Resty client is properly configured
func TestClient_Resty_Configuration(t *testing.T) {
	// 创建测试服务器
	var capturedRequest *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedRequest = r
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    "test_api_key_12345",
		Language:  "ja-JP",
		RateLimit: 40,
	}
	client := NewClient(cfg, logger)
	client.httpClient.SetBaseURL(server.URL)

	// 发送请求
	ctx := context.Background()
	err := client.Ping(ctx)
	require.NoError(t, err)

	// 验证请求配置
	require.NotNil(t, capturedRequest)

	// 验证 API Key 自动添加
	apiKey := capturedRequest.URL.Query().Get("api_key")
	assert.Equal(t, "test_api_key_12345", apiKey)

	// 验证 language 自动添加
	language := capturedRequest.URL.Query().Get("language")
	assert.Equal(t, "ja-JP", language)

	// 验证 User-Agent
	userAgent := capturedRequest.Header.Get("User-Agent")
	assert.Contains(t, userAgent, "tmdb-mcp")
}
