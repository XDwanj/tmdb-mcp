# 简单的 GET 请求

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/get-request.md

对指定的 URL 执行基本的 GET 请求。捕获响应和任何潜在的错误。

```go
res, err := client.R()
    .Get("https://httpbin.org/get")

fmt.Println(err, res)
```

--------------------------------
