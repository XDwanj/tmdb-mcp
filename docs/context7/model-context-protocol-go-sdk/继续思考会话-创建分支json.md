# 继续思考会话 - 创建分支（JSON）

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/sequentialthinking/README.md

在思考会话中创建替代的推理路径或分支。它需要会话 ID 和新分支的思考内容。

```json
{
  "method": "tools/call",
  "params": {
    "name": "continue_thinking",
    "arguments": {
      "sessionId": "architecture_design", 
      "thought": "Alternative approach: Start with a monolith-first strategy and extract services gradually.",
      "createBranch": true
    }
  }
}
```

--------------------------------
