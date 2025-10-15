# 在 go-resty 中配置 SOCKS5 代理

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/socks5-proxy.md

演示如何为 go-resty 客户端设置 SOCKS5 代理。这涉及创建一个新的客户端实例，并使用 SOCKS5 URI 调用 `SetProxy` 方法。示例展示了一个通过配置的代理发出的 GET 请求。

```go
package main

import (
	"fmt"

	"resty.dev/v3"
)

func main() {
	c := resty.New().
		SetProxy("socks5://127.0.0.1:1080")
	defer c.Close()

	res, err := c.R().
		Get("https://httpbin.org/get")

	fmt.Println(err, res)
}
```

--------------------------------
