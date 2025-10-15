# Go Resty 客户端方法重命名和更改

Source: https://github.com/go-resty/docs/blob/main/content/docs/upgrading-to-v3.md

本节详细介绍了 Go Resty 客户端中 getter 方法的重命名，以实现线程安全并与命名约定保持一致。它还涵盖了方法签名和功能的更改。

```APIDOC
Client Methods:

- Getter Methods:
  - `Client.BaseURL()`: 获取基础 URL。
  - `Client.FormData()`: 获取表单数据。
  - `Client.Header()`: 获取标头。
  - `Client.AuthToken()`: 获取身份验证令牌。
  - `Client.Client()`: 获取底层客户端实例。

- Method Signature Changes:
  - `Client.SetDebugBodyLimit(limit int)`: 从 `int64` 更改为 `int`。
  - `Client.ResponseBodyLimit()`: 从 `int` 更改为 `int64`。
  - `Client.SetAllowMethodGetPayload(allow bool)`: 从 `SetAllowGetMethodPayload` 重命名。
  - `Client.Clone(ctx context.Context)`: 添加了 `context.Context` 参数。
  - `Client.EnableGenerateCurlCmd()`: 从 `EnableGenerateCurlOnDebug` 重命名。
  - `Client.DisableGenerateCurlCmd()`: 从 `DisableGenerateCurlOnDebug` 重命名。
  - `Client.SetRootCertificates(certs ...*x509.Certificate)`: 从 `SetRootCertificate` 重命名。
  - `Client.SetClientRootCertificates(certs ...*x509.Certificate)`: 从 `SetClientRootCertificate` 重命名。
  - `Client.IsDebug()`: 从 `Debug` 重命名。
  - `Client.IsDisableWarn()`: 从 `DisableWarn` 重命名。
  - `Client.AddRetryConditions(conditions ...RetryConditionFunc)`: 从 `AddRetryCondition` 重命名。
  - `Client.AddRetryHooks(hooks ...RetryHookFunc)`: 从 `AddRetryHook` 重命名。
  - `Client.SetRetryStrategy(strategy RetryStrategyFunc)`: 从 `SetRetryAfter` 重命名。
  - `Client.HTTPTransport()`: 新方法，返回 `http.Transport`。
  - `Client.AddRequestMiddleware(middleware RequestMiddlewareFunc)`: 从 `OnBeforeRequest` 重命名。
  - `Client.AddResponseMiddleware(middleware ResponseMiddlewareFunc)`: 从 `OnAfterResponse` 重命名。
```

--------------------------------
