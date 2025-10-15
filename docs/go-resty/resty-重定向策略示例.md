# Resty 重定向策略示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/redirect-policy.md

演示如何在 Resty 中设置内置的重定向策略。这包括设置具有最大次数的灵活重定向策略，以及组合多个策略，如灵活策略和域名检查策略。

```go
client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(5))

client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(5),
    resty.DomainCheckRedirectPolicy("host1.com", "host2.org", "host3.net"))
```

--------------------------------
