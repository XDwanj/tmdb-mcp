# Go：向服务器添加具有显式架构的工具

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

演示如何向具有显式定义的 `Tool` 结构（包括其名称、描述和输入架构）的服务器添加工具。

```Go
t := &Tool{Name: ..., Description: ..., InputSchema: &jsonschema.Schema{...}}
server.AddTool(t, myHandler)
```

--------------------------------
