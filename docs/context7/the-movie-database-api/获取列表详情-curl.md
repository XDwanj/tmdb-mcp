# 获取列表详情 (cURL)

来源: https://developer.themoviedb.org/reference/list-details

如何使用 cURL 检索特定列表详情的示例。它指定了 GET 请求、带有列表 ID 占位符的 API 端点以及必需的标头。

```Shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/list/list_id?language=en-US&page=1' \
     --header 'accept: application/json'
```

--------------------------------
