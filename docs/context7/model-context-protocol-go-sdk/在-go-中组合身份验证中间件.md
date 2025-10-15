# 在 Go 中组合身份验证中间件

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/auth-middleware/README.md

演示了如何将身份验证中间件与其他自定义中间件处理程序在 Go 中进行组合。这种模式允许在到达主 MCP 处理程序之前分层安全检查。

```Go
// Combine authentication middleware with custom middleware
authenticatedHandler := authMiddleware(customMiddleware(mcpHandler))
```

--------------------------------
