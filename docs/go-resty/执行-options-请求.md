# 执行 OPTIONS 请求

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/options-head-trace-request.md

演示如何使用 Resty 执行 OPTIONS 请求。这通常用于确定允许的 HTTP 方法或用于 CORS 预检请求。它包括设置身份验证令牌。

```go
res, err := client.R().
    SetAuthToken("bc594900518b4f7eac75bd37f019e08fbc594900518b4f7eac75bd37f019e08f").
    Options("https://myapp.com/servers/nyc-dc-01")

fmt.Println(err, res)
fmt.Println(res.Header())
```

--------------------------------
