# Model Context Protocol Go SDK 文档

## 📦 库信息

- **Library ID**: `/modelcontextprotocol/go-sdk`
- **GitHub**: https://github.com/modelcontextprotocol/go-sdk
- **包路径**: `github.com/modelcontextprotocol/go-sdk/mcp`
- **描述**: Model Context Protocol (MCP) 的官方 Go SDK
- **Trust Score**: 7.8
- **代码示例数**: 81 个
- **可用版本**: v0.2.0, v0.4.0

## 🎯 概述

Model Context Protocol Go SDK 提供了构建和使用 MCP 客户端和服务器的 API，包括 JSON Schema 和 JSON RPC 实现。它是 MCP 协议的官方 Go 语言实现。

## 📚 包结构

```go
github.com/modelcontextprotocol/go-sdk/mcp           // 核心 MCP API
github.com/modelcontextprotocol/go-sdk/jsonschema    // JSON Schema 支持
github.com/modelcontextprotocol/go-sdk/internal/jsonrpc2  // JSON-RPC 2.0 实现
```

## 🚀 快速开始

### 安装

```bash
go get github.com/modelcontextprotocol/go-sdk
```

### 创建最简单的 MCP 服务器

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
	// 创建服务器
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)

	// 添加工具
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)

	// 运行服务器（通过 stdin/stdout）
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
```

### 创建 MCP 客户端

```go
package main

import (
	"context"
	"log"
	"os/exec"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()

	// 创建客户端
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

	// 通过 stdin/stdout 连接到服务器
	transport := &mcp.CommandTransport{Command: exec.Command("myserver")}
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// 调用工具
	params := &mcp.CallToolParams{
		Name:      "greet",
		Arguments: map[string]any{"name": "you"},
	}
	res, err := session.CallTool(ctx, params)
	if err != nil {
		log.Fatalf("CallTool failed: %v", err)
	}

	// 处理响应
	for _, c := range res.Content {
		log.Print(c.(*mcp.TextContent).Text)
	}
}
```

## 🔌 传输方式

MCP Go SDK 支持多种传输协议：

### 1. Stdio Transport（标准输入/输出）

这是最常见的传输方式，适用于本地进程间通信。

```go
// 服务器端
server := mcp.NewServer(&mcp.Implementation{Name: "myserver", Version: "v1.0.0"}, nil)
err := server.Run(context.Background(), &mcp.StdioTransport{})

// 客户端
transport := &mcp.CommandTransport{Command: exec.Command("myserver")}
session, err := client.Connect(ctx, transport, nil)
```

**运行方式**：
```bash
go run .  # 默认使用 stdio
```

### 2. HTTP Transport

适用于网络服务，提供标准的 HTTP API。

```go
// 启动 HTTP 服务器
go run main.go -host 0.0.0.0 -port 8080 server

// 客户端连接
go run main.go client

// 自定义端口
go run main.go -host 0.0.0.0 -port 9000 server
```

### 3. SSE (Server-Sent Events)

适用于需要服务器主动推送的场景。

```go
// 服务器端
type SSEHTTPHandler struct { /* unexported fields */ }

func NewSSEHTTPHandler(getServer func(request *http.Request) *Server) *SSEHTTPHandler

handler := mcp.NewSSEHTTPHandler(func(req *http.Request) *mcp.Server {
	return server
})

func (*SSEHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request)
func (*SSEHTTPHandler) Close() error

// 客户端
type SSEClientTransport struct {
	Endpoint   string
	HTTPClient *http.Client
}

transport := &mcp.SSEClientTransport{
	Endpoint: "http://localhost:8080",
}
```

**SSE Server Transport**:
```go
// 表示一个逻辑 SSE 会话，通过挂起的 GET 请求建立
type SSEServerTransport struct {
    Endpoint string
    Response http.ResponseWriter
}

func (*SSEServerTransport) ServeHTTP(w http.ResponseWriter, req *http.Request)
func (*SSEServerTransport) Connect(context.Context) (Connection, error)
```

### 4. Streamable Transport

支持会话恢复和事件存储的流式传输。

```go
// 服务器端
type StreamableServerTransport struct {
	SessionID  string      // 会话 ID
	EventStore EventStore  // 事件存储，支持流恢复
}

handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
	return server
}, nil)

func (*StreamableServerTransport) ServeHTTP(w http.ResponseWriter, req *http.Request)
func (*StreamableServerTransport) Connect(context.Context) (Connection, error)

// 客户端
type StreamableClientTransport struct {
	Endpoint         string
	HTTPClient       *http.Client
	ReconnectOptions *StreamableReconnectOptions
}

// 客户端会透明地处理重连
func (*StreamableClientTransport) Connect(context.Context) (Connection, error)
```

**Streamable HTTP Handler**:
```go
type StreamableHTTPHandler struct { /* unexported fields */ }

func NewStreamableHTTPHandler(getServer func(request *http.Request) *Server) *StreamableHTTPHandler

func (*StreamableHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request)
func (*StreamableHTTPHandler) Close() error
```

### 5. 自定义 Transport

#### InMemoryTransport
进程内通信，使用换行符分隔的 JSON。

```go
type InMemoryTransport struct { /* ... */ }

// 创建两个相互连接的 InMemoryTransports
func NewInMemoryTransports() (*InMemoryTransport, *InMemoryTransport)
```

#### LoggingTransport
中间件传输，记录 RPC 详情到 io.Writer。

```go
type LoggingTransport struct {
	Delegate Transport
	Writer   io.Writer
}

// 示例：记录到 stdout
serverTransport, clientTransport := NewInMemoryTransports()
logger := os.Stdout
loggingTransport := &LoggingTransport{
	Delegate: serverTransport,
	Writer:   logger
}
```

## 🔧 工具（Tools）管理

### 添加工具

```go
// 定义工具的输入输出结构
type AddParams struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// 工具处理函数
func addHandler(ctx context.Context, req *mcp.ServerRequest[*mcp.CallToolParamsFor[AddParams]]) (*mcp.CallToolResultFor[int], error) {
	return &mcp.CallToolResultFor[int]{
		StructuredContent: req.Params.Arguments.X + req.Params.Arguments.Y
	}, nil
}

// 添加工具
mcp.AddTool(server, &mcp.Tool{
	Name:        "add",
	Description: "add numbers"
}, addHandler)

// 也可以使用方法
server.AddTool(&mcp.Tool{
	Name:        "subtract",
	Description: "subtract numbers"
}, subHandler)
```

### 工具定义

```go
type Tool struct {
	Annotations *ToolAnnotations   `json:"annotations,omitempty"`
	Description string             `json:"description,omitempty"`
	InputSchema *jsonschema.Schema `json:"inputSchema"`
	Name        string             `json:"name"`
}

type CallToolParamsFor[In any] struct {
	Meta      Meta   `json:"_meta,omitempty"`
	Arguments In     `json:"arguments,omitempty"`
	Name      string `json:"name"`
}
```

### 移除工具

```go
server.RemoveTools("add", "subtract")
```

### 泛型工具处理器

SDK 提供了类型安全的泛型处理器：

```go
// 泛型工具处理器
func AddTool[In, Out any](s *Server, t *Tool, h ToolHandlerFor[In, Out])

// MethodHandler 处理 MCP 消息
// 对于方法：恰好一个返回值必须为 nil
// 对于通知：两个都必须为 nil
type MethodHandler func(ctx context.Context, method string, req Request) (result Result, err error)
```

## 📦 资源（Resources）管理

### 添加资源

```go
// 添加单个资源
server.AddResource(&mcp.Resource{
	URI:         "file:///example.txt",
	Name:        "Example",
	Description: "An example resource",
}, resourceHandler)

// 添加资源模板（支持参数化 URI）
server.AddResourceTemplate(&mcp.ResourceTemplate{
	URITemplate: "file:///{path}",
	Name:        "File Resource",
	Description: "Dynamic file resource",
}, templateHandler)
```

### 移除资源

```go
// 移除特定资源
server.RemoveResources("file:///example.txt", "file:///another.txt")

// 移除资源模板
server.RemoveResourceTemplates("file:///{path}")
```

### 资源订阅

客户端可以订阅资源更新：

```go
// 客户端订阅资源
err := session.Subscribe(ctx, &mcp.SubscribeParams{
	URI: "resource://example",
})

// 取消订阅
err := session.Unsubscribe(ctx, &mcp.UnsubscribeParams{
	URI: "resource://example",
})

// 客户端处理资源更新通知
type ClientOptions struct {
	ResourceUpdatedHandler func(context.Context, *ClientSession, *ResourceUpdatedNotificationParams)
}
```

### 服务器端订阅处理

```go
type ServerOptions struct {
	// 客户端订阅资源时调用
	SubscribeHandler func(context.Context, *ServerRequest[*SubscribeParams]) error

	// 客户端取消订阅时调用
	UnsubscribeHandler func(context.Context, *ServerRequest[*UnsubscribeParams]) error
}

// 通知客户端资源已更新
err := server.ResourceUpdated(ctx, &mcp.ResourceUpdatedNotificationParams{
	URI: "resource://example",
})
```

## 💬 提示（Prompts）管理

### 添加提示

```go
type codeReviewArgs struct {
	Code string `json:"code"`
}

func codeReviewHandler(context.Context, *ServerSession, *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	// 处理提示逻辑
	return &mcp.GetPromptResult{
		Messages: []mcp.Message{
			{
				Role: "user",
				Content: mcp.TextContent{
					Text: "Please review this code...",
				},
			},
		},
	}, nil
}

server.AddPrompt(
	&mcp.Prompt{
		Name:        "code_review",
		Description: "review code",
	},
	codeReviewHandler,
)
```

### 移除提示

```go
server.RemovePrompts("code_review")
```

## 🔐 认证和授权

### JWT 认证

```go
// 实现 JWT 验证器
func jwtVerifier(ctx context.Context, tokenString string) (*auth.TokenInfo, error) {
	// JWT 验证逻辑
	// 成功返回: TokenInfo
	// 失败返回: auth.ErrInvalidToken
	return &auth.TokenInfo{
		UserID: "user123",
		Scopes: []string{"read", "write"},
	}, nil
}

// 创建认证中间件
authMiddleware := auth.RequireBearerToken(jwtVerifier, &auth.RequireBearerTokenOptions{
	Scopes: []string{"read", "write"},
})

// 创建 MCP handler
handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
	return server
}, nil)

// 应用认证中间件
authenticatedHandler := authMiddleware(customMiddleware(handler))
```

### 在工具中访问认证信息

```go
func MyTool(ctx context.Context, req *mcp.CallToolRequest, args MyArgs) (*mcp.CallToolResult, any, error) {
	// 从请求中提取认证信息
	userInfo := req.Extra.TokenInfo

	// 检查权限范围
	if !slices.Contains(userInfo.Scopes, "read") {
		return nil, nil, fmt.Errorf("insufficient permissions: read scope required")
	}

	// 执行工具逻辑
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "Tool executed successfully"},
		},
	}, nil, nil
}
```

### 测试认证

```bash
# 生成 JWT token
curl 'http://localhost:8080/generate-token?user_id=alice&scopes=read,write'

# 使用 JWT 调用工具
curl -H 'Authorization: Bearer <token>' \
     -H 'Content-Type: application/json' \
     -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"say_hi","arguments":{}}}' \
     http://localhost:8080/mcp/jwt

# 生成 API Key
curl -X POST 'http://localhost:8080/generate-api-key?user_id=bob&scopes=read'

# 使用 API Key 调用工具
curl -H 'Authorization: Bearer <api_key>' \
     -H 'Content-Type: application/json' \
     -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"get_user_info","arguments":{"user_id":"test"}}}' \
     http://localhost:8080/mcp/apikey

# 测试权限范围限制
curl -H 'Authorization: Bearer <token_with_write_scope>' \
     -H 'Content-Type: application/json' \
     -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"create_resource","arguments":{"name":"test","description":"test resource","content":"test content"}}}' \
     http://localhost:8080/mcp/jwt
```

### 中间件组合

```go
// 组合认证中间件和自定义中间件
authenticatedHandler := authMiddleware(customMiddleware(mcpHandler))
```

## 🔄 中间件系统

### 中间件定义

```go
// Middleware 是从 MethodHandlers 到 MethodHandlers 的函数
type Middleware func(MethodHandler) MethodHandler

// 添加发送中间件
func (c *Client) AddSendingMiddleware(middleware ...Middleware)
func (s *Server) AddSendingMiddleware(middleware ...Middleware)

// 添加接收中间件
func (c *Client) AddReceivingMiddleware(middleware ...Middleware)
func (s *Server) AddReceivingMiddleware(middleware ...Middleware)
```

### 中间件执行顺序

中间件从右到左应用，第一个先执行：

```go
// AddMiddleware(m1, m2, m3) 会将处理器增强为：
// m1(m2(m3(handler)))
server.AddSendingMiddleware(m1, m2, m3)
```

## 📊 日志处理

### 使用标准库 slog

```go
type LoggingHandler struct {}

type LoggingHandlerOptions struct {
	// "logger" 字段的值
	LoggerName string
	// 限制日志消息发送速率
	// 如果为零，则无速率限制
	MinInterval time.Duration
}

// 创建日志处理器
func NewLoggingHandler(ss *ServerSession, opts *LoggingHandlerOptions) *LoggingHandler

// 服务器选项
type ServerOptions struct {
	// "logger" 字段的值
	LoggerName string
	// 日志通知不会比此持续时间更频繁地发送
	LoggingInterval time.Duration
}
```

### 客户端日志处理

```go
type ClientOptions struct {
	// 从服务器接收日志消息时调用
	LoggingMessageHandler func(context.Context, *ClientSession, *LoggingMessageParams)
}
```

## 📈 进度通知

### 请求进度更新

```go
type XXXParams struct {
	Meta Meta
	// ... 其他参数
}

type Meta struct {
	Data          map[string]any
	ProgressToken any // string 或 int
}
```

### 发送进度通知

```go
// 客户端发送进度
func (*ClientSession) NotifyProgress(context.Context, *ProgressNotification)

// 服务器发送进度
func (*ServerSession) NotifyProgress(context.Context, *ProgressNotification)
```

## 🔔 通知处理

### 工具列表变更通知

```go
type ClientOptions struct {
	// 工具列表变更时调用
	ToolListChangedHandler func(context.Context, *ClientRequest[*ToolListChangedParams])

	// 提示列表变更时调用
	PromptListChangedHandler func(context.Context, *ClientRequest[*PromptListChangedParams])

	// 资源和资源模板列表变更时调用
	ResourceListChangedHandler func(context.Context, *ClientRequest[*ResourceListChangedParams])
}
```

## 🌲 Roots 管理（客户端）

```go
// 添加 roots（替换相同 URI 的 roots）
// 并通知所有已连接的服务器
func (*Client) AddRoots(roots ...*Root)

// 移除指定 URI 的 roots
// 如果列表发生变化，则通知已连接的服务器
// 移除不存在的 root 不会报错
func (*Client) RemoveRoots(uris ...string)
```

## 🤖 AI 采样（Sampling）

### 客户端处理采样请求

```go
type ClientOptions struct {
	// 服务器调用 CreateMessage 时调用
	CreateMessageHandler func(context.Context, *ClientSession, *CreateMessageParams) (*CreateMessageResult, error)
}
```

## 🔄 自动补全

### 服务器处理补全请求

```go
type ServerOptions struct {
	// 客户端发送补全请求时调用
	CompletionHandler func(context.Context, *ServerRequest[*CompleteParams]) (*CompleteResult, error)
}
```

## 📄 分页支持

```go
type ServerOptions struct {
	// 定义每页返回的项目数
	PageSize int
}
```

## 🎯 规范方法签名

所有 RPC 方法遵循统一的签名格式：

```go
// 标准规范方法签名
func (*ClientSession) ListTools(context.Context, *ListToolsParams) (*ListToolsResult, error)

// 通用签名模式
func (*Session) MethodName(ctx context.Context, params *MethodParams) (*MethodResult, error)
```

## 🧪 测试

```bash
# 运行所有测试
go test -v

# 运行基准测试
go test -bench=.

# 生成覆盖率报告
go test -cover
```

## 🚀 运行示例

### Sequential Thinking Server

```bash
# Stdio 模式（默认）
go run .

# HTTP 模式
go run . -http :8080
```

### HTTP Server/Client Example

```bash
# 启动服务器
go run main.go server

# 运行客户端
go run main.go client
```

### 认证中间件示例

```bash
cd examples/server/auth-middleware
go mod tidy
go run main.go
```

## 🔗 集成到 Claude Code

```bash
# 添加 HTTP MCP 服务器到 Claude Code
claude mcp add -t http timezone http://localhost:8080
```

### 在 Claude Code 中使用示例

```bash
# 查询支持的时区
> what timezones do you support?

⏺ The timezone tool supports three US cities:
  - NYC (Eastern Time)
  - SF (Pacific Time)
  - Boston (Eastern Time)

# 查询当前时区
> what's the timezone

⏺ I'll get the current time in a major US city for you.

⏺ timezone - cityTime (MCP)(city: "nyc")
  ⎿ The current time in New York City is 7:30:16 PM EDT on Wednesday, July 23, 2025

⏺ The current timezone is EDT (Eastern Daylight Time), and it's 7:30 PM on Wednesday, July 23, 2025.
```

## 🛠️ 高级特性

### Sequential Thinking Tools

SDK 包含一个顺序思考工具的示例实现：

#### 开始思考会话
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

#### 继续思考 - 添加步骤
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

#### 继续思考 - 修订步骤
```json
{
  "method": "tools/call",
  "params": {
    "name": "continue_thinking",
    "arguments": {
      "sessionId": "architecture_design",
      "thought": "Actually, before identifying domains, I should analyze the current system's pain points and requirements.",
      "reviseStep": 1
    }
  }
}
```

#### 继续思考 - 创建分支
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

#### 完成思考
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

#### 回顾思考会话
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

## 📝 贡献指南

### 版权头

所有 Go 源文件必须包含标准版权头：

```go
// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
```

### 开发设置

使用 `go.work` 进行多模块开发：

```bash
# 初始化工作区，同时测试 SDK 更改和本地项目
go work init ./project ./go-sdk
```

## 🔍 最佳实践

1. **选择合适的传输方式**
   - 本地工具：使用 Stdio Transport
   - Web 服务：使用 HTTP Transport
   - 实时推送：使用 SSE Transport
   - 需要会话恢复：使用 Streamable Transport

2. **错误处理**
   - 始终检查错误返回值
   - 使用 context 进行超时控制
   - 提供有意义的错误消息

3. **类型安全**
   - 使用泛型处理器获得类型安全
   - 利用 JSON Schema 标签自动生成模式
   - 定义清晰的输入输出结构

4. **安全性**
   - 实施适当的认证中间件
   - 使用 scope 进行细粒度权限控制
   - 验证所有输入参数

5. **性能优化**
   - 使用上下文进行请求取消
   - 实施日志速率限制
   - 考虑使用分页处理大型结果集

## 📚 相关资源

- **MCP 规范**: https://modelcontextprotocol.io/specification
- **MCP 文档**: https://github.com/modelcontextprotocol/docs
- **Go SDK GitHub**: https://github.com/modelcontextprotocol/go-sdk
- **MCP 初学者教程**: https://github.com/microsoft/mcp-for-beginners

## 🌐 其他语言 SDK

- **TypeScript SDK**: `/modelcontextprotocol/typescript-sdk`
- **Python SDK**: `/modelcontextprotocol/python-sdk`
- **C# SDK**: `/modelcontextprotocol/csharp-sdk`
- **Java SDK**: `/modelcontextprotocol/java-sdk`
- **Swift SDK**: `/modelcontextprotocol/swift-sdk`
- **Kotlin SDK**: `/modelcontextprotocol/kotlin-sdk`
- **Ruby SDK**: `/modelcontextprotocol/ruby-sdk`
- **PHP SDK**: `/modelcontextprotocol/php-sdk`

## 🎓 学习路径

1. **入门**: 从最简单的 Stdio 服务器开始
2. **进阶**: 实现 HTTP 传输和认证
3. **高级**: 使用中间件、订阅和流式传输
4. **专家**: 构建复杂的工具链和集成

---

*最后更新: 2025-10-10*
*基于 Context7 查询结果整理*
