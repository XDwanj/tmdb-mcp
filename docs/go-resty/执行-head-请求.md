# 执行 HEAD 请求

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/options-head-trace-request.md

演示如何使用 Resty 执行 HEAD 请求。此方法可用于在不获取响应正文的情况下检索与资源相关的标头。示例中包含一个身份验证令牌。

```go
res, err = client.R().
    SetAuthToken("bc594900518b4f7eac75bd37f019e08fbc594900518b4f7eac75bd37f019e08f").
    Head("https://myapp.com/videos/hi-res-video")

fmt.Println(err, res)
fmt.Println(res.Header())
```

--------------------------------
