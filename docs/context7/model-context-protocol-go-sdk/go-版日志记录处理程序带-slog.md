# Go 版日志记录处理程序，带 slog

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

为 MCP 定义了一个 slog.Handler，允许服务器作者使用标准的 slog 包进行日志记录。它包括日志记录器名称和速率限制的选项。

```Go
package mcp

import (
	"context"
	"time"
	"log/slog"
)

// A LoggingHandler is a [slog.Handler] for MCP.
type LoggingHandler struct {}

// LoggingHandlerOptions are options for a LoggingHandler.
type LoggingHandlerOptions struct {
	// The value for the "logger" field of logging notifications.
	LoggerName string
	// Limits the rate at which log messages are sent.
	// If zero, there is no rate limiting.
	MinInterval time.Duration
}

// NewLoggingHandler creates a [LoggingHandler] that logs to the given [ServerSession] using a
// [slog.JSONHandler].
func NewLoggingHandler(ss *ServerSession, opts *LoggingHandlerOptions) *LoggingHandler
```

--------------------------------
