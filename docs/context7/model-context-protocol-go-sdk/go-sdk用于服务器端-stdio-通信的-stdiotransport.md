# Go SDK：用于服务器端 Stdio 通信的 StdioTransport

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了 stdio 传输服务器端 `StdioTransport`。此传输通过绑定到 os.Stdin 和 os.Stdout 来连接，并通过换行符分隔的 JSON 进行通信。

```Go
// A StdioTransport is a [Transport] that communicates using newline-delimited
// JSON over stdin/stdout.
type StdioTransport struct { }

func (t *StdioTransport) Connect(context.Context) (Connection, error)
```

--------------------------------
