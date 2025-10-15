# 用于日志记录和分页的 Go 服务器选项

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

配置服务器端的日志记录和分页行为。包括日志记录器名称、日志记录间隔和页面大小的选项。

```Go
package mcp

import "time"

type ServerOptions struct {
  // ...
  // The value for the "logger" field of the notification.
  LoggerName string
  // Log notifications to a single ClientSession will not be
  // sent more frequently than this duration.
  LoggingInterval time.Duration
  // PageSize defines the number of items to return per page.
  PageSize int
}
```

--------------------------------
