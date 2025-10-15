# Go SDK 进度通知参数

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了 Go SDK 中用于请求和发送进度通知的结构。`Meta` 包含 `ProgressToken` 以指示请求进度更新。

```Go
type XXXParams struct { // where XXX is each type of call
  Meta Meta
  ...
}

type Meta struct {
  Data          map[string]any // arbitrary data
  ProgressToken any // string or int
}
```

--------------------------------
