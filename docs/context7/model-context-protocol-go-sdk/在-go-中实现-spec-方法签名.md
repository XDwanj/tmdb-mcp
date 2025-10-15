# 在 Go 中实现 Spec 方法签名

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

此 Go 代码片段演示了规范中定义的 RPC 方法的标准签名。它包括上下文和参数指针，返回结果指针和错误。为了向后兼容 Spec 更改，保留了此签名。

```Go
func (*ClientSession) ListTools(context.Context, *ListToolsParams) (*ListToolsResult, error)
```

--------------------------------
