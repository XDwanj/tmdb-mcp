# Go SDK：用于 Stdio 通信的 CommandTransport

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了 stdio 传输的 CommandTransport。此类型通过启动命令并通过其 stdin 和 stdout 流式传输 JSON-RPC 消息进行连接，使用换行符分隔的 JSON 进行通信。

```Go
// A CommandTransport is a [Transport] that runs a command and communicates
// with it over stdin/stdout, using newline-delimited JSON.
type CommandTransport struct { Command *exec.Command }

// Connect starts the command, and connects to it over stdin/stdout.
func (*CommandTransport) Connect(ctx context.Context) (Connection, error) {
```

--------------------------------
