# 用于 MCP 服务器的 SSE HTTP 处理程序

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

提供了一个 `http.Handler`，用于处理基于服务器发送事件 (SSE) 的 MCP 会话。它接受一个回调函数来绑定新会话的服务器，允许每个连接使用无状态或有状态的服务器实例。处理程序管理会话的生命周期，包括关闭活动会话。

```Go
package mcp

import (
	"net/http"
)

// SSEHTTPHandler is an http.Handler that serves SSE-based MCP sessions as defined by
// the 2024-11-05 version of the MCP protocol.
type SSEHTTPHandler struct { /* unexported fields */ }

// NewSSEHTTPHandler returns a new [SSEHTTPHandler] that is ready to serve HTTP.
//
// The getServer function is used to bind created servers for new sessions. It
// is OK for getServer to return the same server multiple times.
func NewSSEHTTPHandler(getServer func(request *http.Request) *Server) *SSEHTTPHandler

func (*SSEHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request)

// Close prevents the SSEHTTPHandler from accepting new sessions, closes active
// sessions, and awaits their graceful termination.
func (*SSEHTTPHandler) Close() error
```

--------------------------------
