# Go：添加工具的示例

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

提供使用 `AddTool` 函数和特定工具定义和处理程序将“add”和“subtract”工具添加到 MCP 服务器的具体示例。

```Go
mcp.AddTool(server, &mcp.Tool{Name: "add", Description: "add numbers"}, addHandler)
mcp.AddTool(server, &mcp.Tool{Name: "subtract", Description: "subtract numbers"}, subHandler)
```

--------------------------------
