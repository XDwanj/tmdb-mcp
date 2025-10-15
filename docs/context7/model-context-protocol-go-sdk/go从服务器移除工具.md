# Go：从服务器移除工具

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

展示了如何通过指定名称使用 `RemoveTools` 方法从 MCP 服务器中删除工具。

```Go
server.RemoveTools("add", "subtract")
```

--------------------------------
