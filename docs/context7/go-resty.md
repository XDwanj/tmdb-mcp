================
代码片段
================
## 入门 Server-Sent Events

Source: https://github.com/go-resty/docs/blob/main/content/docs/server-sent-events.md

演示如何使用 Resty 初始化和连接到 Server-Sent Events 流的基本示例。它设置了 URL 并定义了一个用于接收消息的处理程序。

```go
es := resty.NewEventSource().
    SetURL("https://sse.dev/test").
    OnMessage(func(e any) {
        fmt.Println(e.(*resty.Event))
    }, nil)

err := es.Get()
fmt.Println(err)
```

--------------------------------

## 创建 Resty 客户端

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/get-request.md

初始化一个新的 Resty 客户端实例。建议在完成后关闭客户端以释放资源。

```go
client := resty.New()
defer client.Close()
```

--------------------------------

## 在 go-resty 中配置 SOCKS5 代理

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

## 带路径参数的 GET 请求

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/get-request.md

使用路径参数构建 GET 请求，以指定 URL 的动态部分，例如用户 ID 或帐户标识符。示例还包含了一个身份验证令牌。

```go
res, err := client.R()
    .SetPathParams(map[string]string{
		"userId":       "sample@sample.com",
		"subAccountId": "100002",
	})
    .SetAuthToken("bc594900518b4f7eac75bd37f019e08fbc594900518b4f7eac75bd37f019e08f")
    .Get("/v1/users/{userId}/{subAccountId}/details")

fmt.Println(err, res)
```

--------------------------------

## 简单的 GET 请求

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/get-request.md

对指定的 URL 执行基本的 GET 请求。捕获响应和任何潜在的错误。

```go
res, err := client.R()
    .Get("https://httpbin.org/get")

fmt.Println(err, res)
```

--------------------------------

## 断路器配置示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/circuit-breaker.md

演示如何创建和配置一个具有自定义超时、失败阈值和成功阈值的新断路器，然后将其设置在 Resty 客户端上。

```go
import (
	"time"
	"github.com/go-resty/resty/v2"
)

// 使用必需的值创建断路器，根据需要覆盖
cb := resty.NewCircuitBreaker().
	SetTimeout(15 * time.Second).
	SetFailureThreshold(10).
	SetSuccessThreshold(5)

// 创建 Resty 客户端
c := resty.New().
    SetCircuitBreaker(cb)
defer c.Close()

// 开始使用客户端...
```

--------------------------------

## 带查询参数的 GET 请求

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/get-request.md

执行带有多个查询参数、自定义标头和身份验证令牌的 GET 请求。这对于过滤或分页搜索结果非常有用。

```go
res, err := client.R()
    .SetQueryParams(map[string]string{
        "page_no": "1",
        "limit":   "20",
        "sort":    "name",
        "order":   "asc",
        "random":  strconv.FormatInt(time.Now().Unix(), 10),
    })
    .SetHeader("Accept", "application/json")
    .SetAuthToken("bc594900518b4f7eac75bd37f019e08fbc594900518b4f7eac75bd37f019e08f")
    .Get("/search_result")

fmt.Println(err, res)
```

--------------------------------

## HTTP 客户端用法

Source: https://github.com/go-resty/docs/blob/main/content/_index.md

演示 Resty HTTP 客户端的基本用法，以发出 GET 请求。它展示了如何创建客户端、启用跟踪以及处理响应。依赖项包括 Resty 库。

```go
client := resty.New()
defer client.Close()

res, err := client.R().
    EnableTrace().
    Get("https://httpbin.org/get")
fmt.Println(err, res)
fmt.Println(res.Request.TraceInfo())
```

--------------------------------

## Resty v3 轮询负载均衡器示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/load-balancer-and-service-discovery.md

演示如何创建和使用 Resty v3 轮询负载均衡器。这涉及使用多个 URL 初始化负载均衡器，然后将其设置在 Resty 客户端上。

```go
rr, err := resty.NewRoundRobin(
    "https://example1.com",
    "https://example2.com",
    "https://example3.com",
)
if err != nil {
    log.Printf("ERROR %v", err)
    return
}

c := resty.New().
    SetLoadBalancer(rr)
defer c.Close()

// 开始使用客户端...
```

--------------------------------

## Go Resty 请求方法更新

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

## Go Resty 请求跟踪示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/request-tracing.md

演示如何在 Go Resty 中启用请求跟踪，并访问详细的跟踪信息，例如 DNS 查找时间、连接时间、服务器时间和响应时间。

```go
client := resty.New()
defer client.Close()

res, _ = client.R().
    EnableTrace().
    Get("https://httpbin.org/get")

ti := res.Request.TraceInfo()

// 探索跟踪信息
fmt.Println("Request Trace Info:")
fmt.Println("  DNSLookup     :", ti.DNSLookup)
fmt.Println("  ConnTime      :", ti.ConnTime)
fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
fmt.Println("  ServerTime    :", ti.ServerTime)
fmt.Println("  ResponseTime  :", ti.ResponseTime)
fmt.Println("  TotalTime     :", ti.TotalTime)
fmt.Println("  IsConnReused  :", ti.IsConnReused)
fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
fmt.Println("  RequestAttempt:", ti.RequestAttempt)
fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
```

--------------------------------

## Go Resty 响应和多部分更改

Source: https://github.com/go-resty/docs/blob/main/content/docs/upgrading-to-v3.md

详细介绍了 Go Resty v3 中 Response 和 MultipartField 对象的更改，包括方法重命名。

```APIDOC
Response Methods:

- `Response.Duration()`: 从 `Time` 重命名。

MultipartField Methods:

- `MultipartField.Name()`: 从 `Param` 重命名。
```

--------------------------------

## Go Resty TraceInfo 和包级更改

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

## 执行 HEAD 请求

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/options-head-trace-request.md

演示如何使用 Resty 执行 HEAD 请求。此方法可用于在不获取响应正文的情况下检索与资源相关的标头。示例中包含一个身份验证令牌。

```go
res, err = client.R().
    SetAuthToken("bc594900518b4f7eac75bd37f019e08fbc594900518b4f7eac75bd37f019e08f").
    Head("https://myapp.com/videos/hi-res-video")

fmt.Println(err, res)
fmt.Println(res.Header())
```

--------------------------------

## Resty 重定向策略示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/redirect-policy.md

演示如何在 Resty 中设置内置的重定向策略。这包括设置具有最大次数的灵活重定向策略，以及组合多个策略，如灵活策略和域名检查策略。

```go
client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(5))

client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(5),
    resty.DomainCheckRedirectPolicy("host1.com", "host2.org", "host3.net"))
```

--------------------------------

## 执行 TRACE 请求

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/options-head-trace-request.md

展示如何使用 Resty 发送 TRACE 请求。TRACE 方法通常用于诊断目的，将请求回显给客户端。示例包括设置身份验证令牌。

```go
res, err = client.R().
    SetAuthToken("bc594900518b4f7eac75bd37f019e08fbc594900518b4f7eac75bd37f019e08f").
    Trace("https://myapp.com/test")

fmt.Println(err, res)
```

--------------------------------

## Resty v3 加权轮询负载均衡器示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/load-balancer-and-service-discovery.md

展示如何设置 Resty v3 加权轮询负载均衡器。这允许为不同的主机指定权重以控制请求分发。它还涵盖了设置恢复持续时间。

```go
wrr, err := resty.NewWeightedRoundRobin(
    3*time.Second, // 恢复持续时间
    []*resty.Host{
        {
            BaseURL: "https://example1.com",
            Weight:  50, // 确定到此主机的请求百分比
        },
        {BaseURL: "https://example2.com", Weight: 30},
        {BaseURL: "https://example3.com", Weight: 20},
    }...,
)
if err != nil {
    log.Printf("ERROR %v", err)
    return
}

// 默认情况下，恢复持续时间为 120 秒，可以按如下方式更改
// wrr.SetRecoveryDuration(3 * time.Minute)

c := resty.New().
    SetLoadBalancer(wrr)
defer c.Close()

// 开始使用客户端...
```

--------------------------------

## 生成 Curl 命令示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/curl-command.md

演示如何为 Resty 请求启用 curl 命令生成并检索生成的命令字符串。这涉及在请求上设置 `GenerateCurlCmd` 选项。

```go
c := resty.New()
deffer c.Close()

res, _ := c.R().
    SetGenerateCurlCmd(true).
    SetBody(map[string]string{
        "name": "Resty",
    }).
    Post("https://httpbin.org/post")

curlCmdStr := res.Request.CurlCmd()
fmt.Println(curlCmdStr)

// Result:
// curl -X POST -H 'Accept-Encoding: gzip, deflate' -H 'Content-Type: application/json' -H 'User-Agent: go-resty/3.0.0 (https://resty.dev)' -d '{"name":"Resty"}' https://httpbin.org/post
```

--------------------------------

## 管理连接事件 (OnOpen, OnError)

Source: https://github.com/go-resty/docs/blob/main/content/docs/server-sent-events.md

示例演示如何处理 Server-Sent Events 的连接级别事件。它设置了 `OnOpen` 的处理程序以确认连接，以及 `OnError` 的处理程序以记录任何连接错误。

```go
es := resty.NewEventSource().
    SetURL("https://sse.dev/test").
    OnMessage(
        func(e any) {
            fmt.Println(e.(*resty.Event))
        },
        nil,
    ).
    OnError(
        func(err error) {
			fmt.Println("Error occurred:", err)
		},
    ).
    OnOpen(
        func(url string) {
			fmt.Println("I'm connected:", url)
		},
    )

err := es.Get()
fmt.Println(err)

// Output:
//  I'm connected: https://sse.dev/test
//  &{  {"testing":true,"sse_dev":"is great","msg":"It works!","now":1737510458794}}
//  &{  {"testing":true,"sse_dev":"is great","msg":"It works!","now":1737510460794}}
//  &{  {"testing":true,"sse_dev":"is great","msg":"It works!","now":1737510462794}}
//  ...

```

--------------------------------

## 实现并设置自定义调试日志格式化程序

Source: https://github.com/go-resty/docs/blob/main/content/docs/debug-log.md

提供了一个创建自定义调试日志格式化函数并将其设置在 Resty 客户端上的示例。

```go
// 实现自定义调试日志格式化程序
func DebugLogCustomFormatter(dl *DebugLog) string {
    logContent := ""

    // 在此处执行日志操作逻辑

	return logContent
}

// 设置自定义调试日志格式化程序
c := resty.New().
    SetDebugLogFormatter(DebugLogCustomFormatter)
```

--------------------------------

## 带负载的 GET 请求

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/get-request.md

演示如何发送带有 JSON 负载的 GET 请求。这是通过将 `AllowMethodGetPayload` 设置为 true 并使用 `SetBody` 提供负载来实现的。

```go
res, err := client.R()
    .SetAllowMethodGetPayload(true) // 客户端级别的选项可用
    .SetContentType("application/json")
    .SetBody(`{
        "query": {
            "simple_query_string" : {
                "query": "\"fried eggs\" +(eggplant | potato) -frittata",
                "fields": ["title^5", "body"],
                "default_operator": "and"
            }
        }
    }`) // 这是字符串形式的请求正文
    .SetAuthToken("bc594900518b4f7eac75bd37f019e08fbc594900518b4f7eac75bd37f019e08f")
    .Get("/_search")

fmt.Println(err, res)
```

--------------------------------

## Go Resty 包级方法可用性

Source: https://github.com/go-resty/docs/blob/main/content/docs/upgrading-to-v3.md

本节列出了 Go Resty 中包级别的实用方法。这些方法涵盖了各种功能，例如检查字符串是否为空、内容类型检测和反序列化。

```APIDOC
IsStringEmpty
  - 检查字符串是否为空。

IsJSONType
  - 检测内容是否为 JSON 类型。

IsXMLType
  - 检测内容是否为 XML 类型。

DetectContentType
  - 检测数据的内容类型。

Unmarshalc
  - 带有内容类型检测的反序列化数据。

Backoff
  - 为重试实现退避策略。
```

--------------------------------

## Go Resty 客户端方法更新

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

## Go Resty TraceInfo 方法

Source: https://github.com/go-resty/docs/blob/main/content/docs/new-features-and-enhancements.md

本节记录了 Go Resty 库中用于格式化和访问跟踪信息的 GetTraceInfo 方法。它包括将跟踪信息作为字符串或 JSON 获取的方法。

```APIDOC
TraceInfo Methods:

String() string
  - 将跟踪信息作为格式化字符串返回。

JSON() string
  - 将跟踪信息作为 JSON 字符串返回。
```

--------------------------------

## 更新 go.mod 以使用 Resty v3

Source: https://github.com/go-resty/docs/blob/main/content/docs/upgrading-to-v3.md

此代码片段显示如何更新 go.mod 文件以使用 Resty 版本 3。它指定了最低要求的 Go 版本和要使用的 Resty 版本。

```bash
require resty.dev/v3 {{% param Resty.V3.Version %}}
```

--------------------------------

## Go Resty 请求方法重命名和更改

Source: https://github.com/go-resty/docs/blob/main/content/docs/upgrading-to-v3.md

本节详细介绍了 Go Resty Request 对象中方法的重命名和功能更改，以帮助迁移到 v3。

```APIDOC
Request Methods:

- `Request.QueryParams()`: 从 `QueryParam` 重命名。
- `Request.AuthToken()`: 从 `Token` 重命名。
- `Request.DoNotParseResponse()`: 从 `NotParseResponse` 重命名。
- `Request.SetExpectResponseContentType(contentType string)`: 从 `ExpectContentType` 重命名。
- `Request.SetForceResponseContentType(contentType string)`: 从 `ForceContentType` 重命名。
- `Request.SetOutputFileName(filename string)`: 从 `SetOutput` 重命名。
- `Request.EnableGenerateCurlCmd()`: 从 `EnableGenerateCurlOnDebug` 重命名。
- `Request.DisableGenerateCurlCmd()`: 从 `DisableGenerateCurlOnDebug` 重命名。
- `Request.CurlCmd()`: 从 `GenerateCurlCommand` 重命名。
- `Request.AddRetryConditions(conditions ...RetryConditionFunc)`: 从 `AddRetryCondition` 重命名。
```

--------------------------------

## Go Resty 客户端方法重命名和更改

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

## Go Resty：设置单个路径参数

Source: https://github.com/go-resty/docs/blob/main/content/docs/request-path-params.md

演示在 Go Resty 中为 GET 请求设置单个动态路径参数。参数值会自动进行 URL 编码。

```go
c := resty.New()
defere c.Close()

c.R().
    SetPathParam("userId", "sample@sample.com").
    Get("/v1/users/{userId}/details")

// Result:
//     /v1/users/sample@sample.com/details
```

--------------------------------

## Resty v3 SRV 加权轮询负载均衡器示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/load-balancer-and-service-discovery.md

演示 Resty v3 的 SRV 加权轮询负载均衡器的用法。此方法通过 SRV 记录发现服务并应用加权轮询分发。它包括设置刷新和恢复持续时间的选项。

```go
swrr, err := resty.NewSRVWeightedRoundRobin(
    "_sample-server",
    "tcp", // 默认协议为 tcp
    "example.com",
    "https", // 默认方案为 https
)
if err != nil {
    log.Printf("ERROR %v", err)
    return
}

// 默认情况下，SRV 记录刷新持续时间为 180 秒，可以按如下方式更改
// swrr.SetRefreshDuration(1 * time.Hour)

// 默认情况下，恢复持续时间为 120 秒，可以按如下方式更改
// swrr.SetRecoveryDuration(3 * time.Minute)

c := resty.New().
    SetLoadBalancer(swrr)
defer c.Close()

// 开始使用客户端...
```

--------------------------------

## 获取 GitHub Release Tag

Source: https://github.com/go-resty/docs/blob/main/layouts/shortcodes/restyrelease.html

此代码片段演示如何使用 Go 模板函数构建 GitHub release tag URL。它检索 GitHub 存储库路径和版本号以创建指向特定版本的链接。

```go
package main

import "fmt"

func main() {
	// 示例用法：
	// 假设 $gh_repo 是 "https://github.com/go-resty/resty"
	// 假设 $version 是 "v1.14.0"
	ghRepo := "https://github.com/go-resty/resty"
	version := "v1.14.0"
	releaseTagURL := fmt.Sprintf("%s/releases/tag/%s", ghRepo, version)
	fmt.Println(releaseTagURL)
}

```

--------------------------------

## 无限响应体读取示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/unlimited-response-body-reads.md

演示如何在 Go Resty 中启用无限响应体读取、自动解析响应并将响应保存到文件。此功能会将响应体保留在内存中以供多次读取。

```go
loginResponseFile := "login-response.txt"

res, err := client.R().
    SetHeader(hdrContentTypeKey, "application/json").
    SetBody(&User{Username: "testuser", Password: "testpass"}).
    SetResponseBodyUnlimitedReads(true).
    SetResult(&LoginResponse{}).
    SetOutputFileName(loginResponseFile).
    Post("/login")

fmt.Println(err)

fmt.Println("")
loginResponse := res.Result().(*LoginResponse)
fmt.Println("ID:", loginResponse.ID)
fmt.Println("Message:", loginResponse.Message)

fmt.Println("")
loginResponseCnt, _ := os.ReadFile(loginResponseFile)
fmt.Println("File Content:", string(loginResponseCnt))
```

--------------------------------

## Go Resty SetResult 和 SetError 用法示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/response-auto-parse.md

演示使用 SetResult 的不同方法，包括传递内联指针、非指针值和预先声明的指针。SetError 的原理相同。

```go
// 设置
client.R().SetResult(&LoginResponse{})

// 访问
fmt.Println(res.Result().(*LoginResponse))
```

```go
// 设置
client.R().SetResult(LoginResponse{})

// 访问
fmt.Println(res.Result().(*LoginResponse))
```

```go
loginResponse := &LoginResponse{}
// 设置
client.R().SetResult(loginResponse)

// 访问
fmt.Println(loginResponse)
```

--------------------------------

## 保存所有响应 - Go Resty

Source: https://github.com/go-resty/docs/blob/main/content/docs/save-response.md

此示例演示如何配置 Resty 客户端将所有 HTTP 响应保存到指定目录。`SetSaveResponse(true)` 选项适用于此客户端进行的所有后续请求。请记住在完成后关闭客户端。

```go
c := resty.New().
    SetOutputDirectory("/path/to/save/all/response").
    SetSaveResponse(true) // 适用于所有请求
defer c.Close()

// 开始使用客户端...
```

--------------------------------

## Resty 客户端的负载允许方法

Source: https://github.com/go-resty/docs/blob/main/content/docs/allow-payload-on.md

Resty 客户端上的这些方法允许您配置 GET 和 DELETE 请求是否可以包含负载。这对于遵循标准 HTTP 动词语义的系统很有用。

```APIDOC
Client.SetAllowMethodGetPayload()
  允许 GET HTTP 动词上的请求负载。

Client.SetAllowMethodDeletePayload()
  允许 DELETE HTTP 动词上的请求负载。
```

--------------------------------

## Go Resty SetMultipartFields 示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/multipart.md

此 Go 代码片段演示如何使用 `SetMultipartFields` 方法构建 multipart/form-data 请求。它展示了添加简单的表单字段、使用文件路径上传文件、包括进度回调、指定文件名和内容类型以及从 `io.Reader` 上传数据。

```go
myImageFile, _ := os.Open("/path/to/image-1.png")
myImageFileStat, _ := myImageFile.Stat()

// 使用各种组合和可能性进行演示
client.R().
    SetMultipartFields(
        []*resty.MultipartField{
            // 添加表单数据，顺序得以保留
            {
                Name:   "field1",
                Values: []string{"field1value1", "field1value2"},
            },
            {
                Name:   "field2",
                Values: []string{"field2value1", "field2value2"},
            },
            // 添加文件上传
            {
                Name:             "myfile_1",
                FilePath:         "/path/to/file-1.txt",
            },
            // 添加带有进度回调的文件上传
            {
                Name:             "myfile_1",
                FilePath:         "/path/to/file-1.txt",
                ProgressCallback: func(mp MultipartFieldProgress) {
    				// 使用进度详细信息
    				},
            },
            // 带有文件名和内容类型
            {
                Name:             "myimage_1",
                FileName:         "image-1.png",
                ContentType:      "image/png",
                FilePath:         "/path/to/image-1.png",
            },
            // 带有 io.Reader 和文件大小
            {
                Name:             "myimage_2",
                FileName:         "image-2.png",
                ContentType:      "image/png",
                Reader:           myImageFile,
                FileSize:         myImageFileStat.Size(),
            },
            // 带有 io.Reader
            {
                Name:        "uploadManifest1",
                FileName:    "upload-file-1.json",
                ContentType: "application/json",
                Reader:      strings.NewReader(`{"input": {"name": "Uploaded document 1", "_filename" : ["file1.txt"]}}`),
            },
            // 带有 io.Reader 和进度回调
            {
                Name:             "image-file1",
                FileName:         "image-file1.png",
                ContentType:      "image/png",
                Reader:           bytes.NewReader(fileBytes),
                ProgressCallback: func(mp MultipartFieldProgress) {
                    // 使用进度详细信息
                },
            },
        }...
    )

```

--------------------------------

## Go Resty 包级类型更改

Source: https://github.com/go-resty/docs/blob/main/content/docs/upgrading-to-v3.md

本节详细介绍了 Go Resty 包中类型的更改。它解释了某些类型如何变得未导出、被新功能替换或集成到增强功能中，如负载均衡和多部分字段。

```APIDOC
User
  - 现在已设为未导出，称为 'credentials'。

SRVRecord
  - 已被支持 SRV 记录查找的新负载均衡器功能取代。

File
  - 已被增强的 MultipartField 功能取代。

RequestLog, ResponseLog
  - 已弃用：请改用 DebugLog。

RequestLogCallback, ResponseLogCallback
  - 已弃用：请改用 DebugLogCallbackFunc。
```

--------------------------------

## 带 JSON 正文和身份验证令牌的 PUT 请求

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/post-put-patch-request.md

演示如何发送带有 JSON 正文和身份验证令牌的 PUT 请求。此示例包括设置请求正文、身份验证以及处理潜在错误。

```go
res, err := client.R().
    SetBody(Article{
        Title: "Resty",
        Content: "This is my article content, oh ya!",
        Author: "Jeevanandam M",
        Tags: []string{"article", "sample", "resty"},
    }). // 默认请求内容类型为 JSON
    SetAuthToken("bc594900518b4f7eac75bd37f019e08fbc594900518b4f7eac75bd37f019e08f").
    SetError(&Error{}). // 或 SetError(Error{}).
    Put("https://myapp.com/articles/123456")

fmt.Println(err, res)
fmt.Println(res.Error().(*Error))
```

--------------------------------

## Go Resty：设置单个原始路径参数

Source: https://github.com/go-resty/docs/blob/main/content/docs/request-path-params.md

演示在 Go Resty 中为 GET 请求设置单个动态原始路径参数。参数值按原样使用，不进行 URL 编码。

```go
c := resty.New()
defere c.Close()

c.R().
    SetRawPathParam("path", "groups/developers").
    Get("/v1/users/{userId}/details")

// Result:
//     /v1/users/groups/developers/details
```

--------------------------------

## 创建 Resty 客户端

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/options-head-trace-request.md

初始化一个新的 Resty 客户端实例。建议延迟关闭客户端以正确释放资源。

```go
client := resty.New()
def client.Close()
```

--------------------------------

## Resty v3 破坏性更改和行为更新

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

## Go Resty 响应方法

Source: https://github.com/go-resty/docs/blob/main/content/docs/new-features-and-enhancements.md

本节记录了用于访问和检查 Go 中 Resty 响应对象的 GetTraceInfo 方法。它包括获取响应正文、字节、检查正文是否已读取、访问错误以及获取重定向历史记录的方法。

```APIDOC
Response Methods:

Body() string
  - 将响应正文作为字符串返回。

Bytes() []byte
  - 将响应正文作为字节切片返回。

IsRead() bool
  - 检查响应正文是否已被读取。

Err() error
  - 返回请求期间遇到的任何错误。

RedirectHistory() []*Request
  - 返回响应的重定向历史记录。
```

--------------------------------

## Resty 请求的负载允许方法

Source: https://github.com/go-resty/docs/blob/main/content/docs/allow-payload-on.md

Resty Request 对象上的这些方法允许您配置特定的 GET 或 DELETE 请求是否可以包含负载。这提供了对单个请求负载包含的精细控制。

```APIDOC
Request.SetAllowMethodGetPayload()
  允许此特定请求的 GET HTTP 动词上的请求负载。

Request.SetAllowMethodDeletePayload()
  允许此特定请求的 DELETE HTTP 动词上的请求负载。
```

--------------------------------

## Go Resty：设置多个原始路径参数

Source: https://github.com/go-resty/docs/blob/main/content/docs/request-path-params.md

演示在 Go Resty 中使用 map 设置多个动态原始路径参数以进行 GET 请求。值按原样使用，不进行 URL 编码。

```go
c := resty.New()
defere c.Close()

c.R().
    SetRawPathParams(map[string]string{
        "userId":       "sample@sample.com",
        "subAccountId": "100002",
        "path":         "groups/developers",
    }).
    Get("/v1/users/{userId}/{subAccountId}/{path}/details")

// Result:
//     /v1/users/sample@sample.com/100002/groups/developers/details
```

--------------------------------

## 自动反序列化 Server-Sent Events

Source: https://github.com/go-resty/docs/blob/main/content/docs/server-sent-events.md

示例展示了如何自动将传入的 SSE 数据反序列化为 Go 结构体。它定义了一个带有 JSON 标签的 `Data` 结构体，并将其与 `OnMessage` 一起使用，以类型安全地处理事件负载。

```go
// https://sse.dev/test 返回
// {"testing":true,"sse_dev":"is great","msg":"It works!","now":1737508994502}
type Data struct {
    Testing bool   `json:"testing"`
    SSEDev  string `json:"sse_dev"`
    Message string `json:"msg"`
    Now     int64  `json:"now"`
}

es := resty.NewEventSource().
    SetURL("https://sse.dev/test").
    OnMessage(
        func(e any) {
            d := e.(*Data)
            fmt.Println("Testing:", d.Testing)
            fmt.Println("SSEDev:", d.SSEDev)
            fmt.Println("Message:", d.Message)
            fmt.Println("Now:", d.Now)
            fmt.Println("")
        },
        Data{},
    )

err := es.Get()
fmt.Println(err)

// Output:
//     Testing: true
//     SSEDev: is great
//     Message: It works!
//     Now: 1737509497652

//     Testing: true
//     SSEDev: is great
//     Message: It works!
//     Now: 1737509499652

//     ...

```

--------------------------------

## 保存单个响应 (从 URL 获取文件名) - Go Resty

Source: https://github.com/go-resty/docs/blob/main/content/docs/save-response.md

此示例演示如何将单个 HTTP 响应保存到文件系统。文件名是从请求的 URL 自动确定的。在这种情况下，图像将保存为“resty-logo.svg”。

```go
client.R().
    SetSaveResponse(true).
    Get("https://resty.dev/svg/resty-logo.svg")
```

--------------------------------

## Go Resty：设置多个路径参数

Source: https://github.com/go-resty/docs/blob/main/content/docs/request-path-params.md

演示在 Go Resty 中使用 map 设置多个动态路径参数以进行 GET 请求。值会被 URL 编码，包括路径段中的特殊字符。

```go
c := resty.New()
defere c.Close()

c.R().
    SetPathParams(map[string]string{
        "userId":       "sample@sample.com",
        "subAccountId": "100002",
        "path":         "groups/developers",
    }).
    Get("/v1/users/{userId}/{subAccountId}/{path}/details)

// Result:
//   /v1/users/sample@sample.com/100002/groups%2Fdevelopers/details
```

--------------------------------

## Resty v2 文档链接

Source: https://github.com/go-resty/docs/blob/main/content/_index.md

提供 Resty v2 文档的链接，供仍在使用或从旧版本迁移的用户参考。

```APIDOC
README.md: https://github.com/go-resty/resty/blob/v2/README.md
```

--------------------------------

## Go Resty API 文档

Source: https://github.com/go-resty/docs/blob/main/layouts/shortcodes/hintreqoverride.html

提供 Go Resty HTTP 客户端的关键方法和配置的文档。这包括客户端创建、设置标头、发出请求和处理响应。

```APIDOC
resty.New()
  - 创建并返回一个新的 Resty 客户端实例。
  - 无参数。
  - 返回值：*resty.Client

Client.R()
  - 创建一个与客户端关联的新 Request 对象。
  - 无参数。
  - 返回值：*resty.Request

Request.SetHeader(name, value string)
  - 设置请求标头。
  - 参数：
    - name: 标头的名称。
    - value: 标头的值。
  - 返回值：*resty.Request (用于链式调用)

Request.Get(url string)
  - 向指定 URL 执行 HTTP GET 请求。
  - 参数：
    - url: 要发送 GET 请求的 URL。
  - 返回值：
    - *resty.Response: HTTP 响应对象。
    - error: 如果请求失败，则为错误。

Response.String() string
  - 将响应正文作为字符串返回。
  - 无参数。
  - 返回值：string

Client-Level Settings:
  - 应用于客户端实例的设置（例如，默认标头、超时）。
  - 可以在请求级别覆盖。

Request-Level Settings:
  - 应用于特定请求的设置（例如，自定义标头、查询参数）。
  - 指定时会覆盖客户端级别设置。

```

--------------------------------

## 使用 quic-go 启用 HTTP3

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

## Go Resty 响应自动解析示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/response-auto-parse.md

演示如何使用 Resty 通过 SetResult 和 SetError 将 JSON 响应自动反序列化为 Go 结构体。它展示了如何设置请求正文、指定结果和错误类型以及访问解析后的结果。

```go
res, err := client.R().
    SetBody(User{
        Username: "testuser",
        Password: "testpass",
    }). // 默认请求内容类型为 JSON
    SetResult(&LoginResponse{}).
    SetError(&LoginError{}).
    Post("https://myapp.com/login")

fmt.Println(err)
fmt.Println(res.Result().(*LoginResponse))
fmt.Println(res.Error().(*LoginError))
```

--------------------------------

## Go Resty 客户端创建

Source: https://github.com/go-resty/docs/blob/main/content/docs/new-features-and-enhancements.md

演示了使用不同配置创建 Resty 客户端的各种方法，包括传输设置和拨号器。

```go
client := resty.New()

// 覆盖所有传输设置和超时值
client.SetTransport(NewWithTransportSettings(settings))

// 使用自定义拨号器创建客户端
client.SetDialer(NewWithDialer(dialer))

// 使用自定义拨号器和传输设置创建客户端
client = resty.NewWithDialerAndTransportSettings(dialer, settings)
```

--------------------------------

## 保存单个响应 (自定义文件名) - Go Resty

Source: https://github.com/go-resty/docs/blob/main/content/docs/save-response.md

此示例演示如何将单个 HTTP 响应保存为自定义文件名。`SetOutputFileName` 方法允许您为保存的文件指定相对或绝对路径。在此示例中，图像将保存为“resty-logo-blue.svg”。

```go
client.R().
    SetSaveResponse(true).
    SetOutputFileName("resty-logo-blue.svg"). // 可以是相对或绝对路径
    Get("https://resty.dev/svg/resty-logo.svg")
```

--------------------------------

## 执行 OPTIONS 请求

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/options-head-trace-request.md

演示如何使用 Resty 执行 OPTIONS 请求。这通常用于确定允许的 HTTP 方法或用于 CORS 预检请求。它包括设置身份验证令牌。

```go
res, err := client.R().
    SetAuthToken("bc594900518b4f7eac75bd37f019e08fbc594900518b4f7eac75bd37f019e08f").
    Options("https://myapp.com/servers/nyc-dc-01")

fmt.Println(err, res)
fmt.Println(res.Header())
```

--------------------------------

## 创建 Resty 客户端

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/post-put-patch-request.md

初始化一个新的 Resty 客户端实例。建议延迟调用客户端的 Close 方法以正确释放资源。

```go
client := resty.New()
def client.Close()
```

--------------------------------

## Go Resty 客户端配置

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

## Resty 客户端根证书方法

Source: https://github.com/go-resty/docs/blob/main/content/docs/root-certificates.md

提供了 Resty 客户端上可用的管理根证书的方法的概述，包括从文件、带监视器的文件和字符串设置。

```APIDOC
Client.SetRootCertificates(paths ...string)
  - 从文件路径添加一个或多个 PEM 编码的根证书。
  - 参数：
    - paths: PEM 编码证书文件的文件路径的可变参数列表。
  - 返回值：Resty 客户端实例，用于链式调用。

Client.SetRootCertificatesWatcher(opts *CertWatcherOptions, paths ...string)
  - 从文件路径添加一个或多个 PEM 编码的根证书，并带有用于动态重新加载的监视器。
  - 参数：
    - opts: 证书监视器的配置选项（例如，PoolInterval）。
    - paths: PEM 编码证书文件的文件路径的可变参数列表。
  - 返回值：Resty 客户端实例，用于链式调用。

Client.SetRootCertificateFromString(cert string)
  - 从字符串添加单个 PEM 编码的根证书。
  - 参数：
    - cert: 包含 PEM 编码证书的字符串。
  - 返回值：Resty 客户端实例，用于链式调用。

CertWatcherOptions struct {
  PoolInterval time.Duration // 检查证书修改的间隔。默认为 24 小时。
}
```

--------------------------------

## 从文件添加客户端根证书

Source: https://github.com/go-resty/docs/blob/main/content/docs/client-root-certificates.md

此 Go 代码片段演示如何使用一个或多个 PEM 文件路径为 Resty 客户端设置客户端根证书。它支持传递单个路径或路径切片。

```go
client.SetClientRootCertificates("/path/to/client-root/pemFile.pem")

client.SetClientRootCertificates(
    "/path/to/client-root/pemFile1.pem",
    "/path/to/client-root/pemFile2.pem"
    "/path/to/client-root/pemFile3.pem"
)

client.SetClientRootCertificates(certs...)
```