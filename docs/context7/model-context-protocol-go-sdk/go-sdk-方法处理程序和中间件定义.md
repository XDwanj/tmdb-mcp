# Go SDK 方法处理程序和中间件定义

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了用于处理 MCP 消息的 `MethodHandler` 类型和用于包装处理程序的 `Middleware` 类型。包括将发送和接收中间件添加到客户端和服务器的函数。

```Go
package mcp

import "context"

// A MethodHandler handles MCP messages.
// For methods, exactly one of the return values must be nil.
// For notifications, both must be nil.
type MethodHandler func(ctx context.Context, method string, req Request) (result Result, err error)

// Middleware is a function from MethodHandlers to MethodHandlers.
type Middleware func(MethodHandler) MethodHandler

// AddMiddleware wraps the client/server's current method handler using the provided
// middleware. Middleware is applied from right to left, so that the first one
// is executed first.
//
// For example, AddMiddleware(m1, m2, m3) augments the server method handler as
// m1(m2(m3(handler))).
func (c *Client) AddSendingMiddleware(middleware ...Middleware)
func (c *Client) AddReceivingMiddleware(middleware ...Middleware)
func (s *Server) AddSendingMiddleware(middleware ...Middleware)
func (s *Server) AddReceivingMiddleware(middleware ...Middleware)
```

--------------------------------
