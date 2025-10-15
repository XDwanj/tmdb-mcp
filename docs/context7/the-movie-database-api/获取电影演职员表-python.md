# 获取电影演职员表 (Python)

来源: https://developer.themoviedb.org/reference/movie-credits

演示如何使用 Python 的 `requests` 库检索电影演职员表。该示例包括设置 TMDB API 调用所需的请求 URL、参数和标头。

```python
import requests

url = "https://api.themoviedb.org/3/movie/movie_id/credits"

params = {
    "language": "en-US"
}

headers = {
    "accept": "application/json"
}

response = requests.get(url, params=params, headers=headers)

print(response.text)
```

--------------------------------
