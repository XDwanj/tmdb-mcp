# 获取公司详情 (Ruby)

来源: https://developer.themoviedb.org/reference/company-details

使用 'httparty' gem 获取公司详情（按 ID）的示例 Ruby 请求。需要 API 密钥。以 JSON 格式返回公司信息。

```ruby
require 'httparty'

url = 'https://api.themoviedb.org/3/company/company_id'

headers = { 'accept' => 'application/json' }

response = HTTParty.get(url, headers: headers)

puts response.parsed_response
```

--------------------------------
