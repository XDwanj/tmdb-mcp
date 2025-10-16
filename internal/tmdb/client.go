package tmdb

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/XDwanj/tmdb-mcp/internal/config"
	"github.com/XDwanj/tmdb-mcp/internal/ratelimit"
	"github.com/XDwanj/tmdb-mcp/pkg/version"
)

const (
	// baseURL is the TMDB API v3 base URL
	baseURL = "https://api.themoviedb.org/3"

	// defaultTimeout is the default HTTP request timeout
	defaultTimeout = 10 * time.Second
)

// Client is the TMDB API client
type Client struct {
	httpClient  *resty.Client
	apiKey      string
	language    string
	logger      *zap.Logger
	rateLimiter *ratelimit.Limiter
}

// NewClient creates a new TMDB API client with configured Resty client
func NewClient(cfg config.TMDBConfig, logger *zap.Logger) *Client {
	// 构建 User-Agent
	userAgent := fmt.Sprintf("tmdb-mcp/%s", version.Version)

	// Create rate limiter first (需要在 middleware 中使用)
	rateLimiter := ratelimit.NewLimiter(cfg, logger)

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
	}
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
