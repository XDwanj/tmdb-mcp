# 执行 TRACE 请求

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/options-head-trace-request.md

展示如何使用 Resty 发送 TRACE 请求。TRACE 方法通常用于诊断目的，将请求回显给客户端。示例包括设置身份验证令牌。

```go
res, err = client.R().
    SetAuthToken("bc594900518b4f7eac75bd37f019e08fbc594900518b4f7eac75bd37f019e08f").
    Trace("https://myapp.com/test")

fmt.Println(err, res)
```

--------------------------------
