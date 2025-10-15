# 获取电视剧集演职员表 (Ruby)

来源: https://developer.themoviedb.org/reference/tv-series-credits

使用 'httparty' gem 发出 GET 请求以获取电视剧集演职员表的示例 Ruby 代码。它展示了如何定义 API 端点、设置必需的标头以及处理响应。

```ruby
require 'httparty'

response = HTTParty.get('https://api.themoviedb.org/3/tv/series_id/credits?language=en-US',
  headers: { 'accept' => 'application/json' })

puts response.body
```

--------------------------------
