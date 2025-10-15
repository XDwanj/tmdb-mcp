# 测试 MCP 工具的作用域限制

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/auth-middleware/README.md

使用 curl 调用需要“write”作用域的“create_resource”MCP 工具，并使用已授予“write”作用域的 JWT 令牌的示例。

```bash
# Access MCP tool requiring write scope
curl -H 'Authorization: Bearer <token_with_write_scope>'
     -H 'Content-Type: application/json'
     -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"create_resource","arguments":{"name":"test","description":"test resource","content":"test content"}}}'
     http://localhost:8080/mcp/jwt
```

--------------------------------
