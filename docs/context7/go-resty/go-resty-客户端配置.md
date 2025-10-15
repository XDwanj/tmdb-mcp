# Go Resty 客户端配置

Source: https://github.com/go-resty/docs/blob/main/content/docs/new-features-and-enhancements.md

本节涵盖了 Go Resty 客户端的核心配置方法。它包括设置中间件、内容编码器/解码器以及管理响应正文的可读性。

```APIDOC
Client.Close()
  - 关闭客户端并释放任何持有的资源。

Client.SetRequestMiddlewares(middlewares ...middleware.Middleware)
  - 设置要在发送请求之前应用于请求的自定义中间件函数。

Client.SetResponseMiddlewares(middlewares ...middleware.Middleware)
  - 设置要在接收响应之后应用于响应的自定义中间件函数。

Client.AddContentTypeEncoder(contentType string, encoder func(interface{}) (string, error))
  - 为特定内容类型（例如 'application/json'）注册自定义编码器。

Client.AddContentTypeDecoder(contentType string, decoder func([]byte) (interface{}, error))
  - 为特定内容类型（例如 'application/json'）注册自定义解码器。

Client.SetResponseBodyUnlimitedReads(unlimited bool)
  - 配置响应正文是否可以被读取多次。默认为 false。
```

--------------------------------
