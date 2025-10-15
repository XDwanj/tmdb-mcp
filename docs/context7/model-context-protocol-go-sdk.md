================
代码片段
================
## 在自定义主机/端口上启动 MCP HTTP 服务器

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/http/README.md

在自定义主机和端口上启动 MCP HTTP 服务器。

```bash
go run main.go -host 0.0.0.0 -port 9000 server
```

--------------------------------

## 设置和运行 MCP 服务器

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/auth-middleware/README.md

用于导航到示例目录、整理 Go 模块并运行启用身份验证中间件的 MCP 服务器的命令。

```bash
cd examples/server/auth-middleware
go mod tidy
go run main.go
```

--------------------------------

## 启动 MCP HTTP 服务器

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/http/README.md

在默认端口 8080 上启动 MCP HTTP 服务器，提供一个“cityTime”工具。

```bash
go run main.go server
```

--------------------------------

## Go 版权头部示例

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

SDK 存储库中使用的标准 Go 文件版权头部示例，遵循 MIT 风格的许可证和归属。

```go
// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
```

--------------------------------

## Go：添加工具的示例

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

提供使用 `AddTool` 函数和特定工具定义和处理程序将“add”和“subtract”工具添加到 MCP 服务器的具体示例。

```Go
mcp.AddTool(server, &mcp.Tool{Name: "add", Description: "add numbers"}, addHandler)
mcp.AddTool(server, &mcp.Tool{Name: "subtract", Description: "subtract numbers"}, subHandler)
```

--------------------------------

## 带有工具示例的 Go SDK 服务器

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

## 测试 MCP 服务器

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/auth-middleware/README.md

用于运行单元测试、基准测试和生成 MCP 服务器代码覆盖率报告的 Go 命令。

```bash
# Run all tests
go test -v

# Run benchmark tests
go test -bench=.

# Generate coverage report
go test -cover
```

--------------------------------

## 在 Claude 代码中列出支持的时区

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/http/README.md

在 Claude 代码中列出“cityTime”工具支持的城市的示例交互。

```bash
> what timezones do you support?

⏺ The timezone tool supports three US cities:
  - NYC (Eastern Time)
  - SF (Pacific Time)
  - Boston (Eastern Time)
```

--------------------------------

## 运行 MCP HTTP 客户端

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/http/README.md

运行 MCP 客户端以连接到服务器、列出工具并为指定的城市调用“cityTime”工具。

```bash
go run main.go client
```

--------------------------------

## 在 Go 中创建和运行 MCP 服务器

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/README.md

演示如何创建 MCP 服务器、添加工具并通过 stdin/stdout 运行它。此示例定义了“SayHi”工具的输入和输出结构。

```go
package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Input struct {
	Name string `json:"name" jsonschema:"the name of the person to greet"`
}

type Output struct {
	Greeting string `json:"greeting" jsonschema:"the greeting to tell to the user"`
}

func SayHi(ctx context.Context, req *mcp.CallToolRequest, input Input) (*mcp.CallToolResult, Output, error) {
	return nil, Output{Greeting: "Hi " + input.Name}, nil
}

func main() {
	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)
	// Run the server over stdin/stdout, until the client disconnects
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
```

--------------------------------

## 在 Go 中创建和运行 MCP 服务器

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

## 使用 go.work 设置 Go SDK 开发环境

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/CONTRIBUTING.md

演示如何初始化 Go 工作区，以便针对本地项目测试 SDK 更改。这对于多模块开发非常有用，其中 SDK 与使用它的项目一起被修改。

```bash
go work init ./project ./go-sdk
```

--------------------------------

## 创建和连接到 Go 中的 MCP 客户端

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

## 将 MCP 服务器添加到 Claude 代码

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/http/README.md

将正在运行的 MCP HTTP 服务器添加到 Claude 代码以进行集成。

```bash
claude mcp add -t http timezone http://localhost:8080
```

--------------------------------

## 在 Claude 代码中查询时区

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/http/README.md

在 Claude 代码中查询时区的示例交互，使用“cityTime”工具。

```bash
> what's the timezone

⏺ I'll get the current time in a major US city for you.

⏺ timezone - cityTime (MCP)(city: "nyc")
  ⎿ The current time in New York City is 7:30:16 PM EDT on Wedn
    esday, July 23, 2025


⏺ The current timezone is EDT (Eastern Daylight Time), and it's
   7:30 PM on Wednesday, July 23, 2025.
```

--------------------------------

## 生成 API 密钥并使用 MCP 工具

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

## Go SDK 客户端连接示例

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

## Go SDK 服务器端日志记录中间件示例

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

## Go：带有类型化参数的示例处理程序

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

提供了一个 Go 处理程序函数 (`addHandler`) 的示例，该函数处理类型化参数 (`AddParams`) 并为 MCP 工具返回类型化结果 (`int`)。

```Go
type AddParams struct {
    X int `json:"x"`
    Y int `json:"y"`
}

func addHandler(ctx context.Context, req *mcp.ServerRequest[*mcp.CallToolParamsFor[AddParams]]) (*mcp.CallToolResultFor[int], error) {
    return &mcp.CallToolResultFor[int]{StructuredContent: req.Params.Arguments.X + req.Params.Arguments.Y}, nil
}
```

--------------------------------

## 生成 JWT 令牌并使用 MCP 工具

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

## 运行顺序思考 MCP 服务器（Go）

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/sequentialthinking/README.md

以标准 I/O 模式启动用于顺序思考的模型上下文协议 (MCP) 服务器。这是运行服务器的默认方式。

```bash
go run .
```

--------------------------------

## 测试 MCP 工具的作用域限制

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/auth-middleware/README.md

使用 curl 调用需要“write”作用域的“create_resource”MCP 工具，并使用已授予“write”作用域的 JWT 令牌的示例。

```bash
# Access MCP tool requiring write scope
curl -H 'Authorization: Bearer <token_with_write_scope>'
     -H 'Content-Type: application/json'
     -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"create_resource","arguments":{"name":"test","description":"test resource","content":"test content"}}}'
     http://localhost:8080/mcp/jwt
```

--------------------------------

## MCP 服务器身份验证中间件实现

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

## 启动思考会话（JSON）

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/sequentialthinking/README.md

为给定的问题启动一个新的顺序思考会话。它接受问题陈述、可选的会话 ID 和步骤的初始估计。

```json
{
  "method": "tools/call",
  "params": {
    "name": "start_thinking",
    "arguments": {
      "problem": "How should I design a scalable microservices architecture?",
      "sessionId": "architecture_design",
      "estimatedSteps": 8
    }
  }
}
```

--------------------------------

## 在 Go 中组合身份验证中间件

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/auth-middleware/README.md

演示了如何将身份验证中间件与其他自定义中间件处理程序在 Go 中进行组合。这种模式允许在到达主 MCP 处理程序之前分层安全检查。

```Go
// Combine authentication middleware with custom middleware
authenticatedHandler := authMiddleware(customMiddleware(mcpHandler))
```

--------------------------------

## 在 MCP 工具中检索和使用身份验证信息（Go）

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/auth-middleware/README.md

演示了如何在 MCP 工具中的传入请求中访问身份验证信息，特别是 `TokenInfo`（包括范围）。它包括在执行工具逻辑之前检查所需范围。

```Go
// Get authentication information in MCP tool
func MyTool(ctx context.Context, req *mcp.CallToolRequest, args MyArgs) (*mcp.CallToolResult, any, error) {
    // Extract authentication info from request 
    userInfo := req.Extra.TokenInfo
    
    // Check scopes
    if !slices.Contains(userInfo.Scopes, "read") {
        return nil, nil, fmt.Errorf("insufficient permissions: read scope required")
    }
    
    // Execute tool logic
    return &mcp.CallToolResult{
        Content: []mcp.Content{
            &mcp.TextContent{Text: "Tool executed successfully"},
        },
    }, nil, nil
}
```

--------------------------------

## 通过 HTTP 运行顺序思考 MCP 服务器（Go）

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/sequentialthinking/README.md

启动用于顺序思考的模型上下文协议 (MCP) 服务器，并通过 HTTP 在指定端口上公开它。这允许与服务器进行基于网络的交互。

```bash
go run . -http :8080
```

--------------------------------

## 调用带有 nil 参数的 Spec 方法

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

此 Go 示例展示了如何通过为当前不需要的参数传递 nil 来调用 Ping 等 Spec 方法。这种方法确保了兼容性，即使 Spec 将来引入新参数。

```Go
err := session.Ping(ctx, nil)
```

--------------------------------

## Go SDK 上下文取消示例

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

演示如何使用上下文取消在 Go SDK 中实现操作取消。创建一个带有取消函数的新的上下文，并在 goroutine 中调用操作。

```Go
import "context"

ctx, cancel := context.WithCancel(ctx)
go session.CallTool(ctx, "slow", map[string]any{}, nil)
cancel()
```

--------------------------------

## 在 Go 中实现 JWT 令牌验证器

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/auth-middleware/README.md

提供了用于验证 JWT 令牌的 Go 函数签名。此函数接受上下文和令牌字符串，在成功时返回令牌信息，在令牌无效时返回错误。

```Go
func jwtVerifier(ctx context.Context, tokenString string) (*auth.TokenInfo, error) {
    // JWT token verification logic
    // On success: Return TokenInfo
    // On failure: Return auth.ErrInvalidToken
}
```

--------------------------------

## Go SDK：用于 Stdio 通信的 CommandTransport

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了 stdio 传输的 CommandTransport。此类型通过启动命令并通过其 stdin 和 stdout 流式传输 JSON-RPC 消息进行连接，使用换行符分隔的 JSON 进行通信。

```Go
// A CommandTransport is a [Transport] that runs a command and communicates
// with it over stdin/stdout, using newline-delimited JSON.
type CommandTransport struct { Command *exec.Command }

// Connect starts the command, and connects to it over stdin/stdout.
func (*CommandTransport) Connect(ctx context.Context) (Connection, error) {
```

--------------------------------

## 用于 Prompt 分页的迭代器方法

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

此 Go 代码为 ListPrompts Spec 方法定义了一个迭代器方法。它自动处理分页，允许遍历所有页面的提示。如果提供了参数，则迭代从指定的游标开始。

```Go
func (*ClientSession) Prompts(context.Context, *ListPromptsParams) iter.Seq2[Prompt, error]
```

--------------------------------

## 用于 Tool 分页的迭代器方法

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

此 Go 代码为 ListTools Spec 方法定义了一个迭代器方法。它自动处理分页，允许遍历所有页面的工具。如果提供了参数，则迭代从指定的游标开始。

```Go
func (*ClientSession) Tools(context.Context, *ListToolsParams) iter.Seq2[Tool, error]
```

--------------------------------

## 用于 Resource 分页的迭代器方法

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

此 Go 代码为 ListResource Spec 方法定义了一个迭代器方法。它自动处理分页，允许遍历所有页面的资源。如果提供了参数，则迭代从指定的游标开始。

```Go
func (*ClientSession) Resources(context.Context, *ListResourceParams) iter.Seq2[Resource, error]
```

--------------------------------

## 用于 Resource Template 分页的迭代器方法

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

此 Go 代码为 ListResourceTemplates Spec 方法定义了一个迭代器方法。它自动处理分页，允许遍历所有页面的资源模板。如果提供了参数，则迭代从指定的游标开始。

```Go
func (*ClientSession) ResourceTemplates(context.Context, *ListResourceTemplatesParams) iter.Seq2[ResourceTemplate, error]
```

--------------------------------

## 用于 MCP 会话的 SSE 服务器传输

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

表示通过挂起 GET 请求建立的逻辑服务器发送事件 (SSE) 会话。此传输处理到会话端点的传入请求并促进连接过程，允许事件流式传输。

```Go
package mcp

import (
	"context"
	"net/http"
)

// A SSEServerTransport is a logical SSE session created through a hanging GET
// request.
type SSEServerTransport struct {
    Endpoint string
    Response http.ResponseWriter
}

// ServeHTTP handles POST requests to the transport endpoint.
func (*SSEServerTransport) ServeHTTP(w http.ResponseWriter, req *http.Request)

// Connect sends the 'endpoint' event to the client.
// See [SSEServerTransport] for more details on the [Connection] implementation.
func (*SSEServerTransport) Connect(context.Context) (Connection, error)
```

--------------------------------

## 在 Go 中连接到 MCP 服务器并调用工具

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/README.md

展示了如何创建 MCP 客户端，通过命令传输（使用 stdin/stdout）连接到服务器，并调用工具。它处理响应并记录输出。

```Go
package main

import (
	"context"
	"log"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()

	// Create a new client, with no features.
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

	// Connect to a server over stdin/stdout
	transport := &mcp.CommandTransport{Command: exec.Command("myserver")}
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Call a tool on the server.
	params := &mcp.CallToolParams{
		Name:      "greet",
		Arguments: map[string]any{"name": "you"},
	}
	res, err := session.CallTool(ctx, params)
	if err != nil {
		log.Fatalf("CallTool failed: %v", err)
	}
	if res.IsError {
		log.Fatal("tool failed")
	}
	for _, c := range res.Content {
		log.Print(c.(*mcp.TextContent).Text)
	}
}
```

--------------------------------

## 用于日志记录和分页的 Go 服务器选项

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

配置服务器端的日志记录和分页行为。包括日志记录器名称、日志记录间隔和页面大小的选项。

```Go
package mcp

import "time"

type ServerOptions struct {
  // ...
  // The value for the "logger" field of the notification.
  LoggerName string
  // Log notifications to a single ClientSession will not be
  // sent more frequently than this duration.
  LoggingInterval time.Duration
  // PageSize defines the number of items to return per page.
  PageSize int
}
```

--------------------------------

## 在 Go 中添加和删除资源/模板

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

演示在服务器上添加和删除资源及资源模板的方法。它展示了 AddResource、AddResourceTemplate、RemoveResources 和 RemoveResourceTemplates 的签名。

```Go
func (*Server) AddResource(*Resource, ResourceHandler)
func (*Server) AddResourceTemplate(*ResourceTemplate, ResourceHandler)

func (s *Server) RemoveResources(uris ...string)
func (s *Server) RemoveResourceTemplates(uriTemplates ...string)
```

--------------------------------

## Go：向服务器添加具有显式架构的工具

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

演示如何向具有显式定义的 `Tool` 结构（包括其名称、描述和输入架构）的服务器添加工具。

```Go
t := &Tool{Name: ..., Description: ..., InputSchema: &jsonschema.Schema{...}}
server.AddTool(t, myHandler)
```

--------------------------------

## Go SDK 服务器 Run 方法

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

提供了 `Server.Run` 方法的签名，这是一个方便的函数，用于处理运行服务器会话直到客户端断开连接的常见情况。

```Go
func (*Server) Run(context.Context, Transport)
```

--------------------------------

## Go SDK 客户端结构

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了 Go SDK 中客户端的结构，包括创建新客户端、连接到传输以及检索活动会话的方法。它还概述了客户端选项和会话详细信息。

```Go
type Client struct { /* ... */ }
func NewClient(impl *Implementation, opts *ClientOptions) *Client
func (*Client) Connect(context.Context, Transport) (*ClientSession, error)
func (*Client) Sessions() iter.Seq[*ClientSession]

type ClientOptions struct { /* ... */ } // described below

type ClientSession struct { /* ... */ }
func (*ClientSession) Client() *Client
func (*ClientSession) Close() error
func (*ClientSession) Wait() error
// Methods for calling through the ClientSession are described below.
// For example: ClientSession.ListTools.
```

--------------------------------

## Go SDK 服务器结构

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了 Go SDK 中服务器的结构，包括创建新服务器、连接到传输以及检索活动会话的方法。它还概述了服务器选项和会话详细信息。

```Go
type Server struct { /* ... */ }
func NewServer(impl *Implementation, opts *ServerOptions) *Server
func (*Server) Connect(context.Context, Transport) (*ServerSession, error)
func (*Server) Sessions() iter.Seq[*ServerSession]
// Methods for adding/removing server features are described below.

type ServerOptions struct { /* ... */ } // described below

type ServerSession struct { /* ... */ }
func (*ServerSession) Server() *Server
func (*ServerSession) Close() error
func (*ServerSession) Wait() error
// Methods for calling through the ServerSession are described below.
// For example: ServerSession.ListRoots.
```

--------------------------------

## Go：向服务器添加工具

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

演示如何使用 `AddTool` 方法或函数向 MCP 服务器添加工具。该函数是通用的，并从处理程序参数推断模式。

```Go
func (s *Server) AddTool(t *Tool, h ToolHandler)
func AddTool[In, Out any](s *Server, t *Tool, h ToolHandlerFor[In, Out])
```

--------------------------------

## Go SDK 版权头部

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/CONTRIBUTING.md

提供了 SDK 中所有 Go 源文件所需的标准版权头部格式。这确保了包含正确的归属和许可信息。

```Go
// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

```

--------------------------------

## 在 Go 中添加和删除 Prompt

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

演示如何使用 Server.AddPrompt 和 Server.RemovePrompts 方法添加带有处理程序的 Prompt，然后删除它。它包括 Prompt 处理程序参数的定义。

```Go
type codeReviewArgs struct {
  Code string `json:"code"`
}

func codeReviewHandler(context.Context, *ServerSession, *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {...}

server.AddPrompt(
  &mcp.Prompt{Name: "code_review", Description: "review code"},
  codeReviewHandler,
)

server.RemovePrompts("code_review")
```

--------------------------------

## Go 中的客户端资源订阅

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

展示了客户端会话如何订阅和取消订阅资源更新。它还详细介绍了 `ClientOptions` 结构，其中包括一个用于处理资源更新通知的回调。

```Go
func (*ClientSession) Subscribe(context.Context, *SubscribeParams) error
func (*ClientSession) Unsubscribe(context.Context, *UnsubscribeParams) error

type ClientOptions struct {
  ...
  ResourceUpdatedHandler func(context.Context, *ClientSession, *ResourceUpdatedNotificationParams)
}
```

--------------------------------

## 查看思考会话（JSON）

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/sequentialthinking/README.md

检索对特定思考会话的完整审查，包括所有思考步骤及其历史记录。它需要会话 ID。

```json
{
  "method": "tools/call", 
  "params": {
    "name": "review_thinking", 
    "arguments": {
      "sessionId": "architecture_design"
    }
  }
}
```

--------------------------------

## Go SDK 包布局

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了 Go SDK 的建议包结构，将核心 MCP API 整合到一个 'mcp' 包中，以提高可发现性和与 Go 习惯的对齐。它还为 JSON Schema 和 JSON-RPC 等相关功能指定了单独的包。

```Go
github.com/modelcontextprotocol/go-sdk/mcp
github.com/modelcontextprotocol/go-sdk/jsonschema
github.com/modelcontextprotocol/go-sdk/internal/jsonrpc2
```

--------------------------------

## Go 中的服务器端订阅处理

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

详细介绍了处理客户端订阅和通知的服务器端实现。它包括用于定义 `SubscribeHandler` 和 `UnsubscribeHandler` 的 `ServerOptions`，以及通知客户端资源更新的方法。

```Go
type ServerOptions struct {
  ...
	// Function called when a client session subscribes to a resource.
	SubscribeHandler func(context.Context, *ServerRequest[*SubscribeParams]) error
	// Function called when a client session unsubscribes from a resource.
	UnsubscribeHandler func(context.Context, *ServerRequest[*UnsubscribeParams]) error
}

func (*Server) ResourceUpdated(context.Context, *ResourceUpdatedNotificationParams) error
```

--------------------------------

## Go：定义 CallToolParams 和 Tool 结构

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了用于类型化工具参数的通用 `CallToolParamsFor` 结构以及用于表示 MCP 工具的 `Tool` 结构，包括模式和描述。

```Go
type CallToolParamsFor[In any] struct {
	Meta      Meta   `json:"_meta,omitempty"`
	Arguments In     `json:"arguments,omitempty"`
	Name      string `json:"name"`
}

type Tool struct {
	Annotations *ToolAnnotations   `json:"annotations,omitempty"`
	Description string             `json:"description,omitempty"`
	InputSchema *jsonschema.Schema `json:"inputSchema"`
	Name string                    `json:"name"`
}
```

--------------------------------

## Go：定义 ToolHandlerFor 签名

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了通用 `ToolHandlerFor` 函数签名，该函数处理带有类型化输入和输出参数的工具调用，接受上下文和服务器请求。

```Go
// A ToolHandlerFor handles a call to tools/call with typed arguments and results.
type ToolHandlerFor[In, Out any] func(context.Context, *ServerRequest[*CallToolParamsFor[In]]) (*CallToolResultFor[Out], error)
```

--------------------------------

## 继续思考会话 - 创建分支（JSON）

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/sequentialthinking/README.md

在思考会话中创建替代的推理路径或分支。它需要会话 ID 和新分支的思考内容。

```json
{
  "method": "tools/call",
  "params": {
    "name": "continue_thinking",
    "arguments": {
      "sessionId": "architecture_design", 
      "thought": "Alternative approach: Start with a monolith-first strategy and extract services gradually.",
      "createBranch": true
    }
  }
}
```

--------------------------------

## 用于 MCP 服务器的 Streamable HTTP 处理程序

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

实现了用于提供 Streamable MCP 会话的 `http.Handler`。与 SSE 处理程序类似，它接受用于创建新会话服务器的回调，并管理会话生命周期，包括关闭活动会话。

```Go
package mcp

import (
	"net/http"
)

// The StreamableHTTPHandler interface is symmetrical to the SSEHTTPHandler.
type StreamableHTTPHandler struct { /* unexported fields */ }
func NewStreamableHTTPHandler(getServer func(request *http.Request) *Server) *StreamableHTTPHandler
func (*StreamableHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request)
func (*StreamableHTTPHandler) Close() error
```

--------------------------------

## 在 Go SDK 中实现 Completion Handler

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了用于处理来自客户端的完成请求的服务器选项。当客户端发送完成请求时，会调用 `CompletionHandler`，允许服务器提供建议。

```Go
type ServerOptions struct {
  ...
  // If non-nil, called when a client sends a completion request.
	CompletionHandler func(context.Context, *ServerRequest[*CompleteParams]) (*CompleteResult, error)
}
```

--------------------------------

## Go SDK 进度通知参数

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了 Go SDK 中用于请求和发送进度通知的结构。`Meta` 包含 `ProgressToken` 以指示请求进度更新。

```Go
type XXXParams struct { // where XXX is each type of call
  Meta Meta
  ...
}

type Meta struct {
  Data          map[string]any // arbitrary data
  ProgressToken any // string or int
}
```

--------------------------------

## 在 Go 中实现 Spec 方法签名

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

此 Go 代码片段演示了规范中定义的 RPC 方法的标准签名。它包括上下文和参数指针，返回结果指针和错误。为了向后兼容 Spec 更改，保留了此签名。

```Go
func (*ClientSession) ListTools(context.Context, *ListToolsParams) (*ListToolsResult, error)
```

--------------------------------

## 用于 MCP 的 SSE 客户端传输

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

使客户端能够连接到基于 SSE 的 MCP 端点。它允许指定目标端点和可选的 `http.Client` 来发出请求。传输处理到服务器的连接建立。

```Go
package mcp

import (
	"context"
	"net/http"
)

type SSEClientTransport struct {
	// Endpoint is the SSE endpoint to connect to.
	Endpoint string
	// HTTPClient is the client to use for making HTTP requests. If nil,
	// http.DefaultClient is used.
	HTTPClient *http.Client
}

// Connect connects through the client endpoint.
func (*SSEClientTransport) Connect(ctx context.Context) (Connection, error)
```

--------------------------------

## Go 版日志记录处理程序，带 slog

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

为 MCP 定义了一个 slog.Handler，允许服务器作者使用标准的 slog 包进行日志记录。它包括日志记录器名称和速率限制的选项。

```Go
package mcp

import (
	"context"
	"time"
	"log/slog"
)

// A LoggingHandler is a [slog.Handler] for MCP.
type LoggingHandler struct {}

// LoggingHandlerOptions are options for a LoggingHandler.
type LoggingHandlerOptions struct {
	// The value for the "logger" field of logging notifications.
	LoggerName string
	// Limits the rate at which log messages are sent.
	// If zero, there is no rate limiting.
	MinInterval time.Duration
}

// NewLoggingHandler creates a [LoggingHandler] that logs to the given [ServerSession] using a
// [slog.JSONHandler].
func NewLoggingHandler(ss *ServerSession, opts *LoggingHandlerOptions) *LoggingHandler
```

--------------------------------

## Go SDK 方法处理程序和中间件定义

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了用于处理 MCP 消息的 `MethodHandler` 类型和用于包装处理程序的 `Middleware` 类型。包括将发送和接收中间件添加到客户端和服务器的函数。

```Go
package mcp

import "context"

// A MethodHandler handles MCP messages.
// For methods, exactly one of the return values must be nil.
// For notifications, both must be nil.
type MethodHandler func(ctx context.Context, method string, req Request) (result Result, err error)

// Middleware is a function from MethodHandlers to MethodHandlers.
type Middleware func(MethodHandler) MethodHandler

// AddMiddleware wraps the client/server's current method handler using the provided
// middleware. Middleware is applied from right to left, so that the first one
// is executed first.
//
// For example, AddMiddleware(m1, m2, m3) augments the server method handler as
// m1(m2(m3(handler))).
func (c *Client) AddSendingMiddleware(middleware ...Middleware)
func (c *Client) AddReceivingMiddleware(middleware ...Middleware)
func (s *Server) AddSendingMiddleware(middleware ...Middleware)
func (s *Server) AddReceivingMiddleware(middleware ...Middleware)
```

--------------------------------

## 在 Go 客户端中添加和删除根目录

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

解释了如何使用 `AddRoots` 和 `RemoveRoots` 方法管理 Go 客户端的根目录。`AddRoots` 添加指定的根目录，替换具有相同 URI 的现有根目录，并通知已连接的服务器。`RemoveRoots` 按 URI 删除根目录，如果根目录不存在则不报错。

```Go
// AddRoots adds the given roots to the client,
// replacing any with the same URIs,
// and notifies any connected servers.
func (*Client) AddRoots(roots ...*Root)

// RemoveRoots removes the roots with the given URIs.
// and notifies any connected servers if the list has changed.
// It is not an error to remove a nonexistent root.
func (*Client) RemoveRoots(uris ...string)
```

--------------------------------

## 用于 MCP 服务器的 SSE HTTP 处理程序

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

提供了一个 `http.Handler`，用于处理基于服务器发送事件 (SSE) 的 MCP 会话。它接受一个回调函数来绑定新会话的服务器，允许每个连接使用无状态或有状态的服务器实例。处理程序管理会话的生命周期，包括关闭活动会话。

```Go
package mcp

import (
	"net/http"
)

// SSEHTTPHandler is an http.Handler that serves SSE-based MCP sessions as defined by
// the 2024-11-05 version of the MCP protocol.
type SSEHTTPHandler struct { /* unexported fields */ }

// NewSSEHTTPHandler returns a new [SSEHTTPHandler] that is ready to serve HTTP.
//
// The getServer function is used to bind created servers for new sessions. It
// is OK for getServer to return the same server multiple times.
func NewSSEHTTPHandler(getServer func(request *http.Request) *Server) *SSEHTTPHandler

func (*SSEHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request)

// Close prevents the SSEHTTPHandler from accepting new sessions, closes active
// sessions, and awaits their graceful termination.
func (*SSEHTTPHandler) Close() error
```

--------------------------------

## 继续思考会话 - 完成（JSON）

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/sequentialthinking/README.md

将当前思考过程标记为已完成，或表示当前不需要进一步的步骤。它需要会话 ID 和最终的思考。

```json
{
  "method": "tools/call",
  "params": {
    "name": "continue_thinking",
    "arguments": {
      "sessionId": "architecture_design",
      "thought": "Based on this analysis, I recommend starting with 3 core services: User Management, Order Processing, and Inventory Management.",
      "nextNeeded": false
    }
  }
}
```

--------------------------------

## Go SDK 进度通知方法

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

提供了 Go SDK 中发送进度通知的方法。如果请求中提供了进度令牌，`NotifyProgress` 会向对等方发送通知。

```Go
import "context"

func (*ClientSession) NotifyProgress(context.Context, *ProgressNotification)
func (*ServerSession) NotifyProgress(context.Context, *ProgressNotification)
```

--------------------------------

## Go：从服务器移除工具

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

展示了如何通过指定名称使用 `RemoveTools` 方法从 MCP 服务器中删除工具。

```Go
server.RemoveTools("add", "subtract")
```

--------------------------------

## Go 客户端/服务器会话中的 Ping 和 KeepAlive

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

演示了 `ClientSession` 和 `ServerSession` 的 `Ping` 方法，用于检查对等连接。它还展示了如何通过 `KeepAlive` 选项配置自动保持活动行为，该选项会在对等方未能响应 ping 时关闭会话。

```Go
func (c *ClientSession) Ping(ctx context.Context, *PingParams) error
func (c *ServerSession) Ping(ctx context.Context, *PingParams) error
```

```Go
type ClientOptions struct {
  ...
  KeepAlive time.Duration
}

type ServerOptions struct {
  ...
  KeepAlive time.Duration
}
```

--------------------------------

## 在 Go 客户端中创建用于采样的消息处理程序

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

演示了如何为 Go 客户端设置 `CreateMessageHandler` 选项，以处理用于采样的服务器调用。当服务器会话调用 `CreateMessage` Spec 方法时，会调用此函数。

```Go
type ClientOptions struct {
  ...
  CreateMessageHandler func(context.Context, *ClientSession, *CreateMessageParams) (*CreateMessageResult, error)
}
```

--------------------------------

## Go 客户端日志记录消息选项

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了用于接收服务器日志记录消息的客户端处理程序。当收到 `LoggingMessageNotification` 时，会调用此回调。

```Go
package mcp

import "context"

type ClientOptions struct {
  // ...
  LoggingMessageHandler func(context.Context, *ClientSession, *LoggingMessageParams)
}
```

--------------------------------

## 在 Go SDK 中处理列表更改通知

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了用于接收工具、Prompt 或资源列表更改通知的客户端选项。当服务器发送更新时，会调用这些处理程序。

```Go
type ClientOptions struct {
  ...
	ToolListChangedHandler      func(context.Context, *ClientRequest[*ToolListChangedParams])
	PromptListChangedHandler    func(context.Context, *ClientRequest[*PromptListChangedParams])
  // For both resources and resource templates.
	ResourceListChangedHandler  func(context.Context, *ClientRequest[*ResourceListChangedParams])
}
```

--------------------------------

## 协议数据类型（Go）

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了模型上下文协议的关键数据结构，这些结构是从其 JSON Schema 生成的。它包括参数、结果和内容表示的类型，利用 Go 的接口和结构来实现灵活性。

```Go
package main

import "encoding/json"

// Meta includes arbitrary data and a progress token.
type Meta struct {
    Data        map[string]any `json:"data,omitempty"`
    ProgressToken string         `json:"progress_token,omitempty"`
}

// ReadResourceParams defines parameters for reading a resource.
type ReadResourceParams struct {
    URI string `json:"uri"`
}

// CallToolResult represents the outcome of a tool call.
type CallToolResult struct {
    Meta    Meta      `json:"_meta,omitempty"`
    Content []Content `json:"content"`
    IsError bool      `json:"isError,omitempty"`
}

// Content is an interface representing different types of content.
// It is implemented by types like TextContent, ImageContent, etc.
type Content interface {
    // (unexported methods)
}

// TextContent represents textual content.
type TextContent struct {
    Text string
}

// ImageContent represents image content (example).
type ImageContent struct {
    // Image data or reference
}

// AudioContent represents audio content (example).
type AudioContent struct {
    // Audio data or reference
}

// EmbeddedResource represents embedded resource content (example).
type EmbeddedResource struct {
    // Resource data or reference
}

// ResourceContents is a struct for representing multiple resource contents, using optional fields for union types.
type ResourceContents struct {
    TextContents     []TextContent     `json:"text_contents,omitempty"`
    ImageContents    []ImageContent    `json:"image_contents,omitempty"`
    AudioContents    []AudioContent    `json:"audio_contents,omitempty"`
    EmbeddedResource []EmbeddedResource `json:"embedded_resources,omitempty"`
}

func main() {
    // Example usage (conceptual):
    // params := ReadResourceParams{URI: "example.com/resource"}
    // text := TextContent{Text: "Hello, world!"}
    // result := CallToolResult{Content: []Content{&text}, IsError: false}
    // _ = json.Marshal(params)
    // _ = json.Marshal(result)
}
```

--------------------------------

## 继续思考会话 - 添加步骤（JSON）

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/sequentialthinking/README.md

将下一个思考或分析步骤添加到正在进行的思考会话中。它需要会话 ID 和思考内容。

```json
{
  "method": "tools/call", 
  "params": {
    "name": "continue_thinking",
    "arguments": {
      "sessionId": "architecture_design",
      "thought": "First, I need to identify the core business domains and their boundaries to determine service decomposition."
    }
  }
}
```

--------------------------------

## Go：Client CallTool 签名

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了客户端 `CallTool` 方法的签名，该方法允许客户端使用类型化参数调用工具，并期望原始 JSON 消息作为参数。

```Go
func (cs *ClientSession) CallTool(context.Context, *CallToolParams[json.RawMessage]) (*CallToolResult, error)
```

--------------------------------

## Go SDK：用于服务器端 Stdio 通信的 StdioTransport

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了 stdio 传输服务器端 `StdioTransport`。此传输通过绑定到 os.Stdin 和 os.Stdout 来连接，并通过换行符分隔的 JSON 进行通信。

```Go
// A StdioTransport is a [Transport] that communicates using newline-delimited
// JSON over stdin/stdout.
type StdioTransport struct { }

func (t *StdioTransport) Connect(context.Context) (Connection, error)
```

--------------------------------

## 用于 MCP 的 Streamable 客户端传输

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

促进客户端连接到 Streamable MCP 端点，并透明地处理重新连接。它需要端点 URL，并允许配置 HTTP 客户端和重新连接行为。

```Go
package mcp

import (
	"context"
	"net/http"
)

// The streamable client handles reconnection transparently to the user.
type StreamableClientTransport struct {
	Endpoint         string
	HTTPClient       *http.Client
	ReconnectOptions *StreamableReconnectOptions
}

func (*StreamableClientTransport) Connect(context.Context) (Connection, error)
```