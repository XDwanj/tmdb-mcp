# cURL 请求示例：验证登录

来源: https://developer.themoviedb.org/reference/authentication-create-session-from-login

使用登录凭据验证请求令牌的示例 cURL 命令。需要 content-type 和 accept 标头。建议使用 HTTPS。

```shell
curl --request POST \
     --url https://api.themoviedb.org/3/authentication/token/validate_with_login \
     --header 'accept: application/json' \
     --header 'content-type: application/json'
```

--------------------------------
