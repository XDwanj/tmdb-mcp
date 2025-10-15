# Go：定义 ToolHandlerFor 签名

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了通用 `ToolHandlerFor` 函数签名，该函数处理带有类型化输入和输出参数的工具调用，接受上下文和服务器请求。

```Go
// A ToolHandlerFor handles a call to tools/call with typed arguments and results.
type ToolHandlerFor[In, Out any] func(context.Context, *ServerRequest[*CallToolParamsFor[In]]) (*CallToolResultFor[Out], error)
```

--------------------------------
