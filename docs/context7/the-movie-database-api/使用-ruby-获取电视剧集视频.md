# 使用 Ruby 获取电视剧集视频

来源: https://developer.themoviedb.org/reference/tv-season-videos

演示如何使用 Ruby 获取电视剧集视频。此示例使用内置的 'net/http' 库发出到 TMDB API 的 GET 请求，并解析 JSON 响应。

```ruby
require 'uri'
require 'net/http'

uri = URI('https://api.themoviedb.org/3/tv/series_id/season/season_number/videos?language=en-US')

Net::HTTP.start(uri.hostname, uri.port, :use_ssl => uri.scheme == 'https') do |http|
  request = Net::HTTP::Get.new(uri)
  request['accept'] = 'application/json'

  response = http.request(request)
  puts response.body
end
```

--------------------------------
