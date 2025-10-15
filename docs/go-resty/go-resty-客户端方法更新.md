# Go Resty 客户端方法更新

Source: https://github.com/go-resty/docs/blob/main/content/docs/upgrading-to-v3.md

本节详细介绍了 Go Resty 客户端中 Client 方法的更改。它重点介绍了已弃用的方法，并建议使用现代替代方法，这些方法通常涉及新的或重构的功能，如内容类型编码/解码和请求中间件。

```APIDOC
Client.SetHostURL
  - 已弃用：请改用 Client.SetBaseURL。

Client.SetJSONMarshaler, Client.SetJSONUnmarshaler, Client.SetXMLMarshaler, Client.SetXMLUnmarshaler
  - 已弃用：请改用 Client.AddContentTypeEncoder 和 Client.AddContentTypeDecoder。

Client.RawPathParams
  - 已弃用：请改用 Client.PathParams()。

Client.SetRetryResetReaders
  - 功能现在是自动的。

Client.SetRetryAfter
  - 已弃用：请改用 Client.SetRetryStrategy 或 Request.SetRetryStrategy。

Client.RateLimiter 和 Client.SetRateLimiter
  - 重试机制会尊重存在的“Retry-After”标头。

Client.AddRetryAfterErrorCondition
  - 已弃用：请改用 Client.AddRetryConditions。

Client.SetPreRequestHook
  - 已弃用：请改用 Client.SetRequestMiddlewares。请参阅请求中间件文档。

Client.OnRequestLog, Client.OnResponseLog
  - 已弃用：请改用 Client.OnDebugLog。
```

--------------------------------
