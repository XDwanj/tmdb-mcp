# 入门 Server-Sent Events

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
