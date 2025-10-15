# 电视剧集关键词 cURL 请求

来源: https://developer.themoviedb.org/reference/tv-series-keywords

此代码片段演示了如何向 TMDB API 发出 GET 请求以获取电视剧集的关键词。它包括必要的 URL 和标头。通常需要通过 API 密钥进行身份验证，但在本通用示例中省略了。

```bash
curl --request GET \
     --url https://api.themoviedb.org/3/tv/series_id/keywords \
     --header 'accept: application/json'
```

--------------------------------
