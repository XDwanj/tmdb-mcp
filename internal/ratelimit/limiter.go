// Package ratelimit provides rate limiting functionality for TMDB API requests.
// It uses the Token Bucket algorithm via golang.org/x/time/rate to ensure
// requests respect TMDB's rate limits (default: 40 requests per 10 seconds).
package ratelimit

import (
	"context"
	"time"

	"github.com/XDwanj/tmdb-mcp/internal/config"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// Limiter wraps golang.org/x/time/rate.Limiter to control TMDB API request rate
type Limiter struct {
	rateLimiter *rate.Limiter
	logger      *zap.Logger
	rateLimit   int
}

// NewLimiter creates a new rate limiter with the specified configuration
func NewLimiter(cfg config.TMDBConfig, logger *zap.Logger) *Limiter {
	// Calculate rate: requests per second
	// TMDB limit: 40 requests / 10 seconds = 4 req/s
	interval := 10 * time.Second
	perRequestInterval := interval / time.Duration(cfg.RateLimit)

	// Create rate limiter with Token Bucket algorithm
	// Burst = cfg.RateLimit allows initial burst of requests
	rateLimiter := rate.NewLimiter(rate.Every(perRequestInterval), cfg.RateLimit)

	logger.Info("Rate Limiter initialized",
		zap.Int("rate_limit", cfg.RateLimit),
		zap.Duration("per_request_interval", perRequestInterval),
		zap.Int("burst", cfg.RateLimit),
		zap.String("component", "rate_limiter"),
	)

	return &Limiter{
		rateLimiter: rateLimiter,
		logger:      logger,
		rateLimit:   cfg.RateLimit,
	}
}

// Wait blocks until a token is available or context is cancelled
func (l *Limiter) Wait(ctx context.Context) error {
	start := time.Now()

	// Wait for token from rate limiter (blocks if no tokens available)
	err := l.rateLimiter.Wait(ctx)

	// Calculate wait time
	elapsed := time.Since(start)

	if err != nil {
		// Log warning when wait is cancelled (e.g., context.Canceled)
		l.logger.Warn("Rate limiter wait cancelled",
			zap.Duration("wait_duration", elapsed),
			zap.Error(err),
			zap.String("component", "rate_limiter"),
		)
		return err
	}

	// Log only if we actually had to wait (avoid log noise)
	// Use 1ms threshold to filter out instant returns
	if elapsed > time.Millisecond {
		l.logger.Debug("Rate limiter wait completed",
			zap.Duration("wait_duration", elapsed),
			zap.Int("rate_limit", l.rateLimit),
			zap.String("component", "rate_limiter"),
		)
	}

	return nil
}
