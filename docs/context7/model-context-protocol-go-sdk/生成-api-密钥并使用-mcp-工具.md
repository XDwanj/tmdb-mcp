# 生成 API 密钥并使用 MCP 工具

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/auth-middleware/README.md

使用 curl 生成具有指定用户 ID 和范围的 API 密钥，然后使用该密钥通过已认证的 API 密钥端点调用“get_user_info”MCP 工具的示例。

```bash
# Generate an API key
curl -X POST 'http://localhost:8080/generate-api-key?user_id=bob&scopes=read'

# Use MCP tool with API key authentication
curl -H 'Authorization: Bearer <generated_api_key>'
     -H 'Content-Type: application/json'
     -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"get_user_info","arguments":{"user_id":"test"}}}'
     http://localhost:8080/mcp/apikey
```

--------------------------------
