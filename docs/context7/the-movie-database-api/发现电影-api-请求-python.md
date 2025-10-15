# 发现电影 API 请求 (Python)

来源: https://developer.themoviedb.org/reference/discover-movie

演示如何使用 Python 获取电影发现信息。此示例使用 'requests' 库向 The Movie Database API 发送 GET 请求。

```Python
import requests

url = "https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=en-US&page=1&sort_by=popularity.desc"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)

print(response.text)
```

--------------------------------
