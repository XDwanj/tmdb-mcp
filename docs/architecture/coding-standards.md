# Coding Standards

**说明**: 以下标准为 **MANDATORY for AI agents**，开发时必须严格遵守。

## Core Standards

**Languages & Runtimes**:
- Go 1.21+
- 使用 `go fmt` 格式化所有代码（提交前必须运行）
- 使用 `go vet` 进行静态检查

**Style & Linting**:
- 遵循 [Effective Go](https://go.dev/doc/effective_go) 规范
- 不使用第三方 linter（golangci-lint 等）
- 包名小写单词，不使用下划线或驼峰

**Test Organization**:
- 测试文件命名：`*_test.go`
- 测试函数命名：`TestFunctionName`（遵循 Go 约定）
- Table-driven tests 用于多场景测试

## Naming Conventions

| Element | Convention | Example |
|---------|-----------|---------|
| 包名 | 小写单词 | `package tmdb` |
| 文件名 | 小写蛇形命名 | `get_details.go` |
| 结构体 | 大驼峰（exported）| `type TMDBClient struct` |
| 接口 | 大驼峰，-er 后缀 | `type Limiter interface` |
| 函数 | 大驼峰（exported）| `func NewClient()` |
| 变量 | 小驼峰 | `var apiKey string` |
| 常量 | 大驼峰或全大写 | `const DefaultPort = 8910` |

## Critical Rules

- **日志规则**: 永远不要使用 `fmt.Println` 或 `log.Println`，始终使用 Zap logger（`logger.Info()`, `logger.Error()` 等）

- **错误处理**: 永远不要忽略 error 返回值，必须检查或明确使用 `_` 表示忽略

- **Context 传递**: 所有需要取消或超时控制的函数必须接受 `context.Context` 作为第一个参数

- **配置管理**: 永远不要硬编码 API Key、Token 等敏感信息，必须从配置文件或环境变量读取

- **JSON Tags**: 所有需要 JSON 序列化的结构体字段必须添加 `json` tag，使用小写蛇形命名

- **依赖注入**: 使用构造函数注入依赖（`NewClient(config, logger)`），不使用全局变量

- **错误包装**: 使用 `fmt.Errorf("context: %w", err)` 包装错误，保留原始错误链

## Go-Specific Guidelines

**Error Handling**:
```go
// ✅ Good
result, err := client.Search(ctx, query)
if err != nil {
    logger.Error("search failed", zap.Error(err))
    return nil, fmt.Errorf("search failed: %w", err)
}

// ❌ Bad
result, _ := client.Search(ctx, query)  // 忽略错误
```

**Context Usage**:
```go
// ✅ Good
func (c *Client) Search(ctx context.Context, query string) ([]Result, error) {
    // context 作为第一个参数
}

// ❌ Bad
func (c *Client) Search(query string, ctx context.Context) ([]Result, error) {
    // context 不是第一个参数
}
```

**Struct Initialization**:
```go
// ✅ Good
client := &Client{
    apiKey:   config.APIKey,
    language: config.Language,
    logger:   logger,
}

// ❌ Bad
client := new(Client)
client.apiKey = config.APIKey  // 逐字段赋值
```

---
