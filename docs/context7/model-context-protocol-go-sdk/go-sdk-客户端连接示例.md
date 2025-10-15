# Go SDK 客户端连接示例

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

演示如何创建客户端实例，使用命令传输（stdin/stdout）建立与服务器的连接，并进行工具调用。它还展示了如何处理错误和关闭会话。

```Go
client := mcp.NewClient(&mcp.Implementation{Name:"mcp-client", Version:"v1.0.0"}, nil)
// Connect to a server over stdin/stdout
transport := &mcp.CommandTransport{
    Command: exec.Command("myserver"),
}
session, err := client.Connect(ctx, transport)
if err != nil { ... }
// Call a tool on the server.
content, err := session.CallTool(ctx, "greet", map[string]any{"name": "you"}, nil)
...
return session.Close()
```

--------------------------------
