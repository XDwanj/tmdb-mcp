# 用于 Tool 分页的迭代器方法

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

此 Go 代码为 ListTools Spec 方法定义了一个迭代器方法。它自动处理分页，允许遍历所有页面的工具。如果提供了参数，则迭代从指定的游标开始。

```Go
func (*ClientSession) Tools(context.Context, *ListToolsParams) iter.Seq2[Tool, error]
```

--------------------------------
