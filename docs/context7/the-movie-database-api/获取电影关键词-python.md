# 获取电影关键词 (Python)

来源: https://developer.themoviedb.org/reference/movie-keywords

展示如何使用 Python 获取电影关键词。它使用 'requests' 库发送 GET 请求，并包含必要的标头。

```python
import requests

url = "https://api.themoviedb.org/3/movie/movie_id/keywords"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)

print(response.text)
```

--------------------------------
