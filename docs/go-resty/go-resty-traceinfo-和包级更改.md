# Go Resty TraceInfo 和包级更改

Source: https://github.com/go-resty/docs/blob/main/content/docs/upgrading-to-v3.md

涵盖了 Go Resty v3 中 TraceInfo 和包级函数类型的更改，包括类型更新和重命名。

```APIDOC
TraceInfo Changes:

- `TraceInfo.RemoteAddr`: 类型从 `net.Addr` 更改为 `string`。

Package Level Types:

- Retry:
  - `RetryHookFunc`: 从 `OnRetryFunc` 重命名。
  - `RetryStrategyFunc`: 从 `RetryStrategyFunc` 重命名（名称不变，但上下文可能不同）。
```

--------------------------------
