# 在 Go 中添加和删除 Prompt

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

演示如何使用 Server.AddPrompt 和 Server.RemovePrompts 方法添加带有处理程序的 Prompt，然后删除它。它包括 Prompt 处理程序参数的定义。

```Go
type codeReviewArgs struct {
  Code string `json:"code"`
}

func codeReviewHandler(context.Context, *ServerSession, *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {...}

server.AddPrompt(
  &mcp.Prompt{Name: "code_review", Description: "review code"},
  codeReviewHandler,
)

server.RemovePrompts("code_review")
```

--------------------------------
