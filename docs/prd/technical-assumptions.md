# Technical Assumptions

## Repository Structure: Monorepo

**决策**: 采用 **Monorepo** 结构

**理由**:
- 单体项目，代码规模可控（MVP 阶段）
- 便于统一管理依赖和构建配置
- 简化开发、测试和部署流程
- 符合简报中提到的 Repository Structure 建议

## Service Architecture

**决策**: 采用 **单进程 Monolith 架构**

**架构设计**:
```
单进程 MCP 服务（通过 stdio 通信）
├─ MCP 协议层（处理 JSON-RPC over stdio）
├─ 工具层（6 个 MCP 工具实现）
├─ TMDB 客户端层（封装 TMDB API v3 调用）
└─ 速率限制层（请求队列管理，40 req/10s）
```

**理由**:
- MVP 阶段功能单一，无需微服务复杂度
- 单进程启动快（< 2 秒要求）
- 通过 stdio 与 MCP 客户端通信，架构简单清晰
- Golang 的并发能力足以处理速率限制范围内的并发请求

## Testing Requirements

**决策**: **Unit + Integration 测试策略**

**测试范围**:
- **Unit Tests**: 覆盖工具层、TMDB 客户端层、速率限制层
- **Integration Tests**: 使用 Mock TMDB API 测试端到端流程
- **Manual Tests**: 使用真实 Claude 客户端进行用户场景测试

**理由**:
- MVP 阶段优先保证核心逻辑正确性
- 集成测试验证 MCP 协议和 TMDB API 交互
- E2E 测试成本高且依赖外部环境，暂不纳入自动化

## Additional Technical Assumptions

1. **编程语言**: Golang 1.21+
   - 类型安全、高性能、优秀的并发模型
   - 官方 MCP SDK 支持（`github.com/modelcontextprotocol/go-sdk`）

2. **核心依赖库**:
   - MCP SDK: `github.com/modelcontextprotocol/go-sdk` (官方 SDK，内置 SSE 支持)
   - HTTP 客户端（TMDB API）: `github.com/go-resty/resty/v2`
   - **HTTP 服务器（SSE）**: `net/http` (标准库 + MCP SDK 的 `SSEHTTPHandler`)
   - 速率限制: `golang.org/x/time/rate`
   - 日志: `go.uber.org/zap`
   - 配置: `github.com/spf13/viper` (支持配置文件、环境变量、命令行 flags)
   - Token 生成: `crypto/rand` (标准库)

3. **配置管理**:
   - 优先级: 命令行 flags > 环境变量 > 配置文件
   - 配置文件路径: `~/.tmdb-mcp/config.yaml`
   - 配置文件格式: YAML
   - 必需配置: TMDB API Key
   - 可选配置: 语言偏好、速率限制、日志级别、SSE Token
   - **SSE Token 管理**:
     - 环境变量 `SSE_TOKEN` 优先级最高（方便 Docker 用户）
     - 配置文件 `server.sse.token` 次之
     - 若两者都未设置且 SSE 模式启用，首次启动时自动生成 256-bit 随机 token 并写入配置文件
     - Docker 用户通过修改环境变量 + 重启容器即可刷新 token

4. **部署方式**:
   - 独立二进制文件（跨平台编译：Linux、macOS、Windows）
   - Docker 容器镜像（发布到 Docker Hub）
   - 分发通过 GitHub Releases

5. **开发环境**:
   - Go Modules 管理依赖
   - 使用 Go 自带工具链：
     - `go build` - 编译
     - `go test` - 测试
     - `go fmt` - 格式化
     - `go vet` - 静态检查
   - **HTTP 服务器配置**:
     - 使用标准库 `net/http` 和 MCP SDK 的 `SSEHTTPHandler`
     - 集成 zap 结构化日志记录 HTTP 请求

6. **安全考虑**:
   - API Key 不得硬编码，仅通过环境变量或配置文件读取
   - 配置文件包含敏感信息时应添加到 `.gitignore`
   - 遵守 TMDB API 使用条款，不滥用 API

7. **监控和调试**:
   - 结构化日志记录所有 API 调用、错误和性能指标
   - 支持日志级别配置（DEBUG、INFO、WARN、ERROR）
   - 关键路径添加性能追踪（响应时间、API 调用次数）

8. **MCP 协议实现**:
   - 严格遵循 MCP 规范（JSON-RPC 2.0 over stdio）
   - 工具描述必须清晰完整，包含参数类型、示例和约束
   - 支持 MCP 协议的 `tools/list` 和 `tools/call` 方法

9. **错误处理策略**:
   - 所有外部 API 调用必须有超时控制（默认 10 秒）
   - 速率限制触发时（429），自动等待并重试（最多 3 次）
   - 网络错误、超时等暂时性错误返回友好提示，建议用户重试

10. **文档要求**:
    - README 包含快速开始指南、配置说明、使用示例
    - 每个 MCP 工具提供清晰的参数说明和示例
    - 提供故障排查指南（常见错误及解决方法）

11. **SSE 访问模式**:
    - 支持两种运行模式：`stdio`（标准 MCP）和 `sse`（Server-Sent Events over HTTP）
    - 可同时启用两种模式（`mode: both`）
    - **HTTP 服务器实现**: 使用标准库 `net/http` + MCP SDK 的 `SSEHTTPHandler`
      - `SSEHTTPHandler` 是 MCP SDK 提供的官方 SSE 处理器
      - 实现了 `http.Handler` 接口,可直接用于 `http.Server`
      - 通过 `mcp.NewSSEHTTPHandler(getServer func(*http.Request) *Server)` 创建
    - SSE 配置：
      - 默认端口：`8910`
      - 默认绑定：`0.0.0.0`（支持远程访问）
      - 认证方式：`Authorization: Bearer <token>` header（使用标准库中间件实现）
      - Token 长期有效，无过期机制
    - SSE 端点：
      - `GET /mcp/sse` - 建立 SSE 连接（需要 Bearer token）
      - `GET /health` - 健康检查端点（无需认证）
    - HTTPS 支持：由用户通过 Nginx/Caddy 等反向代理实现，服务本身仅提供 HTTP
    - 安全建议：
      - Token 应保密，避免提交到版本控制
      - 配置文件权限应设置为 `600`（仅所有者可读写）
      - 公网暴露时强烈建议配置反向代理 + HTTPS

**配置文件示例**:

```yaml
# ~/.tmdb-mcp/config.yaml
tmdb:
  api_key: "your_tmdb_api_key"
  language: "en-US"
  rate_limit: 40  # requests per 10 seconds

server:
  mode: "both"  # stdio, sse, or both
  sse:
    enabled: true
    host: "0.0.0.0"  # 监听所有网络接口（支持远程访问）
    port: 8910
    token: "auto-generated-on-first-run"  # 首次启动自动生成，或通过 SSE_TOKEN 环境变量设置

logging:
  level: "info"  # debug, info, warn, error
```

---
