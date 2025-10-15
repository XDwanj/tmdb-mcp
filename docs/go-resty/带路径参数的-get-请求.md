# 带路径参数的 GET 请求

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/get-request.md

使用路径参数构建 GET 请求，以指定 URL 的动态部分，例如用户 ID 或帐户标识符。示例还包含了一个身份验证令牌。

```go
res, err := client.R()
    .SetPathParams(map[string]string{
		"userId":       "sample@sample.com",
		"subAccountId": "100002",
	})
    .SetAuthToken("bc594900518b4f7eac75bd37f019e08fbc594900518b4f7eac75bd37f019e08f")
    .Get("/v1/users/{userId}/{subAccountId}/details")

fmt.Println(err, res)
```

--------------------------------
