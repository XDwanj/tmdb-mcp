# 启动思考会话（JSON）

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/sequentialthinking/README.md

为给定的问题启动一个新的顺序思考会话。它接受问题陈述、可选的会话 ID 和步骤的初始估计。

```json
{
  "method": "tools/call",
  "params": {
    "name": "start_thinking",
    "arguments": {
      "problem": "How should I design a scalable microservices architecture?",
      "sessionId": "architecture_design",
      "estimatedSteps": 8
    }
  }
}
```

--------------------------------
