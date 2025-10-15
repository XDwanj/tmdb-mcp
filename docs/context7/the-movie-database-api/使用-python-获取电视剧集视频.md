# 使用 Python 获取电视剧集视频

来源: https://developer.themoviedb.org/reference/tv-season-videos

提供使用 Python 获取 TMDB 电视剧集视频的示例。此代码使用 'requests' 库发出 GET 请求并打印 JSON 响应。

```python
import requests

url = "https://api.themoviedb.org/3/tv/series_id/season/season_number/videos?language=en-US"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)

print(response.text)
```

--------------------------------
