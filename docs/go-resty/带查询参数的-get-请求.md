# 带查询参数的 GET 请求

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
