# HTTP 客户端用法

Source: https://github.com/go-resty/docs/blob/main/content/_index.md

演示 Resty HTTP 客户端的基本用法，以发出 GET 请求。它展示了如何创建客户端、启用跟踪以及处理响应。依赖项包括 Resty 库。

```go
client := resty.New()
defer client.Close()

res, err := client.R().
    EnableTrace().
    Get("https://httpbin.org/get")
fmt.Println(err, res)
fmt.Println(res.Request.TraceInfo())
```

--------------------------------
