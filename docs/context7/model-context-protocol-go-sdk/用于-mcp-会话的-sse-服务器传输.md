# 用于 MCP 会话的 SSE 服务器传输

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

表示通过挂起 GET 请求建立的逻辑服务器发送事件 (SSE) 会话。此传输处理到会话端点的传入请求并促进连接过程，允许事件流式传输。

```Go
package mcp

import (
	"context"
	"net/http"
)

// A SSEServerTransport is a logical SSE session created through a hanging GET
// request.
type SSEServerTransport struct {
    Endpoint string
    Response http.ResponseWriter
}

// ServeHTTP handles POST requests to the transport endpoint.
func (*SSEServerTransport) ServeHTTP(w http.ResponseWriter, req *http.Request)

// Connect sends the 'endpoint' event to the client.
// See [SSEServerTransport] for more details on the [Connection] implementation.
func (*SSEServerTransport) Connect(context.Context) (Connection, error)
```

--------------------------------
