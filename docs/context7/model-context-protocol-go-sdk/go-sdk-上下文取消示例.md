# Go SDK 上下文取消示例

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

演示如何使用上下文取消在 Go SDK 中实现操作取消。创建一个带有取消函数的新的上下文，并在 goroutine 中调用操作。

```Go
import "context"

ctx, cancel := context.WithCancel(ctx)
go session.CallTool(ctx, "slow", map[string]any{}, nil)
cancel()
```

--------------------------------
