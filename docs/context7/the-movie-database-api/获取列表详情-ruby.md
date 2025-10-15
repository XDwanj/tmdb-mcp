# 获取列表详情 (Ruby)

来源: https://developer.themoviedb.org/reference/list-details

使用 Ruby 和 'httparty' gem 获取列表详情的示例。它演示了如何使用 API 端点和标头配置 GET 请求。

```Ruby
require 'httparty'

url = 'https://api.themoviedb.org/3/list/list_id?language=en-US&page=1'
headers = { 'accept' => 'application/json' }

response = HTTParty.get(url, headers: headers)

puts response.body
```

--------------------------------
