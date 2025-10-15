# Go：定义 CallToolParams 和 Tool 结构

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了用于类型化工具参数的通用 `CallToolParamsFor` 结构以及用于表示 MCP 工具的 `Tool` 结构，包括模式和描述。

```Go
type CallToolParamsFor[In any] struct {
	Meta      Meta   `json:"_meta,omitempty"`
	Arguments In     `json:"arguments,omitempty"`
	Name      string `json:"name"`
}

type Tool struct {
	Annotations *ToolAnnotations   `json:"annotations,omitempty"`
	Description string             `json:"description,omitempty"`
	InputSchema *jsonschema.Schema `json:"inputSchema"`
	Name string                    `json:"name"`
}
```

--------------------------------
