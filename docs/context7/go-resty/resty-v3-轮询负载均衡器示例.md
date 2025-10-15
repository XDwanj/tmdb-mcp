# Resty v3 轮询负载均衡器示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/load-balancer-and-service-discovery.md

演示如何创建和使用 Resty v3 轮询负载均衡器。这涉及使用多个 URL 初始化负载均衡器，然后将其设置在 Resty 客户端上。

```go
rr, err := resty.NewRoundRobin(
    "https://example1.com",
    "https://example2.com",
    "https://example3.com",
)
if err != nil {
    log.Printf("ERROR %v", err)
    return
}

c := resty.New().
    SetLoadBalancer(rr)
defer c.Close()

// 开始使用客户端...
```

--------------------------------
