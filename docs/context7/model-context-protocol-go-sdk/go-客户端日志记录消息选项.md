# Go 客户端日志记录消息选项

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了用于接收服务器日志记录消息的客户端处理程序。当收到 `LoggingMessageNotification` 时，会调用此回调。

```Go
package mcp

import "context"

type ClientOptions struct {
  // ...
  LoggingMessageHandler func(context.Context, *ClientSession, *LoggingMessageParams)
}
```

--------------------------------
