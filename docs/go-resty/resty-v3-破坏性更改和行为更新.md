# Resty v3 破坏性更改和行为更新

Source: https://github.com/go-resty/docs/blob/main/content/docs/upgrading-to-v3.md

详细介绍了 Resty v3 中的破坏性更改，包括错误消息格式、需要 defer client.Close() 以及内容长度、DELETE 负载和重试机制处理方式的更改。它还涵盖了多部分、重定向、中间件、标头、摘要身份验证和超时。

```APIDOC
Resty v3 升级说明：

错误格式：
  - 所有 Resty 错误现在都以 `resty: ...` 前缀开头。
  - 子功能错误包括功能名称，例如 `resty: digest: ...`。

客户端生命周期：
  - 在创建客户端后添加 `defer client.Close()`。

行为更改：
  - 内容长度：
    - 内容长度选项不再适用于 `io.Reader` 流。
  - DELETE 负载：
    - 默认情况下，HTTP 动词 DELETE 不支持负载。
    - 使用 `Client.AllowMethodDeletePayload` 或 `Request.AllowMethodDeletePayload` 为 DELETE 请求启用负载支持。
  - 重试机制：
    - 请求值从创建时继承自客户端，在重试尝试期间不会刷新。通过 `Response.Request` 在请求实例上更新值。
    - 如果存在 `Retry-After` 标头，则会尊重该标头。
    - 如果支持 `io.ReadSeeker` 接口，则在重试请求时重置读取器。
    - 仅在幂等 HTTP 动词（GET、HEAD、PUT、DELETE、OPTIONS、TRACE）上重试，符合 RFC 9110 和 RFC 5789。
    - 使用 `Client.SetAllowNonIdempotentRetry` 或 `Request.SetAllowNonIdempotentRetry` 允许对非幂等方法进行重试。
    - 应用默认重试条件，可以通过 `Client.SetRetryDefaultConditions` 或 `Request.SetRetryDefaultConditions` 禁用。
  - 多部分：
    - 默认情况下，当在 MultipartField 输入中检测到文件或 `io.Reader` 时，Resty 会在请求正文中流式传输内容。
  - 重定向策略：
    - `NoRedirectPolicy` 返回错误 `http.ErrUseLastResponse`。
  - 响应中间件：
    - 所有响应中间件都会执行，无论错误如何，并将错误向下级联。
    - 检查错误以确定是继续还是跳过逻辑执行。
  - 标头：
    - 默认情况下，Resty 不为请求设置 `Accept` 标头。
  - 摘要身份验证：
    - 仅在客户端级别支持。创建专用客户端以利用它。
  - 超时：
    - 不使用 `http.Client.Timeout`；而是使用带超时的上下文。
  - Curl 命令生成：
    - `curl` 命令生成流是独立的，不需要启用调试或跟踪。
```

--------------------------------
