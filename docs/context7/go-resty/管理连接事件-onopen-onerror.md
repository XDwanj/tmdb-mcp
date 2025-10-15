# 管理连接事件 (OnOpen, OnError)

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
