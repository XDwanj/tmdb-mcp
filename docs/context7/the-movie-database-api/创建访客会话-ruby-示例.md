# 创建访客会话 - Ruby 示例

来源: https://developer.themoviedb.org/reference/authentication-create-guest-session

此 Ruby 代码片段展示了如何使用 'httparty' gem 与 TMDb API 交互来创建访客会话。它向指定的端点发送 GET 请求。

```ruby
require 'httparty'

response = HTTParty.get('https://api.themoviedb.org/3/authentication/guest_session/new',
  headers: { 'accept' => 'application/json' })

puts response.body

```

--------------------------------
