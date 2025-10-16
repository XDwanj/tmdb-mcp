package tmdb

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/XDwanj/tmdb-mcp/internal/config"
	"github.com/XDwanj/tmdb-mcp/internal/ratelimit"
	"github.com/XDwanj/tmdb-mcp/pkg/version"
)

// contextKey is used for context values (避免 key 冲突)
type contextKey string

const (
	// startTimeKey is the context key for request start time
	startTimeKey contextKey = "request_start_time"
)

const (
	// baseURL is the TMDB API v3 base URL
	baseURL = "https://api.themoviedb.org/3"

	// defaultTimeout is the default HTTP request timeout
	defaultTimeout = 10 * time.Second

	// performanceThreshold is the response time threshold for performance alerts
	performanceThreshold = 1 * time.Second
)

// Client is the TMDB API client
type Client struct {
	httpClient  *resty.Client
	apiKey      string
	language    string
	logger      *zap.Logger
	rateLimiter *ratelimit.Limiter
	callCounter *uint64 // API 调用计数器(指针以支持 atomic 操作)
}

// NewClient creates a new TMDB API client with configured Resty client
func NewClient(cfg config.TMDBConfig, logger *zap.Logger) *Client {
	// 构建 User-Agent
	userAgent := fmt.Sprintf("tmdb-mcp/%s", version.Version)

	// Create rate limiter first (需要在 middleware 中使用)
	rateLimiter := ratelimit.NewLimiter(cfg, logger)

	// 初始化 API 调用计数器(在创建 httpClient 之前,以便在 middleware 中引用)
	var counter uint64 = 0

	// 创建 Resty 客户端
	httpClient := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(defaultTimeout).
		SetHeader("User-Agent", userAgent).
		OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
			// 1. 统一处理 rate limiting (阻塞等待)
			if err := rateLimiter.Wait(req.Context()); err != nil {
				logger.Error("rate limit wait failed", zap.Error(err))
				return fmt.Errorf("rate limit wait failed: %w", err)
			}

			// 2. 自动添加 API Key
			req.SetQueryParam("api_key", cfg.APIKey)

			// 3. language 参数仅在请求中未显式设置时使用配置默认值
			if req.QueryParam.Get("language") == "" && cfg.Language != "" {
				req.SetQueryParam("language", cfg.Language)
			}
			return nil
		}).
		// 统一处理请求开始日志和计时
		OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
			// 记录请求开始时间(存储在 context 中)
			ctx := context.WithValue(req.Context(), startTimeKey, time.Now())
			req.SetContext(ctx)

			// 记录请求开始日志
			logger.Debug("Starting TMDB API request",
				zap.String("method", req.Method),
				zap.String("url", req.URL),
			)
			return nil
		}).
		// 统一处理成功响应日志
		OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
			// 递增 API 调用计数器(线程安全)
			atomic.AddUint64(&counter, 1)

			// 从 context 获取开始时间
			responseTime := time.Duration(0)
			if startTime, ok := resp.Request.Context().Value(startTimeKey).(time.Time); ok {
				responseTime = time.Since(startTime)
			}

			if resp.IsSuccess() {
				logger.Info("TMDB API request succeeded",
					zap.String("method", resp.Request.Method),
					zap.String("url", resp.Request.URL),
					zap.Int("status_code", resp.StatusCode()),
					zap.Duration("response_time", responseTime),
				)

				// 性能阈值告警
				if responseTime > performanceThreshold {
					logger.Warn("TMDB API request exceeded performance threshold",
						zap.String("method", resp.Request.Method),
						zap.String("url", resp.Request.URL),
						zap.Int("status_code", resp.StatusCode()),
						zap.Duration("response_time", responseTime),
						zap.Duration("threshold", performanceThreshold),
					)
				}
			}
			return nil
		}).
		// 统一处理错误响应日志
		OnError(func(req *resty.Request, err error) {
			// 从 context 获取开始时间
			responseTime := time.Duration(0)
			if startTime, ok := req.Context().Value(startTimeKey).(time.Time); ok {
				responseTime = time.Since(startTime)
			}

			logger.Error("TMDB API request failed",
				zap.String("method", req.Method),
				zap.String("url", req.URL),
				zap.Duration("response_time", responseTime),
				zap.Error(err),
			)
		}).
		// 配置重试机制
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		// 自定义重试条件：仅对 429/500/502/503 重试
		AddRetryCondition(func(res *resty.Response, err error) bool {
			// 网络错误或 nil response，由 Resty 默认条件处理
			if err != nil || res == nil {
				return false
			}

			// 仅对特定状态码重试
			statusCode := res.StatusCode()
			return statusCode == 429 || statusCode == 500 || statusCode == 502 || statusCode == 503
		}).
		// 添加重试钩子：记录重试日志
		AddRetryHook(func(res *resty.Response, err error) {
			statusCode := 0
			endpoint := ""
			if res != nil {
				statusCode = res.StatusCode()
				endpoint = res.Request.URL
			}

			logger.Warn("Retrying TMDB API request",
				zap.String("endpoint", endpoint),
				zap.Int("status_code", statusCode),
				zap.Error(err),
			)
		})

	logger.Debug("TMDB client initialized",
		zap.String("base_url", baseURL),
		zap.String("language", cfg.Language),
		zap.String("user_agent", userAgent),
		zap.Int("retry_count", 3),
	)

	logger.Debug("Rate Limiter integrated to TMDB Client",
		zap.String("component", "tmdb_client"),
	)

	return &Client{
		httpClient:  httpClient,
		apiKey:      cfg.APIKey,
		language:    cfg.Language,
		logger:      logger,
		rateLimiter: rateLimiter,
		callCounter: &counter,
	}
}

// GetCallCount returns the current API call count (thread-safe)
func (c *Client) GetCallCount() uint64 {
	return atomic.LoadUint64(c.callCounter)
}

// Ping tests the TMDB API Key validity by calling the /configuration endpoint
func (c *Client) Ping(ctx context.Context) error {
	// Rate limiting is handled by OnBeforeRequest middleware
	resp, err := c.httpClient.R().
		SetContext(ctx).
		Get("/configuration")

	if err != nil {
		c.logger.Error("failed to ping TMDB API", zap.Error(err))
		return fmt.Errorf("failed to ping TMDB API: %w", err)
	}

	// 检查错误响应
	if resp.IsError() {
		return handleError(resp)
	}

	c.logger.Info("TMDB API Key validation successful")
	return nil
}
