# Go Resty 响应方法

Source: https://github.com/go-resty/docs/blob/main/content/docs/new-features-and-enhancements.md

本节记录了用于访问和检查 Go 中 Resty 响应对象的 GetTraceInfo 方法。它包括获取响应正文、字节、检查正文是否已读取、访问错误以及获取重定向历史记录的方法。

```APIDOC
Response Methods:

Body() string
  - 将响应正文作为字符串返回。

Bytes() []byte
  - 将响应正文作为字节切片返回。

IsRead() bool
  - 检查响应正文是否已被读取。

Err() error
  - 返回请求期间遇到的任何错误。

RedirectHistory() []*Request
  - 返回响应的重定向历史记录。
```

--------------------------------
