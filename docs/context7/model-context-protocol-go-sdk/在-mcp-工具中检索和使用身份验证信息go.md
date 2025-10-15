# 在 MCP 工具中检索和使用身份验证信息（Go）

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/auth-middleware/README.md

演示了如何在 MCP 工具中的传入请求中访问身份验证信息，特别是 `TokenInfo`（包括范围）。它包括在执行工具逻辑之前检查所需范围。

```Go
// Get authentication information in MCP tool
func MyTool(ctx context.Context, req *mcp.CallToolRequest, args MyArgs) (*mcp.CallToolResult, any, error) {
    // Extract authentication info from request 
    userInfo := req.Extra.TokenInfo
    
    // Check scopes
    if !slices.Contains(userInfo.Scopes, "read") {
        return nil, nil, fmt.Errorf("insufficient permissions: read scope required")
    }
    
    // Execute tool logic
    return &mcp.CallToolResult{
        Content: []mcp.Content{
            &mcp.TextContent{Text: "Tool executed successfully"},
        },
    }, nil, nil
}
```

--------------------------------
