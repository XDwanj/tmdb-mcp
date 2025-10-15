# 获取电影推荐 (cURL)

来源: https://developer.themoviedb.org/reference/movie-recommendations

此代码片段演示了如何使用 cURL 向 movie/movie_id/recommendations 端点发出 GET 请求。它包括基本 URL、路径参数以及语言和页面等查询参数。

```shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/movie/movie_id/recommendations?language=en-US&page=1' \
     --header 'accept: application/json'
```

--------------------------------
