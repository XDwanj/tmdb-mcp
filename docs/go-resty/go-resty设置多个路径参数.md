# Go Resty：设置多个路径参数

Source: https://github.com/go-resty/docs/blob/main/content/docs/request-path-params.md

演示在 Go Resty 中使用 map 设置多个动态路径参数以进行 GET 请求。值会被 URL 编码，包括路径段中的特殊字符。

```go
c := resty.New()
defere c.Close()

c.R().
    SetPathParams(map[string]string{
        "userId":       "sample@sample.com",
        "subAccountId": "100002",
        "path":         "groups/developers",
    }).
    Get("/v1/users/{userId}/{subAccountId}/{path}/details)

// Result:
//   /v1/users/sample@sample.com/100002/groups%2Fdevelopers/details
```

--------------------------------
