# TMDB API 核心接口文档

## 设计原则

- **数据优先**：只包含AI真正需要的6个核心工具
- **消除冗余**：详情接口内嵌常用信息，减少工具调用次数
- **简洁执念**：参数保持最少，覆盖90%使用场景

---

## 认证

所有请求需要在Header中包含API Key：

```bash
Authorization: Bearer YOUR_API_KEY
```

**Base URL**: `https://api.themoviedb.org/3`

---

## 1. 搜索电影

**端点**: `GET /search/movie`

**用途**: 根据关键词搜索电影

**参数**:
- `query` (required): 搜索关键词
- `language` (optional): 语言代码，默认 `en-US`
- `page` (optional): 页码，默认 1

**示例**:
```bash
GET /search/movie?query=Inception&language=en-US
```

**返回**:
```json
{
  "page": 1,
  "results": [
    {
      "id": 27205,
      "title": "Inception",
      "overview": "Cobb, a skilled thief...",
      "release_date": "2010-07-16",
      "vote_average": 8.4,
      "poster_path": "/9gk7adHYeDvHkCSEqAvQNLV5Uge.jpg"
    }
  ],
  "total_results": 42
}
```

---

## 2. 搜索电视剧

**端点**: `GET /search/tv`

**用途**: 根据关键词搜索电视剧

**参数**:
- `query` (required): 搜索关键词
- `language` (optional): 语言代码，默认 `en-US`
- `page` (optional): 页码，默认 1

**示例**:
```bash
GET /search/tv?query=Breaking%20Bad&language=en-US
```

**返回**:
```json
{
  "page": 1,
  "results": [
    {
      "id": 1396,
      "name": "Breaking Bad",
      "overview": "A high school chemistry teacher...",
      "first_air_date": "2008-01-20",
      "vote_average": 8.9,
      "poster_path": "/ggFHVNu6YYI5L9pCfOacjizRGt.jpg"
    }
  ],
  "total_results": 15
}
```

---

## 3. 获取电影详情

**端点**: `GET /movie/{movie_id}`

**用途**: 获取电影完整信息（包含演员表）

**参数**:
- `movie_id` (required): 电影ID
- `language` (optional): 语言代码，默认 `en-US`
- `append_to_response` (optional): 附加信息，推荐 `credits,videos`

**示例**:
```bash
GET /movie/27205?append_to_response=credits,videos&language=en-US
```

**返回**:
```json
{
  "id": 27205,
  "title": "Inception",
  "overview": "Cobb, a skilled thief who commits corporate espionage...",
  "release_date": "2010-07-16",
  "runtime": 148,
  "vote_average": 8.4,
  "vote_count": 35420,
  "poster_path": "/9gk7adHYeDvHkCSEqAvQNLV5Uge.jpg",
  "backdrop_path": "/s3TBrRGB1iav7gFOCNx3H31MoES.jpg",
  "genres": [
    {"id": 28, "name": "Action"},
    {"id": 878, "name": "Science Fiction"}
  ],
  "credits": {
    "cast": [
      {
        "id": 6193,
        "name": "Leonardo DiCaprio",
        "character": "Cobb",
        "profile_path": "/wo2hJpn04vbtmh0B9utCFdsQhxM.jpg"
      }
    ],
    "crew": [
      {
        "id": 525,
        "name": "Christopher Nolan",
        "job": "Director"
      }
    ]
  },
  "videos": {
    "results": [
      {
        "key": "YoHD9XEInc0",
        "name": "Official Trailer",
        "site": "YouTube",
        "type": "Trailer"
      }
    ]
  }
}
```

**图片完整URL构建**:
```
https://image.tmdb.org/t/p/w500{poster_path}
https://image.tmdb.org/t/p/original{backdrop_path}
```

---

## 4. 获取电视剧详情

**端点**: `GET /tv/{tv_id}`

**用途**: 获取电视剧完整信息（包含演员表）

**参数**:
- `tv_id` (required): 电视剧ID
- `language` (optional): 语言代码，默认 `en-US`
- `append_to_response` (optional): 附加信息，推荐 `credits,videos`

**示例**:
```bash
GET /tv/1396?append_to_response=credits,videos&language=en-US
```

**返回**:
```json
{
  "id": 1396,
  "name": "Breaking Bad",
  "overview": "A high school chemistry teacher diagnosed with...",
  "first_air_date": "2008-01-20",
  "last_air_date": "2013-09-29",
  "number_of_seasons": 5,
  "number_of_episodes": 62,
  "vote_average": 8.9,
  "poster_path": "/ggFHVNu6YYI5L9pCfOacjizRGt.jpg",
  "genres": [
    {"id": 18, "name": "Drama"}
  ],
  "credits": {
    "cast": [
      {
        "id": 17419,
        "name": "Bryan Cranston",
        "character": "Walter White"
      }
    ]
  }
}
```

---

## 5. 获取推荐

**端点**: `GET /movie/{movie_id}/recommendations` 或 `GET /tv/{tv_id}/recommendations`

**用途**: 基于指定电影/电视剧获取推荐列表

**参数**:
- `movie_id` / `tv_id` (required): 媒体ID
- `language` (optional): 语言代码，默认 `en-US`
- `page` (optional): 页码，默认 1

**示例**:
```bash
GET /movie/27205/recommendations?language=en-US
```

**返回**:
```json
{
  "page": 1,
  "results": [
    {
      "id": 155,
      "title": "The Dark Knight",
      "overview": "Batman raises the stakes...",
      "release_date": "2008-07-16",
      "vote_average": 8.5,
      "poster_path": "/qJ2tW6WMUDux911r6m7haRef0WH.jpg"
    }
  ],
  "total_results": 20
}
```

---

## 6. 获取趋势

**端点**: `GET /trending/{media_type}/{time_window}`

**用途**: 获取当前趋势内容

**参数**:
- `media_type` (required): `movie` 或 `tv`
- `time_window` (required): `day` 或 `week`
- `language` (optional): 语言代码，默认 `en-US`

**示例**:
```bash
GET /trending/movie/week?language=en-US
```

**返回**:
```json
{
  "page": 1,
  "results": [
    {
      "id": 872585,
      "title": "Oppenheimer",
      "overview": "The story of J. Robert Oppenheimer's role...",
      "release_date": "2023-07-21",
      "vote_average": 8.3,
      "poster_path": "/8Gxv8gSFCU0XGDykEGv7zR1n2ua.jpg"
    }
  ],
  "total_results": 20
}
```

---

## MCP工具映射

| MCP工具 | TMDB端点 | 说明 |
|---------|---------|------|
| `search_movies(query)` | `GET /search/movie` | 搜索电影 |
| `search_tv(query)` | `GET /search/tv` | 搜索电视剧 |
| `get_movie_details(movie_id)` | `GET /movie/{id}?append_to_response=credits,videos` | 电影详情+演员 |
| `get_tv_details(tv_id)` | `GET /tv/{id}?append_to_response=credits,videos` | 电视剧详情+演员 |
| `get_recommendations(movie_id)` | `GET /movie/{id}/recommendations` | 电影推荐 |
| `get_trending(media_type, time_window)` | `GET /trending/{media_type}/{time_window}` | 趋势榜单 |

---

## 错误处理

**常见错误码**:
- `401`: API Key无效或缺失
- `404`: 资源不存在
- `429`: 请求频率超限

**错误响应示例**:
```json
{
  "status_code": 7,
  "status_message": "Invalid API key: You must be granted a valid key.",
  "success": false
}
```

---

## 速率限制

- **免费API**: 40 requests/10 seconds
- 建议实现请求队列和重试机制
