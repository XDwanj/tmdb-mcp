# 按关键词获取电影 (Ruby)

来源: https://developer.themoviedb.org/reference/keyword-movies

此 Ruby 代码片段展示了如何使用 `httparty` gem 检索与关键词相关的电影。它使用适当的参数和标头构建 API 请求。确保已安装 `httparty` 并配置了 API 密钥。

```ruby
require 'httparty'

url = 'https://api.themoviedb.org/3/keyword/keyword_id/movies'

options = {
  query: {
    include_adult: 'false',
    language: 'en-US',
    page: '1'
  },
  headers: {
    'accept' => 'application/json'
  }
}

response = HTTParty.get(url, options)

puts response.body
```

--------------------------------
