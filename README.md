<div align="center">

# TMDB MCP Server

[![GitHub Stars](https://img.shields.io/github/stars/XDwanj/tmdb-mcp?style=social)](https://github.com/XDwanj/tmdb-mcp)
![License](https://img.shields.io/github/license/XDwanj/tmdb-mcp)
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go)](https://go.dev/)
[![Build](https://github.com/XDwanj/tmdb-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/XDwanj/tmdb-mcp/actions)
[![Release](https://img.shields.io/github/release/XDwanj/tmdb-mcp)](https://github.com/XDwanj/tmdb-mcp/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/xdwanj/tmdb-mcp)](https://github.com/XDwanj/tmdb-mcp/pkgs/container/tmdb-mcp)

An MCP (Model Context Protocol) server for The Movie Database (TMDB): search, details, discovery, trending, and recommendations — ready for stdio and SSE modes, with Docker support.

[English](README.md) | [简体中文](README.zh-CN.md)

</div>

## Overview

TMDB MCP Server exposes six MCP tools so AI coding assistants and MCP clients (e.g., Claude Code) can search and retrieve movie/TV data from TMDB.

## Quick Start

### Prerequisites
- Go 1.21+
- A valid TMDB API Key (https://www.themoviedb.org/settings/api)

### Run (stdio mode)

```bash
TMDB_API_KEY=your_api_key go run ./cmd/tmdb-mcp --server-mode stdio
```

This starts the server over stdio. Connect with an MCP-compatible client via stdio transport.

### Run (SSE mode)

```bash
# Option A: Let the app generate a token and persist to ~/.tmdb-mcp/config.yaml
TMDB_API_KEY=your_api_key \
SERVER_MODE=sse \
go run ./cmd/tmdb-mcp

# Option B: Provide your own token
TMDB_API_KEY=your_api_key \
SERVER_MODE=sse \
SSE_TOKEN=$(openssl rand -base64 32) \
go run ./cmd/tmdb-mcp

# Health check
curl -s http://localhost:8910/health
```

### Docker

#### Option 1: Use pre-built image (Recommended)

```bash
docker pull ghcr.io/xdwanj/tmdb-mcp:latest
docker run --rm -p 8910:8910 \
  -e TMDB_API_KEY=your_api_key \
  -e SERVER_MODE=sse \
  -e SSE_TOKEN=$(openssl rand -base64 32) \
  ghcr.io/xdwanj/tmdb-mcp:latest
```

#### Option 2: Build from source

```bash
docker build -t tmdb-mcp .
docker run --rm -p 8910:8910 \
  -e TMDB_API_KEY=your_api_key \
  -e SERVER_MODE=sse \
  -e SSE_TOKEN=$(openssl rand -base64 32) \
  tmdb-mcp
```

Or use docker-compose: see `examples/docker-compose.yml`.

### Binary Downloads

Pre-compiled binaries are available on the [Releases page](https://github.com/XDwanj/tmdb-mcp/releases).

Supported platforms:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

Example:
```bash
# Linux amd64
wget https://github.com/XDwanj/tmdb-mcp/releases/latest/download/tmdb-mcp-linux-amd64.tar.gz
tar -xzf tmdb-mcp-linux-amd64.tar.gz
sudo mv tmdb-mcp /usr/local/bin/
```

## Features

Exposed MCP tools:
- `search` — Search movies/TV by query
- `get_details` — Get details by media type and ID
- `discover_movies` — Discover movies with rich filters
- `discover_tv` — Discover TV with rich filters
- `get_trending` — Trending items by media type and window
- `get_recommendations` — Recommendations based on a movie/TV ID

Typical flows:
- search → get_details
- discover_movies → get_recommendations
- get_trending → get_details

## Configuration

Configuration sources and priority: CLI flags > Environment variables > Config file.

Flags (subset):
- `--tmdb-api-key`, `--tmdb-language`, `--tmdb-rate-limit`
- `--server-mode`, `--sse-host`, `--sse-port`, `--sse-token`
- `--logging-level`

Environment variables (when flags are not provided):
- `TMDB_API_KEY`, `TMDB_LANGUAGE`, `TMDB_RATE_LIMIT`
- `SERVER_MODE`, `SERVER_SSE_HOST`, `SERVER_SSE_PORT`, `SSE_TOKEN`
- `LOGGING_LEVEL`

Config file (fallback): `~/.tmdb-mcp/config.yaml`
See `examples/config.yaml` for a complete example.

Key fields:
- `tmdb.api_key`, `tmdb.language`, `tmdb.rate_limit`
- `server.mode` (stdio|sse|both), `server.sse.host`, `server.sse.port`, `server.sse.token`
- `logging.level`

## Deployment

### Quick Deployment Options

1. **Pre-built Docker Image (Recommended):**
   ```bash
   docker run -d --name tmdb-mcp -p 8910:8910 \
     -e TMDB_API_KEY=your_api_key \
     -e SERVER_MODE=sse \
     -e SSE_TOKEN=your_token \
     ghcr.io/xdwanj/tmdb-mcp:latest
   ```

2. **Download Binary:**
   ```bash
   # Download the latest release for your platform
   wget https://github.com/XDwanj/tmdb-mcp/releases/latest/download/tmdb-mcp-linux-amd64.tar.gz
   tar -xzf tmdb-mcp-linux-amd64.tar.gz
   ./tmdb-mcp --server-mode sse
   ```

3. **Build from Source:**
   ```bash
   go build -o tmdb-mcp ./cmd/tmdb-mcp
   ./tmdb-mcp --server-mode sse
   ```

### Production Deployment

- **Docker Compose:** See `examples/docker-compose.yml`
- **Kubernetes:** Expose the container at port 8910 and configure `SSE_TOKEN`/`TMDB_API_KEY` via secrets
- **Docker Registry:** Use `ghcr.io/xdwanj/tmdb-mcp:{tag}` for specific versions

### Version Management

- Use semantic versioning (e.g., `v1.0.0`, `v1.0.1`)
- `latest` tag always points to the most recent stable release
- For production, pin to a specific version tag (e.g., `ghcr.io/xdwanj/tmdb-mcp:v1.0.0`)

## Use with Claude Code (MCP)

When running in SSE mode (for example at `http://localhost:8910/mcp/sse`) with an `SSE_TOKEN`, you can connect Claude Code via either of the following methods.

- Option 1: Configure `~/.claude.json`

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

- Option 2: Use the Claude CLI

```sh
claude mcp add --transport sse tmdb-mcp-sse \
  http://localhost:8910/mcp/sse \
  --header "Authorization: Bearer xxxxxxxxxxxxxxxxxxxxx" \
  --scope user
```

- Verify installation

```sh
claude mcp list
```

Replace the placeholder token with your actual `SSE_TOKEN`, and adjust the URL if the server runs on a different host/port.

## Development

- Format: `go fmt ./...`
- Vet: `go vet ./...`
- Test: `go test ./...`
- Lint (Markdown): run your preferred linter locally (e.g., markdownlint)
- Code style and structure: see `docs/architecture/coding-standards.md` and `docs/architecture/source-tree.md`

## Contributing

1) Open an issue to discuss changes
2) Fork and create a feature branch
3) Add tests and docs
4) Submit a PR for review

## License

TBD. Please add a LICENSE file (e.g., MIT) to this repository.

## Acknowledgements

- TMDB for the API and data
- MCP Go SDK https://github.com/modelcontextprotocol/go-sdk
