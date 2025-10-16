package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/XDwanj/tmdb-mcp/internal/config"
	"github.com/XDwanj/tmdb-mcp/internal/logger"
	"github.com/XDwanj/tmdb-mcp/internal/mcp"
	"github.com/XDwanj/tmdb-mcp/internal/tmdb"
	"github.com/XDwanj/tmdb-mcp/pkg/version"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

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

	// 设置信号处理
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 启动 MCP Server（在 goroutine 中）
	errChan := make(chan error, 1)
	go func() {
		log.Info("Starting MCP Server in stdio mode")
		transport := &mcpsdk.StdioTransport{}
		if err := mcpServer.Run(ctx, transport); err != nil {
			errChan <- err
		}
	}()

	// 等待信号或错误
	select {
	case <-sigChan:
		log.Info("Received shutdown signal")
		cancel()
	case err := <-errChan:
		log.Error("MCP Server failed", zap.Error(err))
		cancel()
	}

	// 优雅退出
	log.Info("TMDB MCP Service shutdown complete")
}
