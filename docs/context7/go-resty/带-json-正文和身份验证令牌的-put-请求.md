# 带 JSON 正文和身份验证令牌的 PUT 请求

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/post-put-patch-request.md

演示如何发送带有 JSON 正文和身份验证令牌的 PUT 请求。此示例包括设置请求正文、身份验证以及处理潜在错误。

```go
res, err := client.R().
    SetBody(Article{
        Title: "Resty",
        Content: "This is my article content, oh ya!",
        Author: "Jeevanandam M",
        Tags: []string{"article", "sample", "resty"},
    }). // 默认请求内容类型为 JSON
    SetAuthToken("bc594900518b4f7eac75bd37f019e08fbc594900518b4f7eac75bd37f019e08f").
    SetError(&Error{}). // 或 SetError(Error{}).
    Put("https://myapp.com/articles/123456")

fmt.Println(err, res)
fmt.Println(res.Error().(*Error))
```

--------------------------------
