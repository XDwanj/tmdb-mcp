package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"

	"github.com/XDwanj/tmdb-mcp/internal/config"
	"github.com/XDwanj/tmdb-mcp/internal/logger"
	"github.com/XDwanj/tmdb-mcp/internal/mcp"
	"github.com/XDwanj/tmdb-mcp/internal/server/middleware"
	"github.com/XDwanj/tmdb-mcp/internal/tmdb"
	"github.com/XDwanj/tmdb-mcp/pkg/version"
)

// healthHandler returns server health status
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid configuration: %v\n", err)
		os.Exit(1)
	}

	// 初始化 logger
	log, err := logger.InitLogger(cfg.Logging)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	// 记录程序启动
	log.Info("TMDB MCP Service starting",
		zap.String("version", version.Version),
		zap.String("mode", cfg.Server.Mode),
	)

	// 记录配置加载成功
	log.Info("Configuration loaded successfully",
		zap.String("language", cfg.TMDB.Language),
		zap.Int("rate_limit", cfg.TMDB.RateLimit),
		zap.String("logging_level", cfg.Logging.Level),
	)

	// 如果 SSE 模式启用，显示 Token 信息
	if cfg.Server.Mode == "sse" || cfg.Server.Mode == "both" {
		if cfg.TokenGenerated {
			// Token 是自动生成的，显示完整 token（用户需要复制配置客户端）
			log.Info("SSE Token auto-generated (save this for client configuration)",
				zap.String("token", cfg.Server.SSE.Token),
				zap.String("config_file", "~/.tmdb-mcp/config.yaml"),
			)
		} else {
			// Token 是从配置或环境变量加载的，仅显示前 8 个字符
			tokenPrefix := cfg.Server.SSE.Token
			if len(tokenPrefix) > 8 {
				tokenPrefix = tokenPrefix[:8] + "..."
			}
			log.Info("SSE Token loaded from configuration",
				zap.String("token_prefix", tokenPrefix),
			)
		}
	}

	// 创建 TMDB Client
	tmdbClient := tmdb.NewClient(cfg.TMDB, log)
	log.Info("TMDB Client created")

	// 启动时性能基准测试: 验证 TMDB API Key 有效性并记录响应时间
	log.Info("Running TMDB API baseline check...")
	ctx := context.Background()
	baselineStart := time.Now()
	if err := tmdbClient.Ping(ctx); err != nil {
		baselineTime := time.Since(baselineStart)
		log.Error("TMDB API baseline check failed",
			zap.Error(err),
			zap.Duration("baseline_response_time", baselineTime),
			zap.String("status", "failed"),
		)
		os.Exit(1)
	}
	baselineTime := time.Since(baselineStart)
	log.Info("TMDB API baseline check completed",
		zap.Duration("baseline_response_time", baselineTime),
		zap.String("status", "success"),
	)

	// 创建 MCP Server
	mcpServer := mcp.NewServer(tmdbClient, log)

	// 根据配置模式启动服务
	switch cfg.Server.Mode {
	case "stdio":
		RunStdioModeServer(ctx, mcpServer, log)
	case "sse":
		RunSSEModeServer(ctx, mcpServer, cfg, log)
	case "both":
		RunBothModeServer(ctx, mcpServer, cfg, log)
	default:
		log.Fatal("Invalid server mode", zap.String("mode", cfg.Server.Mode))
	}
}

func RunBothModeServer(ctx context.Context, mcpServer *mcp.Server, cfg *config.Config, log *zap.Logger) {
	// 同时运行 stdio 和 SSE 模式
	log.Info("Starting MCP server in both stdio and SSE modes")
	go func() {
		RunSSEModeServer(ctx, mcpServer, cfg, log)
	}()
	RunStdioModeServer(ctx, mcpServer, log)
}

func RunStdioModeServer(ctx context.Context, mcpServer *mcp.Server, log *zap.Logger) {
	// stdio 模式：通过标准输入输出通信
	log.Info("Starting MCP server in stdio mode")
	transport := &mcpsdk.StdioTransport{}
	if err := mcpServer.Run(ctx, transport); err != nil {
		log.Fatal("MCP server failed", zap.Error(err))
	}
}

func RunSSEModeServer(ctx context.Context, mcpServer *mcp.Server, cfg *config.Config, log *zap.Logger) {
	// SSE 模式：通过 HTTP SSE 通信
	log.Info("Starting MCP server in SSE mode")

	// 设置 SSE 处理器
	handler := mcpServer.GetSSEHandler()
	handler = middleware.AuthMiddleware(cfg.Server.SSE.Token, handler)

	// 设置路由
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.Handle("/mcp/sse", handler)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.SSE.Host, cfg.Server.SSE.Port)
	log.Info("Starting HTTP server", zap.String("addr", addr))
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("HTTP server failed", zap.Error(err))
	}
}
