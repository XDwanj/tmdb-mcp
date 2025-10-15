# Go：带有类型化参数的示例处理程序

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

提供了一个 Go 处理程序函数 (`addHandler`) 的示例，该函数处理类型化参数 (`AddParams`) 并为 MCP 工具返回类型化结果 (`int`)。

```Go
type AddParams struct {
    X int `json:"x"`
    Y int `json:"y"`
}

func addHandler(ctx context.Context, req *mcp.ServerRequest[*mcp.CallToolParamsFor[AddParams]]) (*mcp.CallToolResultFor[int], error) {
    return &mcp.CallToolResultFor[int]{StructuredContent: req.Params.Arguments.X + req.Params.Arguments.Y}, nil
}
```

--------------------------------
