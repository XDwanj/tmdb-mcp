package mcp

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// LoggingMiddleware creates a middleware that logs all MCP method calls
// It records:
// - Method name and session ID
// - Tool name and arguments (for CallToolRequest)
// - Execution duration
// - Success/failure status
// - Result information (for CallToolResult)
func LoggingMiddleware(logger *zap.Logger) mcp.Middleware {
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(
			ctx context.Context,
			method string,
			req mcp.Request,
		) (mcp.Result, error) {
			// 1. Record request started
			logger.Info("MCP method started",
				zap.String("method", method),
				zap.String("session_id", req.GetSession().ID()),
				zap.Bool("has_params", req.GetParams() != nil))

			// 2. For CallToolRequest, log additional tool information
			if ctr, ok := req.(*mcp.CallToolRequest); ok {
				logger.Info("Calling tool",
					zap.String("name", ctr.Params.Name),
					zap.Any("args", ctr.Params.Arguments))
			}

			// 3. Execute actual method and measure duration
			start := time.Now()
			result, err := next(ctx, method, req)
			duration := time.Since(start)

			// 4. Record result based on error status
			if err != nil {
				logger.Error("MCP method failed",
					zap.String("method", method),
					zap.String("session_id", req.GetSession().ID()),
					zap.Duration("duration", duration),
					zap.Error(err))
			} else {
				logger.Info("MCP method completed",
					zap.String("method", method),
					zap.String("session_id", req.GetSession().ID()),
					zap.Duration("duration", duration),
					zap.Bool("has_result", result != nil))

				// 5. For CallToolResult, log additional result information
				if ctr, ok := result.(*mcp.CallToolResult); ok {
					logger.Info("Tool result",
						zap.Bool("isError", ctr.IsError),
						zap.Int("content_count", len(ctr.Content)))
				}
			}

			return result, err
		}
	}
}
