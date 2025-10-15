# 获取电视剧集演职员表 (cURL)

来源: https://developer.themoviedb.org/reference/tv-series-credits

获取电视剧集最新季演职员表的示例 cURL 请求。它指定了 GET 方法、带有系列 ID 和语言占位符的 API 端点 URL 以及必需的 Accept 标头。

```shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/tv/series_id/credits?language=en-US' \
     --header 'accept: application/json'
```

--------------------------------
