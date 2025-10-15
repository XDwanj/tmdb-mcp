# GET /configuration

来源: https://developer.themoviedb.org/reference/configuration-details

获取 API 配置详情，其中包括图像的基 URL 和支持的图像大小。

```APIDOC
# GET /configuration

# 描述
检索 API 配置详情，包括图像基 URL 和支持的图像大小。

# 方法
GET

# 端点
https://api.themoviedb.org/3/configuration

# 参数
## 查询参数
无

## 请求正文
无

# 响应
## 成功响应 (200)
- **images** (object) - 包含图像配置详情。
  - **base_url** (string) - 图像的基 URL。
  - **secure_base_url** (string) - 图像的安全基 URL。
  - **backdrop_sizes** (string 数组) - 支持的背景图像大小。
  - **logo_sizes** (string 数组) - 支持的 Logo 图像大小。
  - **poster_sizes** (string 数组) - 支持的海报图像大小。
  - **profile_sizes** (string 数组) - 支持的个人资料图像大小。
  - **still_sizes** (string 数组) - 支持的剧照图像大小。
- **change_keys** (string 数组) - 触发 API 更改的键。

## 响应示例
```json
{
  "images": {
    "base_url": "http://image.tmdb.org/t/p/",
    "secure_base_url": "https://image.tmdb.org/t/p/",
    "backdrop_sizes": [
      "w300",
      "w780",
      "w1280",
      "original"
    ],
    "logo_sizes": [
      "w45",
      "w92",
      "w154",
      "w185",
      "w300",
      "w500"
    ],
    "poster_sizes": [
      "w92",
      "w154",
      "w185",
      "w342",
      "w500",
      "w780",
      "original"
    ],
    "profile_sizes": [
      "w45",
      "w185",
      "h632"
    ],
    "still_sizes": [
      "w92",
      "w185",
      "w300",
      "original"
    ]
  },
  "change_keys": [
    "adult",
    "air_date",
    "also_known_as",
    "backdrop_path",
    "biography",
    "birthday",
    "budget",
    "character",
    "created_by",
    "crew",
    ""}
```
```

--------------------------------
