# 生成 Curl 命令示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/curl-command.md

演示如何为 Resty 请求启用 curl 命令生成并检索生成的命令字符串。这涉及在请求上设置 `GenerateCurlCmd` 选项。

```go
c := resty.New()
deffer c.Close()

res, _ := c.R().
    SetGenerateCurlCmd(true).
    SetBody(map[string]string{
        "name": "Resty",
    }).
    Post("https://httpbin.org/post")

curlCmdStr := res.Request.CurlCmd()
fmt.Println(curlCmdStr)

// Result:
// curl -X POST -H 'Accept-Encoding: gzip, deflate' -H 'Content-Type: application/json' -H 'User-Agent: go-resty/3.0.0 (https://resty.dev)' -d '{"name":"Resty"}' https://httpbin.org/post
```

--------------------------------
