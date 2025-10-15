# 保存所有响应 - Go Resty

Source: https://github.com/go-resty/docs/blob/main/content/docs/save-response.md

此示例演示如何配置 Resty 客户端将所有 HTTP 响应保存到指定目录。`SetSaveResponse(true)` 选项适用于此客户端进行的所有后续请求。请记住在完成后关闭客户端。

```go
c := resty.New().
    SetOutputDirectory("/path/to/save/all/response").
    SetSaveResponse(true) // 适用于所有请求
defer c.Close()

// 开始使用客户端...
```

--------------------------------
