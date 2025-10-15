# 获取 API 配置详情 (Ruby)

来源: https://developer.themoviedb.org/reference/configuration-details

此 Ruby 代码片段使用 'httparty' gem 获取 API 配置详情。它向配置端点发出 GET 请求并打印 JSON 响应，其中包含图像大小和基 URL 等必要详细信息，用于 API 集成。

```ruby
require 'httparty'

response = HTTParty.get('https://api.themoviedb.org/3/configuration', headers: { 'accept' => 'application/json' })

puts response.body
```

--------------------------------
