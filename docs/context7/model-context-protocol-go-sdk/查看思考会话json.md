# 查看思考会话（JSON）

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/sequentialthinking/README.md

检索对特定思考会话的完整审查，包括所有思考步骤及其历史记录。它需要会话 ID。

```json
{
  "method": "tools/call", 
  "params": {
    "name": "review_thinking", 
    "arguments": {
      "sessionId": "architecture_design"
    }
  }
}
```

--------------------------------
