# Resty v3 SRV 加权轮询负载均衡器示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/load-balancer-and-service-discovery.md

演示 Resty v3 的 SRV 加权轮询负载均衡器的用法。此方法通过 SRV 记录发现服务并应用加权轮询分发。它包括设置刷新和恢复持续时间的选项。

```go
swrr, err := resty.NewSRVWeightedRoundRobin(
    "_sample-server",
    "tcp", // 默认协议为 tcp
    "example.com",
    "https", // 默认方案为 https
)
if err != nil {
    log.Printf("ERROR %v", err)
    return
}

// 默认情况下，SRV 记录刷新持续时间为 180 秒，可以按如下方式更改
// swrr.SetRefreshDuration(1 * time.Hour)

// 默认情况下，恢复持续时间为 120 秒，可以按如下方式更改
// swrr.SetRecoveryDuration(3 * time.Minute)

c := resty.New().
    SetLoadBalancer(swrr)
defer c.Close()

// 开始使用客户端...
```

--------------------------------
