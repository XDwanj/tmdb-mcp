# 用于 MCP 服务器的 Streamable HTTP 处理程序

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

实现了用于提供 Streamable MCP 会话的 `http.Handler`。与 SSE 处理程序类似，它接受用于创建新会话服务器的回调，并管理会话生命周期，包括关闭活动会话。

```Go
package mcp

import (
	"net/http"
)

// The StreamableHTTPHandler interface is symmetrical to the SSEHTTPHandler.
type StreamableHTTPHandler struct { /* unexported fields */ }
func NewStreamableHTTPHandler(getServer func(request *http.Request) *Server) *StreamableHTTPHandler
func (*StreamableHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request)
func (*StreamableHTTPHandler) Close() error
```

--------------------------------
