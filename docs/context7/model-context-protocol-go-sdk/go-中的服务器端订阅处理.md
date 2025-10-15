# Go 中的服务器端订阅处理

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

详细介绍了处理客户端订阅和通知的服务器端实现。它包括用于定义 `SubscribeHandler` 和 `UnsubscribeHandler` 的 `ServerOptions`，以及通知客户端资源更新的方法。

```Go
type ServerOptions struct {
  ...
	// Function called when a client session subscribes to a resource.
	SubscribeHandler func(context.Context, *ServerRequest[*SubscribeParams]) error
	// Function called when a client session unsubscribes from a resource.
	UnsubscribeHandler func(context.Context, *ServerRequest[*UnsubscribeParams]) error
}

func (*Server) ResourceUpdated(context.Context, *ResourceUpdatedNotificationParams) error
```

--------------------------------
