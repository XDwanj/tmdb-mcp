# 获取已评分电影 (cURL)

来源: https://developer.themoviedb.org/reference/account-rated-movies

如何使用 cURL 获取用户已评分电影的示例。它演示了 GET 请求、带有语言、页面和排序顺序等参数的 URL 构建以及必需的 Accept 标头。

```shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/account/null/rated/movies?language=en-US&page=1&sort_by=created_at.asc' \
     --header 'accept: application/json'
```

--------------------------------
