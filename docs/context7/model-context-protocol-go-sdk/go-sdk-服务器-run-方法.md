# Go SDK 服务器 Run 方法

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

提供了 `Server.Run` 方法的签名，这是一个方便的函数，用于处理运行服务器会话直到客户端断开连接的常见情况。

```Go
func (*Server) Run(context.Context, Transport)
```

--------------------------------
