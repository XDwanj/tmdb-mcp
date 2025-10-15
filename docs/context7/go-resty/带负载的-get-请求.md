# 带负载的 GET 请求

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
