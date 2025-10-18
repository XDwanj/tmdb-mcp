# TMDB MCP 服务产品需求文档 (PRD)

**项目名称**: TMDB MCP 服务
**版本**: v1.0
**日期**: 2025-10-10
**状态**: 已批准 ✅
**PRD 完整度**: 92/100

---

## 目录

1. [Goals and Background Context](#goals-and-background-context)
2. [Requirements](#requirements)
3. [Technical Assumptions](#technical-assumptions)
4. [Epic List](#epic-list)
5. [Epic Details](#epic-details)
6. [Checklist Results Report](#checklist-results-report)
7. [Next Steps](#next-steps)

---

## Goals and Background Context

### Goals

- 通过 MCP 协议将 TMDB API 与 LLM 连接，实现自然语言交互方式查询影视数据
- 解决自建流媒体用户面临的文件命名混乱和内容发现困难问题
- 降低 TMDB 数据库的使用门槛，让普通用户能够通过对话方式使用强大的影视数据
- 实现智能文件重命名功能，自动识别混乱的文件名并标准化
- 提供对话式内容探索和个性化推荐，解决用户"片荒"痛点
- 建立高性能、类型安全的 Golang MCP 服务，支持并发处理和速率限制

### Background Context

TMDB MCP 服务是一个创新性项目，旨在利用 LLM 作为"超级胶水"将传统 TMDB API 转化为对话式操作接口。当前自建流媒体服务器（Jellyfin、Emby、Plex）用户在管理影视资源时面临严重的文件命名混乱问题——下载的文件通常包含复杂格式导致刮削器无法识别。同时，用户在探索新内容时需要在 TMDB 网站进行繁琐的多次搜索和页面跳转。本项目通过 Model Context Protocol (MCP) 将 TMDB 的强大数据能力暴露给 Claude 等 LLM，使用户能够用一句话完成原本需要多个步骤的复杂操作。

该项目采用场景驱动设计理念，精选 6 个核心 MCP 工具（搜索、详情、发现电影、发现电视剧、趋势、推荐），而非完整映射 TMDB API。通过智能合并策略将 10 个原子 API 合并为 2 个统一工具，大幅降低 LLM 的认知负担。项目使用 Golang 实现，确保高性能、类型安全和天然并发支持，满足 TMDB 免费 API 的 40 req/10s 速率限制要求。

### Change Log

| Date | Version | Description | Author |
|------|---------|-------------|--------|
| 2025-10-10 | 1.0 | 初始 PRD 创建，基于项目简报 v1.0 | John (PM) |

---

## Requirements

### Functional Requirements

1. **FR1**: 系统必须提供统一的搜索工具（search），支持通过自然语言查询电影、电视剧、人物和多媒体内容，并返回匹配结果列表（包括标题、年份、评分等关键信息）

2. **FR2**: 系统必须提供详情获取工具（get_details），能够根据 TMDB ID 获取电影、电视剧或人物的完整信息，并自动追加演职人员表（credits）和相关视频（videos）数据

3. **FR3**: 系统必须提供电影发现工具（discover_movies），支持按类型、年份、评分、语言等条件筛选电影，并返回符合条件的电影列表

4. **FR4**: 系统必须提供电视剧发现工具（discover_tv），支持按类型、首播年份、评分、状态等条件筛选电视剧，并返回符合条件的剧集列表

5. **FR5**: 系统必须提供趋势工具（get_trending），能够获取当前热门的电影、电视剧或人物，支持按日榜（day）或周榜（week）查询

6. **FR6**: 系统必须提供推荐工具（get_recommendations），能够基于指定的电影或电视剧 ID 获取相似内容推荐列表

7. **FR7**: 系统必须实现完整的错误处理机制，能够识别并优雅处理 TMDB API 返回的 401（未授权）、404（未找到）、429（速率限制）等错误，并向用户返回清晰的错误信息

8. **FR8**: 系统必须支持 TMDB API Key 和语言参数的配置管理:
   - API Key 通过环境变量或配置文件读取
   - 语言参数支持**两级优先级**:
     1. **工具级参数**: MCP 工具调用时传入的 `language` 参数(最高优先级)
     2. **配置默认值**: 环境变量/配置文件中的 `language` 设置(未指定时默认 `en-US`)
   - 不实现自动语言检测
   - 所有 MCP 工具必须统一遵循此优先级模型

9. **FR9**: 系统必须支持两种运行模式：
   - **stdio 模式**：通过标准输入/输出与 MCP 客户端通信（遵循 MCP 协议规范）
   - **SSE 模式**：通过 HTTP Server-Sent Events 提供远程访问，默认端口 8910，使用 `Authorization: Bearer <token>` 进行认证
   - 支持同时启用两种模式

10. **FR10**: (已删除 - LLM 负责文件名解析)

11. **FR11**: SSE 模式必须提供以下端点：
    - `GET /mcp/sse` - 建立 MCP over SSE 连接（需要有效的 Bearer token）
    - `GET /health` - 健康检查端点（返回服务状态，无需认证）

### Non-Functional Requirements

1. **NFR1**: API 响应时间的 P95 必须小于 500 毫秒（不包含 TMDB API 自身延迟）

2. **NFR2**: 系统启动时间必须小于 2 秒，确保用户快速开始使用

3. **NFR3**: 系统必须实现速率限制机制，支持通过配置文件或命令行 flags 设置请求速率（默认 40 requests/10 seconds），命令行参数优先级高于配置文件，确保不触发 TMDB API 限流

4. **NFR4**: 系统的 API 调用错误率（不含 TMDB 自身错误）必须小于 1%

5. **NFR5**: 系统必须支持跨平台运行（Linux、macOS、Windows），以独立二进制文件或 Docker 容器形式分发

6. **NFR6**: 系统必须使用 Golang 1.21+ 实现，利用其类型安全和天然并发能力，确保代码质量和性能

7. **NFR7**: 系统必须支持并发处理多个请求，在速率限制范围内最大化吞吐量

8. **NFR8**: 系统日志必须使用结构化日志框架（如 zap），记录关键操作、错误和性能指标，便于调试和监控

9. **NFR9**: 系统配置（API Key、语言偏好、速率限制）必须通过环境变量、配置文件或命令行 flags 管理（优先级：命令行 > 环境变量 > 配置文件），不得在代码中硬编码敏感信息

10. **NFR10**: 系统必须遵守 TMDB API 使用条款，包括在返回数据中保留归属声明和遵守速率限制

11. **NFR11**: SSE 模式的 Token 管理必须满足：
    - 优先使用环境变量 `SSE_TOKEN`（方便 Docker 部署）
    - 其次从配置文件 `~/.tmdb-mcp/config.yaml` 读取
    - 若两者都未设置，首次启动时自动生成 256-bit 随机 token 并持久化到配置文件
    - Token 验证失败时返回 `401 Unauthorized`

12. **NFR12**: 配置文件 `~/.tmdb-mcp/config.yaml` 创建时必须设置权限为 `600`（仅所有者可读写），防止 token 泄露

---

## Technical Assumptions

### Repository Structure: Monorepo

**决策**: 采用 **Monorepo** 结构

**理由**:
- 单体项目，代码规模可控（MVP 阶段）
- 便于统一管理依赖和构建配置
- 简化开发、测试和部署流程
- 符合简报中提到的 Repository Structure 建议

### Service Architecture

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

### Testing Requirements

**决策**: **Unit + Integration 测试策略**

**测试范围**:
- **Unit Tests**: 覆盖工具层、TMDB 客户端层、速率限制层
- **Integration Tests**: 使用 Mock TMDB API 测试端到端流程
- **Manual Tests**: 使用真实 Claude 客户端进行用户场景测试

**理由**:
- MVP 阶段优先保证核心逻辑正确性
- 集成测试验证 MCP 协议和 TMDB API 交互
- E2E 测试成本高且依赖外部环境，暂不纳入自动化

### Additional Technical Assumptions

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

## Epic List

### Epic 1: Foundation & Search
**目标**: 建立项目基础设施并实现第一个可工作的 MCP 工具（search），支持通过 stdio 模式搜索电影、电视剧和人物。

### Epic 2: Details & Discovery Tools
**目标**: 实现内容详情获取和发现功能，支持智能文件重命名场景和内容筛选场景，完善错误处理机制。核心工具集功能完备，自动化测试验证所有使用场景的可行性。

### Epic 3: Trending & Recommendations
**目标**: 实现趋势和推荐工具，完成所有 6 个核心 MCP 工具，优化速率限制和性能监控。通过全面的自动化集成测试验证所有工具协同工作，确保性能指标符合 NFR 要求。此 Epic 完成后，stdio 模式功能完整可用，达到 MVP 核心功能完整的里程碑。可选地准备演示材料用于后续文档和社区宣传。

### Epic 4: SSE Remote Access Mode
**目标**: 使用 MCP SDK 的 SSEHTTPHandler 和标准库 net/http 实现 SSE 远程访问模式和 Token 认证，支持 stdio + sse 双模式运行，完成 Docker 镜像构建和部署。

### Epic 5: Documentation, Examples & Community Launch
**目标**: 完善项目文档（README、配置指南、使用示例、故障排查），提供真实场景的示例配置和脚本，准备并发布 GitHub Release、Docker Hub 镜像，向社区宣传（r/selfhosted、r/jellyfin）并收集早期用户反馈。

---

## Epic Details

### Epic 1: Foundation & Search

**Epic Goal**: 建立 TMDB MCP 服务的核心技术基础设施，包括配置管理、日志系统、TMDB API 客户端封装、速率限制机制和 MCP stdio 协议集成。实现并交付第一个可工作的 MCP 工具（search），使用户能够通过 Claude 等 LLM 客户端搜索电影、电视剧和人物，验证整个技术栈的可行性并为后续功能打下坚实基础。

#### Story 1.1: Project Initialization and Configuration Management

**As a** developer,
**I want** to initialize the Go project structure and implement a flexible configuration management system,
**so that** I can manage TMDB API Key, language preferences, and other settings through multiple sources (config file, environment variables, command-line flags) with proper priority.

**Acceptance Criteria**:

1. 项目使用 Go Modules 初始化（`go mod init github.com/[username]/tmdb-mcp`）
2. 目录结构遵循 Go 标准布局：`cmd/tmdb-mcp/`, `internal/`, `pkg/`, `.gitignore`
3. 集成 `github.com/spf13/viper` 实现配置管理，支持从 `~/.tmdb-mcp/config.yaml` 读取配置，支持环境变量和命令行 flags，优先级：命令行 > 环境变量 > 配置文件
4. 定义配置结构体，包含 `tmdb.api_key`, `tmdb.language`, `tmdb.rate_limit`, `logging.level`
5. 首次运行时，如果 `~/.tmdb-mcp/` 目录不存在，自动创建
6. 如果缺少必需配置（TMDB API Key），程序输出清晰的错误提示并退出
7. 提供配置验证功能，启动时检查配置有效性

#### Story 1.2: Structured Logging System

**As a** developer,
**I want** to integrate a structured logging system using zap,
**so that** I can record key operations, errors, and performance metrics in a structured format for debugging and monitoring.

**Acceptance Criteria**:

1. 集成 `go.uber.org/zap` 日志库
2. 实现日志初始化函数，根据配置文件中的 `logging.level` 设置日志级别
3. 日志输出格式：开发模式使用 `zap.NewDevelopment()`，生产模式使用 `zap.NewProduction()`（JSON 格式）
4. 提供全局 logger 实例，可在整个项目中使用
5. 记录关键事件：程序启动、配置加载、程序退出
6. 日志字段包含上下文信息（时间戳、日志级别、caller 信息）
7. 确保日志不会泄露敏感信息（API Key 应被遮盖）

#### Story 1.3: TMDB API Client Foundation

**As a** developer,
**I want** to create a TMDB API client wrapper using resty,
**so that** I can make authenticated HTTP requests to TMDB API with proper error handling and response parsing.

**Acceptance Criteria**:

1. 集成 `github.com/go-resty/resty/v2` 作为 HTTP 客户端
2. 创建 `TMDBClient` 结构体，封装 TMDB API Key、Base URL、Language preference、Resty client 实例
3. 实现 `NewTMDBClient(apiKey, language string)` 构造函数
4. 配置 Resty client：设置 Base URL、自动添加 API Key、设置超时时间（10 秒）、设置 User-Agent
5. 实现通用错误处理函数，解析 TMDB API 错误响应（401/404/429）
6. 实现测试方法 `Ping()`，调用 `/configuration` 端点验证 API Key 有效性
7. 编写单元测试，使用 mock 验证 API Key 正确添加、错误响应被正确解析

#### Story 1.4: Rate Limiting Mechanism

**As a** developer,
**I want** to implement a rate limiting layer using `golang.org/x/time/rate`,
**so that** I can ensure all TMDB API requests respect the rate limit (40 requests per 10 seconds) and avoid triggering API throttling.

**Acceptance Criteria**:

1. 集成 `golang.org/x/time/rate` 包
2. 创建 `RateLimiter` 包装器，使用 `rate.NewLimiter(rate.Every(10*time.Second/40), 40)` 配置速率，支持通过配置文件自定义
3. 将 `RateLimiter` 集成到 `TMDBClient`，每次 API 调用前调用 `Wait(ctx)` 方法
4. 实现可观测性：记录速率限制等待事件到日志（debug 级别）
5. 编写单元测试验证：在 10 秒内最多允许 40 个请求
6. 编写集成测试：模拟 50 个快速请求，验证请求被正确限流

#### Story 1.5: MCP Protocol Integration via stdio

**As a** developer,
**I want** to integrate the official MCP Go SDK and implement stdio transport,
**so that** the service can communicate with MCP clients (like Claude) using JSON-RPC over standard input/output.

**Acceptance Criteria**:

1. 集成 `github.com/modelcontextprotocol/go-sdk` (官方 MCP Go SDK)
2. 实现 MCP 服务器初始化：创建 MCP server 实例、配置 stdio transport、注册服务器信息
3. 实现 `tools/list` 方法，返回可用工具列表（当前为空）
4. 实现 `tools/call` 方法框架，支持调用已注册的工具
5. 实现主程序入口：加载配置、初始化日志、创建 TMDB client、启动 MCP server、优雅退出
6. 验证程序可以编译和运行
7. 手动测试：使用 MCP 客户端连接，验证 `tools/list` 返回空列表

#### Story 1.6: Implement Search Tool

**As a** user,
**I want** to search for movies, TV shows, and people using natural language queries through the MCP search tool,
**so that** I can find TMDB content without knowing TMDB IDs or using complex web interfaces.

**Acceptance Criteria**:

1. 实现 `search` 工具，映射到 TMDB API `/search/multi` 端点
2. 工具定义：Name: `search`, Description: "Search for movies, TV shows, and people on TMDB using a query string", Parameters: `query` (string, required), `page` (integer, optional)
3. 实现 TMDB client 的 `Search(query string, page int)` 方法
4. 返回结果包含：`id`, `media_type`, `title`/`name`, `release_date`/`first_air_date`, `vote_average`, `overview`
5. 错误处理：query 为空返回错误、TMDB API 错误返回友好消息、无结果返回空数组
6. 在 MCP server 中注册 `search` 工具
7. 编写单元测试：Mock TMDB API 响应、验证查询参数、验证结果解析
8. 编写集成测试：搜索 "Inception"、搜索 "Christopher Nolan"、搜索不存在内容

#### Story 1.7: Automated End-to-End Integration Testing

**As a** developer,
**I want** to implement automated integration tests using MCP SDK's InMemoryTransports,
**so that** I can continuously verify the search tool works correctly without manual intervention and ensure the entire MCP protocol stack is functioning properly.

**Acceptance Criteria**:

1. **自动化集成测试框架**（必需）：
   - 创建 `cmd/tmdb-mcp/integration_test.go` 使用 InMemoryTransports
   - 使用 `mcp.NewInMemoryTransports()` 创建 client-server 通信对
   - 在同一进程内模拟完整的 MCP 协议交互
   - 无需启动外部进程或 Claude Code 客户端

2. **测试用例覆盖**（必需）：
   - ✅ 成功场景：搜索流行电影（"Inception"）、搜索电视剧（"Breaking Bad"）、搜索人物（"Christopher Nolan"）
   - ✅ 边界场景：空查询、不存在的内容（返回空结果）、分页测试
   - ✅ 错误场景：无效参数、TMDB API 错误模拟
   - ✅ 结果验证：检查返回数据结构、字段完整性、数据类型正确性

3. **性能验证**（必需）：
   - 每次搜索调用的响应时间 < 3 秒
   - 记录并验证 API 调用次数
   - 使用 Go testing 的 benchmark 功能测试吞吐量

4. **速率限制验证**（必需）：
   - 快速执行 10 次搜索请求
   - 验证没有触发 429 错误
   - 验证 RateLimiter 正确工作（通过日志或计数器）

5. **测试覆盖率**（必需）：
   - 使用 `go test -cover` 检查覆盖率
   - 目标：`internal/tools` 包覆盖率 ≥ 70%
   - 目标：`internal/tmdb` 包覆盖率 ≥ 70%

6. **CI/CD 集成**（必需）：
   - 测试可以通过 `go test ./...` 运行
   - 无需外部依赖（使用 Mock TMDB API 或环境变量控制）
   - 测试结果输出清晰，失败时提供有用的错误信息

7. **手动验证**（可选，作为补充）：
   - 在 `.ai/epic1-e2e-test-results.md` 记录使用真实 Claude Code 的手动测试结果
   - 验证用户体验和自然语言交互效果
   - 截图和日志作为文档参考

**实现参考**（基于官方 MCP SDK）：
```go
func TestSearchTool_Integration(t *testing.T) {
    ctx := context.Background()

    // 创建内存传输对
    clientTransport, serverTransport := mcp.NewInMemoryTransports()

    // 初始化 server
    server := setupMCPServer(t) // 包含 search tool
    serverSession, _ := server.Connect(ctx, serverTransport, nil)
    defer serverSession.Close()

    // 初始化 client
    client := mcp.NewClient(&mcp.Implementation{Name: "test-client"}, nil)
    clientSession, _ := client.Connect(ctx, clientTransport, nil)
    defer clientSession.Close()

    // 测试搜索功能
    start := time.Now()
    result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
        Name: "search",
        Arguments: map[string]any{"query": "Inception", "page": 1},
    })
    duration := time.Since(start)

    // 验证结果
    assert.NoError(t, err)
    assert.Less(t, duration, 3*time.Second)
    // ... 更多验证
}
```

---

### Epic 2: Details & Discovery Tools

**Epic Goal**: 实现内容详情获取工具（get_details）和内容发现工具（discover_movies、discover_tv），使用户能够获取电影/电视剧/人物的完整信息，并通过自然语言表达的筛选条件探索新内容。完善错误处理机制，优雅处理 TMDB API 的各类错误（401/404/429）。此 Epic 完成后，核心工具集功能完备，自动化测试验证所有使用场景的可行性。

#### Story 2.1: Implement get_details Tool

**As a** user,
**I want** to retrieve detailed information about a movie, TV show, or person using their TMDB ID,
**so that** I can get complete metadata (including cast, crew, and videos) for file renaming or content exploration.

**Acceptance Criteria**:

1. 实现 `get_details` 工具，支持三种 media_type：movie, tv, person
2. 工具定义：Name: `get_details`, Parameters: `media_type` (string, required), `id` (integer, required)
3. 自动追加功能：电影/电视剧详情自动追加 `append_to_response=credits,videos`，人物详情自动追加 `combined_credits`
4. 实现 TMDB client 方法：`GetMovieDetails()`, `GetTVDetails()`, `GetPersonDetails()`
5. 返回结果包含核心字段（根据 media_type）
6. 错误处理：media_type 无效、ID 不存在（404）、TMDB API 错误
7. 在 MCP server 中注册 `get_details` 工具
8. 编写单元测试：Mock TMDB API 响应、验证 append_to_response 参数
9. 编写集成测试：获取《盗梦空间》、《权力的游戏》、克里斯托弗·诺兰详情

#### Story 2.2: Implement discover_movies Tool

**As a** user,
**I want** to discover movies using filters like genre, year, rating, and language,
**so that** I can find movies matching my preferences without manually browsing TMDB website.

**Acceptance Criteria**:

1. 实现 `discover_movies` 工具，映射到 `/discover/movie` 端点
2. 工具定义：Name: `discover_movies`, Parameters: `with_genres`, `primary_release_year`, `vote_average.gte/lte`, `with_original_language`, `sort_by`, `page`
3. 实现 TMDB client 的 `DiscoverMovies(params DiscoverMoviesParams)` 方法
4. 返回结果字段：`id`, `title`, `release_date`, `vote_average`, `overview`, `genre_ids`, `popularity`
5. 参数验证：vote_average 范围 0-10、sort_by 支持的值
6. 默认行为：所有参数为空时返回最流行的电影
7. 工具描述中提供示例
8. 在 MCP server 中注册工具
9. 编写单元测试：Mock API 响应、验证参数映射
10. 编写集成测试：查找 2020 年后的高分科幻片、评分最高的动作片

#### Story 2.3: Implement discover_tv Tool

**As a** user,
**I want** to discover TV shows using filters like genre, year, rating, and status,
**so that** I can find TV series matching my preferences.

**Acceptance Criteria**:

1. 实现 `discover_tv` 工具，映射到 `/discover/tv` 端点
2. 工具定义：Name: `discover_tv`, Parameters: `with_genres`, `first_air_date_year`, `vote_average.gte/lte`, `with_original_language`, `with_status`, `sort_by`, `page`
3. 实现 TMDB client 的 `DiscoverTV(params DiscoverTVParams)` 方法
4. 返回结果字段：`id`, `name`, `first_air_date`, `vote_average`, `overview`, `genre_ids`, `origin_country`
5. 参数验证（同 discover_movies）
6. 工具描述中提供示例
7. 在 MCP server 中注册工具
8. 编写单元测试
9. 编写集成测试：查找高分犯罪剧、正在播出的科幻剧

#### Story 2.4: Enhanced Error Handling and Retry Logic

**As a** developer,
**I want** to implement comprehensive error handling for all TMDB API calls,
**so that** users receive clear error messages and the service can gracefully handle rate limiting, network issues, and API errors.

**Acceptance Criteria**:

1. 标准化错误响应结构：创建 `TMDBError` 类型
2. 401 Unauthorized 处理：立即返回 "Invalid or missing TMDB API Key"、记录 ERROR 日志、不重试
3. 404 Not Found 处理：返回 "Resource not found"、记录 INFO 日志、不重试
4. 429 Rate Limit Exceeded 处理：解析 `Retry-After` header、等待后重试（最多 3 次）、记录 WARN 日志
5. 网络超时处理：返回 "Request timeout"、记录 WARN 日志
6. 其他 HTTP 错误（500, 502, 503）：返回错误消息、记录 ERROR 日志
7. JSON 解析错误：返回 "Failed to parse response"
8. MCP 工具层错误处理：转换为 MCP 错误响应格式
9. 日志记录增强：记录 endpoint, parameters, response_time, error_type
10. 编写单元测试：Mock 各类错误响应、验证重试逻辑
11. 编写集成测试：使用无效 API Key 触发 401、请求不存在的 ID 触发 404

#### Story 2.5: (已删除)

**说明**: 此 Story 已移除。原 Story 2.5 依赖外部工具（Claude Code）进行手动测试，导致测试不可靠且无法自动化。文件重命名场景的测试已由 Story 1.7（自动化集成测试）和 Story 3.4（综合集成测试）充分覆盖。

**文件重命名场景的自动化测试策略**:
- Story 1.7 已使用 InMemoryTransports 模拟完整 MCP 协议交互
- Story 3.4 提供多工具组合测试（search → get_details）
- 文档示例将在 Epic 5 (Story 5.2) 中提供

---

### Epic 3: Trending & Recommendations

**Epic Goal**: 实现最后两个 MCP 工具（get_trending 和 get_recommendations），完成所有 6 个核心工具的功能集。添加性能监控和指标记录，优化 API 调用效率和响应时间。通过全面的自动化集成测试验证所有工具协同工作，确保性能指标符合 NFR 要求。此 Epic 完成后，stdio 模式功能完整可用，达到 MVP 核心功能完整的里程碑。可选地准备演示材料用于后续文档和社区宣传。

#### Story 3.1: Implement get_trending Tool

**As a** user,
**I want** to get trending movies, TV shows, or people for a specific time window (day or week),
**so that** I can quickly discover currently popular content without browsing TMDB website.

**Acceptance Criteria**:

1. 实现 `get_trending` 工具，映射到 `/trending/{media_type}/{time_window}` 端点
2. 工具定义：Name: `get_trending`, Parameters: `media_type` (movie/tv/person), `time_window` (day/week), `page`
3. 实现 TMDB client 的 `GetTrending(mediaType, timeWindow string, page int)` 方法
4. 返回结果字段（根据 media_type）
5. 参数验证：media_type 和 time_window 必须是有效值
6. 工具描述中提供示例
7. 在 MCP server 中注册工具
8. 编写单元测试
9. 编写集成测试：获取今日热门电影、本周热门电视剧、热门人物

#### Story 3.2: Implement get_recommendations Tool

**As a** user,
**I want** to get movie or TV show recommendations based on a specific title I like,
**so that** I can discover similar content matching my preferences.

**Acceptance Criteria**:

1. 实现 `get_recommendations` 工具，映射到 `/movie/{id}/recommendations` 和 `/tv/{id}/recommendations` 端点
2. 工具定义：Name: `get_recommendations`, Parameters: `media_type` (movie/tv), `id` (integer), `page`
3. 实现 TMDB client 方法：`GetMovieRecommendations()`, `GetTVRecommendations()`
4. 返回结果字段：`id`, `title`/`name`, `release_date`/`first_air_date`, `vote_average`, `overview`, `popularity`
5. 参数验证：media_type 和 id 有效性
6. 错误处理：ID 不存在（404）、无推荐结果返回空数组
7. 工具描述中提供示例
8. 在 MCP server 中注册工具
9. 编写单元测试
10. 编写集成测试：基于《盗梦空间》获取电影推荐、基于《绝命毒师》获取电视剧推荐

#### Story 3.3: Performance Monitoring and Metrics

**As a** developer,
**I want** to add performance monitoring and metrics recording,
**so that** I can track API response times, call counts, and identify performance bottlenecks.

**Acceptance Criteria**:

1. 响应时间记录：为每个 TMDB API 调用记录响应时间，使用 zap 结构化日志
2. API 调用计数：在内存中维护计数器（使用 sync/atomic）
3. 性能阈值告警：当响应时间超过 1 秒时，记录 WARN 级别日志
4. 速率限制观测性：记录速率限制等待事件（DEBUG 级别）
5. 启动时性能基准：调用 `/configuration` 端点记录响应时间作为基准
6. 定期统计日志（可选）：每 100 次 API 调用后，输出统计摘要
7. 编写单元测试：验证响应时间记录、计数器递增、阈值告警
8. 集成到所有现有工具

#### Story 3.4: Comprehensive Integration Testing

**As a** developer,
**I want** to perform comprehensive integration tests covering all 6 MCP tools,
**so that** I can verify they work correctly both individually and in combination.

**Acceptance Criteria**:

1. 单工具集成测试（使用真实 TMDB API）：每个工具至少 3 个测试用例
2. 多工具组合测试：search → get_details、discover_movies → get_recommendations、get_trending → get_details
3. 性能测试：顺序调用所有 6 个工具，总耗时 < 10 秒，验证无 429 错误
4. 并发测试：使用 goroutines 并发调用多个工具，验证速率限制正确工作，验证无数据竞争（`go test -race`）
5. 错误场景测试：无效 API Key、不存在的 ID、无效参数
6. 测试覆盖率：使用 `go test -cover` 检查覆盖率，目标：核心业务逻辑覆盖率 ≥ 70%
7. 测试结果文档：记录到 `.ai/epic3-integration-tests.md`

#### Story 3.5: End-to-End Scenario Validation (Optional Documentation)

**As a** user,
**I want** to prepare demonstration materials and document real-world usage scenarios,
**so that** potential users can understand the value of tmdb-mcp in practical contexts.

**注意**: 此 Story 为**可选**，主要用于准备演示材料和用户文档，**不作为 Epic 3 完成的阻塞条件**。技术验证已由 Story 1.7 和 3.4 的自动化测试完成。

**Acceptance Criteria**:

1. 使用 Claude Code 执行 4 个核心场景：智能文件重命名、片荒推荐、关联探索、智能推荐
2. 额外组合场景：热门内容探索 + 详情查看、发现 + 推荐链条
3. 性能验证：每个场景端到端响应时间（包括 LLM 推理）< 10 秒，复杂场景 < 15 秒
4. 用户体验验证：Claude 的回复是否自然有用、工具选择是否准确、返回数据是否满足需求
5. 错误恢复验证：故意提供模糊或错误输入，验证 Claude 能够引导用户
6. 测试结果文档：记录截图/日志到 `.ai/epic3-e2e-scenarios.md`，记录用户体验评分
7. 问题修复：记录所有问题、修复阻塞性问题
8. 交付物：演示材料（截图/录屏）记录到 `.ai/epic3-e2e-scenarios.md`，用于 Epic 5 文档和社区宣传

**里程碑确认标准（移至 Story 3.4）**:
- ✅ 所有 6 个工具的自动化集成测试通过
- ✅ 多工具组合测试成功（Story 3.4）
- ✅ 性能指标符合 NFR（Story 3.3 + 3.4）
- ✅ 测试覆盖率 ≥ 70%（Story 3.4）

---

### Epic 4: SSE Remote Access Mode

**Epic Goal**: 在现有 stdio 模式的基础上，使用 MCP Go SDK 提供的 `SSEHTTPHandler` 实现 Server-Sent Events (SSE) 远程访问模式，提供 HTTP API 端点供远程客户端连接。实现 Bearer Token 认证中间件（基于标准库 `net/http`）保护 SSE 端点，支持通过环境变量和配置文件管理 Token。支持 stdio 和 sse 双模式同时运行，并完成 Docker 镜像构建和多平台二进制文件编译，使服务可以方便地部署到远程服务器或容器环境中。

#### Story 4.1: HTTP Server Setup with Standard Library

**As a** developer,
**I want** to set up a basic HTTP server using standard library `net/http`,
**so that** I can provide HTTP endpoints for SSE connections and health checks.

**Acceptance Criteria**:

1. 在 `cmd/tmdb-mcp/main.go` 中实现 HTTP 服务器启动逻辑，使用标准库 `net/http`
2. 创建 `/health` 端点（无需认证）：返回 `{"status": "ok"}`，使用 `http.HandlerFunc`
3. 实现 SSE 模式运行函数 `RunSSEModeServer()`，配置监听地址和端口
4. 使用 `http.ServeMux` 注册路由：`/health` 和 `/mcp/sse`
5. 使用 `http.ListenAndServe()` 启动服务器（阻塞式）
6. 更新配置结构体，添加 SSE 相关配置（host, port, token）
7. 验证 `/health` 端点返回 200 OK
8. 验证服务器可正常启动并接受 HTTP 请求

#### Story 4.2: Token Generation and Management

**As a** developer,
**I want** to implement SSE Token 自动生成和管理机制,
**so that** users can securely access SSE endpoints with minimal configuration.

**Acceptance Criteria**:

1. Token 生成逻辑：使用 `crypto/rand` 生成 256-bit (32 bytes) 随机 token，编码为 hex string（64 字符）
2. Token 加载优先级：环境变量 `SSE_TOKEN` > 配置文件 `server.sse.token` > 自动生成
3. Token 持久化：新生成的 token 必须写入配置文件，确保配置文件权限为 `0600`
4. Token 显示：启动时，如果自动生成则显示完整 token，如果加载则显示前 8 个字符
5. 配置验证：如果 SSE 模式启用但 token 为空，返回错误
6. 编写单元测试：测试 token 生成长度、随机性、加载优先级
7. 编写集成测试：模拟首次启动、使用环境变量启动、验证配置文件权限

#### Story 4.3: Bearer Token Authentication Middleware

**As a** developer,
**I want** to implement Bearer Token authentication middleware using standard library `net/http`,
**so that** only authorized clients can access the SSE endpoint.

**Acceptance Criteria**:

1. 实现标准库中间件 `AuthMiddleware(expectedToken string, next http.Handler) http.Handler`
2. 认证逻辑：提取 `Authorization` header、验证格式 `Bearer <token>`、比对 token（使用字符串比较）
   - **注意**: 生产环境建议使用 `crypto/subtle.ConstantTimeCompare` 防止时序攻击
3. 认证成功：调用 `next.ServeHTTP(w, r)`
4. 认证失败：返回 `401 Unauthorized`、JSON 响应 `{"error": "unauthorized"}`
5. 错误场景处理：缺少 header、格式错误、token 不匹配
6. 将中间件应用到 SSE 路由（不应用到 `/health`）
7. 手动验证：使用有效/无效 token 访问 SSE 端点，验证认证逻辑正确

#### Story 4.4: Implement SSE Endpoint with MCP SDK

**As a** user,
**I want** to connect to the MCP service via SSE over HTTP using MCP SDK's built-in support,
**so that** I can access TMDB tools remotely from any device on the network.

**Acceptance Criteria**:

1. 使用 MCP SDK 创建 SSE handler：`sseHandler := mcp.NewSSEHTTPHandler(func(req *http.Request) *mcp.Server { return mcpServer })`
2. 实现 `/mcp/sse` 端点（需要认证）：
   - 方法 GET
   - 应用 `AuthMiddleware` 包装 `sseHandler`
   - `SSEHTTPHandler` 自动处理 SSE 连接、Content-Type 和必需的 headers
3. SSE 连接处理（由 `SSEHTTPHandler` 自动处理）：
   - 自动设置正确的 SSE headers（Content-Type: text/event-stream 等）
   - 保持连接打开
   - 内置心跳机制
4. MCP over SSE 协议（SDK 自动处理）：
   - 客户端通过 SSE 发送 JSON-RPC 请求
   - 服务器处理 MCP 请求（复用 stdio 模式的工具实现）
   - 通过 SSE 事件返回响应
5. 连接管理：记录活跃连接数、记录连接建立/断开日志
6. 错误处理：MCP 请求解析失败、工具调用失败、连接异常断开
7. 编写单元测试：测试 SSE handler 创建、认证中间件集成
8. 编写集成测试：建立 SSE 连接、发送 `tools/list`、发送 `tools/call`、验证响应格式

#### Story 4.5: Dual Mode Support (stdio + sse)

**As a** developer,
**I want** to support running both stdio and sse modes simultaneously,
**so that** users can choose their preferred connection method or use both.

**Acceptance Criteria**:

1. 模式配置：`server.mode` 支持三个值：`stdio`, `sse`, `both`（默认）
2. stdio 模式实现：启动 MCP server，监听 stdin/stdout，阻塞主 goroutine
3. sse 模式实现：启动 HTTP server（阻塞式，使用 `http.ListenAndServe`），监听端口 8910
4. both 模式实现：HTTP server 在 goroutine 中运行，stdio 在主 goroutine 中运行，共享 TMDB client 和工具实现
5. 验证三种模式可正常启动：`stdio`, `sse`, `both`
6. 验证 `both` 模式下 stdio 和 SSE 同时工作
7. 日志记录：启动时记录启用的模式
8. 配置验证：如果 mode="sse" 但 `enabled=false`，返回错误

#### Story 4.6: Docker Image and Multi-Platform Build

**As a** user,
**I want** to run tmdb-mcp in a Docker container,
**so that** I can easily deploy it to any server or cloud environment.

**Acceptance Criteria**:

1. 创建 Dockerfile（多阶段构建）：Build stage（golang:1.21-alpine）+ Runtime stage（alpine:latest）
2. 创建 `.dockerignore`：排除 `.git`, `*.md`, `.ai/`, `config.yaml`
3. 支持环境变量配置：`TMDB_API_KEY`, `SSE_TOKEN`, `SERVER_MODE`, `SERVER_SSE_HOST`, `SERVER_SSE_PORT`, `LOGGING_LEVEL`
4. 配置文件挂载支持：支持挂载 `/root/.tmdb-mcp/config.yaml`
5. 构建多平台镜像：使用 Docker Buildx 构建 `linux/amd64`, `linux/arm64`, `linux/arm/v7`
6. 创建 docker-compose.yml 示例
7. 健康检查：Dockerfile 添加 HEALTHCHECK
8. 文档：在 README 添加 Docker 部署章节
9. 测试：本地构建 Docker 镜像、运行容器并验证健康检查、端点可访问、工具调用正常
10. 多平台二进制编译（Bonus）：使用 `go build` 编译多平台二进制

#### Story 4.7: SSE Mode End-to-End Testing

**As a** user,
**I want** to verify that all MCP tools work correctly via SSE remote access,
**so that** I can confidently use the service remotely.

**Acceptance Criteria**:

1. 测试环境准备：启动服务（sse 或 both 模式）、记录 SSE Token、确认 HTTP server 运行
2. 手动 HTTP 客户端测试：测试健康检查（无需认证）、SSE 连接（无 token）、SSE 连接（有效 token）
3. MCP 工具调用测试：通过 SSE 调用所有 6 个工具
4. 并发连接测试：同时建立 5 个 SSE 连接，验证无相互干扰
5. 长连接稳定性测试：保持 SSE 连接 5 分钟，验证心跳消息、连接稳定
6. Docker 容器测试：使用 docker-compose 启动、从宿主机访问、验证环境变量配置
7. 远程访问测试（如果有远程服务器）：部署到远程服务器、从本地访问
8. 性能验证：SSE 连接建立时间 < 1 秒、工具调用响应时间与 stdio 模式相当
9. 测试结果文档：记录到 `.ai/epic4-sse-tests.md`，包含 curl 命令、响应示例
10. 里程碑确认：SSE 模式所有功能正常、Token 认证有效、Docker 镜像可用、双模式正常工作

---

### Epic 5: Documentation, Examples & Community Launch

**Epic Goal**: 完善项目文档体系，包括 README、安装指南、配置说明、使用示例和故障排查指南，使新用户能够在 5 分钟内快速上手。创建真实场景的配置示例和演示脚本，展示核心使用价值。准备并发布 GitHub Release（v1.0.0），包含多平台二进制文件和详细 Release Notes。发布 Docker Hub 镜像，简化部署流程。向目标社区（r/selfhosted、r/jellyfin）宣传项目，收集早期用户反馈，达成 MVP 成功指标（500+ 活跃用户、GitHub Stars 1000+）。

#### Story 5.1: Core Documentation and README

**As a** new user,
**I want** to read clear and comprehensive documentation in the README,
**so that** I can quickly understand what the project does and how to get started.

**Acceptance Criteria**:

1. README.md 结构包含：项目介绍、快速开始、功能特性、使用场景、配置说明、部署方式、开发、贡献指南、许可证、致谢
2. 添加徽章（Badges）：GitHub Stars、License、Go Version、Docker Pulls、Build Status
3. 添加截图/GIF：Claude Code 中使用演示、配置文件示例
4. 多语言支持（可选）：提供中文版 README（`README.zh-CN.md`）
5. 文档质量检查：使用 Markdown linter、确保链接有效、代码示例可运行、请他人审阅

#### Story 5.2: Usage Examples and Scenario Demonstrations

**As a** new user,
**I want** to see real-world usage examples and scenario demonstrations,
**so that** I can understand how to apply the tool to my specific needs.

**Acceptance Criteria**:

1. 创建 `examples/` 目录，包含：基础配置文件、完整配置文件、stdio 模式配置、SSE 模式配置、Docker Compose、Docker Compose with Nginx
2. 创建演示脚本 `examples/demo.sh`：自动化演示、调用所有 6 个工具
3. 创建文件重命名脚本示例 `examples/rename-movies.md`
4. 创建 Claude Code 配置示例 `examples/claude-code-config.json`
5. 创建故障排查场景示例 `examples/troubleshooting-scenarios.md`
6. 文档说明：在 README 添加 "Examples" 章节

#### Story 5.3: Configuration Guide and Troubleshooting

**As a** user,
**I want** to have a comprehensive configuration guide and troubleshooting documentation,
**so that** I can solve common problems independently.

**Acceptance Criteria**:

1. 创建配置指南 `docs/configuration.md`：配置文件详解、环境变量、命令行参数、配置优先级、常见配置场景
2. 创建故障排查指南 `docs/troubleshooting.md`：常见问题 FAQ（401/429/SSE 连接失败等）、日志分析、性能问题排查、获取帮助
3. 创建 API 参考文档 `docs/api-reference.md`（可选）：每个 MCP 工具的详细 API 文档
4. 文档质量保证：所有错误场景实际测试、解决方法确认有效、请早期用户审阅

#### Story 5.4: GitHub Release Preparation and Publishing

**As a** project maintainer,
**I want** to prepare and publish a GitHub Release (v1.0.0),
**so that** users can easily download and install the software.

**Acceptance Criteria**:

1. 版本标记：在代码中添加版本常量、更新 README 版本号、Git tag: `v1.0.0`
2. 编译多平台二进制文件：Linux AMD64/ARM64、macOS AMD64/ARM64、Windows AMD64
3. 打包发布文件：为每个平台创建 tar.gz/zip 压缩包
4. 计算校验和：为每个压缩包生成 SHA256 校验和、创建 `checksums.txt`
5. 编写 Release Notes（`RELEASE_NOTES.md`）：Highlights、Features、Quick Start、Known Issues、Acknowledgments、Full Changelog
6. 创建 CHANGELOG.md：遵循 Keep a Changelog 格式
7. 发布到 GitHub Releases：创建 Release v1.0.0、上传所有二进制压缩包和 checksums.txt、标记为 "Latest release"
8. 验证下载链接：测试每个平台的下载、验证校验和、测试二进制可运行

#### Story 5.5: Docker Hub Image Publishing

**As a** user,
**I want** to easily pull and run the tmdb-mcp Docker image from Docker Hub,
**so that** I can quickly deploy without building from source.

**Acceptance Criteria**:

1. 注册 Docker Hub 账号（如未有）、创建 repository: `username/tmdb-mcp`
2. 构建多平台 Docker 镜像：使用 Docker Buildx 构建 `linux/amd64`, `linux/arm64`, `linux/arm/v7`
3. 镜像标签策略：`latest`, `v1.0.0`, `v1.0`, `v1`
4. 更新 Docker Hub 描述：项目简介、快速开始命令、链接到 README、环境变量说明
5. 添加 README.md 到 Docker Hub
6. 测试镜像拉取和运行：`docker pull`、`docker run`、验证健康检查
7. 验证多平台支持：在 AMD64 和 ARM64 机器上测试、验证镜像大小合理（< 50MB）
8. 文档更新：在 README 添加 Docker Hub 徽章、更新 Docker 安装命令

#### Story 5.6: Community Launch and Early User Recruitment

**As a** project maintainer,
**I want** to launch the project to relevant communities and recruit early users,
**so that** I can gather feedback and build a user base.

**Acceptance Criteria**:

1. 准备社区发布内容：Reddit 发布帖模板（r/selfhosted, r/jellyfin）
2. 准备演示材料：录制简短演示视频（1-2 分钟）或 GIF
3. 社区发布计划：第 1 天 Reddit 发布、第 2-3 天论坛发布（Jellyfin Forum、Hacker News）、第 4-7 天博客和社交媒体
4. 设置反馈收集机制：GitHub Discussions 启用、创建 Discussion 类别
5. 早期用户招募：在发布帖中征集、提供 "Early Adopters" 标签、承诺快速响应
6. 监控和响应：第 1 周每天检查 GitHub Issues/Reddit 评论、及时回复、记录常见问题
7. 衡量指标跟踪（第 1 个月）：GitHub Stars、Docker Hub Pulls、Issues/Discussions 活跃度、实际用户反馈数量
8. 迭代计划：根据反馈识别 Top 3 优先功能、创建 v1.1.0 里程碑
9. 创建 "感谢早期用户" 文档 `CONTRIBUTORS.md`
10. 准备 MVP 成功报告：总结关键指标、识别成功点和改进空间、规划长期路线图

---

## Checklist Results Report

### Executive Summary

**Overall PRD Completeness**: **92%**

**MVP Scope Appropriateness**: **Just Right** ✅
- 功能范围聚焦核心价值（6 个 MCP 工具）
- 双模式支持（stdio + sse）合理，不过度设计
- 明确排除了缓存、批量处理等非核心功能
- 5 个 Epic 的时间估算现实（3-4 周 MVP）

**Readiness for Architecture Phase**: **READY** ✅

**Most Critical Gaps**:
- ✅ **无阻塞性缺陷**
- ⚠️ Minor: 缺少用户研究证据（基于项目简报的假设）
- ⚠️ Minor: 测试策略可以更详细（Epic 1-3 中已部分覆盖）

### Category Analysis Table

| Category                         | Status     | Critical Issues | Notes |
| -------------------------------- | ---------- | --------------- | ----- |
| 1. Problem Definition & Context  | **PASS**   | None            | 清晰定义了自建流媒体用户的文件管理痛点 |
| 2. MVP Scope Definition          | **PASS**   | None            | 6 个工具 + 双模式，边界清晰，Out of Scope 明确 |
| 3. User Experience Requirements  | **PASS**   | None            | CLI/MCP 服务，无传统 UI，已正确跳过 UI Design Goals |
| 4. Functional Requirements       | **PASS**   | None            | 9 个 FR 详细且可测试，覆盖所有核心功能 |
| 5. Non-Functional Requirements   | **PASS**   | None            | 12 个 NFR 全面，包括性能、安全、配置管理 |
| 6. Epic & Story Structure        | **PASS**   | None            | 5 个 Epic，30+ Stories，AC 详细且可执行 |
| 7. Technical Guidance            | **PASS**   | None            | 技术栈明确（Golang, net/http, MCP SDK, Resty, Viper, Zap） |
| 8. Cross-Functional Requirements | **PASS**   | None            | 集成（TMDB API）、运维（Docker）、监控（日志）已覆盖 |
| 9. Clarity & Communication       | **PASS**   | None            | 结构清晰，术语一致，中文文档流畅 |

**Overall Status**: **9/9 PASS** 🎉

### Technical Readiness

**Identified Technical Risks**:

1. **MCP Go SDK 的 SSE 支持** (MEDIUM RISK)
   - 风险：官方 SDK 可能不原生支持 SSE transport
   - 影响：Epic 4 的 Story 4.4 实现复杂度增加
   - 缓解：架构师优先调研、如不支持需设计 SSE 适配器、备选方案：WebSocket

2. **LLM 工具理解能力** (LOW RISK)
   - 风险：Claude 等 LLM 可能无法有效理解 discover 工具的复杂参数
   - 影响：内容发现场景效果不佳
   - 缓解：工具描述中提供清晰示例、Epic 3, Story 3.5 中验证

3. **TMDB API 稳定性** (LOW RISK)
   - 风险：TMDB API 可能变更或限制加严
   - 影响：服务不可用
   - 缓解：已实现错误处理和速率限制、监控 TMDB 官方公告

### Final Decision

## ✅ **READY FOR ARCHITECT**

**Summary**: PRD 和 Epic 定义**全面、结构清晰、可执行性强**，已准备好进入架构设计阶段。

**Key Strengths**:
- ✅ 问题定义清晰，目标用户具体
- ✅ MVP 范围精准，聚焦核心价值
- ✅ 功能和非功能需求完整且可测试
- ✅ Epic 和 Story 结构合理，大小适中
- ✅ 技术栈明确，约束清晰
- ✅ 无阻塞性缺陷

**Critical Success Factors**:
1. 架构师需优先调研 MCP Go SDK 的 SSE 支持
2. Epic 3, Story 3.5（E2E 测试）中验证 LLM 工具理解能力
3. Epic 5 社区发布后，收集用户反馈验证假设

**PRD Completeness Score: 92/100** 🎉

---

## Next Steps

### UX Expert Prompt

**N/A - No Traditional UI Required**

本项目为 CLI/MCP 服务，无传统用户界面。用户体验通过以下方式实现：

1. **对话式交互**：通过 Claude Code 等 MCP 客户端，用户使用自然语言与服务交互
2. **CLI 输出**：结构化日志和错误消息提供反馈
3. **文档 UX**：README、配置指南、示例确保用户快速上手

**用户体验责任**：
- **PM**（已完成）：定义 4 个核心使用场景和用户流程
- **Dev**（待实现）：确保错误消息清晰、日志有用、配置简单
- **文档作者**（Epic 5）：编写清晰的快速开始指南和故障排查文档

### Architect Prompt

**开始架构设计，使用本 PRD 作为输入。**

---

## 📐 架构师任务：TMDB MCP 服务架构设计

### 背景

你将为 **TMDB MCP 服务** 设计技术架构。这是一个基于 Golang 的 MCP (Model Context Protocol) 服务器，将 TMDB 电影数据库与 LLM（如 Claude）连接，使用户能够通过自然语言查询影视内容。

**核心文档**：
- 📄 PRD：`docs/prd.md`（本文档）
- 📄 项目简报：`docs/brief.md`
- 📄 TMDB API 文档：`docs/tmdb-api.md`

### 你的任务

设计完整的技术架构，包括：

1. **系统架构设计**
   - 分层架构（MCP 层、工具层、TMDB 客户端层、速率限制层）
   - stdio 和 SSE 模式的代码共享策略
   - 项目结构（cmd/, internal/, pkg/）

2. **关键技术调研**
   - MCP SDK 的 `SSEHTTPHandler` 使用模式和最佳实践
   - 标准库 `net/http` 中间件模式实现 Bearer Token 认证
   - Viper 的配置优先级实现（CLI > ENV > File）

3. **核心组件设计**
   - TMDB API 客户端（Resty、速率限制、错误处理）
   - MCP 工具注册和调用机制
   - Token 生成和认证中间件（标准库）
   - 配置管理（多源、优先级、持久化）
   - 日志系统（Zap、结构化日志、性能监控）

4. **数据流设计**
   - stdio 模式：stdin/stdout → MCP handler → 工具 → TMDB API
   - SSE 模式：HTTP request → 认证中间件 → `SSEHTTPHandler` → MCP handler → 工具 → TMDB API

5. **错误处理策略**
   - TMDB API 错误（401/404/429）的统一处理
   - 重试逻辑（429 自动重试）
   - MCP 错误响应格式

6. **性能和并发设计**
   - 速率限制器实现（40 req/10s）
   - 并发请求处理（goroutines）
   - 性能监控（响应时间、调用计数）

7. **安全设计**
   - Token 生成（crypto/rand）
   - 配置文件权限（600）
   - API Key 管理（不硬编码）

8. **部署设计**
   - Docker 镜像（多阶段构建）
   - 多平台二进制编译
   - 配置文件路径（`~/.tmdb-mcp/config.yaml`）

### 关键约束

**必须遵守**：
- ✅ Golang 1.21+
- ✅ TMDB API 速率限制：40 requests / 10 seconds
- ✅ MCP 协议：stdio + SSE 双模式
- ✅ 跨平台支持：Linux / macOS / Windows
- ✅ 性能目标：P95 响应时间 < 500ms，启动时间 < 2 秒
- ✅ 精简原则：不使用 Makefile、golangci-lint，仅用 Go 原生工具链

**技术栈**：
- MCP SDK：`github.com/modelcontextprotocol/go-sdk`（内置 SSE 支持）
- HTTP 服务器：标准库 `net/http`
- HTTP 客户端：Resty
- 配置管理：Viper
- 日志：Zap
- 速率限制：`golang.org/x/time/rate`

### 关键风险点

**请优先关注**：

1. **SSEHTTPHandler 与标准库中间件集成** ⚠️ MEDIUM RISK
   - 风险：需要确保认证中间件能正确包装 `SSEHTTPHandler`
   - 影响：Epic 4 的 Story 4.3 和 4.4 实现复杂度
   - 缓解：参考 MCP SDK 示例中的认证中间件模式、早期原型验证

2. **LLM 工具理解能力** (LOW RISK)
   - 风险：Claude 等 LLM 可能无法有效理解 discover 工具的复杂参数
   - 影响：内容发现场景效果不佳
   - 缓解：工具描述中提供清晰示例、Epic 3, Story 3.5 中验证

3. **TMDB API 稳定性** (LOW RISK)
   - 风险：TMDB API 可能变更或限制加严
   - 影响：服务不可用
   - 缓解：已实现错误处理和速率限制、监控 TMDB 官方公告

4. **配置优先级实现**
   - 风险：Viper 的配置优先级逻辑需要仔细设计
   - 影响：用户配置体验
   - 缓解：参考 Viper 官方文档、编写完整的配置管理测试

### 交付物

请创建以下文档：

1. **架构文档** (`docs/architecture.md` 或 `docs/architecture/`)：
   - 系统架构图（分层架构）
   - 数据流图（stdio 和 SSE 模式）
   - 目录结构设计
   - 关键设计决策和理由

2. **SSE 集成方案**：
   - `SSEHTTPHandler` 使用模式
   - 认证中间件与 `SSEHTTPHandler` 的集成方案
   - 连接管理和心跳机制说明

3. **接口设计**：
   - 核心 struct 和 interface 定义
   - 配置结构体
   - TMDB API 客户端接口

4. **部署方案**：
   - Dockerfile 设计
   - 多平台编译脚本（如需）
   - 环境变量清单

### 成功标准

✅ 架构设计清晰，开发者可直接按设计实现
✅ SSE 集成方案明确，利用 MCP SDK 的 `SSEHTTPHandler`
✅ 性能和安全要求已体现在设计中
✅ 所有技术风险已识别并有缓解方案
✅ 架构文档完整，可传递给开发团队

### 参考资源

- MCP 协议规范：https://spec.modelcontextprotocol.io/
- MCP Go SDK：https://github.com/modelcontextprotocol/go-sdk
- **MCP Go SDK 文档**：`docs/mcp-go-sdk.md`（本地文档，包含 SSE 支持详情）
- TMDB API v3：https://developers.themoviedb.org/3
- 项目简报：`docs/brief.md`

---

**开始架构设计吧！如有任何问题或需要澄清 PRD 内容，请随时询问。**

---

*文档生成日期：2025-10-10*
*作者：John (Product Manager)*
*使用 BMAD™ 方法创建*
