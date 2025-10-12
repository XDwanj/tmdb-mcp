# 📐 架构师任务：TMDB MCP 服务架构设计

## 背景

你将为 **TMDB MCP 服务** 设计技术架构。这是一个基于 Golang 的 MCP (Model Context Protocol) 服务器，将 TMDB 电影数据库与 LLM（如 Claude）连接，使用户能够通过自然语言查询影视内容。

**核心文档**：
- 📄 PRD：`docs/prd.md`（本文档）
- 📄 项目简报：`docs/brief.md`
- 📄 TMDB API 文档：`docs/tmdb-api.md`

## 你的任务

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

## 关键约束

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

## 关键风险点

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

## 交付物

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

## 成功标准

✅ 架构设计清晰，开发者可直接按设计实现
✅ SSE 集成方案明确，利用 MCP SDK 的 `SSEHTTPHandler`
✅ 性能和安全要求已体现在设计中
✅ 所有技术风险已识别并有缓解方案
✅ 架构文档完整，可传递给开发团队

## 参考资源

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
