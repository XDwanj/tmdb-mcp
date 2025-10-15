# 带有工具示例的 Go SDK 服务器

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

演示如何设置一个带有特定工具（“greet”）的服务器并通过 stdin/stdout 运行它。此服务器可以处理来自客户端的“greet”工具的请求。

```Go
// Create a server with a single tool.
server := mcp.NewServer(&mcp.Implementation{Name:"greeter", Version:"v1.0.0"}, nil)
mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)
// Run the server over stdin/stdout, until the client disconnects.
if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
    log.Fatal(err)
}
```

--------------------------------
