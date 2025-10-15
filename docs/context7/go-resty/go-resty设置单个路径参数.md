# Go Resty：设置单个路径参数

Source: https://github.com/go-resty/docs/blob/main/content/docs/request-path-params.md

演示在 Go Resty 中为 GET 请求设置单个动态路径参数。参数值会自动进行 URL 编码。

```go
c := resty.New()
defere c.Close()

c.R().
    SetPathParam("userId", "sample@sample.com").
    Get("/v1/users/{userId}/details")

// Result:
//     /v1/users/sample@sample.com/details
```

--------------------------------
