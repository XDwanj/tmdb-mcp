# Components

基于分层架构，以下是核心组件及其职责定义：

## Main Application (cmd/tmdb-mcp/main.go)

**Responsibility**: 应用程序入口，负责初始化所有组件并启动服务

**Key Interfaces**:
- `func main()` - 程序入口函数
- 信号处理（SIGINT/SIGTERM）- 优雅关闭

**Dependencies**: Config, Logger, TMDB Client, MCP Server, HTTP Server

**Technology Stack**: Go 标准库（`os`, `os/signal`, `context`）

**Implementation Details**:
```go
func main() {
    // 1. 加载配置（Viper）
    config := loadConfig()

    // 2. 初始化日志（Zap）
    logger := initLogger(config.Logging)

    // 3. 创建 TMDB Client（Resty + Rate Limiter）
    tmdbClient := tmdb.NewClient(config.TMDB, logger)

    // 4. 创建 MCP Server（MCP SDK）
    mcpServer := mcp.NewServer(tmdbClient, logger)

    // 5. 根据模式启动服务
    if config.Server.Mode == "stdio" || config.Server.Mode == "both" {
        go mcpServer.ServeStdio()
    }
    if config.Server.Mode == "sse" || config.Server.Mode == "both" {
        go serveSSE(mcpServer, config.Server.SSE, logger)
    }

    // 6. 等待信号并优雅关闭
    waitForShutdown(mcpServer, logger)
}
```

## Config Manager (internal/config)

**Responsibility**: 管理配置加载、验证和优先级控制（CLI > ENV > File）

**Key Interfaces**:
- `func LoadConfig() (*Config, error)` - 加载配置
- `func ValidateConfig(*Config) error` - 验证配置有效性
- `func GenerateSSEToken() (string, error)` - 生成 SSE Token

**Dependencies**: Viper, crypto/rand

**Technology Stack**: `github.com/spf13/viper`, `crypto/rand`, `os`

**Configuration Priority**:
1. 命令行 flags（最高优先级）
2. 环境变量（`TMDB_API_KEY`, `SSE_TOKEN` 等）
3. 配置文件 `~/.tmdb-mcp/config.yaml`（默认优先级）

**Validation Rules**:
- TMDB API Key 必须存在且非空
- SSE 模式启用时，Token 必须存在或自动生成
- Rate Limit 必须 > 0（默认 40）
- Server Mode 必须是 "stdio", "sse", "both" 之一

## TMDB API Client (internal/tmdb)

**Responsibility**: 封装 TMDB API v3 调用，处理 HTTP 请求、错误和响应解析

**Key Interfaces**:
- `func Search(ctx context.Context, query string, page int) ([]SearchResult, error)`
- `func GetMovieDetails(ctx context.Context, id int) (*MovieDetails, error)`
- `func GetTVDetails(ctx context.Context, id int) (*TVDetails, error)`
- `func GetPersonDetails(ctx context.Context, id int) (*PersonDetails, error)`
- `func DiscoverMovies(ctx context.Context, params DiscoverMoviesParams) ([]MovieResult, error)`
- `func DiscoverTV(ctx context.Context, params DiscoverTVParams) ([]TVResult, error)`
- `func GetTrending(ctx context.Context, mediaType, timeWindow string, page int) ([]TrendingResult, error)`
- `func GetMovieRecommendations(ctx context.Context, id, page int) ([]MovieResult, error)`
- `func GetTVRecommendations(ctx context.Context, id, page int) ([]TVResult, error)`

**Dependencies**: Resty, Rate Limiter, Logger

**Technology Stack**: `github.com/go-resty/resty/v2`, `golang.org/x/time/rate`, Zap

**Implementation Pattern**:
```go
type Client struct {
    httpClient  *resty.Client
    rateLimiter *rate.Limiter
    apiKey      string
    language    string
    logger      *zap.Logger
}

func (c *Client) Search(ctx context.Context, query string, page int) ([]SearchResult, error) {
    // 1. 等待速率限制
    if err := c.rateLimiter.Wait(ctx); err != nil {
        return nil, err
    }

    // 2. 发起 HTTP 请求
    resp, err := c.httpClient.R().
        SetContext(ctx).
        SetQueryParam("query", query).
        SetQueryParam("page", strconv.Itoa(page)).
        Get("/search/multi")

    // 3. 错误处理（401/404/429）
    if err := c.handleError(resp); err != nil {
        return nil, err
    }

    // 4. 解析响应
    var result SearchResponse
    if err := json.Unmarshal(resp.Body(), &result); err != nil {
        return nil, err
    }

    return result.Results, nil
}
```

**Error Handling Strategy**:
- **401 Unauthorized**: 立即返回错误，记录 ERROR 日志
- **404 Not Found**: 返回空结果，记录 INFO 日志
- **429 Rate Limit**: 解析 `Retry-After` header，等待后重试（最多 3 次）
- **Network Timeout**: 返回错误，记录 WARN 日志
- **5xx Server Error**: 返回错误，记录 ERROR 日志

## Rate Limiter (internal/ratelimit)

**Responsibility**: 控制 TMDB API 调用频率，防止触发 429 错误

**Key Interfaces**:
- `func NewLimiter(ratePerSecond float64, burst int) *rate.Limiter` - 创建限制器
- `func Wait(ctx context.Context) error` - 等待获取令牌

**Dependencies**: golang.org/x/time/rate

**Technology Stack**: `golang.org/x/time/rate`

**Configuration**:
- **Rate**: 40 requests / 10 seconds = 4 req/s
- **Burst**: 40（允许短时突发）
- **Algorithm**: Token Bucket

**Usage**:
```go
// 创建限制器：每秒 4 个请求，突发 40
limiter := rate.NewLimiter(rate.Every(10*time.Second/40), 40)

// 每次 API 调用前等待
if err := limiter.Wait(ctx); err != nil {
    return err
}
```

## MCP Server (internal/mcp)

**Responsibility**: MCP 协议实现，工具注册和调度，stdio 和 SSE 模式支持

**Key Interfaces**:
- `func NewServer(tmdbClient *tmdb.Client, logger *zap.Logger) *mcp.Server` - 创建服务器
- `func ServeStdio() error` - 启动 stdio 模式
- `func GetHTTPHandler() http.Handler` - 获取 SSE HTTP Handler

**Dependencies**: MCP SDK, TMDB Client, Tools, Logger

**Technology Stack**: `github.com/modelcontextprotocol/go-sdk`, Zap

**Tool Registration**:
```go
func NewServer(tmdbClient *tmdb.Client, logger *zap.Logger) *Server {
    // Create server options
    opts := &mcp.ServerOptions{
        Instructions: "TMDB Movie Database MCP Server - provides tools for searching and retrieving movie information",
    }

    // Create MCP server with implementation info
    mcpServer := mcp.NewServer(&mcp.Implementation{
        Name:    "tmdb-mcp",
        Version: "1.0.0",
    }, opts)

    // Create and register search tool
    searchTool := tools.NewSearchTool(tmdbClient, logger)
    mcp.AddTool(mcpServer, &mcp.Tool{
        Name:        searchTool.Name(),
        Description: searchTool.Description(),
    }, searchTool.Handler())

    // 其他 5 个工具同理注册...

    return &Server{
        mcpServer:  mcpServer,
        tmdbClient: tmdbClient,
        logger:     logger,
    }
}
```

**设计说明**:
- **Server 层职责**: 仅负责工具注册,不包含业务逻辑(50 行代码)
- **Tools 层职责**: 提供 `Handler()` 工厂方法,返回符合 MCP SDK 签名的处理函数
- **闭包模式**: `Handler()` 返回的闭包自动捕获 `tmdbClient` 和 `logger` 依赖
- **SDK 自动处理**: jsonschema 标签自动生成 InputSchema,SDK 自动解析和验证参数

## MCP Tools (internal/tools)

**Responsibility**: 实现 6 个 MCP 工具，处理参数验证和结果转换

**Key Interfaces**:
- `func (t *SearchTool) Name() string` - 返回工具名称
- `func (t *SearchTool) Description() string` - 返回工具描述
- `func (t *SearchTool) Handler() func(...)` - 返回处理函数（工厂方法）

**Dependencies**: TMDB Client, Logger

**Technology Stack**: MCP SDK, Zap

**Tool List**:
1. **SearchTool** - 搜索电影/电视剧/人物
2. **GetDetailsTool** - 获取详情（movie/tv/person）
3. **DiscoverMoviesTool** - 发现电影
4. **DiscoverTVTool** - 发现电视剧
5. **GetTrendingTool** - 获取热门内容
6. **GetRecommendationsTool** - 获取推荐

**Implementation Pattern**:
```go
type SearchTool struct {
    tmdbClient *tmdb.Client
    logger     *zap.Logger
}

func (t *SearchTool) Name() string {
    return "search"
}

func (t *SearchTool) Description() string {
    return "Search for movies, TV shows, and people on TMDB using a query string"
}

// Handler 返回一个符合 MCP SDK 签名的处理函数
// 这允许工具在 MCP Server 中注册,同时将业务逻辑封装在 SearchTool 中
func (t *SearchTool) Handler() func(context.Context, *mcp.CallToolRequest, SearchParams) (*mcp.CallToolResult, SearchResponse, error) {
    return func(ctx context.Context, req *mcp.CallToolRequest, params SearchParams) (*mcp.CallToolResult, SearchResponse, error) {
        // 1. 设置默认值
        if params.Page == 0 {
            params.Page = 1
        }

        // 2. 记录请求日志
        t.logger.Info("Search request received",
            zap.String("query", params.Query),
            zap.Int("page", params.Page),
        )

        // 3. 调用 TMDB Client (参数验证在 Client 层完成)
        results, err := t.tmdbClient.Search(ctx, params.Query, params.Page)
        if err != nil {
            t.logger.Error("Search failed",
                zap.Error(err),
                zap.String("query", params.Query),
            )
            return nil, SearchResponse{}, err
        }

        t.logger.Info("Search completed",
            zap.String("query", params.Query),
            zap.Int("results", len(results.Results)),
        )

        // 4. 返回空的 CallToolResult 和结构化响应
        return &mcp.CallToolResult{}, SearchResponse{Results: results.Results}, nil
    }
}
```

**设计说明**:
- **Handler() 工厂方法**: 返回闭包函数,自动捕获 `tmdbClient` 和 `logger` 依赖
- **类型安全参数**: MCP SDK 自动从 SearchParams 的 jsonschema 标签生成 InputSchema 并验证参数
- **职责分离**: Tools 层只处理日志和默认值,参数验证在 Client 层统一完成
- **闭包模式**: 避免全局变量,每个工具实例独立维护依赖

## HTTP Server (cmd/tmdb-mcp/main.go)

**Responsibility**: 在主程序中直接提供 SSE HTTP 服务，集成 MCP SDK 的 `SSEHTTPHandler`

**Key Functions**:
- `func RunSSEModeServer(ctx, mcpServer, config, logger)` - 启动 SSE 模式服务器
- `func RunStdioModeServer(ctx, mcpServer, logger)` - 启动 stdio 模式服务器
- `func RunBothModeServer(ctx, mcpServer, config, logger)` - 同时启动两种模式
- `func healthHandler(w, r)` - 健康检查端点

**Dependencies**: MCP Server, net/http, AuthMiddleware, Logger

**Technology Stack**: `net/http` (标准库), MCP SDK 的 `SSEHTTPHandler`, Zap

**Endpoints**:
- `GET /mcp/sse` - MCP over SSE 连接（需要认证）
- `GET /health` - 健康检查（无需认证）

**Server Setup**:
```go
func RunSSEModeServer(ctx context.Context, mcpServer *mcp.Server, cfg *config.Config, log *zap.Logger) {
    // 设置 SSE 处理器
    handler := mcpServer.GetSSEHandler()
    handler = middleware.AuthMiddleware(cfg.Server.SSE.Token, handler)

    // 设置路由
    mux := http.NewServeMux()
    mux.HandleFunc("/health", healthHandler)
    mux.Handle("/mcp/sse", handler)

    // 启动服务器（阻塞）
    addr := fmt.Sprintf("%s:%d", cfg.Server.SSE.Host, cfg.Server.SSE.Port)
    log.Info("Starting HTTP server", zap.String("addr", addr))
    if err := http.ListenAndServe(addr, mux); err != nil {
        log.Fatal("HTTP server failed", zap.Error(err))
    }
}
```

**Authentication Middleware** (`internal/server/middleware/auth.go`):
```go
func AuthMiddleware(token string, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        auth := r.Header.Get("Authorization")
        if auth != "Bearer "+token {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte(`{"error":"unauthorized"}`))
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

**Design Notes**:
- HTTP 服务器逻辑直接在 `main.go` 中实现，无独立组件
- 使用 `http.ServeMux` 和 `http.ListenAndServe()` 启动
- 认证中间件封装 SSE handler
- 无优雅关闭机制（简化设计）
- AuthMiddleware 使用字符串直接比较（MVP 阶段，生产环境建议使用 `crypto/subtle.ConstantTimeCompare`）

## Logger (internal/logger)

**Responsibility**: 提供结构化日志功能，记录关键事件和错误

**Key Interfaces**:
- `func InitLogger(config LogConfig) (*zap.Logger, error)` - 初始化 logger
- Zap logger 标准接口（`Info`, `Error`, `Debug`, `Warn`）

**Dependencies**: Zap

**Technology Stack**: `go.uber.org/zap`

**Configuration**:
- **Development Mode**: 使用 `zap.NewDevelopment()`，输出到 console，彩色输出
- **Production Mode**: 使用 `zap.NewProduction()`，JSON 格式，结构化输出
- **Log Level**: 从配置文件读取（debug/info/warn/error）

**Usage**:
```go
logger.Info("Starting TMDB MCP Service",
    zap.String("mode", config.Server.Mode),
    zap.Int("port", config.Server.SSE.Port),
)

logger.Error("TMDB API call failed",
    zap.String("endpoint", "/search/multi"),
    zap.Error(err),
    zap.Duration("duration", elapsed),
)
```

## Component Diagrams

```mermaid
graph TB
    subgraph "Application Layer"
        Main[Main<br/>程序入口]
    end

    subgraph "Protocol Layer"
        StdioTransport[stdio Transport<br/>MCP SDK]
        SSETransport[SSE Transport<br/>HTTP + MCP SDK]
        AuthMW[Auth Middleware<br/>Bearer Token]
    end

    subgraph "Service Layer"
        MCPServer[MCP Server<br/>工具注册和调度]
        Tools[MCP Tools x6<br/>search, details, discover...]
    end

    subgraph "Data Access Layer"
        TMDBClient[TMDB Client<br/>Resty HTTP]
        RateLimiter[Rate Limiter<br/>Token Bucket]
    end

    subgraph "Cross-Cutting"
        Config[Config Manager<br/>Viper]
        Logger[Logger<br/>Zap]
    end

    Main --> Config
    Main --> Logger
    Main --> MCPServer
    Main --> StdioTransport
    Main --> SSETransport

    StdioTransport --> MCPServer
    SSETransport --> AuthMW
    AuthMW --> MCPServer

    MCPServer --> Tools
    Tools --> TMDBClient
    TMDBClient --> RateLimiter

    Config -.-> Main
    Config -.-> TMDBClient
    Logger -.-> Main
    Logger -.-> MCPServer
    Logger -.-> Tools
    Logger -.-> TMDBClient
```

## Component Interaction Sequence

**stdio 模式请求流程**:
```mermaid
sequenceDiagram
    participant Claude as Claude Code
    participant Stdio as stdio Transport
    participant MCP as MCP Server
    participant Tool as Search Tool
    participant Client as TMDB Client
    participant RL as Rate Limiter
    participant TMDB as TMDB API

    Claude->>Stdio: JSON-RPC request<br/>tools/call "search"
    Stdio->>MCP: Parse and route
    MCP->>Tool: Call(params)
    Tool->>Client: Search(query, page)
    Client->>RL: Wait()
    RL-->>Client: OK (token granted)
    Client->>TMDB: GET /search/multi
    TMDB-->>Client: 200 OK + JSON
    Client-->>Tool: []SearchResult
    Tool-->>MCP: Result
    MCP-->>Stdio: JSON-RPC response
    Stdio-->>Claude: Results
```

**SSE 模式请求流程**:
```mermaid
sequenceDiagram
    participant Browser as Browser/App
    participant Main as main.go
    participant Auth as Auth Middleware
    participant SSE as SSE Handler
    participant MCP as MCP Server
    participant Tool as Tool
    participant Client as TMDB Client

    Browser->>Main: GET /mcp/sse<br/>Authorization: Bearer TOKEN
    Main->>Auth: Verify token
    Auth->>Auth: String Compare
    Auth-->>Main: OK
    Main->>SSE: Handle SSE connection
    SSE->>MCP: MCP request
    MCP->>Tool: Call tool
    Tool->>Client: TMDB API call
    Client-->>Tool: Result
    Tool-->>MCP: Result
    MCP-->>SSE: MCP response
    SSE-->>Browser: SSE event
```

---
