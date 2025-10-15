# 获取用户列表 (Ruby)

来源: https://developer.themoviedb.org/reference/account-lists

此 Ruby 示例使用 'httparty' gem 从 TMDb 获取用户的自定义列表。它发送带有必要参数和标头的 GET 请求，并返回包含列表详情的 JSON 响应。

```ruby
require 'httparty'

url = 'https://api.themoviedb.org/3/account/{account_id}/lists'

response = HTTParty.get(url, query: { page: '1' }, headers: { 'accept' => 'application/json' })

puts response.body
```

--------------------------------
