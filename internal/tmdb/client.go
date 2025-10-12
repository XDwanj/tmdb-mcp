package tmdb

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/XDwanj/tmdb-mcp/internal/config"
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
	httpClient *resty.Client
	apiKey     string
	language   string
	logger     *zap.Logger
}

// NewClient creates a new TMDB API client with configured Resty client
func NewClient(cfg config.TMDBConfig, logger *zap.Logger) *Client {
	// 构建 User-Agent
	userAgent := fmt.Sprintf("tmdb-mcp/%s", version.Version)

	// 创建 Resty 客户端
	httpClient := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(defaultTimeout).
		SetHeader("User-Agent", userAgent).
		OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
			// 自动添加 API Key 和 language 查询参数
			req.SetQueryParam("api_key", cfg.APIKey)
			if cfg.Language != "" {
				req.SetQueryParam("language", cfg.Language)
			}
			return nil
		})

	logger.Debug("TMDB client initialized",
		zap.String("base_url", baseURL),
		zap.String("language", cfg.Language),
		zap.String("user_agent", userAgent),
	)

	return &Client{
		httpClient: httpClient,
		apiKey:     cfg.APIKey,
		language:   cfg.Language,
		logger:     logger,
	}
}

// Ping tests the TMDB API Key validity by calling the /configuration endpoint
func (c *Client) Ping(ctx context.Context) error {
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
