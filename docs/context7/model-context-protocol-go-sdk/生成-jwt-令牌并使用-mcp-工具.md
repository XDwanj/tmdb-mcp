# 生成 JWT 令牌并使用 MCP 工具

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/auth-middleware/README.md

使用 curl 生成具有指定用户 ID 和范围的 JWT 令牌，然后使用该令牌通过已认证的 JWT 端点调用“say_hi”MCP 工具的示例。

```bash
# Generate a token
curl 'http://localhost:8080/generate-token?user_id=alice&scopes=read,write'

# Use MCP tool with JWT authentication
curl -H 'Authorization: Bearer <generated_token>'
     -H 'Content-Type: application/json'
     -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"say_hi","arguments":{}}}'
     http://localhost:8080/mcp/jwt
```

--------------------------------
