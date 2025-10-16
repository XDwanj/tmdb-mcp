# Tech Stack

## Cloud Infrastructure

**Provider**: N/A - 自托管部署

**说明**:
- 本项目为 **自托管工具**，无需云服务商
- 用户在本地机器或自己的服务器上运行
- 部署方式：独立二进制文件或 Docker 容器
- 无云服务依赖，无运行时成本

**可选云部署**:
- 用户可自行选择任何云服务商（AWS, GCP, Azure, DigitalOcean 等）
- 仅需支持 Linux 容器的环境即可
- 建议配置：1 vCPU, 512MB RAM（轻量级部署）

## Technology Stack Table

基于 PRD Technical Assumptions 和精简原则，以下是最终技术栈选择：

| Category              | Technology                               | Version    | Purpose         | Rationale                                              |
| --------------------- | ---------------------------------------- | ---------- | --------------- | ------------------------------------------------------ |
| **Language**          | Go                                       | 1.21+      | 主开发语言      | 类型安全、高性能、优秀并发模型、编译为独立二进制       |
| **Runtime**           | Go Runtime                               | 1.21+      | 程序运行环境    | 跨平台支持（Linux/macOS/Windows）、静态链接、启动快    |
| **MCP SDK**           | `github.com/modelcontextprotocol/go-sdk` | 最新稳定版 | MCP 协议实现    | 官方 Go SDK，内置 stdio 和 SSE 支持，遵循规范          |
| **HTTP Client**       | `github.com/go-resty/resty/v2`           | v2.11.0+   | TMDB API 调用   | 链式 API、自动重试、中间件支持、超时控制               |
| **HTTP Server**       | `net/http` (标准库)                      | Go 1.21+   | SSE HTTP 服务器 | 标准库稳定可靠、零依赖、配合 MCP SDK 的 SSEHTTPHandler |
| **Rate Limiter**      | `golang.org/x/time/rate`                 | v0.5.0+    | API 速率限制    | 官方扩展库、Token Bucket 算法、并发安全                |
| **Logging**           | `go.uber.org/zap`                        | v1.26.0+   | 结构化日志      | 高性能、零分配、JSON 输出、日志级别控制                |
| **Configuration**     | `github.com/spf13/viper`                 | v1.18.0+   | 配置管理        | 多源支持（文件/ENV/CLI）、优先级控制、热重载           |
| **Testing**           | `testing` (标准库)                       | Go 1.21+   | 单元测试        | Go 原生测试框架，`go test` 命令                        |
| **Testing - Mocking** | `github.com/stretchr/testify`            | v1.8.4+    | Mock 和断言     | 丰富的断言、Mock 支持、与标准库兼容                    |
| **Security - Token**  | `crypto/rand` (标准库)                   | Go 1.21+   | SSE Token 生成  | 加密安全的随机数生成、标准库零依赖                     |
| **Security - Auth**   | `crypto/subtle` (标准库)                 | Go 1.21+   | Token 比对      | 常量时间比对，防止时序攻击                             |
| **Build Tool**        | `go build`                               | Go 1.21+   | 编译二进制      | Go 原生工具，无需 Makefile                             |
| **Formatting**        | `go fmt`                                 | Go 1.21+   | 代码格式化      | Go 官方格式化工具                                      |
| **Static Analysis**   | `go vet`                                 | Go 1.21+   | 静态检查        | Go 官方静态分析工具                                    |
| **Containerization**  | Docker                                   | 24.0+      | 容器化部署      | 多平台镜像、轻量 Alpine 基础镜像、Docker Compose 支持  |

## 关键技术决策说明

**1. 为什么选择 Go 1.21+？**
- **类型安全**: 编译时捕获错误，减少运行时故障
- **并发模型**: Goroutines 和 channels 天然支持并发处理 TMDB API 调用
- **部署简单**: 编译为单个二进制文件，无需运行时依赖
- **性能优异**: 满足 P95 < 500ms 的响应时间要求
- **启动快速**: 满足 < 2 秒启动时间要求

**2. 为什么选择标准库 `net/http` 而非 Gin/Echo？**
- **精简原则**: PRD 明确要求避免不必要的框架
- **MCP SDK 兼容**: MCP SDK 的 `SSEHTTPHandler` 实现了 `http.Handler` 接口，直接兼容标准库
- **零依赖**: 标准库稳定可靠，无版本冲突风险
- **足够功能**: 仅需 `/mcp/sse` 和 `/health` 两个端点，标准库完全满足

**3. 为什么选择 Resty 而非标准库 `net/http` 客户端？**
- **链式 API**: 更简洁的代码风格（`client.R().SetQueryParam().Get()`）
- **自动重试**: 内置重试逻辑，处理 TMDB API 临时故障
- **中间件支持**: 方便统一添加日志、指标记录，**支持 OnBeforeRequest 统一处理 rate limiting**
- **超时控制**: 更好的超时和取消控制

**4. 为什么选择 Viper 配置管理？**
- **多源支持**: 文件、环境变量、命令行 flags 三种来源
- **优先级控制**: CLI > ENV > File，符合 12-factor app 原则
- **生态成熟**: Go 社区广泛使用，文档完善

**5. 为什么不使用数据库？**
- **纯 API 转发**: 所有数据来自 TMDB API，无需持久化
- **无状态设计**: 每次请求独立，符合 MCP 协议特性
- **简化运维**: 用户无需管理数据库，降低部署门槛

**6. 为什么不使用缓存（Redis/Memcached）？**
- **实时性要求**: PRD 明确要求实时数据，避免过期内容
- **MVP 范围**: 第一版聚焦核心功能，缓存属于优化范畴
- **TMDB 速率限制**: 40 req/10s 足够支持单用户场景，无需缓存

## Rate Limiting 架构设计

**实现模式**: 使用 Resty `OnBeforeRequest` middleware 统一处理

**关键设计决策**：
- **集中式处理**: Rate limiting 在 `client.go` 的 `OnBeforeRequest` 中统一处理，避免在每个 API 方法中重复代码
- **阻塞式等待**: 使用 `rateLimiter.Wait(ctx)` 而非 `Allow()` 模式，提供更好的用户体验
- **Context 支持**: 尊重 context 取消和超时，防止资源泄漏
- **透明集成**: API 方法无需显式调用 rate limiting，middleware 自动处理

**实现示例**（`internal/tmdb/client.go`）：
```go
httpClient := resty.New().
    SetBaseURL(baseURL).
    OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
        // 1. 统一处理 rate limiting (阻塞等待)
        if err := rateLimiter.Wait(req.Context()); err != nil {
            logger.Error("rate limit wait failed", zap.Error(err))
            return fmt.Errorf("rate limit wait failed: %w", err)
        }

        // 2. 自动添加 API Key
        req.SetQueryParam("api_key", cfg.APIKey)

        // 3. language 参数仅在请求中未显式设置时使用配置默认值
        if req.QueryParam.Get("language") == "" && cfg.Language != "" {
            req.SetQueryParam("language", cfg.Language)
        }
        return nil
    })
```

**优势**：
- ✅ **DRY 原则**: 消除了 6+ 个 API 方法中的重复代码
- ✅ **维护性**: 修改 rate limiting 逻辑只需修改一处
- ✅ **一致性**: 所有 API 调用自动应用相同的 rate limiting 策略
- ✅ **可观测性**: `Wait(ctx)` 允许记录实际等待时间，便于监控

[Source: docs/stories/1.4.rate-limiting-mechanism.md#Architecture Decision Record]


## MCP Middleware 架构设计

**实现模式**: 使用 MCP SDK `AddReceivingMiddleware()` 统一处理日志记录

**背景问题**:
- `internal/tools/` 下的 6+ 个工具文件中存在 20+ 处重复的日志记录代码
- 每个工具的 `Handler()` 都包含相同的 "request received"、"failed"、"completed" 日志
- 违反 DRY 原则，维护成本高，容易出现格式不一致

**关键设计决策**：
- **集中式拦截**: Middleware 在 `server.go` 的 `NewServer()` 中统一注册，拦截所有 MCP 工具调用
- **分层日志**: 通用日志（请求/响应）由 middleware 处理，业务特定日志（如 "No movies found"）保留在工具层
- **完整追踪**: 记录 method、session_id、tool name、arguments、duration、success/error 状态
- **透明集成**: 工具代码无需显式调用日志，middleware 自动处理

**实现示例**（`internal/mcp/middleware.go`）：
```go
func LoggingMiddleware(logger *zap.Logger) mcp.Middleware {
    return func(next mcp.MethodHandler) mcp.MethodHandler {
        return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
            // 1. 记录请求开始
            logger.Info("MCP method started",
                zap.String("method", method),
                zap.String("session_id", req.GetSession().ID()),
                zap.Bool("has_params", req.GetParams() != nil))

            // 2. 对 CallToolRequest，记录工具名称和参数
            if ctr, ok := req.(*mcp.CallToolRequest); ok {
                logger.Info("Calling tool",
                    zap.String("name", ctr.Params.Name),
                    zap.Any("args", ctr.Params.Arguments))
            }

            // 3. 执行实际方法 + 计时
            start := time.Now()
            result, err := next(ctx, method, req)
            duration := time.Since(start)

            // 4. 记录结果（成功或失败）
            if err != nil {
                logger.Error("MCP method failed",
                    zap.String("method", method),
                    zap.Duration("duration", duration),
                    zap.Error(err))
            } else {
                logger.Info("MCP method completed",
                    zap.String("method", method),
                    zap.Duration("duration", duration))
            }

            return result, err
        }
    }
}
```

**注册方式**（`internal/mcp/server.go`）：
```go
mcpServer := mcp.NewServer(&mcp.Implementation{
    Name:    "tmdb-mcp",
    Version: "1.0.0",
}, opts)

// Add logging middleware (must be added before registering tools)
mcpServer.AddReceivingMiddleware(LoggingMiddleware(logger))
```

**优势**：
- ✅ **DRY 原则**: 消除了 6+ 个工具文件中的 20+ 处重复日志代码
- ✅ **维护性**: 修改日志格式只需修改一处（middleware）
- ✅ **一致性**: 所有 MCP 工具自动应用相同的日志格式和字段
- ✅ **可观测性**: 统一记录 session_id、duration，便于追踪和性能监控
- ✅ **关注点分离**: 工具层专注业务逻辑，横切关注点（日志）由 middleware 处理

**对比项目现有实践**:
本设计模式与项目中 TMDB HTTP Client 的 Resty middleware 一致：
- HTTP Client 层：使用 `OnBeforeRequest` middleware 统一处理 rate limiting
- MCP Server 层：使用 `AddReceivingMiddleware` 统一处理日志记录
- 统一的 middleware 模式降低了开发人员的认知负担

[Source: Official MCP SDK Middleware Example]


## 精简原则落地

**不使用的工具/库**:
- ❌ Makefile - 仅使用 `go build`, `go test`
- ❌ golangci-lint - 仅使用 `go vet`
- ❌ 其他 linter（staticcheck, errcheck 等）
- ❌ Web 框架（Gin, Echo, Fiber）
- ❌ ORM 框架（无数据库需求）
- ❌ 依赖注入框架（wire, dig）- 手动构造函数注入

**构建命令示例**:
```bash
# 编译
go build -o tmdb-mcp ./cmd/tmdb-mcp

# 测试
go test ./...

# 格式化
go fmt ./...

# 静态检查
go vet ./...

# 运行
./tmdb-mcp
```

---
