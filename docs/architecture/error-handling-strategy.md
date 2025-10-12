# Error Handling Strategy

## General Approach

**Error Model**: Go 标准错误处理 + 自定义错误类型

**Exception Hierarchy**:
- `TMDBError` - TMDB API 错误
- `ConfigError` - 配置错误
- `ValidationError` - 参数验证错误

**Error Propagation**: 使用 `error` 返回值，层层向上传递，在适当层次记录日志

## Logging Standards

**Library**: `go.uber.org/zap` v1.26.0+

**Format**:
- **Development**: Console encoder，彩色输出
- **Production**: JSON encoder，结构化日志

**Levels**:
- **DEBUG**: 详细调试信息（API 请求参数、速率限制等待）
- **INFO**: 正常运行信息（服务启动、配置加载成功）
- **WARN**: 警告信息（速率限制触发、重试操作）
- **ERROR**: 错误信息（API 调用失败、配置验证失败）

**Required Context**:
- **Correlation ID**: 无需（单用户场景，无分布式追踪需求）
- **Service Context**: 服务名 `tmdb-mcp`，版本号
- **User Context**: 无需（无用户身份，MCP 客户端透明）

**Log Filtering**:
- ❌ **禁止记录**: 完整 API Key（仅显示前 8 个字符）
- ❌ **禁止记录**: SSE Token 明文
- ✅ **允许记录**: TMDB API 请求 URL（不含 api_key 参数）
- ✅ **允许记录**: 错误消息、响应时间、HTTP 状态码

## Error Handling Patterns

### External API Errors (TMDB API)

**Retry Policy**:
- **401 Unauthorized**: 不重试，立即返回错误
- **404 Not Found**: 不重试，返回空结果
- **429 Rate Limit**: 重试最多 3 次，等待 `Retry-After` header 指定时间
- **500/502/503**: 重试最多 3 次，指数退避（1s, 2s, 4s）
- **Network Timeout**: 不重试，返回错误

**Circuit Breaker**: 不实现（MVP 范围外，单用户场景无需熔断）

**Timeout Configuration**:
- **HTTP 请求超时**: 10 秒
- **Rate Limiter 等待超时**: 无限制（由 context 控制）

**Error Translation**:
```go
func (c *Client) handleError(resp *resty.Response) error {
    if resp.IsSuccess() {
        return nil
    }

    var tmdbErr TMDBError
    _ = json.Unmarshal(resp.Body(), &tmdbErr)
    tmdbErr.StatusCode = resp.StatusCode()

    switch resp.StatusCode() {
    case 401:
        return fmt.Errorf("TMDB API authentication failed: invalid API key")
    case 404:
        return fmt.Errorf("resource not found")
    case 429:
        retryAfter := resp.Header().Get("Retry-After")
        return fmt.Errorf("rate limit exceeded, retry after %s seconds", retryAfter)
    default:
        return &tmdbErr
    }
}
```

### Business Logic Errors

**Custom Exceptions**:
- `ErrInvalidMediaType` - media_type 参数无效
- `ErrInvalidTimeWindow` - time_window 参数无效
- `ErrMissingRequiredParam` - 缺少必需参数

**User-Facing Errors**: 返回清晰的英文错误消息（MCP 协议规定）

**Error Codes**: 使用 HTTP 状态码语义（400 Bad Request, 404 Not Found, 500 Internal Server Error）

### Data Consistency

**Transaction Strategy**: 无事务需求（无数据库，纯 API 转发）

**Compensation Logic**: 不需要（无副作用操作）

**Idempotency**: 所有 MCP 工具调用都是幂等的（只读操作）

---
