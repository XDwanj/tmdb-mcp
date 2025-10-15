# 创建和连接到 Go 中的 MCP 客户端

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/internal/readme/README.src.md

演示如何创建 `mcp.Client` 来与 MCP 服务器通信。此示例假定服务器正在运行并通过 stdin/stdout 可访问，演示了如何建立连接并在服务器上调用功能。

```Go
package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Create a new transport that uses stdin and stdout.
	transport := mcp.NewStdinStdoutTransport(os.Stdin, os.Stdout)

	// Create a new MCP client connected to the server via the transport.
	client := mcp.NewClient(transport)

	// Call the 'greet' feature on the server with the argument "World".
	response, err := client.CallFeature(context.Background(), "greet", "World")
	if err != nil {
		log.Fatalf("Failed to call feature: %v", err)
	}

	// Print the response from the server.
	log.Printf("Server response: %s\n", response)
}

```

--------------------------------
