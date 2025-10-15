# 在 Go 客户端中创建用于采样的消息处理程序

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

演示了如何为 Go 客户端设置 `CreateMessageHandler` 选项，以处理用于采样的服务器调用。当服务器会话调用 `CreateMessage` Spec 方法时，会调用此函数。

```Go
type ClientOptions struct {
  ...
  CreateMessageHandler func(context.Context, *ClientSession, *CreateMessageParams) (*CreateMessageResult, error)
}
```

--------------------------------
