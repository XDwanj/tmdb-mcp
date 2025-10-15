# Go Resty 响应自动解析示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/response-auto-parse.md

演示如何使用 Resty 通过 SetResult 和 SetError 将 JSON 响应自动反序列化为 Go 结构体。它展示了如何设置请求正文、指定结果和错误类型以及访问解析后的结果。

```go
res, err := client.R().
    SetBody(User{
        Username: "testuser",
        Password: "testpass",
    }). // 默认请求内容类型为 JSON
    SetResult(&LoginResponse{}).
    SetError(&LoginError{}).
    Post("https://myapp.com/login")

fmt.Println(err)
fmt.Println(res.Result().(*LoginResponse))
fmt.Println(res.Error().(*LoginError))
```

--------------------------------
