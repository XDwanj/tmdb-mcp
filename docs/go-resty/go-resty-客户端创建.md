# Go Resty 客户端创建

Source: https://github.com/go-resty/docs/blob/main/content/docs/new-features-and-enhancements.md

演示了使用不同配置创建 Resty 客户端的各种方法，包括传输设置和拨号器。

```go
client := resty.New()

// 覆盖所有传输设置和超时值
client.SetTransport(NewWithTransportSettings(settings))

// 使用自定义拨号器创建客户端
client.SetDialer(NewWithDialer(dialer))

// 使用自定义拨号器和传输设置创建客户端
client = resty.NewWithDialerAndTransportSettings(dialer, settings)
```

--------------------------------
