# Source Tree

项目采用 **标准 Go 项目布局**，遵循 Go 社区最佳实践：

```
tmdb-mcp/
├── cmd/
│   └── tmdb-mcp/                 # 主应用程序
│       └── main.go               # 程序入口
│
├── internal/                     # 私有应用代码（不可被外部导入）
│   ├── config/                   # 配置管理
│   │   ├── config.go             # 配置加载和验证
│   │   └── token.go              # SSE Token 生成
│   │
│   ├── tmdb/                     # TMDB API 客户端
│   │   ├── client.go             # HTTP 客户端封装
│   │   ├── search.go             # 搜索相关 API
│   │   ├── details.go            # 详情相关 API
│   │   ├── discover.go           # 发现相关 API
│   │   ├── trending.go           # 热门相关 API
│   │   ├── recommendations.go    # 推荐相关 API
│   │   ├── error.go              # 错误处理
│   │   └── models.go             # TMDB 响应模型
│   │
│   ├── ratelimit/                # 速率限制
│   │   └── limiter.go            # Token Bucket 限制器
│   │
│   ├── mcp/                      # MCP 服务器
│   │   └── server.go             # MCP Server 初始化和工具注册
│   │
│   ├── tools/                    # MCP 工具实现
│   │   ├── search.go             # search 工具
│   │   ├── get_details.go        # get_details 工具
│   │   ├── discover_movies.go    # discover_movies 工具
│   │   ├── discover_tv.go        # discover_tv 工具
│   │   ├── get_trending.go       # get_trending 工具
│   │   ├── get_recommendations.go # get_recommendations 工具
│   │   └── params.go             # 参数模型定义
│   │
│   ├── server/                   # HTTP Server 相关组件
│   │   └── middleware/           # HTTP 中间件
│   │       └── auth.go           # Bearer Token 认证中间件
│   │
│   └── logger/                   # 日志系统
│       └── logger.go             # Zap Logger 初始化
│
├── pkg/                          # 公共库代码（可被外部导入）
│   └── version/                  # 版本信息
│       └── version.go            # 版本常量
│
├── docs/                         # 项目文档
│   ├── architecture.md           # 架构文档（本文档）
│   ├── prd.md                    # 产品需求文档
│   ├── brief.md                  # 项目简报
│   └── tmdb-api.md               # TMDB API 文档
│
├── examples/                     # 示例配置和脚本
│   ├── config.yaml               # 示例配置文件
│   ├── docker-compose.yml        # Docker Compose 示例
│   └── claude-code-config.json   # Claude Code 配置示例
│
├── Dockerfile                    # Docker 镜像构建文件
├── .dockerignore                 # Docker 忽略文件
├── go.mod                        # Go Modules 依赖定义
├── go.sum                        # Go Modules 依赖校验和
├── .gitignore                    # Git 忽略文件
└── README.md                     # 项目说明文档
```

## 目录结构说明

**`cmd/`** - 主要应用程序

- 包含可执行程序的入口代码
- 每个子目录名对应一个可执行文件名
- `main.go` 应简洁，主要负责初始化和启动

**`internal/`** - 私有应用代码

- Go 的特殊目录，代码不可被外部项目导入
- 存放所有业务逻辑和实现细节
- 按功能分包：`config`, `tmdb`, `mcp`, `tools`, `server`, `logger`

**`pkg/`** - 公共库代码

- 可被外部项目导入的库代码
- 本项目仅包含 `version` 包，未来可扩展
- 如果未来开源为库，考虑将通用部分移到 `pkg/`

**`docs/`** - 文档

- 架构文档、PRD、API 文档等
- 不包含代码，仅包含 Markdown 文档

**`examples/`** - 示例和模板

- 示例配置文件、Docker Compose、脚本
- 帮助用户快速上手

## 文件命名约定

- **Go 文件**: 小写蛇形命名（`snake_case.go`）
- **包名**: 小写单词（`package tmdb`）
- **测试文件**: `_test.go` 后缀（`client_test.go`）
- **配置文件**: YAML 格式（`config.yaml`）

## 包组织原则

1. **按功能分包**: 而非按层次（避免 `models/`, `services/` 这种通用命名）
2. **包内聚**: 每个包有明确单一职责
3. **包独立**: 减少包间循环依赖
4. **internal 优先**: 默认放在 `internal/`，除非需要被外部导入

---
