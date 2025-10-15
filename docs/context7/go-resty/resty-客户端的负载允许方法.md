# Resty 客户端的负载允许方法

Source: https://github.com/go-resty/docs/blob/main/content/docs/allow-payload-on.md

Resty 客户端上的这些方法允许您配置 GET 和 DELETE 请求是否可以包含负载。这对于遵循标准 HTTP 动词语义的系统很有用。

```APIDOC
Client.SetAllowMethodGetPayload()
  允许 GET HTTP 动词上的请求负载。

Client.SetAllowMethodDeletePayload()
  允许 DELETE HTTP 动词上的请求负载。
```

--------------------------------
