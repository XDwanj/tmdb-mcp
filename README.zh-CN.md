<div align="center">

# TMDB MCP 服务

[![GitHub Stars](https://img.shields.io/github/stars/XDwanj/tmdb-mcp?style=social)](https://github.com/XDwanj/tmdb-mcp)
[![License](https://img.shields.io/github/license/XDwanj/tmdb-mcp)]
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go)](https://go.dev/)
[![Build](https://github.com/XDwanj/tmdb-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/XDwanj/tmdb-mcp/actions)
[![Docker Pulls](https://img.shields.io/docker/pulls/xdwanj/tmdb-mcp)](#)

基于 TMDB 的 MCP（Model Context Protocol）服务：提供搜索、详情、发现、趋势与推荐等六个工具，支持 stdio 与 SSE 两种模式，内置 Docker 支持。

[English](README.md) | [简体中文](README.zh-CN.md)

</div>

## 项目介绍

TMDB MCP 服务向 MCP 客户端（如 Claude Code）暴露 6 个工具，便于 AI 助手检索与获取 TMDB 的影视数据。

## 快速开始

### 先决条件
- Go 1.21+
- 有效的 TMDB API Key（https://www.themoviedb.org/settings/api）

### 运行（stdio 模式）

```bash
TMDB_API_KEY=your_api_key go run ./cmd/tmdb-mcp --server-mode stdio
```

使用 stdio 作为传输层，MCP 客户端以 stdio 方式连接。

### 运行（SSE 模式）

```bash
# 方式 A：由程序自动生成 Token 并写入 ~/.tmdb-mcp/config.yaml
TMDB_API_KEY=your_api_key \
SERVER_MODE=sse \
go run ./cmd/tmdb-mcp

# 方式 B：自备 token
TMDB_API_KEY=your_api_key \
SERVER_MODE=sse \
SSE_TOKEN=$(openssl rand -base64 32) \
go run ./cmd/tmdb-mcp

# 健康检查
curl -s http://localhost:8910/health
```

### Docker

```bash
docker build -t tmdb-mcp .
docker run --rm -p 8910:8910 \
  -e TMDB_API_KEY=your_api_key \
  -e SERVER_MODE=sse \
  -e SSE_TOKEN=$(openssl rand -base64 32) \
  tmdb-mcp
```

或参考 `examples/docker-compose.yml`。

## 功能特性

暴露的 MCP 工具：
- `search` — 关键字搜索电影/剧集
- `get_details` — 按媒体类型与 ID 获取详情
- `discover_movies` — 按丰富条件发现电影
- `discover_tv` — 按丰富条件发现剧集
- `get_trending` — 按媒体类型与窗口获取趋势
- `get_recommendations` — 基于电影/剧集 ID 获取推荐

典型使用串联：
- 搜索 → 详情
- 发现（电影）→ 推荐
- 趋势 → 详情

## 配置说明

配置来源与优先级：命令行（flags）> 环境变量 > 配置文件。

命令行 flags（部分）：
- `--tmdb-api-key`, `--tmdb-language`, `--tmdb-rate-limit`
- `--server-mode`, `--sse-host`, `--sse-port`, `--sse-token`
- `--logging-level`

环境变量（当未提供 flags 时生效）：
- `TMDB_API_KEY`, `TMDB_LANGUAGE`, `TMDB_RATE_LIMIT`
- `SERVER_MODE`, `SERVER_SSE_HOST`, `SERVER_SSE_PORT`, `SSE_TOKEN`
- `LOGGING_LEVEL`

配置文件（兜底）：`~/.tmdb-mcp/config.yaml`（完整示例见 `examples/config.yaml`）

关键字段：
- `tmdb.api_key`, `tmdb.language`, `tmdb.rate_limit`
- `server.mode`（stdio|sse|both）, `server.sse.host`, `server.sse.port`, `server.sse.token`
- `logging.level`

## 部署方式

- 本地二进制：`go build -o tmdb-mcp ./cmd/tmdb-mcp`
- Docker：`docker build -t tmdb-mcp .`
- Docker Compose：`examples/docker-compose.yml`
- K8s：以 8910 端口暴露容器，`SSE_TOKEN`/`TMDB_API_KEY` 使用 Secret 管理

## 截图/GIF（占位）

Claude Code 使用演示：

![Claude Demo](docs/images/claude-demo.svg)

配置文件示例：

![Config Example](docs/images/config-example.svg)

## 开发与贡献指南

- 代码格式：`go fmt ./...`
- 静态检查：`go vet ./...`
- 运行测试：`go test ./...`
- Markdown 规范：建议本地使用 markdownlint 等工具
- 代码风格/目录规范：见 `docs/architecture/coding-standards.md` 与 `docs/architecture/source-tree.md`

贡献流程：Issue → PR → Review → Merge

## 许可证

待定（TBD）。请在仓库中添加 LICENSE 文件（例如 MIT）。

## 致谢

- TMDB 提供的 API 与数据
- MCP Go SDK https://github.com/modelcontextprotocol/go-sdk
