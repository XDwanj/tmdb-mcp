# Go 中的客户端资源订阅

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

展示了客户端会话如何订阅和取消订阅资源更新。它还详细介绍了 `ClientOptions` 结构，其中包括一个用于处理资源更新通知的回调。

```Go
func (*ClientSession) Subscribe(context.Context, *SubscribeParams) error
func (*ClientSession) Unsubscribe(context.Context, *UnsubscribeParams) error

type ClientOptions struct {
  ...
  ResourceUpdatedHandler func(context.Context, *ClientSession, *ResourceUpdatedNotificationParams)
}
```

--------------------------------
