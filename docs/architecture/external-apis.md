# External APIs

## TMDB API v3

**Purpose**: 提供电影、电视剧和人物数据，是本服务的唯一外部数据源

**Documentation**: https://developers.themoviedb.org/3

**Base URL**: `https://api.themoviedb.org/3`

**Authentication**: API Key (Query Parameter `api_key=xxx`)

**Rate Limits**: 40 requests per 10 seconds（免费账户）

**Key Endpoints Used**:

| Endpoint                               | Method | Purpose              | 映射工具            |
| -------------------------------------- | ------ | -------------------- | ------------------- |
| `/search/multi`                        | GET    | 搜索电影/电视剧/人物 | search              |
| `/movie/{id}`                          | GET    | 获取电影详情         | get_details         |
| `/tv/{id}`                             | GET    | 获取电视剧详情       | get_details         |
| `/person/{id}`                         | GET    | 获取人物详情         | get_details         |
| `/discover/movie`                      | GET    | 发现电影（筛选）     | discover_movies     |
| `/discover/tv`                         | GET    | 发现电视剧（筛选）   | discover_tv         |
| `/trending/{media_type}/{time_window}` | GET    | 获取热门内容         | get_trending        |
| `/movie/{id}/recommendations`          | GET    | 获取电影推荐         | get_recommendations |
| `/tv/{id}/recommendations`             | GET    | 获取电视剧推荐       | get_recommendations |

**Integration Notes**:

1. **API Key 管理**:
   - 从环境变量 `TMDB_API_KEY` 或配置文件读取
   - 启动时验证 API Key 有效性（调用 `/configuration` 端点）
   - 不在日志中打印完整 API Key（仅显示前 8 个字符）

2. **速率限制处理**:
   - 使用 Token Bucket 限制器确保不超过 40 req/10s
   - 每次请求前调用 `rateLimiter.Wait(ctx)`
   - 收到 429 错误时，解析 `Retry-After` header 并等待后重试（最多 3 次）

3. **自动追加参数**:
   - 电影/电视剧详情：自动追加 `append_to_response=credits,videos`
   - 人物详情：自动追加 `append_to_response=combined_credits`
   - 减少 API 调用次数，提升性能

4. **语言偏好**:
   - 从配置文件读取 `language` 参数（默认 `en-US`）
   - 所有 API 请求自动添加 `language` 查询参数
   - PRD 明确不实现自动语言检测

5. **错误处理**:
   - **401 Unauthorized**: API Key 无效或过期，立即返回错误
   - **404 Not Found**: 资源不存在，返回友好消息
   - **429 Rate Limit**: 触发限流，自动重试
   - **500/502/503**: TMDB 服务器错误，返回错误并记录日志
   - **Network Timeout**: 10 秒超时，返回错误

6. **归属声明**:
   - 遵守 TMDB API 使用条款
   - 在返回数据中保留 TMDB 归属信息
   - 文档中注明数据来源："This product uses the TMDB API but is not endorsed or certified by TMDB."

---
