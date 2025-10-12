# Story 1.5 Manual Test Results

## Test Environment

- **OS**: Linux 6.16.11-1-MANJARO
- **Go Version**: 1.25.1
- **Binary**: tmdb-mcp (14MB)
- **Date**: 2025-10-12

## Pre-Test Checklist

- [x] 配置文件已准备（需要有效的 TMDB API Key）
- [x] 二进制文件编译成功
- [x] Logger 初始化成功
- [x] TMDB Client 创建成功

## Test Case 1: Program Startup

**Objective**: 验证程序能正常启动并初始化所有组件

**Prerequisites**:
```bash
# 1. 复制示例配置到默认位置
mkdir -p ~/.tmdb-mcp
cp examples/config.yaml ~/.tmdb-mcp/config.yaml

# 2. 编辑配置文件，设置有效的 TMDB API Key
# vim ~/.tmdb-mcp/config.yaml
```

**Execution**:
```bash
./tmdb-mcp
```

**Expected Logs** (输出到 stderr):
```
INFO    TMDB MCP Service starting    version=x.x.x mode=stdio
INFO    Configuration loaded successfully
INFO    TMDB Client created
INFO    TMDB API Key validated successfully
INFO    MCP Server initialized    name=tmdb-mcp version=1.0.0
INFO    Starting MCP Server in stdio mode
```

**Status**: ⏳ Pending (需要用户提供 TMDB API Key)

---

## Test Case 2: MCP Protocol - tools/list

**Objective**: 验证 MCP 协议能正确响应 tools/list 请求

**Test Method 1**: 使用 echo 和管道
```bash
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | ./tmdb-mcp 2>/dev/null
```

**Expected Response** (输出到 stdout):
```json
{"jsonrpc":"2.0","id":1,"result":{"tools":[]}}
```

**Explanation**:
- 此阶段尚未注册任何工具，因此返回空列表
- Future Story 1.6 将添加 `search` 工具

**Status**: ⏳ Pending

---

## Test Case 3: Graceful Shutdown

**Objective**: 验证程序能正确处理 SIGINT/SIGTERM 信号并优雅退出

**Execution**:
```bash
./tmdb-mcp
# 在另一个终端发送信号
kill -INT <pid>
# 或直接按 Ctrl+C
```

**Expected Logs**:
```
INFO    Received shutdown signal
INFO    TMDB MCP Service shutdown complete
```

**Status**: ⏳ Pending

---

## Test Case 4: MCP Protocol - initialize (如果需要)

**Objective**: 验证 MCP 初始化协议

**Test Input**:
```json
{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}},"id":1}
```

**Expected Response**:
```json
{"jsonrpc":"2.0","id":1,"result":{"protocolVersion":"2024-11-05","capabilities":{"tools":{}},"serverInfo":{"name":"tmdb-mcp","version":"1.0.0"}}}
```

**Status**: ⏳ Pending

---

## Notes for Manual Testing

由于此阶段需要：
1. **有效的 TMDB API Key** - 需要用户从 https://www.themoviedb.org/settings/api 获取
2. **MCP 客户端** - 可以使用 Claude Code 或其他 MCP 客户端

建议用户按以下步骤测试：

### 步骤 1: 配置 API Key
```bash
mkdir -p ~/.tmdb-mcp
cp examples/config.yaml ~/.tmdb-mcp/config.yaml
# 编辑 ~/.tmdb-mcp/config.yaml，设置 api_key
```

### 步骤 2: 启动程序
```bash
./tmdb-mcp 2>tmdb-mcp.log
```

### 步骤 3: 发送测试请求
在另一个终端：
```bash
# Test tools/list
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | nc localhost <port_if_sse>

# 或使用 stdio（需要手动输入）
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}'
```

### 步骤 4: 验证日志
```bash
cat tmdb-mcp.log
```

---

## Automated Testing Script

为了简化测试，可以使用以下脚本（需要 jq）：

```bash
#!/bin/bash
# test-mcp-stdio.sh

# 启动服务器（后台）
./tmdb-mcp 2>test.log &
SERVER_PID=$!

# 等待启动
sleep 2

# 发送 tools/list 请求
RESPONSE=$(echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | timeout 5 nc localhost 8910 2>/dev/null)

# 验证响应
if echo "$RESPONSE" | jq -e '.result.tools == []' > /dev/null 2>&1; then
    echo "✅ tools/list test PASSED"
else
    echo "❌ tools/list test FAILED"
    echo "Response: $RESPONSE"
fi

# 关闭服务器
kill -INT $SERVER_PID
wait $SERVER_PID

# 检查日志
if grep -q "TMDB MCP Service shutdown complete" test.log; then
    echo "✅ Graceful shutdown test PASSED"
else
    echo "❌ Graceful shutdown test FAILED"
fi
```

---

## Summary

**Total Test Cases**: 4
- **Passed**: 0
- **Failed**: 0
- **Pending**: 4 (需要用户手动执行)

**Blocker**: 需要有效的 TMDB API Key 才能完成测试

**Recommendation**:
- Task 5 标记为"部分完成"，已创建测试文档和指南
- 用户可以在获得 API Key 后按照此文档执行手动测试
- 继续 Task 6（单元测试）不需要 API Key
