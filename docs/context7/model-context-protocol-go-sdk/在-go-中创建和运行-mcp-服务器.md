# 在 Go 中创建和运行 MCP 服务器

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/internal/readme/README.src.md

演示如何创建 `mcp.Server` 实例，添加功能（如简单工具），并通过 `mcp.Transport` 运行它，特别是使用 stdin/stdout 进行通信。这是设置 MCP 服务器的基础示例。

```Go
package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Create a new MCP server.
	server := mcp.NewServer(
		"my-simple-tool", // Name of the tool
		"A simple tool that says hello.", // Description of the tool
	)

	// Add a feature to the server. This feature is a simple function that takes a name and returns a greeting.
	server.AddFeature(
		"greet", // Name of the feature
		"Greets the user by name.", // Description of the feature
		func(ctx context.Context, name string) (string, error) {
			return "Hello, " + name + "!", nil
		},
	)

	// Create a new transport that uses stdin and stdout.
	transport := mcp.NewStdinStdoutTransport(os.Stdin, os.Stdout)

	// Run the server over the transport.
	log.Fatal(server.Run(context.Background(), transport))
}

```

--------------------------------
