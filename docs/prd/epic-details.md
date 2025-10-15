# Epic Details

## Epic 1: Foundation & Search

**Epic Goal**: 建立 TMDB MCP 服务的核心技术基础设施，包括配置管理、日志系统、TMDB API 客户端封装、速率限制机制和 MCP stdio 协议集成。实现并交付第一个可工作的 MCP 工具（search），使用户能够通过 Claude 等 LLM 客户端搜索电影、电视剧和人物，验证整个技术栈的可行性并为后续功能打下坚实基础。

### Story 1.1: Project Initialization and Configuration Management

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

### Story 1.2: Structured Logging System

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

### Story 1.3: TMDB API Client Foundation

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

### Story 1.4: Rate Limiting Mechanism

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

### Story 1.5: MCP Protocol Integration via stdio

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

### Story 1.6: Implement Search Tool

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

### Story 1.7: Automated End-to-End Integration Testing

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

## Epic 2: Details & Discovery Tools

**Epic Goal**: 实现内容详情获取工具（get_details）和内容发现工具（discover_movies、discover_tv），使用户能够获取电影/电视剧/人物的完整信息，并通过自然语言表达的筛选条件探索新内容。完善错误处理机制，优雅处理 TMDB API 的各类错误（401/404/429）。此 Epic 完成后，核心工具集功能完备，自动化测试验证所有使用场景的可行性。

### Story 2.1: Implement get_details Tool

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

### Story 2.2: Implement discover_movies Tool

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

### Story 2.3: Implement discover_tv Tool

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

### Story 2.4: Enhanced Error Handling and Retry Logic

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

### Story 2.5: (已删除)

**说明**: 此 Story 已移除。原 Story 2.5 依赖外部工具（Claude Code）进行手动测试，导致测试不可靠且无法自动化。文件重命名场景的测试已由 Story 1.7（自动化集成测试）和 Story 3.4（综合集成测试）充分覆盖。

**文件重命名场景的自动化测试策略**:
- Story 1.7 已使用 InMemoryTransports 模拟完整 MCP 协议交互
- Story 3.4 提供多工具组合测试（search → get_details）
- 文档示例将在 Epic 5 (Story 5.2) 中提供

---

## Epic 3: Trending & Recommendations

**Epic Goal**: 实现最后两个 MCP 工具（get_trending 和 get_recommendations），完成所有 6 个核心工具的功能集。添加性能监控和指标记录，优化 API 调用效率和响应时间。通过全面的自动化集成测试验证所有工具协同工作，确保性能指标符合 NFR 要求。此 Epic 完成后，stdio 模式功能完整可用，达到 MVP 核心功能完整的里程碑。可选地准备演示材料用于后续文档和社区宣传。

### Story 3.1: Implement get_trending Tool

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

### Story 3.2: Implement get_recommendations Tool

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

### Story 3.3: Performance Monitoring and Metrics

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

### Story 3.4: Comprehensive Integration Testing

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

### Story 3.5: End-to-End Scenario Validation (Optional Documentation)

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

## Epic 4: SSE Remote Access Mode

**Epic Goal**: 在现有 stdio 模式的基础上，使用 MCP Go SDK 提供的 `SSEHTTPHandler` 实现 Server-Sent Events (SSE) 远程访问模式，提供 HTTP API 端点供远程客户端连接。实现 Bearer Token 认证中间件（基于标准库 `net/http`）保护 SSE 端点，支持通过环境变量和配置文件管理 Token。支持 stdio 和 sse 双模式同时运行，并完成 Docker 镜像构建和多平台二进制文件编译，使服务可以方便地部署到远程服务器或容器环境中。

### Story 4.1: HTTP Server Setup with Standard Library

**As a** developer,
**I want** to set up a basic HTTP server using standard library `net/http`,
**so that** I can provide HTTP endpoints for SSE connections and health checks.

**Acceptance Criteria**:

1. 创建 `internal/server` 包，实现 HTTP 服务器，结构体 `HTTPServer` 包含 `http.Server`、配置、MCP server 引用
2. 实现 `NewHTTPServer(config Config, mcpServer *mcp.Server)` 构造函数
3. 配置 `http.Server`：设置监听地址、读写超时、集成 zap logger 记录 HTTP 请求
4. 实现 `/health` 端点（无需认证）：返回 `{"status": "ok", "version": "1.0.0", "mode": "sse"}`，使用标准 `http.HandlerFunc`
5. 实现服务器启动和优雅关闭：`Start()` 和 `Stop(ctx)` 方法，支持 SIGINT/SIGTERM 信号处理
6. 更新配置结构体，添加 SSE 相关配置（host, port, token）
7. 编写单元测试：测试 server 启动/停止、`/health` 端点
8. 编写集成测试：启动服务器，调用 `/health`，验证 200 OK

### Story 4.2: Token Generation and Management

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

### Story 4.3: Bearer Token Authentication Middleware

**As a** developer,
**I want** to implement Bearer Token authentication middleware using standard library `net/http`,
**so that** only authorized clients can access the SSE endpoint.

**Acceptance Criteria**:

1. 实现标准库中间件 `AuthMiddleware(expectedToken string) func(http.Handler) http.Handler`
2. 认证逻辑：提取 `Authorization` header、验证格式 `Bearer <token>`、使用 `crypto/subtle.ConstantTimeCompare` 比对 token
3. 认证成功：调用 `next.ServeHTTP(w, r)`、记录 DEBUG 日志
4. 认证失败：返回 `401 Unauthorized`、JSON 响应 `{"error": "unauthorized"}`、记录 WARN 日志
5. 错误场景处理：缺少 header、格式错误、token 不匹配
6. 将中间件应用到 SSE 路由（不应用到 `/health`）
7. 编写单元测试：测试有效/无效 token、缺少 header、`/health` 不需要认证
8. 编写集成测试：使用正确/错误 token 访问 SSE 端点

### Story 4.4: Implement SSE Endpoint with MCP SDK

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

### Story 4.5: Dual Mode Support (stdio + sse)

**As a** developer,
**I want** to support running both stdio and sse modes simultaneously,
**so that** users can choose their preferred connection method or use both.

**Acceptance Criteria**:

1. 模式配置：`server.mode` 支持三个值：`stdio`, `sse`, `both`（默认）
2. stdio 模式实现：启动 MCP server，监听 stdin/stdout，阻塞主 goroutine
3. sse 模式实现：启动 HTTP server（非阻塞，使用 goroutine），监听端口 8910
4. both 模式实现：同时启动 stdio 和 HTTP server，共享 TMDB client 和工具实现
5. 优雅关闭：捕获 SIGINT/SIGTERM 信号、同时关闭两个 server、等待活跃连接完成（最多 10 秒超时）
6. 日志记录：启动时记录启用的模式
7. 配置验证：如果 mode="sse" 但 `enabled=false`，返回错误
8. 编写单元测试：测试每种模式的启动逻辑
9. 编写集成测试：启动 stdio 模式、sse 模式、both 模式，验证优雅关闭

### Story 4.6: Docker Image and Multi-Platform Build

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

### Story 4.7: SSE Mode End-to-End Testing

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

## Epic 5: Documentation, Examples & Community Launch

**Epic Goal**: 完善项目文档体系，包括 README、安装指南、配置说明、使用示例和故障排查指南，使新用户能够在 5 分钟内快速上手。创建真实场景的配置示例和演示脚本，展示核心使用价值。准备并发布 GitHub Release（v1.0.0），包含多平台二进制文件和详细 Release Notes。发布 Docker Hub 镜像，简化部署流程。向目标社区（r/selfhosted、r/jellyfin）宣传项目，收集早期用户反馈，达成 MVP 成功指标（500+ 活跃用户、GitHub Stars 1000+）。

### Story 5.1: Core Documentation and README

**As a** new user,
**I want** to read clear and comprehensive documentation in the README,
**so that** I can quickly understand what the project does and how to get started.

**Acceptance Criteria**:

1. README.md 结构包含：项目介绍、快速开始、功能特性、使用场景、配置说明、部署方式、开发、贡献指南、许可证、致谢
2. 添加徽章（Badges）：GitHub Stars、License、Go Version、Docker Pulls、Build Status
3. 添加截图/GIF：Claude Code 中使用演示、配置文件示例
4. 多语言支持（可选）：提供中文版 README（`README.zh-CN.md`）
5. 文档质量检查：使用 Markdown linter、确保链接有效、代码示例可运行、请他人审阅

### Story 5.2: Usage Examples and Scenario Demonstrations

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

### Story 5.3: Configuration Guide and Troubleshooting

**As a** user,
**I want** to have a comprehensive configuration guide and troubleshooting documentation,
**so that** I can solve common problems independently.

**Acceptance Criteria**:

1. 创建配置指南 `docs/configuration.md`：配置文件详解、环境变量、命令行参数、配置优先级、常见配置场景
2. 创建故障排查指南 `docs/troubleshooting.md`：常见问题 FAQ（401/429/SSE 连接失败等）、日志分析、性能问题排查、获取帮助
3. 创建 API 参考文档 `docs/api-reference.md`（可选）：每个 MCP 工具的详细 API 文档
4. 文档质量保证：所有错误场景实际测试、解决方法确认有效、请早期用户审阅

### Story 5.4: GitHub Release Preparation and Publishing

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

### Story 5.5: Docker Hub Image Publishing

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

### Story 5.6: Community Launch and Early User Recruitment

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
