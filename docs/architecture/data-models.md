# Data Models

**说明**: 本项目为无状态 API 转发服务，无数据库和持久化需求。以下定义的是**内部数据结构**（Go structs），用于配置管理、TMDB API 交互和 MCP 工具实现。

## Configuration Model

**Purpose**: 管理服务配置，支持从文件、环境变量、命令行 flags 三种来源加载

**Go Struct 定义**:
```go
type Config struct {
    TMDB      TMDBConfig   `mapstructure:"tmdb"`
    Server    ServerConfig `mapstructure:"server"`
    Logging   LogConfig    `mapstructure:"logging"`
}

type TMDBConfig struct {
    APIKey    string `mapstructure:"api_key"`    // TMDB API Key（必需）
    Language  string `mapstructure:"language"`   // 语言偏好（默认 "en-US"）
    RateLimit int    `mapstructure:"rate_limit"` // 速率限制（默认 40 req/10s）
}

type ServerConfig struct {
    Mode string    `mapstructure:"mode"` // "stdio", "sse", "both"（默认 "both"）
    SSE  SSEConfig `mapstructure:"sse"`
}

type SSEConfig struct {
    Enabled bool   `mapstructure:"enabled"` // 是否启用 SSE 模式
    Host    string `mapstructure:"host"`    // 监听地址（默认 "0.0.0.0"）
    Port    int    `mapstructure:"port"`    // 监听端口（默认 8910）
    Token   string `mapstructure:"token"`   // Bearer Token（自动生成或手动设置）
}

type LogConfig struct {
    Level string `mapstructure:"level"` // "debug", "info", "warn", "error"
}
```

**Key Attributes**:
- `TMDB.APIKey`: string - TMDB API 认证密钥，从环境变量 `TMDB_API_KEY` 或配置文件读取
- `Server.Mode`: string - 运行模式，支持 stdio、sse、both 三种值
- `SSE.Token`: string - SSE 认证 Token，优先级：ENV > 配置文件 > 自动生成
- `Logging.Level`: string - 日志级别，控制 Zap logger 的输出详细程度

**Relationships**:
- Config 是顶层结构，包含 TMDB、Server、Logging 三个子配置
- 通过 Viper 加载，支持 `mapstructure` tag 映射 YAML 字段

## TMDB API Response Models

**Purpose**: 封装 TMDB API 的响应数据，提供类型安全的数据访问

**关键模型**:

### SearchResult（搜索结果）

```go
type SearchResult struct {
    ID           int     `json:"id"`
    MediaType    string  `json:"media_type"`    // "movie", "tv", "person"
    Title        string  `json:"title"`         // 电影标题
    Name         string  `json:"name"`          // 电视剧/人物名称
    ReleaseDate  string  `json:"release_date"`  // 上映日期
    FirstAirDate string  `json:"first_air_date"` // 首播日期
    VoteAverage  float64 `json:"vote_average"`  // 评分
    Overview     string  `json:"overview"`      // 简介
}
```

### MovieDetails（电影详情）

```go
type MovieDetails struct {
    ID          int      `json:"id"`
    Title       string   `json:"title"`
    ReleaseDate string   `json:"release_date"`
    Runtime     int      `json:"runtime"`
    VoteAverage float64  `json:"vote_average"`
    Overview    string   `json:"overview"`
    Genres      []Genre  `json:"genres"`
    Credits     Credits  `json:"credits"`      // 自动追加
    Videos      Videos   `json:"videos"`       // 自动追加
}

type Genre struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type Credits struct {
    Cast []CastMember `json:"cast"`
    Crew []CrewMember `json:"crew"`
}

type CastMember struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    Character string `json:"character"`
}
```

### TVDetails（电视剧详情）

```go
type TVDetails struct {
    ID            int      `json:"id"`
    Name          string   `json:"name"`
    FirstAirDate  string   `json:"first_air_date"`
    NumberOfSeasons int    `json:"number_of_seasons"`
    VoteAverage   float64  `json:"vote_average"`
    Overview      string   `json:"overview"`
    Genres        []Genre  `json:"genres"`
    Credits       Credits  `json:"credits"` // 自动追加
    Videos        Videos   `json:"videos"`  // 自动追加
}
```

### PersonDetails（人物详情）

```go
type PersonDetails struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Birthday    string `json:"birthday"`
    Biography   string `json:"biography"`
    KnownForDepartment string `json:"known_for_department"`
    CombinedCredits    CombinedCredits `json:"combined_credits"` // 自动追加
}
```

**Relationships**:
- SearchResult 是统一搜索结果，通过 `MediaType` 区分电影/电视剧/人物
- Details 模型通过 `Credits` 和 `Videos` 嵌套关联演职员表和视频数据
- 所有模型使用 JSON tags 与 TMDB API 响应字段映射

## MCP Tool Parameter Models

**Purpose**: 定义 MCP 工具的输入参数结构，用于参数验证和类型转换

**jsonschema 标签**: MCP SDK v1.0+ 支持自动从 jsonschema 标签生成 InputSchema,无需手动定义。这遵循官方 SDK 的最佳实践,显著简化了工具定义。

### SearchParams

```go
type SearchParams struct {
    Query string `json:"query" jsonschema:"Search query for movies, TV shows, and people"` // 搜索关键词（必需）
    Page  int    `json:"page" jsonschema:"Page number (default: 1)"`                       // 页码（可选，默认 1）
}
```

### GetDetailsParams

```go
type GetDetailsParams struct {
    MediaType string `json:"media_type"` // "movie", "tv", "person"（必需）
    ID        int    `json:"id"`         // TMDB ID（必需）
}
```

### DiscoverMoviesParams

```go
type DiscoverMoviesParams struct {
    WithGenres         string  `json:"with_genres"`          // 类型 ID（逗号分隔）
    PrimaryReleaseYear int     `json:"primary_release_year"` // 上映年份
    VoteAverageGte     float64 `json:"vote_average.gte"`     // 最低评分
    VoteAverageLte     float64 `json:"vote_average.lte"`     // 最高评分
    WithOriginalLanguage string `json:"with_original_language"` // 语言代码
    SortBy             string  `json:"sort_by"`              // 排序方式
    Page               int     `json:"page"`                 // 页码
}
```

### DiscoverTVParams

```go
type DiscoverTVParams struct {
    WithGenres         string  `json:"with_genres"`
    FirstAirDateYear   int     `json:"first_air_date_year"`
    VoteAverageGte     float64 `json:"vote_average.gte"`
    VoteAverageLte     float64 `json:"vote_average.lte"`
    WithOriginalLanguage string `json:"with_original_language"`
    WithStatus         string  `json:"with_status"` // "Returning Series", "Ended" 等
    SortBy             string  `json:"sort_by"`
    Page               int     `json:"page"`
}
```

### GetTrendingParams

```go
type GetTrendingParams struct {
    MediaType  string `json:"media_type"`  // "movie", "tv", "person"（必需）
    TimeWindow string `json:"time_window"` // "day", "week"（必需）
    Page       int    `json:"page"`        // 页码（可选）
}
```

### GetRecommendationsParams

```go
type GetRecommendationsParams struct {
    MediaType string `json:"media_type"` // "movie", "tv"（必需）
    ID        int    `json:"id"`         // TMDB ID（必需）
    Page      int    `json:"page"`       // 页码（可选）
}
```

**Relationships**:
- 每个 Params 结构对应一个 MCP 工具
- 使用 JSON tags 与 MCP 协议的 JSON-RPC 参数映射
- 工具层负责将这些 Params 转换为 TMDB API 查询参数

## Error Model

**Purpose**: 统一错误处理和错误响应格式

```go
type TMDBError struct {
    StatusCode int    `json:"status_code"` // HTTP 状态码
    StatusMessage string `json:"status_message"` // TMDB API 错误消息
}

func (e *TMDBError) Error() string {
    return fmt.Sprintf("TMDB API Error %d: %s", e.StatusCode, e.StatusMessage)
}
```

**Key Attributes**:
- `StatusCode`: int - HTTP 状态码（401, 404, 429 等）
- `StatusMessage`: string - TMDB API 返回的错误消息

**Error Handling Flow**:
1. TMDB API 返回错误响应
2. 解析为 `TMDBError` 结构
3. 根据 `StatusCode` 决定重试或直接返回
4. 转换为 MCP 错误响应格式返回给客户端

---

**设计注意事项**:

1. **无 ORM 映射**: 所有结构直接映射 JSON，无需数据库注解
2. **类型安全**: 使用 Go 类型系统确保编译时类型检查
3. **JSON Tags**: 使用小写蛇形命名与 TMDB API 保持一致
4. **可扩展性**: 结构设计支持未来添加新字段（TMDB API 更新）
5. **内存效率**: 避免不必要的指针和嵌套，减少 GC 压力

---
