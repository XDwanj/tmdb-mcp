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
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

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

// TestClient_CallCounter tests the API call counter functionality
func TestClient_CallCounter(t *testing.T) {
	tests := []struct {
		name      string
		callCount int
		wantCount uint64
	}{
		{
			name:      "zero_calls",
			callCount: 0,
			wantCount: 0,
		},
		{
			name:      "single_call",
			callCount: 1,
			wantCount: 1,
		},
		{
			name:      "multiple_calls",
			callCount: 5,
			wantCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试服务器
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"results":[]}`))
			}))
			defer server.Close()

			// 创建客户端
			logger := zap.NewNop()
			cfg := config.TMDBConfig{
				APIKey:    "test_api_key",
				Language:  "en-US",
				RateLimit: 40,
			}
			client := NewClient(cfg, logger)
			client.httpClient.SetBaseURL(server.URL)

			// 验证初始计数器为 0
			initialCount := client.GetCallCount()
			assert.Equal(t, uint64(0), initialCount, "Initial counter should be 0")

			// 执行指定次数的 API 调用
			ctx := context.Background()
			for i := 0; i < tt.callCount; i++ {
				_, err := client.httpClient.R().SetContext(ctx).Get("/test")
				require.NoError(t, err)
			}

			// 验证计数器值
			finalCount := client.GetCallCount()
			assert.Equal(t, tt.wantCount, finalCount, "Counter should match call count")
		})
	}
}

// TestClient_CallCounter_Concurrent tests counter thread safety
func TestClient_CallCounter_Concurrent(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":[]}`))
	}))
	defer server.Close()

	// 创建客户端
	logger := zap.NewNop()
	cfg := config.TMDBConfig{
		APIKey:    "test_api_key",
		Language:  "en-US",
		RateLimit: 40,
	}
	client := NewClient(cfg, logger)
	client.httpClient.SetBaseURL(server.URL)

	// 并发执行 API 调用
	const goroutines = 10
	const callsPerGoroutine = 5
	const expectedTotal = goroutines * callsPerGoroutine

	done := make(chan bool, goroutines)
	ctx := context.Background()

	for i := 0; i < goroutines; i++ {
		go func() {
			for j := 0; j < callsPerGoroutine; j++ {
				client.httpClient.R().SetContext(ctx).Get("/test")
			}
			done <- true
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < goroutines; i++ {
		<-done
	}

	// 验证最终计数正确
	finalCount := client.GetCallCount()
	assert.Equal(t, uint64(expectedTotal), finalCount, "Counter should be thread-safe and match total calls")
}

// TestClient_PerformanceThreshold tests performance threshold alerting
func TestClient_PerformanceThreshold(t *testing.T) {
	tests := []struct {
		name          string
		delay         time.Duration
		expectWarning bool
	}{
		{
			name:          "fast_response_no_warning",
			delay:         100 * time.Millisecond,
			expectWarning: false,
		},
		{
			name:          "slow_response_with_warning",
			delay:         1500 * time.Millisecond,
			expectWarning: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建带延迟的测试服务器
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(tt.delay)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"results":[]}`))
			}))
			defer server.Close()

			// 创建 observer logger 捕获日志 (使用 DebugLevel 捕获所有级别的日志)
			core, logs := observer.New(zapcore.DebugLevel)
			logger := zap.New(core)

			// 创建客户端
			cfg := config.TMDBConfig{
				APIKey:    "test_api_key",
				Language:  "en-US",
				RateLimit: 40,
			}
			client := NewClient(cfg, logger)
			client.httpClient.SetBaseURL(server.URL)

			// 执行 API 调用
			ctx := context.Background()
			_, err := client.httpClient.R().SetContext(ctx).Get("/test")
			require.NoError(t, err)

			// 验证日志
			if tt.expectWarning {
				// 应该有 WARN 日志
				warnLogs := logs.FilterMessage("TMDB API request exceeded performance threshold").All()
				assert.NotEmpty(t, warnLogs, "Should log warning for slow response")

				if len(warnLogs) > 0 {
					// 验证日志级别
					assert.Equal(t, zapcore.WarnLevel, warnLogs[0].Level, "Should be WARN level")

					// 验证日志字段
					fields := warnLogs[0].ContextMap()
					assert.Contains(t, fields, "response_time", "Should include response_time field")
					assert.Contains(t, fields, "threshold", "Should include threshold field")
					assert.Contains(t, fields, "method", "Should include method field")
					assert.Contains(t, fields, "url", "Should include url field")
					assert.Contains(t, fields, "status_code", "Should include status_code field")
				}
			} else {
				// 不应该有 WARN 日志
				warnLogs := logs.FilterMessage("TMDB API request exceeded performance threshold").All()
				assert.Empty(t, warnLogs, "Should not log warning for fast response")
			}
		})
	}
}
