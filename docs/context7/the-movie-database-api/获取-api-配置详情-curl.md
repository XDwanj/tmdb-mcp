# 获取 API 配置详情 (cURL)

来源: https://developer.themoviedb.org/reference/configuration-details

此 cURL 命令从 The Movie Database API 获取 API 配置详情。它指定了 GET 请求方法、配置端点 URL 和预期的 JSON 响应格式。

```shell
curl --request GET \
     --url https://api.themoviedb.org/3/configuration \
     --header 'accept: application/json'
```

--------------------------------
