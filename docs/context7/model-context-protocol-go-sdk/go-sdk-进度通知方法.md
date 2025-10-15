# Go SDK 进度通知方法

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

提供了 Go SDK 中发送进度通知的方法。如果请求中提供了进度令牌，`NotifyProgress` 会向对等方发送通知。

```Go
import "context"

func (*ClientSession) NotifyProgress(context.Context, *ProgressNotification)
func (*ServerSession) NotifyProgress(context.Context, *ProgressNotification)
```

--------------------------------
