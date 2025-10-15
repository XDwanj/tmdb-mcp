# MCP 服务器身份验证中间件实现

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/auth-middleware/README.md

展示了创建 MCP 服务器、应用 `auth.RequireBearerToken` 中间件以及集成自定义中间件以处理已认证请求的 Go 代码片段。

```go
// Create MCP server
server := mcp.NewServer(&mcp.Implementation{Name: "authenticated-mcp-server"}, nil)

// Create authentication middleware
authMiddleware := auth.RequireBearerToken(verifier, &auth.RequireBearerTokenOptions{
    Scopes: []string{"read", "write"},
})

// Create MCP handler
handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
    return server
}, nil)

// Apply authentication middleware to MCP handler
authenticatedHandler := authMiddleware(customMiddleware(handler))
```

--------------------------------
