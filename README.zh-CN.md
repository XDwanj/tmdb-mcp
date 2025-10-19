<div align="center">

# TMDB MCP 服务器

[![GitHub Stars](https://img.shields.io/github/stars/XDwanj/tmdb-mcp?style=social)](https://github.com/XDwanj/tmdb-mcp)
![License](https://img.shields.io/github/license/XDwanj/tmdb-mcp)
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go)](https://go.dev/)
[![Build](https://github.com/XDwanj/tmdb-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/XDwanj/tmdb-mcp/actions)
[![Release](https://img.shields.io/github/release/XDwanj/tmdb-mcp)](https://github.com/XDwanj/tmdb-mcp/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/xdwanj/tmdb-mcp)](https://github.com/XDwanj/tmdb-mcp/pkgs/container/tmdb-mcp)

The Movie Database (TMDB) 的 MCP (Model Context Protocol) 服务器：搜索、详情、发现、热门、推荐 — 支持 stdio 和 SSE 模式，并提供 Docker 支持。

[English](README.md) | [简体中文](README.zh-CN.md)

</div>

## 概览

TMDB MCP 服务器暴露了六个 MCP 工具，以便 AI 编码助手和 MCP 客户端（例如 Claude Code）可以从 TMDB 搜索和检索电影/电视数据。

## 快速入门

### 先决条件
- Go 1.21+
- 有效的 TMDB API 密钥 (https://www.themoviedb.org/settings/api)

### 运行 (stdio 模式)

```bash
TMDB_API_KEY=your_api_key go run ./cmd/tmdb-mcp --server-mode stdio
```

这将通过 stdio 启动服务器。通过 stdio 传输连接到 MCP 兼容的客户端。

### 运行 (SSE 模式)

```bash
# 选项 A：让应用程序生成一个 token 并持久化到 ~/.tmdb-mcp/config.yaml
TMDB_API_KEY=your_api_key \
SERVER_MODE=sse \
go run ./cmd/tmdb-mcp

# 选项 B：提供您自己的 token
TMDB_API_KEY=your_api_key \
SERVER_MODE=sse \
SSE_TOKEN=$(openssl rand -base64 32) \
go run ./cmd/tmdb-mcp

# 健康检查
curl -s http://localhost:8910/health
```

### Docker

#### 选项 1：使用预构建镜像（推荐）

```bash
docker pull ghcr.io/xdwanj/tmdb-mcp:latest
docker run --rm -p 8910:8910 \
  -e TMDB_API_KEY=your_api_key \
  -e SERVER_MODE=sse \
  -e SSE_TOKEN=$(openssl rand -base64 32) \
  ghcr.io/xdwanj/tmdb-mcp:latest
```

#### 选项 2：从源码构建

```bash
docker build -t tmdb-mcp .
docker run --rm -p 8910:8910 \
  -e TMDB_API_KEY=your_api_key \
  -e SERVER_MODE=sse \
  -e SSE_TOKEN=$(openssl rand -base64 32) \
  tmdb-mcp
```

或者使用 docker-compose：请参阅 `examples/docker-compose.yml`。

### 二进制文件下载

预编译的二进制文件可在 [Releases 页面](https://github.com/XDwanj/tmdb-mcp/releases) 上获取。

支持的平台：
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

示例：
```bash
# Linux amd64
wget https://github.com/XDwanj/tmdb-mcp/releases/latest/download/tmdb-mcp-linux-amd64.tar.gz
tar -xzf tmdb-mcp-linux-amd64.tar.gz
sudo mv tmdb-mcp /usr/local/bin/
```

## 功能

暴露的 MCP 工具：
- `search` — 按查询搜索电影/电视
- `get_details` — 按媒体类型和 ID 获取详情
- `discover_movies` — 使用丰富的过滤器发现电影
- `discover_tv` — 使用丰富的过滤器发现电视
- `get_trending` — 按媒体类型和时间窗口获取热门内容
- `get_recommendations` — 基于电影/电视 ID 获取推荐

典型流程：
- search → get_details
- discover_movies → get_recommendations
- get_trending → get_details

## 配置

配置来源和优先级：CLI 标志 > 环境变量 > 配置文件。

标志（部分）：
- `--tmdb-api-key`, `--tmdb-language`, `--tmdb-rate-limit`
- `--server-mode`, `--sse-host`, `--sse-port`, `--sse-token`
- `--logging-level`

环境变量（未提供标志时）：
- `TMDB_API_KEY`, `TMDB_LANGUAGE`, `TMDB_RATE_LIMIT`
- `SERVER_MODE`, `SERVER_SSE_HOST`, `SERVER_SSE_PORT`, `SSE_TOKEN`
- `LOGGING_LEVEL`

配置文件（回退）：`~/.tmdb-mcp/config.yaml`
请参阅 `examples/config.yaml` 获取完整示例。

关键字段：
- `tmdb.api_key`, `tmdb.language`, `tmdb.rate_limit`
- `server.mode` (stdio|sse|both), `server.sse.host`, `server.sse.port`, `server.sse.token`
- `logging.level`

## 部署

### 快速部署选项

1. **预构建 Docker 镜像（推荐）：**
   ```bash
   docker run -d --name tmdb-mcp -p 8910:8910 \
     -e TMDB_API_KEY=your_api_key \
     -e SERVER_MODE=sse \
     -e SSE_TOKEN=your_token \
     ghcr.io/xdwanj/tmdb-mcp:latest
   ```

2. **下载二进制文件：**
   ```bash
   # 下载适用于您平台的最新版本
   wget https://github.com/XDwanj/tmdb-mcp/releases/latest/download/tmdb-mcp-linux-amd64.tar.gz
   tar -xzf tmdb-mcp-linux-amd64.tar.gz
   ./tmdb-mcp --server-mode sse
   ```

3. **从源码构建：**
   ```bash
   go build -o tmdb-mcp ./cmd/tmdb-mcp
   ./tmdb-mcp --server-mode sse
   ```

### 生产环境部署

- **Docker Compose：** 请参阅 `examples/docker-compose.yml`
- **Kubernetes：** 在端口 8910 暴露容器，并通过 secrets 配置 `SSE_TOKEN`/`TMDB_API_KEY`
- **Docker 注册表：** 使用 `ghcr.io/xdwanj/tmdb-mcp:{tag}` 指定特定版本

### 版本管理

- 使用语义化版本控制（例如 `v1.0.0`, `v1.0.1`）
- `latest` 标签始终指向最新的稳定发布版本
- 生产环境请固定到特定版本标签（例如 `ghcr.io/xdwanj/tmdb-mcp:v1.0.0`）

## 与 Claude Code (MCP) 配合使用

当在 SSE 模式下运行（例如在 `http://localhost:8910/mcp/sse`）并提供 `SSE_TOKEN` 时，您可以通过以下任一方式连接 Claude Code。

- 选项 1：配置 `~/.claude.json`

```json5
{
  "mcpServers": {
    "tmdb-mcp-sse": {
      "type": "sse",
      "url": "http://localhost:8910/mcp/sse",
      "headers": {
        "Authorization": "Bearer xxxxxxxxxxxxxxxxxxxxx"
      }
    }
  }
}
```

- 选项 2：使用 Claude CLI

```sh
claude mcp add --transport sse tmdb-mcp-sse \
  http://localhost:8910/mcp/sse \
  --header "Authorization: Bearer xxxxxxxxxxxxxxxxxxxxx" \
  --scope user
```

- 验证安装

```sh
claude mcp list
```

将占位符 token 替换为您实际的 `SSE_TOKEN`，如果服务器运行在不同的主机/端口，请调整 URL。

## 开发

- 格式化：`go fmt ./...`
- 检查：`go vet ./...`
- 测试：`go test ./...`
- Lint (Markdown)：在本地运行您喜欢的 linter（例如 markdownlint）
- 代码风格和结构：请参阅 `docs/architecture/coding-standards.md` 和 `docs/architecture/source-tree.md`

## 贡献

1) 打开一个 issue 讨论更改
2) Fork 并创建一个特性分支
3) 添加测试和文档
4) 提交 PR 进行审查

## 许可证

待定。请在此仓库中添加一个 LICENSE 文件（例如 MIT）。

## 致谢

- TMDB 的 API 和数据
- MCP Go SDK https://github.com/modelcontextprotocol/go-sdk