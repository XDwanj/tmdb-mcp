# 调用带有 nil 参数的 Spec 方法

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

此 Go 示例展示了如何通过为当前不需要的参数传递 nil 来调用 Ping 等 Spec 方法。这种方法确保了兼容性，即使 Spec 将来引入新参数。

```Go
err := session.Ping(ctx, nil)
```

--------------------------------
