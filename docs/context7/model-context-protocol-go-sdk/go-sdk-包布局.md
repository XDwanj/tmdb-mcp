# Go SDK 包布局

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了 Go SDK 的建议包结构，将核心 MCP API 整合到一个 'mcp' 包中，以提高可发现性和与 Go 习惯的对齐。它还为 JSON Schema 和 JSON-RPC 等相关功能指定了单独的包。

```Go
github.com/modelcontextprotocol/go-sdk/mcp
github.com/modelcontextprotocol/go-sdk/jsonschema
github.com/modelcontextprotocol/go-sdk/internal/jsonrpc2
```

--------------------------------
