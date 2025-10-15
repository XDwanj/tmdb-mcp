# Go：Client CallTool 签名

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了客户端 `CallTool` 方法的签名，该方法允许客户端使用类型化参数调用工具，并期望原始 JSON 消息作为参数。

```Go
func (cs *ClientSession) CallTool(context.Context, *CallToolParams[json.RawMessage]) (*CallToolResult, error)
```

--------------------------------
