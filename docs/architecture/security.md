# Security

## Input Validation

**Validation Library**: Go 标准库（无第三方库）

**Validation Location**: Tools 层（每个 MCP 工具的 `Call` 方法入口）

**Required Rules**:
- 所有 MCP 工具参数必须验证（query 非空、id > 0、media_type 有效值等）
- 验证在 API 边界（工具入口），处理前验证
- 使用白名单验证（media_type 必须是 "movie"/"tv"/"person" 之一）

## Authentication & Authorization

**Auth Method**: Bearer Token（SSE 模式）

**Session Management**: 无会话（无状态服务）

**Required Patterns**:
- SSE 端点必须验证 `Authorization: Bearer <token>` header
- 使用 `crypto/subtle.ConstantTimeCompare` 防止时序攻击
- Token 验证失败返回 `401 Unauthorized`

## Secrets Management

**Development**: 环境变量（`TMDB_API_KEY`, `SSE_TOKEN`）

**Production**: 配置文件 `~/.tmdb-mcp/config.yaml`（权限 600）

**Code Requirements**:
- ❌ 永远不要硬编码 API Key 或 Token
- ✅ 仅通过配置服务访问敏感信息
- ❌ 不在日志或错误消息中打印完整密钥（仅显示前 8 个字符）

## API Security

**Rate Limiting**: 40 req/10s（Token Bucket）

**CORS Policy**: N/A（无浏览器客户端）

**Security Headers**:
- `X-Content-Type-Options: nosniff`（SSE 响应）
- HTTPS 由用户反向代理实现（Nginx/Caddy）

**HTTPS Enforcement**: 用户自行配置反向代理，服务本身仅提供 HTTP

## Data Protection

**Encryption at Rest**: N/A（无持久化数据）

**Encryption in Transit**:
- TMDB API 调用：HTTPS
- SSE 连接：HTTP（建议用户配置反向代理 + HTTPS）

**PII Handling**: 无 PII 数据（仅 TMDB 公开数据）

**Logging Restrictions**:
- ❌ 不记录完整 API Key
- ❌ 不记录 SSE Token 明文
- ✅ 可记录 API 请求 URL（不含 api_key 参数）

## Dependency Security

**Scanning Tool**: `go list -m all` 查看依赖，手动审查

**Update Policy**: 每 3 个月更新依赖（除非有安全漏洞）

**Approval Process**:
- 新依赖必须是知名库（GitHub Stars > 1000）
- 新依赖必须有理由说明（无法用标准库实现）

## Security Testing

**SAST Tool**: 不使用（MVP 范围外）

**DAST Tool**: 不使用（MVP 范围外）

**Penetration Testing**: 不进行（单用户自托管工具）

---

**架构文档核心部分已完成！** 🎉

这份文档涵盖了：
✅ 高层架构设计（分层单体架构）
✅ 技术栈选择（Go 1.21, MCP SDK, Resty, Viper, Zap）
✅ 数据模型（配置、TMDB API 响应、MCP 参数）
✅ 组件设计（8 个核心组件及其职责）
✅ 外部 API 集成（TMDB API v3）
✅ 项目目录结构（标准 Go 布局）
✅ 基础设施和部署（Docker, 多平台二进制）
✅ 错误处理策略（重试、日志、错误转换）
✅ 编码标准（强制规则，AI 代理必须遵守）
✅ 测试策略（单元/集成/E2E 测试）
✅ 安全设计（Token 认证、密钥管理、输入验证）

让我保存文档并显示完成状态！
