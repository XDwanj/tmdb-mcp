# Go Resty API 文档

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
