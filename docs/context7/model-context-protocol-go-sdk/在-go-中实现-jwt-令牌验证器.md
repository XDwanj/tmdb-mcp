# 在 Go 中实现 JWT 令牌验证器

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/examples/server/auth-middleware/README.md

提供了用于验证 JWT 令牌的 Go 函数签名。此函数接受上下文和令牌字符串，在成功时返回令牌信息，在令牌无效时返回错误。

```Go
func jwtVerifier(ctx context.Context, tokenString string) (*auth.TokenInfo, error) {
    // JWT token verification logic
    // On success: Return TokenInfo
    // On failure: Return auth.ErrInvalidToken
}
```

--------------------------------
