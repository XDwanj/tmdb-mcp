# 创建 Resty 客户端

Source: https://github.com/go-resty/docs/blob/main/content/docs/example/post-put-patch-request.md

初始化一个新的 Resty 客户端实例。建议延迟调用客户端的 Close 方法以正确释放资源。

```go
client := resty.New()
def client.Close()
```

--------------------------------
