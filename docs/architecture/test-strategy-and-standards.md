# Test Strategy and Standards

## Testing Philosophy

**Approach**: **Test-After**（非 TDD）

- 先实现功能，后编写测试
- 测试覆盖核心业务逻辑
- 集成测试验证端到端流程

**Coverage Goals**:
- **Unit Tests**: 核心业务逻辑覆盖率 ≥ 70%
- **Integration Tests**: 覆盖所有 6 个 MCP 工具
- **E2E Tests**: 手动测试 4 个核心使用场景

**Test Pyramid**:
- 70% Unit Tests（快速，隔离）
- 25% Integration Tests（真实 API 或 Mock）
- 5% E2E Tests（手动，真实 Claude Code 客户端）

## Test Types and Organization

### Unit Tests

**Framework**: `testing` (标准库)

**File Convention**: `*_test.go`（与被测试文件同目录）

**Location**: 与源代码同目录（`internal/tmdb/client_test.go`）

**Mocking Library**: `github.com/stretchr/testify/mock`

**Coverage Requirement**: ≥ 70% for `internal/tmdb`, `internal/tools`, `internal/config`

**AI Agent Requirements**:
- 为所有公共函数生成测试
- 覆盖正常情况、边界情况和错误情况
- 使用 Table-driven tests 处理多场景
- Mock 所有外部依赖（HTTP 调用、文件 I/O）

**Example**:
```go
func TestClient_Search(t *testing.T) {
    tests := []struct {
        name    string
        query   string
        page    int
        want    int  // 期望结果数量
        wantErr bool
    }{
        {"valid query", "Inception", 1, 10, false},
        {"empty query", "", 1, 0, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Mock HTTP client
            // ... 测试逻辑
        })
    }
}
```

### Integration Tests

**Scope**: 测试多个组件协同工作（Client + Rate Limiter + Real/Mock TMDB API）

**Location**: `internal/tmdb/integration_test.go`

**Test Infrastructure**:
- **TMDB API**: 使用真实 TMDB API（需要有效 API Key）或 HTTP Mock Server
- **Rate Limiter**: 真实 Rate Limiter，验证速率控制
- **MCP Server**: 真实 MCP SDK，验证协议兼容性

**Example**:
```go
// +build integration

func TestSearchIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    config := loadTestConfig()  // 从环境变量加载测试配置
    client := tmdb.NewClient(config, logger)

    results, err := client.Search(context.Background(), "Inception", 1)
    assert.NoError(t, err)
    assert.NotEmpty(t, results)
}
```

### End-to-End Tests

**Framework**: 手动测试（使用真实 Claude Code 客户端）

**Scope**: 4 个核心使用场景

**Environment**: 真实 TMDB MCP 服务 + Claude Code

**Test Data**: 真实 TMDB 数据

**Test Scenarios**:
1. 智能文件重命名（search + get_details）
2. 片荒推荐（get_trending + get_details）
3. 关联探索（discover_movies + get_recommendations）
4. 智能推荐（基于特定电影的推荐链）

## Test Data Management

**Strategy**: 使用真实 TMDB API 数据（集成测试）或 Mock 数据（单元测试）

**Fixtures**: `testdata/` 目录存放 Mock API 响应 JSON

**Factories**: 不需要（无数据库，无复杂对象构造）

**Cleanup**: 不需要（只读操作，无副作用）

## Continuous Testing

**CI Integration**:
- GitHub Actions 在每次 Push 和 PR 时运行测试
- 单元测试：`go test ./...`
- 集成测试：`go test -tags=integration ./...`（需要 TMDB API Key）

**Performance Tests**: 不在 MVP 范围内（未来可添加 Benchmark）

**Security Tests**: 不在 MVP 范围内（未来可添加 SAST 扫描）

---
