# 获取电视剧集推荐 (cURL)

来源: https://developer.themoviedb.org/reference/tv-series-recommendations

演示如何使用 cURL 向 /tv/{series_id}/recommendations 端点发出 GET 请求。它包括设置 API 端点、语言、页码和 accept 标头。

```shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/tv/series_id/recommendations?language=en-US&page=1' \
     --header 'accept: application/json'
```

--------------------------------
