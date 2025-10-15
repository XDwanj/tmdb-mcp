# Go Resty 请求跟踪示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/request-tracing.md

演示如何在 Go Resty 中启用请求跟踪，并访问详细的跟踪信息，例如 DNS 查找时间、连接时间、服务器时间和响应时间。

```go
client := resty.New()
defer client.Close()

res, _ = client.R().
    EnableTrace().
    Get("https://httpbin.org/get")

ti := res.Request.TraceInfo()

// 探索跟踪信息
fmt.Println("Request Trace Info:")
fmt.Println("  DNSLookup     :", ti.DNSLookup)
fmt.Println("  ConnTime      :", ti.ConnTime)
fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
fmt.Println("  ServerTime    :", ti.ServerTime)
fmt.Println("  ResponseTime  :", ti.ResponseTime)
fmt.Println("  TotalTime     :", ti.TotalTime)
fmt.Println("  IsConnReused  :", ti.IsConnReused)
fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
fmt.Println("  RequestAttempt:", ti.RequestAttempt)
fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
```

--------------------------------
