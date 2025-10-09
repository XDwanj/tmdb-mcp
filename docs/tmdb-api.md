# TMDB API 核心接口文档

## 设计原则

- **场景驱动**：基于真实使用场景设计，而非完整映射 TMDB API
- **智能合并**：相似功能通过参数统一（搜索4合1、详情3合1），降低 LLM 认知负担
- **最小化工具数**：6个精简工具覆盖核心场景（搜索、详情、发现、趋势、推荐）
- **自动优化**：详情接口自动追加 credits/videos，减少额外调用
- **参数精简**：保留高频参数，覆盖90%使用场景

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

## 2. 搜索人物

**端点**: `GET /search/person`

**用途**: 根据关键词搜索演员、导演等人物

**参数**:
- `query` (required): 搜索关键词
- `language` (optional): 语言代码，默认 `en-US`
- `page` (optional): 页码，默认 1

**示例**:
```bash
GET /search/person?query=Leonardo%20DiCaprio&language=en-US
```

**返回**:
```json
{
  "page": 1,
  "results": [
    {
      "id": 6193,
      "name": "Leonardo DiCaprio",
      "known_for_department": "Acting",
      "profile_path": "/wo2hJpn04vbtmh0B9utCFdsQhxM.jpg",
      "known_for": [
        {
          "id": 27205,
          "title": "Inception",
          "media_type": "movie"
        }
      ]
    }
  ],
  "total_results": 15
}
```

---

## 3. 多媒体搜索

**端点**: `GET /search/multi`

**用途**: 同时搜索电影、电视剧和人物

**参数**:
- `query` (required): 搜索关键词
- `language` (optional): 语言代码，默认 `en-US`
- `page` (optional): 页码，默认 1

**示例**:
```bash
GET /search/multi?query=Inception&language=en-US
```

**返回**:
```json
{
  "page": 1,
  "results": [
    {
      "id": 27205,
      "media_type": "movie",
      "title": "Inception",
      "release_date": "2010-07-16",
      "vote_average": 8.4
    },
    {
      "id": 6193,
      "media_type": "person",
      "name": "Leonardo DiCaprio",
      "known_for_department": "Acting"
    }
  ],
  "total_results": 58
}
```

**注意**: 返回结果包含 `media_type` 字段区分类型（movie/tv/person）

---

## 4. 搜索电视剧

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

## 5. 获取人物详情

**端点**: `GET /person/{person_id}`

**用途**: 获取人物完整信息（包含作品列表）

**参数**:
- `person_id` (required): 人物ID
- `language` (optional): 语言代码，默认 `en-US`
- `append_to_response` (optional): 附加信息，推荐 `combined_credits`

**示例**:
```bash
GET /person/6193?append_to_response=combined_credits&language=en-US
```

**返回**:
```json
{
  "id": 6193,
  "name": "Leonardo DiCaprio",
  "birthday": "1974-11-11",
  "place_of_birth": "Los Angeles, California, USA",
  "biography": "Leonardo Wilhelm DiCaprio is an American actor and film producer...",
  "profile_path": "/wo2hJpn04vbtmh0B9utCFdsQhxM.jpg",
  "known_for_department": "Acting",
  "combined_credits": {
    "cast": [
      {
        "id": 27205,
        "title": "Inception",
        "character": "Cobb",
        "media_type": "movie",
        "release_date": "2010-07-16",
        "vote_average": 8.4
      },
      {
        "id": 597,
        "title": "Titanic",
        "character": "Jack Dawson",
        "media_type": "movie",
        "release_date": "1997-11-18",
        "vote_average": 7.9
      }
    ],
    "crew": [
      {
        "id": 466272,
        "title": "The Revenant",
        "job": "Producer",
        "media_type": "movie",
        "release_date": "2015-12-25"
      }
    ]
  }
}
```

**注意**: `combined_credits` 包含电影和电视剧作品，通过 `media_type` 区分

---

## 6. 获取推荐

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

## 8. 发现电影

**端点**: `GET /discover/movie`

**用途**: 按条件筛选和发现电影

**参数**:
- `language` (optional): 语言代码，默认 `en-US`
- `page` (optional): 页码，默认 1
- `sort_by` (optional): 排序方式，默认 `popularity.desc`
  - 可选值: `popularity.desc`, `vote_average.desc`, `release_date.desc`, `revenue.desc`
- `with_genres` (optional): 类型ID（逗号分隔），如 `28,12` (动作,冒险)
- `primary_release_year` (optional): 首映年份，如 `2023`
- `vote_average.gte` (optional): 最低评分，如 `7.0`
- `vote_count.gte` (optional): 最低评分数，如 `100`
- `with_runtime.gte` (optional): 最短时长（分钟）
- `with_runtime.lte` (optional): 最长时长（分钟）

**常用类型ID**:
- 28: Action (动作)
- 12: Adventure (冒险)
- 16: Animation (动画)
- 35: Comedy (喜剧)
- 80: Crime (犯罪)
- 18: Drama (剧情)
- 14: Fantasy (奇幻)
- 27: Horror (恐怖)
- 10749: Romance (爱情)
- 878: Science Fiction (科幻)
- 53: Thriller (惊悚)

**示例**:
```bash
# 查找2020年后的高分科幻片
GET /discover/movie?with_genres=878&primary_release_year=2020&vote_average.gte=7.0&sort_by=vote_average.desc&language=en-US
```

**返回**:
```json
{
  "page": 1,
  "results": [
    {
      "id": 872585,
      "title": "Oppenheimer",
      "overview": "The story of J. Robert Oppenheimer...",
      "release_date": "2023-07-21",
      "vote_average": 8.3,
      "vote_count": 5420,
      "poster_path": "/8Gxv8gSFCU0XGDykEGv7zR1n2ua.jpg",
      "genre_ids": [18, 36]
    }
  ],
  "total_results": 142,
  "total_pages": 8
}
```

---

## 9. 发现电视剧

**端点**: `GET /discover/tv`

**用途**: 按条件筛选和发现电视剧

**参数**:
- `language` (optional): 语言代码，默认 `en-US`
- `page` (optional): 页码，默认 1
- `sort_by` (optional): 排序方式，默认 `popularity.desc`
  - 可选值: `popularity.desc`, `vote_average.desc`, `first_air_date.desc`
- `with_genres` (optional): 类型ID（逗号分隔）
- `first_air_date_year` (optional): 首播年份
- `vote_average.gte` (optional): 最低评分
- `vote_count.gte` (optional): 最低评分数
- `with_status` (optional): 状态筛选
  - 0: Returning Series (续播中)
  - 1: Planned (计划中)
  - 2: In Production (制作中)
  - 3: Ended (已完结)
  - 4: Cancelled (已取消)
  - 5: Pilot (试播)

**示例**:
```bash
# 查找高分剧情剧
GET /discover/tv?with_genres=18&vote_average.gte=8.0&sort_by=vote_average.desc&language=en-US
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
      "vote_count": 12450,
      "poster_path": "/ggFHVNu6YYI5L9pCfOacjizRGt.jpg",
      "genre_ids": [18, 80]
    }
  ],
  "total_results": 58,
  "total_pages": 3
}
```

---

## 10. MCP工具映射

### 设计理念
- **合并相似功能**：减少 LLM 认知负担
- **参数区分类型**：通过 `type` 参数统一接口
- **保持独立性**：功能差异大的接口保持独立

### 工具列表

| MCP工具 | TMDB端点 | 说明 |
|---------|---------|------|
| `search(type, query)` | `GET /search/{type}` | **统一搜索**<br>type: "movie" \| "tv" \| "person" \| "multi"<br>合并4个搜索接口 |
| `get_details(type, id)` | `GET /{type}/{id}` | **获取详情**<br>type: "movie" \| "tv" \| "person"<br>自动追加 credits/combined_credits<br>合并3个详情接口 |
| `discover_movies(params)` | `GET /discover/movie` | **发现电影**<br>支持类型、年份、评分等筛选 |
| `discover_tv(params)` | `GET /discover/tv` | **发现电视剧**<br>支持类型、年份、评分等筛选 |
| `get_trending(media_type, time_window)` | `GET /trending/{media_type}/{time_window}` | **趋势榜单**<br>media_type: "movie" \| "tv"<br>time_window: "day" \| "week" |
| `get_recommendations(type, id)` | `GET /{type}/{id}/recommendations` | **获取推荐**<br>type: "movie" \| "tv" |

### 工具参数详解

#### 1. search(type, query, page?)
```json
{
  "type": "movie | tv | person | multi",
  "query": "搜索关键词",
  "page": 1  // optional
}
```

#### 2. get_details(type, id, language?)
```json
{
  "type": "movie | tv | person",
  "id": 27205,
  "language": "en-US"  // optional, 默认 en-US
}
```
**自动追加**：
- movie/tv: `append_to_response=credits,videos`
- person: `append_to_response=combined_credits`

#### 3. discover_movies(params)
```json
{
  "sort_by": "vote_average.desc",  // optional
  "with_genres": "878,12",  // optional, 类型ID逗号分隔
  "primary_release_year": 2023,  // optional
  "vote_average.gte": 7.0,  // optional
  "vote_count.gte": 100,  // optional
  "page": 1  // optional
}
```

#### 4. discover_tv(params)
```json
{
  "sort_by": "vote_average.desc",  // optional
  "with_genres": "18,80",  // optional
  "first_air_date_year": 2020,  // optional
  "vote_average.gte": 8.0,  // optional
  "with_status": 3,  // optional, 0-5
  "page": 1  // optional
}
```

#### 5. get_trending(media_type, time_window, language?)
```json
{
  "media_type": "movie | tv",
  "time_window": "day | week",
  "language": "en-US"  // optional
}
```

#### 6. get_recommendations(type, id, page?)
```json
{
  "type": "movie | tv",
  "id": 27205,
  "page": 1  // optional
}
```

---

## 使用场景示例

### 场景1：智能文件重命名
```
用户文件: "The.Irishman.2019.HDR.1080p.WEBRip.X265-MEGABOX.mp4"

步骤：
1. LLM 提取关键字 "The Irishman 2019"
2. 调用 search("movie", "The Irishman")
3. 获取结果列表，找到 id=398978
4. 调用 get_details("movie", 398978)
5. 获取官方标题 "The Irishman" 和年份 "2019"
6. 重命名为 "The Irishman (2019).mp4"
```

### 场景2：片荒找电影
```
用户: "推荐一些2020年后的高分科幻片"

步骤：
1. 调用 discover_movies({
     "with_genres": "878",
     "primary_release_year": 2020,
     "vote_average.gte": 7.5,
     "sort_by": "vote_average.desc"
   })
2. 返回结果列表给用户
```

### 场景3：导演作品查询
```
用户: "诺兰还拍过什么电影？"

步骤：
1. 调用 search("person", "Christopher Nolan")
2. 获取 person_id=525
3. 调用 get_details("person", 525)
4. 从 combined_credits.crew 中筛选 job="Director" 的电影
5. 返回导演作品列表
```

### 场景4：推荐相似电影
```
用户: "我喜欢《盗梦空间》，推荐类似的"

步骤：
1. 调用 search("movie", "Inception")
2. 获取 movie_id=27205
3. 调用 get_recommendations("movie", 27205)
4. 返回推荐列表
```

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
