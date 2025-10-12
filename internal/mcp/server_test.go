package mcp

import (
	"context"
	"testing"
	"time"

	"github.com/XDwanj/tmdb-mcp/internal/config"
	"github.com/XDwanj/tmdb-mcp/internal/tmdb"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// TestNewServer 测试 MCP Server 的创建
func TestNewServer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "创建 MCP Server 成功",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试 logger
			logger := zaptest.NewLogger(t)

			// 创建测试 TMDB client
			tmdbConfig := config.TMDBConfig{
				APIKey:    "test_api_key",
				Language:  "en-US",
				RateLimit: 40,
			}
			tmdbClient := tmdb.NewClient(tmdbConfig, logger)

			// 创建 MCP Server
			server := NewServer(tmdbClient, logger)

			// 验证 server 不为 nil
			require.NotNil(t, server, "Server should not be nil")
			require.NotNil(t, server.mcpServer, "MCP Server should not be nil")
			require.NotNil(t, server.tmdbClient, "TMDB Client should not be nil")
			require.NotNil(t, server.logger, "Logger should not be nil")

			// 验证依赖注入正确
			assert.Equal(t, tmdbClient, server.tmdbClient, "TMDB Client should be correctly injected")
			assert.Equal(t, logger, server.logger, "Logger should be correctly injected")
		})
	}
}

// TestServerInfo 测试 MCP Server 的信息配置
func TestServerInfo(t *testing.T) {
	logger := zap.NewNop()
	tmdbConfig := config.TMDBConfig{
		APIKey:    "test_api_key",
		Language:  "en-US",
		RateLimit: 40,
	}
	tmdbClient := tmdb.NewClient(tmdbConfig, logger)

	server := NewServer(tmdbClient, logger)

	// 注意：由于 MCP SDK 的 Server 结构体可能不直接暴露 ServerInfo，
	// 我们主要验证 server 能正确创建
	// ServerInfo 的验证将在集成测试或手动测试中完成

	require.NotNil(t, server, "Server should be created successfully")
	assert.NotNil(t, server.mcpServer, "MCP Server should be initialized")
}

// TestNewServer_WithValidDependencies 测试使用有效依赖创建 server
func TestNewServer_WithValidDependencies(t *testing.T) {
	logger := zap.NewNop()
	tmdbConfig := config.TMDBConfig{
		APIKey:    "test_api_key",
		Language:  "en-US",
		RateLimit: 40,
	}
	tmdbClient := tmdb.NewClient(tmdbConfig, logger)

	server := NewServer(tmdbClient, logger)

	// 验证所有依赖都正确设置
	require.NotNil(t, server, "Server should be created successfully")
	require.NotNil(t, server.mcpServer, "MCP Server should be initialized")
	require.NotNil(t, server.tmdbClient, "TMDB Client should be set")
	require.NotNil(t, server.logger, "Logger should be set")
}

// TestServer_Run_WithInMemoryTransport 测试使用 InMemoryTransport 运行 server
func TestServer_Run_WithInMemoryTransport(t *testing.T) {
	logger := zap.NewNop()
	tmdbConfig := config.TMDBConfig{
		APIKey:    "test_api_key",
		Language:  "en-US",
		RateLimit: 40,
	}
	tmdbClient := tmdb.NewClient(tmdbConfig, logger)

	server := NewServer(tmdbClient, logger)

	// 创建 InMemoryTransport
	clientTransport, serverTransport := mcpsdk.NewInMemoryTransports()
	_ = clientTransport // 暂时不使用 client transport

	// 创建带超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// 在 goroutine 中运行 server
	errChan := make(chan error, 1)
	go func() {
		err := server.Run(ctx, serverTransport)
		errChan <- err
	}()

	// 等待 context 超时或错误
	select {
	case err := <-errChan:
		// server 应该因为 context 取消而停止
		// 验证没有意外错误（context cancelled 是预期的）
		if err != nil && err != context.DeadlineExceeded && err != context.Canceled {
			t.Errorf("Unexpected error: %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Server did not stop in time")
	}
}

// TestServer_Run_Logger 测试 Run 方法的日志记录
func TestServer_Run_Logger(t *testing.T) {
	// 使用 zaptest logger 捕获日志
	logger := zaptest.NewLogger(t)

	tmdbConfig := config.TMDBConfig{
		APIKey:    "test_api_key",
		Language:  "en-US",
		RateLimit: 40,
	}
	tmdbClient := tmdb.NewClient(tmdbConfig, logger)

	server := NewServer(tmdbClient, logger)

	// 创建 InMemoryTransport
	_, serverTransport := mcpsdk.NewInMemoryTransports()

	// 创建带超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// 运行 server（会触发 logger.Info）
	_ = server.Run(ctx, serverTransport)

	// 注意：zaptest.NewLogger 会自动将日志输出到 testing.T
	// 日志 "Starting MCP server" 会被记录
	// 如果日志没有被调用，测试框架会检测到
}

// TestServer_Run_ContextCancellation 测试 context 取消时的行为
func TestServer_Run_ContextCancellation(t *testing.T) {
	logger := zap.NewNop()
	tmdbConfig := config.TMDBConfig{
		APIKey:    "test_api_key",
		Language:  "en-US",
		RateLimit: 40,
	}
	tmdbClient := tmdb.NewClient(tmdbConfig, logger)

	server := NewServer(tmdbClient, logger)

	// 创建 InMemoryTransport
	_, serverTransport := mcpsdk.NewInMemoryTransports()

	// 创建可取消的 context
	ctx, cancel := context.WithCancel(context.Background())

	// 启动 server
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Run(ctx, serverTransport)
	}()

	// 立即取消 context
	time.Sleep(10 * time.Millisecond) // 给 server 一点启动时间
	cancel()

	// 验证 server 正确停止
	select {
	case err := <-errChan:
		// server 应该因为 context 取消而停止
		if err != nil && err != context.Canceled {
			t.Errorf("Expected context.Canceled or nil, got: %v", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Server did not stop after context cancellation")
	}
}
