# Go Resty：设置单个原始路径参数

Source: https://github.com/go-resty/docs/blob/main/content/docs/request-path-params.md

演示在 Go Resty 中为 GET 请求设置单个动态原始路径参数。参数值按原样使用，不进行 URL 编码。

```go
c := resty.New()
defere c.Close()

c.R().
    SetRawPathParam("path", "groups/developers").
    Get("/v1/users/{userId}/details")

// Result:
//     /v1/users/groups/developers/details
```

--------------------------------
