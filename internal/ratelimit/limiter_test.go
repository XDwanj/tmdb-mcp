package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/XDwanj/tmdb-mcp/internal/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// TestNewLimiter tests Limiter initialization
func TestNewLimiter(t *testing.T) {
	cfg := config.TMDBConfig{
		APIKey:    "test-key",
		Language:  "en-US",
		RateLimit: 40,
	}
	logger := zap.NewNop()

	limiter := NewLimiter(cfg, logger)

	assert.NotNil(t, limiter)
	assert.NotNil(t, limiter.rateLimiter)
	assert.NotNil(t, limiter.logger)
	assert.Equal(t, 40, limiter.rateLimit)
}

// TestLimiter_Wait_FastRequests tests that fast consecutive requests are rate limited
func TestLimiter_Wait_FastRequests(t *testing.T) {
	cfg := config.TMDBConfig{
		APIKey:    "test-key",
		Language:  "en-US",
		RateLimit: 40,
	}
	logger := zap.NewNop()
	limiter := NewLimiter(cfg, logger)
	ctx := context.Background()

	start := time.Now()

	// First 40 requests should complete quickly (burst capacity)
	for i := 0; i < 40; i++ {
		err := limiter.Wait(ctx)
		assert.NoError(t, err)
	}

	// First 40 requests should complete in < 100ms
	elapsed := time.Since(start)
	assert.Less(t, elapsed, 100*time.Millisecond, "First 40 requests should complete quickly using burst capacity")

	// 41st request should be delayed (no more tokens in bucket)
	start41 := time.Now()
	err := limiter.Wait(ctx)
	assert.NoError(t, err)
	elapsed41 := time.Since(start41)

	// Should wait at least 200ms (close to 250ms per request interval)
	assert.Greater(t, elapsed41, 200*time.Millisecond, "41st request should be delayed due to rate limiting")
}

// TestLimiter_Wait_ContextCancellation tests that Wait returns error when context is cancelled
func TestLimiter_Wait_ContextCancellation(t *testing.T) {
	cfg := config.TMDBConfig{
		APIKey:    "test-key",
		Language:  "en-US",
		RateLimit: 40,
	}
	logger := zap.NewNop()
	limiter := NewLimiter(cfg, logger)

	// Exhaust burst capacity
	ctx := context.Background()
	for i := 0; i < 40; i++ {
		err := limiter.Wait(ctx)
		assert.NoError(t, err)
	}

	// Create context with timeout
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// This should timeout and return error
	err := limiter.Wait(ctxTimeout)
	assert.Error(t, err, "Wait should return error when context is cancelled")
	// Note: rate.Limiter wraps context errors, so we just check that an error occurred
}

// TestLimiter_Wait_MultipleSlowRequests tests rate limiting over time
func TestLimiter_Wait_MultipleSlowRequests(t *testing.T) {
	// Use smaller rate limit for faster test
	cfg := config.TMDBConfig{
		APIKey:    "test-key",
		Language:  "en-US",
		RateLimit: 10, // 10 requests per 10 seconds = 1 req/s
	}
	logger := zap.NewNop()
	limiter := NewLimiter(cfg, logger)
	ctx := context.Background()

	start := time.Now()

	// Request 15 times (10 burst + 5 delayed)
	for i := 0; i < 15; i++ {
		err := limiter.Wait(ctx)
		assert.NoError(t, err)
	}

	elapsed := time.Since(start)

	// Should take at least 4 seconds (5 requests * 1s interval)
	// But less than 6 seconds (allowing some tolerance)
	assert.Greater(t, elapsed, 4*time.Second, "15 requests with rate limit 10/10s should take at least 4s")
	assert.Less(t, elapsed, 6*time.Second, "Test should complete within reasonable time")
}

// TestLimiter_Wait_ConcurrentRequests tests thread safety
func TestLimiter_Wait_ConcurrentRequests(t *testing.T) {
	cfg := config.TMDBConfig{
		APIKey:    "test-key",
		Language:  "en-US",
		RateLimit: 10,
	}
	logger := zap.NewNop()
	limiter := NewLimiter(cfg, logger)
	ctx := context.Background()

	// Launch 20 concurrent goroutines
	done := make(chan bool, 20)
	for i := 0; i < 20; i++ {
		go func() {
			err := limiter.Wait(ctx)
			assert.NoError(t, err)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 20; i++ {
		<-done
	}

	// No panic or data race should occur (verified with go test -race)
}

// TestLimiter_WaitObservability tests rate limiter wait observability
func TestLimiter_WaitObservability(t *testing.T) {
	// 创建 observer logger 捕获日志
	core, logs := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)

	// 创建低速率限制器 (1 req/s)
	cfg := config.TMDBConfig{
		APIKey:    "test-key",
		Language:  "en-US",
		RateLimit: 1, // 1 request per 10 seconds
	}
	limiter := NewLimiter(cfg, logger)
	ctx := context.Background()

	// 第一次调用应该立即返回(使用 burst token)
	err := limiter.Wait(ctx)
	assert.NoError(t, err)

	// 第二次调用应该等待,并记录 DEBUG 日志
	err = limiter.Wait(ctx)
	assert.NoError(t, err)

	// 验证日志
	debugLogs := logs.FilterMessage("Rate limiter wait completed").All()
	assert.NotEmpty(t, debugLogs, "Should log DEBUG message for wait")

	if len(debugLogs) > 0 {
		// 验证日志级别
		assert.Equal(t, zapcore.DebugLevel, debugLogs[0].Level, "Should be DEBUG level")

		// 验证日志字段
		fields := debugLogs[0].ContextMap()
		assert.Contains(t, fields, "wait_duration", "Should include wait_duration field")
		assert.Contains(t, fields, "rate_limit", "Should include rate_limit field")
		assert.Contains(t, fields, "component", "Should include component field")

		// 验证 wait_duration > 0
		waitDuration, ok := fields["wait_duration"].(time.Duration)
		assert.True(t, ok, "wait_duration should be time.Duration type")
		assert.Greater(t, waitDuration, time.Duration(0), "wait_duration should be greater than 0")
	}
}

// TestLimiter_WaitCancelled tests that Wait logs warning when cancelled
func TestLimiter_WaitCancelled(t *testing.T) {
	// 创建 observer logger 捕获日志
	core, logs := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)

	// 创建低速率限制器
	cfg := config.TMDBConfig{
		APIKey:    "test-key",
		Language:  "en-US",
		RateLimit: 1, // 1 request per 10 seconds
	}
	limiter := NewLimiter(cfg, logger)

	// 耗尽 burst 容量
	ctx := context.Background()
	err := limiter.Wait(ctx)
	assert.NoError(t, err)

	// 创建会立即取消的 context
	ctxCancel, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Wait 应该被取消并记录 WARN 日志
	err = limiter.Wait(ctxCancel)
	assert.Error(t, err, "Wait should return error when context is cancelled")

	// 验证 WARN 日志
	warnLogs := logs.FilterMessage("Rate limiter wait cancelled").All()
	assert.NotEmpty(t, warnLogs, "Should log WARN message when wait is cancelled")

	if len(warnLogs) > 0 {
		// 验证日志级别
		assert.Equal(t, zapcore.WarnLevel, warnLogs[0].Level, "Should be WARN level")

		// 验证日志字段
		fields := warnLogs[0].ContextMap()
		assert.Contains(t, fields, "wait_duration", "Should include wait_duration field")
		assert.Contains(t, fields, "component", "Should include component field")
	}
}
