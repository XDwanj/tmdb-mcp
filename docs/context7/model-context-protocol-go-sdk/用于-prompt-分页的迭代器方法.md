# 用于 Prompt 分页的迭代器方法

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

此 Go 代码为 ListPrompts Spec 方法定义了一个迭代器方法。它自动处理分页，允许遍历所有页面的提示。如果提供了参数，则迭代从指定的游标开始。

```Go
func (*ClientSession) Prompts(context.Context, *ListPromptsParams) iter.Seq2[Prompt, error]
```

--------------------------------
