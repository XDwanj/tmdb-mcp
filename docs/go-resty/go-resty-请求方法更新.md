# Go Resty 请求方法更新

Source: https://github.com/go-resty/docs/blob/main/content/docs/upgrading-to-v3.md

本节介绍 Go Resty 中 Request 方法的修改。它指出了已弃用或替换的方法，引导用户采用推荐的路径参数、身份验证和服务发现实践。

```APIDOC
Request.RawPathParams
  - 已弃用：请改用 Request.PathParams。

Request.SRV 和 Request.SetSRV
  - 已弃用：请改用 NewSRVWeightedRoundRobin 和 Client.SetLoadBalancer。新实现会考虑 SRV 记录权重。

Request.SetDigestAuth
  - 已弃用：请改用 Client.SetDigestAuth。
```

--------------------------------
