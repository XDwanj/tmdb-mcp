# Streamable HTTP 远程访问模式 - Brownfield Enhancement

## Epic Goal
在保持现有 stdio/SSE 通信能力和部署流程不变的前提下，引入 MCP Streamable HTTP 传输模式（基于 `StreamableHTTPHandler` 与 `StreamableServerTransport`），让远程客户端能够通过单一 HTTP 端点获取可恢复的增量消息流，并支持断线重连与会话恢复。

## Epic Description

**Existing System Context:**
- 当前 MCP 服务已经实现 stdio 与 SSE 两种连接方式，所有远程访问均复用标准库 `net/http` 与 Bearer Token 中间件。
- 技术栈基于 Go 1.21、`github.com/modelcontextprotocol/go-sdk` 官方 SDK、`zap` 日志、`viper` 配置、`resty` TMDB 客户端以及 `time/rate` 限流器。
- 现有远程模式通过 `/mcp/sse` 提供事件推送，配套健康检查、配置验证、Docker 镜像与 CI 任务。

**Enhancement Details:**
- 新增 Streamable HTTP 端点（建议 `/mcp/stream`），使用 `mcp.NewStreamableHTTPHandler` 暴露 `StreamableServerTransport` 会话，支持 `StreamableClientTransport` 客户端的断线重连能力。
- 复用现有认证模块：沿用 Bearer Token 校验逻辑，扩展配置（`server.streamable.enabled`、`server.streamable.token` 等）与 CLI flag/env 变量映射。
- 增强会话持久化：利用 SDK 提供的事件存储接口或轻量内存缓存保存最近会话事件，以满足 Streamable HTTP 的重放需求。
- 扩展部署与监控：在 Docker 镜像、Compose、文档与健康检查中加入 Streamable 模式说明，同时更新 Story、PRD 与使用示例。

## Stories
1. **Story 6.1:** 实现 `/mcp/stream` 端点与 `StreamableHTTPHandler` 集成，并复用 SSE Bearer Token 中间件保障访问控制。
2. **Story 6.2:** 扩展配置、CLI flag 与自动化测试，确保三种模式（stdio/SSE/Streamable）可独立或组合启用，并覆盖重连与异常处理场景。
3. **Story 6.3:** 更新文档、示例与监控指标，提供使用指南、Docker/Compose 配置片段，以及端到端集成测试验证会话恢复与多客户端并发。

## Compatibility Requirements
- [ ] 保持现有 stdio 与 SSE API 行为、配置项与默认值不变。
- [ ] 新增配置项保持向后兼容，默认禁用 Streamable 模式。
- [ ] 认证与限流模块复用既有实现，不引入额外依赖。
- [ ] 端点命名、日志与指标遵循既有命名约定，避免破坏监控脚本。

## Risk Mitigation
- **Primary Risk:** Streamable 会话重放或并发处理不当导致内存增长或连接泄漏。
- **Mitigation:** 采用 SDK 建议的事件存储接口，设置最大事件保留与超时回收策略；为 `StreamableHTTPHandler` 增加分级日志与指标监控。
- **Rollback Plan:** 如 Streamable 模式导致问题，可通过配置开关禁用该端点，并回滚到仅 stdio + SSE 的部署；保留独立镜像标签以便恢复。

## Definition of Done
- [ ] 所有相关故事完成，具备端到端 Streamable HTTP 功能与断线重连验证。
- [ ] CI 中新增/更新的集成测试稳定通过，覆盖成功、失败与重连路径。
- [ ] 文档（README、配置指南、使用示例）更新完毕，并含故障排查章节。
- [ ] Docker 镜像、Compose 样例更新，通过实际部署验证。
- [ ] 监控与日志指标新增 Streamable 维度或复用既有结构，确保可观测性。

## Validation Checklist

**Scope Validation:**
- [ ] Epic 范围可在 1-3 个故事内完成，不触发架构重写。
- [ ] 不改变核心 TMDB 工具逻辑，仅扩展传输层。
- [ ] 需要的 SDK 能力已在 `modelcontextprotocol/go-sdk` 中提供并确认可行。

**Risk Assessment:**
- [ ] 评估并记录内存、连接数、事件队列等资源指标。
- [ ] 为生产部署提供推荐的超时、并发与重试配置。
- [ ] 设计手动验证流程，模拟断线重连与并发负载。

**Completeness Check:**
- [ ] Epic 目标、成功指标与依赖被明确记录。
- [ ] 故事定义与验收标准在 Story 阶段可直接引用。
- [ ] 所有新配置项附带默认值、示例与回退策略。

---

**Story Manager Handoff:**  
请据此 Epic 编写详细用户故事。关键注意事项：  
- 复用 Go 官方 MCP SDK (`StreamableHTTPHandler`、`StreamableServerTransport`、`StreamableClientTransport`) 与现有 `net/http` 服务结构。  
- 远程模式需同时兼容 stdio、SSE 与 Streamable，配置层需保证向后兼容。  
- 认证与限流逻辑沿用现有 Token 中间件和速率限制，实现统一日志与指标。  
- 每个故事需验证 Streamable 会话的初始化、断线重连、错误处理与资源回收，并更新相应文档和示例。

