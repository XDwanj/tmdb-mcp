# Resty v3 加权轮询负载均衡器示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/load-balancer-and-service-discovery.md

展示如何设置 Resty v3 加权轮询负载均衡器。这允许为不同的主机指定权重以控制请求分发。它还涵盖了设置恢复持续时间。

```go
wrr, err := resty.NewWeightedRoundRobin(
    3*time.Second, // 恢复持续时间
    []*resty.Host{
        {
            BaseURL: "https://example1.com",
            Weight:  50, // 确定到此主机的请求百分比
        },
        {BaseURL: "https://example2.com", Weight: 30},
        {BaseURL: "https://example3.com", Weight: 20},
    }...,
)
if err != nil {
    log.Printf("ERROR %v", err)
    return
}

// 默认情况下，恢复持续时间为 120 秒，可以按如下方式更改
// wrr.SetRecoveryDuration(3 * time.Minute)

c := resty.New().
    SetLoadBalancer(wrr)
defer c.Close()

// 开始使用客户端...
```

--------------------------------
