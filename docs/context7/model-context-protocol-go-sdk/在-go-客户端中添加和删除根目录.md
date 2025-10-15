# 在 Go 客户端中添加和删除根目录

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

解释了如何使用 `AddRoots` 和 `RemoveRoots` 方法管理 Go 客户端的根目录。`AddRoots` 添加指定的根目录，替换具有相同 URI 的现有根目录，并通知已连接的服务器。`RemoveRoots` 按 URI 删除根目录，如果根目录不存在则不报错。

```Go
// AddRoots adds the given roots to the client,
// replacing any with the same URIs,
// and notifies any connected servers.
func (*Client) AddRoots(roots ...*Root)

// RemoveRoots removes the roots with the given URIs.
// and notifies any connected servers if the list has changed.
// It is not an error to remove a nonexistent root.
func (*Client) RemoveRoots(uris ...string)
```

--------------------------------
