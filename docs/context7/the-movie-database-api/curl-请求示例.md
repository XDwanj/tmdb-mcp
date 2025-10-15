# cURL 请求示例

来源: https://developer.themoviedb.org/reference/authentication-create-session

如何使用 cURL 向创建会话端点发出 POST 请求的示例。它包括 URL、必需的标头以及指定的内容类型。

```shell
curl --request POST \
     --url https://api.themoviedb.org/3/authentication/session/new \
     --header 'accept: application/json' \
     --header 'content-type: application/json'
```

--------------------------------
