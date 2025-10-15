# 使用 quic-go 启用 HTTP3

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/enable-http3.md

演示如何配置和使用 quic-go 传输与 Go Resty 来启用 HTTP3。它强调了需要社区包，因为 HTTP3 尚未包含在 Go 的标准库中，并指向 quic-go 文档以进行进一步定制。

```go
import (
    "crypto/tls"
    "time"

    "http3" "github.com/quic-go/http3"
    "quic" "github.com/quic-go/quic-go"
    "github.com/go-resty/resty/v2"
)

// 请参阅 quic-go 文档进行自定义配置
// https://quic-go.net/docs/
http3Transport := &http3.Transport{
    TLSClientConfig: &tls.Config{}, // 如果需要，设置 TLS 客户端配置
    QUICConfig: &quic.Config{ // QUIC 连接选项
        MaxIdleTimeout:  45 * time.Second,
        KeepAlivePeriod: 30 * time.Second,
    },
}
defer http3Transport.Close()

c := resty.New().
    SetTransport(http3Transport)
defer c.Close()

// 您现在可以使用 HTTP3 与 Resty

```

--------------------------------
