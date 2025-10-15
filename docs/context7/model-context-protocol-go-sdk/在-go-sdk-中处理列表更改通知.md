# 在 Go SDK 中处理列表更改通知

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了用于接收工具、Prompt 或资源列表更改通知的客户端选项。当服务器发送更新时，会调用这些处理程序。

```Go
type ClientOptions struct {
  ...
	ToolListChangedHandler      func(context.Context, *ClientRequest[*ToolListChangedParams])
	PromptListChangedHandler    func(context.Context, *ClientRequest[*PromptListChangedParams])
  // For both resources and resource templates.
	ResourceListChangedHandler  func(context.Context, *ClientRequest[*ResourceListChangedParams])
}
```

--------------------------------
