# 在 Go SDK 中实现 Completion Handler

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了用于处理来自客户端的完成请求的服务器选项。当客户端发送完成请求时，会调用 `CompletionHandler`，允许服务器提供建议。

```Go
type ServerOptions struct {
  ...
  // If non-nil, called when a client sends a completion request.
	CompletionHandler func(context.Context, *ServerRequest[*CompleteParams]) (*CompleteResult, error)
}
```

--------------------------------
