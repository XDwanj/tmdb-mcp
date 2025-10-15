# 发现电影 API 请求 (cURL)

来源: https://developer.themoviedb.org/reference/discover-movie

演示如何使用 cURL 向 /discover/movie 端点发出 GET 请求。它包括语言、排序以及成人和视频内容过滤等常用参数。

```Shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=en-US&page=1&sort_by=popularity.desc' \
     --header 'accept: application/json'
```

--------------------------------
