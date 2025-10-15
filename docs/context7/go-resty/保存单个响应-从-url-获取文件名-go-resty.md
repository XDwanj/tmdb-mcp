# 保存单个响应 (从 URL 获取文件名) - Go Resty

Source: https://github.com/go-resty/docs/blob/main/content/docs/save-response.md

此示例演示如何将单个 HTTP 响应保存到文件系统。文件名是从请求的 URL 自动确定的。在这种情况下，图像将保存为“resty-logo.svg”。

```go
client.R().
    SetSaveResponse(true).
    Get("https://resty.dev/svg/resty-logo.svg")
```

--------------------------------
