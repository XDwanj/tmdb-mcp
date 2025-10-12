# Infrastructure and Deployment

## Infrastructure as Code

**Tool**: N/A - 本项目无需 IaC 工具

**说明**:
- 自托管部署，无云基础设施
- 用户直接运行二进制文件或 Docker 容器
- 无需 Terraform、Pulumi 等 IaC 工具

## Deployment Strategy

**Strategy**: **独立二进制 + Docker 容器**

**CI/CD Platform**: GitHub Actions（用于自动构建和发布）

**Pipeline Configuration**: `.github/workflows/release.yml`

**Deployment Modes**:

1. **本地二进制部署**（推荐用于开发和个人使用）
   - 下载对应平台的二进制文件
   - 配置 `~/.tmdb-mcp/config.yaml`
   - 直接运行：`./tmdb-mcp`

2. **Docker 容器部署**（推荐用于服务器和远程访问）
   - 拉取 Docker 镜像：`docker pull username/tmdb-mcp:latest`
   - 通过环境变量或挂载配置文件运行
   - 支持 Docker Compose 一键启动

## Environments

本项目不区分传统意义上的"环境"（dev/staging/prod），而是通过配置文件和命令行参数控制行为：

- **开发模式**: `logging.level=debug`，控制台输出，详细日志
- **生产模式**: `logging.level=info`，JSON 格式日志，性能优化
- **测试模式**: 使用 Mock TMDB API，单元测试和集成测试

## Dockerfile 设计

**多阶段构建**（减小镜像大小）：

```dockerfile
# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /build

# 复制依赖文件并下载
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码并编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o tmdb-mcp ./cmd/tmdb-mcp

# Stage 2: Runtime
FROM alpine:latest

# 安装 CA 证书（HTTPS 请求需要）
RUN apk --no-cache add ca-certificates

WORKDIR /app

# 从 builder 复制二进制文件
COPY --from=builder /build/tmdb-mcp .

# 创建配置目录
RUN mkdir -p /root/.tmdb-mcp

# 暴露端口（SSE 模式）
EXPOSE 8910

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8910/health || exit 1

# 运行程序
CMD ["./tmdb-mcp"]
```

**镜像大小**: < 20MB（Alpine base + Go 静态二进制）

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `TMDB_API_KEY` | TMDB API 密钥 | - | ✅ Yes |
| `SSE_TOKEN` | SSE 认证 Token | 自动生成 | ❌ No |
| `SERVER_MODE` | 运行模式（stdio/sse/both） | `both` | ❌ No |
| `SERVER_SSE_HOST` | SSE 监听地址 | `0.0.0.0` | ❌ No |
| `SERVER_SSE_PORT` | SSE 监听端口 | `8910` | ❌ No |
| `TMDB_LANGUAGE` | 语言偏好 | `en-US` | ❌ No |
| `TMDB_RATE_LIMIT` | 速率限制（req/10s） | `40` | ❌ No |
| `LOGGING_LEVEL` | 日志级别 | `info` | ❌ No |

## Docker Compose 示例

```yaml
version: '3.8'

services:
  tmdb-mcp:
    image: username/tmdb-mcp:latest
    container_name: tmdb-mcp
    environment:
      - TMDB_API_KEY=your_api_key_here
      - SSE_TOKEN=your_secure_token_here
      - SERVER_MODE=sse
      - LOGGING_LEVEL=info
    ports:
      - "8910:8910"
    volumes:
      - ./config:/root/.tmdb-mcp  # 可选：挂载配置文件
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8910/health"]
      interval: 30s
      timeout: 3s
      retries: 3
```

## Multi-Platform Binary Compilation

使用 GitHub Actions 自动编译多平台二进制：

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o tmdb-mcp-linux-amd64 ./cmd/tmdb-mcp

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o tmdb-mcp-linux-arm64 ./cmd/tmdb-mcp

# macOS AMD64 (Intel)
GOOS=darwin GOARCH=amd64 go build -o tmdb-mcp-darwin-amd64 ./cmd/tmdb-mcp

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o tmdb-mcp-darwin-arm64 ./cmd/tmdb-mcp

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o tmdb-mcp-windows-amd64.exe ./cmd/tmdb-mcp
```

## Rollback Strategy

**Primary Method**: Git tag + Docker 镜像版本

**Trigger Conditions**:
- 严重 Bug 导致服务不可用
- API 不兼容变更导致客户端无法连接
- 性能严重下降（响应时间 > 5 秒）

**Rollback Steps**:
1. 停止当前版本服务
2. 拉取前一版本 Docker 镜像或下载前一版本二进制
3. 启动前一版本服务
4. 验证功能正常

**Recovery Time Objective**: < 5 分钟（手动回滚）

---
