# 获取电影关键词 (Ruby)

来源: https://developer.themoviedb.org/reference/movie-keywords

演示使用 Ruby 获取电影关键词。此示例使用 'httparty' gem 向 API 端点发出 GET 请求。

```ruby
require 'httparty'

response = HTTParty.get('https://api.themoviedb.org/3/movie/movie_id/keywords',
  headers: { 'accept' => 'application/json' })

puts response.body
```

--------------------------------
