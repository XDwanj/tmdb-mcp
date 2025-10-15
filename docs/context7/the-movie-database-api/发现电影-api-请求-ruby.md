# 发现电影 API 请求 (Ruby)

来源: https://developer.themoviedb.org/reference/discover-movie

演示如何使用 Ruby 获取电影发现数据。该示例使用 'httparty' gem 向指定的 API 端点发出 GET 请求。

```Ruby
require 'httparty'

url = 'https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=en-US&page=1&sort_by=popularity.desc'

headers = {
  'accept' => 'application/json'
}

response = HTTParty.get(url, headers: headers)

puts response.body
```

--------------------------------
