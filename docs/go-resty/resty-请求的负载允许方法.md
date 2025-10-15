# Resty 请求的负载允许方法

Source: https://github.com/go-resty/docs/blob/main/content/docs/allow-payload-on.md

Resty Request 对象上的这些方法允许您配置特定的 GET 或 DELETE 请求是否可以包含负载。这提供了对单个请求负载包含的精细控制。

```APIDOC
Request.SetAllowMethodGetPayload()
  允许此特定请求的 GET HTTP 动词上的请求负载。

Request.SetAllowMethodDeletePayload()
  允许此特定请求的 DELETE HTTP 动词上的请求负载。
```

--------------------------------
