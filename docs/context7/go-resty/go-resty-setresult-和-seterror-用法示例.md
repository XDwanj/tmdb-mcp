# Go Resty SetResult 和 SetError 用法示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/response-auto-parse.md

演示使用 SetResult 的不同方法，包括传递内联指针、非指针值和预先声明的指针。SetError 的原理相同。

```go
// 设置
client.R().SetResult(&LoginResponse{})

// 访问
fmt.Println(res.Result().(*LoginResponse))
```

```go
// 设置
client.R().SetResult(LoginResponse{})

// 访问
fmt.Println(res.Result().(*LoginResponse))
```

```go
loginResponse := &LoginResponse{}
// 设置
client.R().SetResult(loginResponse)

// 访问
fmt.Println(loginResponse)
```

--------------------------------
