# 用于 MCP 的 Streamable 客户端传输

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

促进客户端连接到 Streamable MCP 端点，并透明地处理重新连接。它需要端点 URL，并允许配置 HTTP 客户端和重新连接行为。

```Go
package mcp

import (
	"context"
	"net/http"
)

// The streamable client handles reconnection transparently to the user.
type StreamableClientTransport struct {
	Endpoint         string
	HTTPClient       *http.Client
	ReconnectOptions *StreamableReconnectOptions
}

func (*StreamableClientTransport) Connect(context.Context) (Connection, error)
```