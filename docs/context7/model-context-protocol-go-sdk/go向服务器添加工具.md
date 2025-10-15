# Go：向服务器添加工具

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

演示如何使用 `AddTool` 方法或函数向 MCP 服务器添加工具。该函数是通用的，并从处理程序参数推断模式。

```Go
func (s *Server) AddTool(t *Tool, h ToolHandler)
func AddTool[In, Out any](s *Server, t *Tool, h ToolHandlerFor[In, Out])
```

--------------------------------
