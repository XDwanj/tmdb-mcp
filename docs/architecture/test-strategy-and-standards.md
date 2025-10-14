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

**Test Pyramid** (基于 MCP SDK InMemoryTransports):
- 60% Unit Tests（快速，隔离，Mock 外部依赖）
- 35% Integration Tests（InMemoryTransports，真实/Mock TMDB API）
- 5% Manual E2E Tests（真实 Claude Code 客户端，用户体验验证）

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

### Integration Tests - MCP Protocol Testing

**Scope**: 测试完整的 MCP 协议栈（Client ↔ Server ↔ Tools ↔ TMDB API）

**Framework**: `testing` (标准库) + MCP SDK 的 `InMemoryTransports`

**Location**: `cmd/tmdb-mcp/integration_test.go`

**Key Pattern**: 使用 InMemoryTransports 在同一进程内模拟 client-server 通信

**Advantages**:
- ✅ 完全自动化，无需外部进程
- ✅ 快速执行（纯内存，无网络开销）
- ✅ 易于集成到 CI/CD
- ✅ 精确验证 MCP 协议消息格式
- ✅ 可以 Mock TMDB API 进行隔离测试

**Test Structure**:
```go
// 标准测试模式
func TestSearchTool_Integration(t *testing.T) {
    // 1. 创建传输对
    ct, st := mcp.NewInMemoryTransports()

    // 2. 启动 server
    server := NewMCPServer(testConfig, testLogger)
    ss, _ := server.Connect(context.Background(), st, nil)
    defer ss.Close()

    // 3. 连接 client
    client := mcp.NewClient(&mcp.Implementation{Name: "test"}, nil)
    cs, _ := client.Connect(context.Background(), ct, nil)
    defer cs.Close()

    // 4. 执行工具调用
    result, err := cs.CallTool(ctx, &mcp.CallToolParams{
        Name: "search",
        Arguments: map[string]any{"query": "Inception"},
    })

    // 5. 验证结果
    assert.NoError(t, err)
    assert.NotEmpty(t, result.Content)
    // 验证返回数据结构、字段完整性等
}
```

**Test Scenarios**:

1. **单工具测试**（Epic 1, Story 1.7）:
   - 搜索成功场景（流行电影、电视剧、人物）
   - 边界场景（空查询、不存在内容、大页码）
   - 错误场景（无效参数、TMDB API 错误）
   - 性能验证（响应时间 < 3 秒）
   - 速率限制验证（10 次快速请求）

2. **多工具协作测试**（Epic 2-3）:
   - `search` → `get_details` 链式调用
   - `discover_movies` → `get_recommendations` 组合
   - `get_trending` → `get_details` 场景

3. **并发测试**:
   - 使用 goroutines 并发调用多个工具
   - 验证速率限制器在并发场景下正确工作
   - 运行 `go test -race` 检测数据竞争

4. **协议一致性测试**:
   - 验证 `tools/list` 返回正确的工具列表
   - 验证 `tools/call` 参数解析和验证
   - 验证错误响应格式符合 MCP 规范

**Mock Strategy**:
- **真实 TMDB API**（需要有效 API Key）：用于验证真实集成
- **Mock TMDB API**（httptest）：用于错误场景和速率限制测试
- 通过环境变量 `TMDB_API_KEY` 控制是否使用真实 API

**Example - Mock TMDB API**:
```go
func TestSearchTool_WithMockTMDB(t *testing.T) {
    // 创建 Mock TMDB API server
    mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(mockSearchResponse)
    }))
    defer mockServer.Close()

    // 使用 Mock server 创建 TMDB client
    config := testConfig
    config.TMDB.BaseURL = mockServer.URL

    // ... 运行测试
}
```

**Performance Testing**:
```go
func BenchmarkSearchTool(b *testing.B) {
    // Setup
    ct, st := mcp.NewInMemoryTransports()
    // ...

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cs.CallTool(ctx, &mcp.CallToolParams{
            Name: "search",
            Arguments: map[string]any{"query": "Inception"},
        })
    }
}
```

**Coverage Requirements**:
- `internal/tools` 包：≥ 70%
- `internal/tmdb` 包：≥ 70%
- `internal/mcp` 包：≥ 60%（MCP server 初始化逻辑）

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
