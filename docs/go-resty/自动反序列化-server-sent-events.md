# 自动反序列化 Server-Sent Events

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
