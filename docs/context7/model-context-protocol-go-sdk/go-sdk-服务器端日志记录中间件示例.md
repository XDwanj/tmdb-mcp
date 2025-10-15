# Go SDK 服务器端日志记录中间件示例

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

提供了一个在 Go 中实现服务器端日志记录中间件的示例。`withLogging` 函数包装了一个 MethodHandler 来记录传入的请求和传出的响应。

```Go
import (
    "context"
    "log"

    mcp "path/to/mcp"
)

func withLogging(h mcp.MethodHandler) mcp.MethodHandler{
    return func(ctx context.Context, method string, req mcp.Request) (res mcp.Result, err error) {
        log.Printf("request: %s %v", method, params)
        defer func() { log.Printf("response: %v, %v", res, err) }()
        return h(ctx, s , method, params)
    }
}

server.AddReceivingMiddleware(withLogging)
```

--------------------------------
