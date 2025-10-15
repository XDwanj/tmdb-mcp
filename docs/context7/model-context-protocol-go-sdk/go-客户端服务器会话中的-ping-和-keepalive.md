# Go 客户端/服务器会话中的 Ping 和 KeepAlive

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

演示了 `ClientSession` 和 `ServerSession` 的 `Ping` 方法，用于检查对等连接。它还展示了如何通过 `KeepAlive` 选项配置自动保持活动行为，该选项会在对等方未能响应 ping 时关闭会话。

```Go
func (c *ClientSession) Ping(ctx context.Context, *PingParams) error
func (c *ServerSession) Ping(ctx context.Context, *PingParams) error
```

```Go
type ClientOptions struct {
  ...
  KeepAlive time.Duration
}

type ServerOptions struct {
  ...
  KeepAlive time.Duration
}
```

--------------------------------
