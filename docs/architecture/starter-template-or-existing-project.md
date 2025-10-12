# Starter Template or Existing Project

**项目类型**: Greenfield 项目(全新创建)

**决策**: 不使用 Starter Template

**理由**:
1. **Go 项目特性**: Golang 项目通常采用标准项目布局,不需要复杂脚手架
2. **精简原则**: PRD 明确要求精简工具链,仅使用 Go 原生工具(`go build`, `go test`, `go fmt`, `go vet`)
3. **技术栈明确**: 所有核心依赖已在 PRD 中定义(MCP SDK, Resty, Viper, Zap 等)
4. **项目结构清晰**: 采用标准 Go 布局:`cmd/tmdb-mcp/`, `internal/`, `pkg/`

**项目初始化方式**:
```bash
go mod init github.com/[username]/tmdb-mcp
mkdir -p cmd/tmdb-mcp internal pkg
```
