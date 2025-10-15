# 保存单个响应 (自定义文件名) - Go Resty

Source: https://github.com/go-resty/docs/blob/main/content/docs/save-response.md

此示例演示如何将单个 HTTP 响应保存为自定义文件名。`SetOutputFileName` 方法允许您为保存的文件指定相对或绝对路径。在此示例中，图像将保存为“resty-logo-blue.svg”。

```go
client.R().
    SetSaveResponse(true).
    SetOutputFileName("resty-logo-blue.svg"). // 可以是相对或绝对路径
    Get("https://resty.dev/svg/resty-logo.svg")
```

--------------------------------
