# 获取电影演职员表 (Ruby)

来源: https://developer.themoviedb.org/reference/movie-credits

提供使用 `httparty` gem 发出 GET 请求以获取电影演职员表的 Ruby 示例。它展示了如何使用 URL、查询参数和标头配置请求。

```ruby
require 'httparty'

url = 'https://api.themoviedb.org/3/movie/movie_id/credits'
options = {
  query: {
    language: 'en-US'
  },
  headers: {
    'accept': 'application/json'
  }
}

response = HTTParty.get(url, options)

puts response.body
```

--------------------------------
