package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/XDwanj/tmdb-mcp/internal/config"
	mcpServer "github.com/XDwanj/tmdb-mcp/internal/mcp"
	"go.uber.org/zap"
)

// HTTPServer represents an HTTP server for SSE connections
type HTTPServer struct {
	server            *http.Server
	mcpServer         *mcpServer.Server
	config            *config.Config
	logger            *zap.Logger
	mux               *http.ServeMux
	activeConnections atomic.Int32 // Thread-safe counter for active SSE connections
}

// NewHTTPServer creates a new HTTP server instance
func NewHTTPServer(cfg *config.Config, mcpSrv *mcpServer.Server, logger *zap.Logger) *HTTPServer {
	mux := http.NewServeMux()

	httpServer := &HTTPServer{
		mcpServer: mcpSrv,
		config:    cfg,
		logger:    logger,
		mux:       mux,
	}

	// Register health check endpoint (public - no authentication required)
	mux.HandleFunc("/health", httpServer.healthHandler())

	// Register SSE endpoint with authentication (Story 4.4)
	// Get SSE handler from MCP server
	sseHandler := mcpSrv.GetSSEHandler()

	// Apply middleware chain: Auth → ConnectionTracking → SSEHandler
	protectedSSEHandler := AuthMiddleware(cfg.Server.SSE.Token, logger)(sseHandler)
	trackedSSEHandler := ConnectionTrackingMiddleware(httpServer)(protectedSSEHandler)

	// Register the protected SSE endpoint with connection tracking
	mux.Handle("/mcp/sse", trackedSSEHandler)

	// Apply middleware chain (recovery first, then logging)
	// Note: Auth middleware is applied per-endpoint, not globally
	handler := RecoveryMiddleware(logger)(LoggingMiddleware(logger)(mux))

	// Create HTTP server with middleware-wrapped handler
	// Note: WriteTimeout is set to 0 for SSE (long-lived connections)
	httpServer.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.SSE.Host, cfg.Server.SSE.Port),
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 0, // No write timeout for SSE (long-lived connections)
		IdleTimeout:  120 * time.Second,
	}

	return httpServer
}

// healthHandler returns the health check handler
func (s *HTTPServer) healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"status":    "ok",
			"version":   "1.0.0",
			"mode":      "sse",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			s.logger.Error("Failed to encode health response",
				zap.Error(err))
		}
	}
}

// Start starts the HTTP server
func (s *HTTPServer) Start() error {
	s.logger.Info("Starting HTTP server",
		zap.String("addr", s.server.Addr),
		zap.String("mode", "sse"),
	)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server failed to start: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the HTTP server
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server",
		zap.Duration("timeout", 10*time.Second),
	)

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP server shutdown failed: %w", err)
	}

	s.logger.Info("HTTP server shutdown complete")
	return nil
}
