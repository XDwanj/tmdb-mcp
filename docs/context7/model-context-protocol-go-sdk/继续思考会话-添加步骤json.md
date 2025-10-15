# 继续思考会话 - 添加步骤（JSON）

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/sequentialthinking/README.md

将下一个思考或分析步骤添加到正在进行的思考会话中。它需要会话 ID 和思考内容。

```json
{
  "method": "tools/call", 
  "params": {
    "name": "continue_thinking",
    "arguments": {
      "sessionId": "architecture_design",
      "thought": "First, I need to identify the core business domains and their boundaries to determine service decomposition."
    }
  }
}
```

--------------------------------
