# 用于 MCP 的 SSE 客户端传输

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

使客户端能够连接到基于 SSE 的 MCP 端点。它允许指定目标端点和可选的 `http.Client` 来发出请求。传输处理到服务器的连接建立。

```Go
package mcp

import (
	"context"
	"net/http"
)

type SSEClientTransport struct {
	// Endpoint is the SSE endpoint to connect to.
	Endpoint string
	// HTTPClient is the client to use for making HTTP requests. If nil,
	// http.DefaultClient is used.
	HTTPClient *http.Client
}

// Connect connects through the client endpoint.
func (*SSEClientTransport) Connect(ctx context.Context) (Connection, error)
```

--------------------------------
