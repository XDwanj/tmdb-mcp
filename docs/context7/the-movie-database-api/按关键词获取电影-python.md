# 按关键词获取电影 (Python)

来源: https://developer.themoviedb.org/reference/keyword-movies

此 Python 代码片段展示了如何使用 `requests` 库获取与关键词相关的电影。它使用 API 端点、查询参数和必要的标头配置 GET 请求。请确保已安装 `requests` 库。

```python
import requests

url = "https://api.themoviedb.org/3/keyword/keyword_id/movies"

params = {
    "include_adult": "false",
    "language": "en-US",
    "page": "1"
}

headers = {
    "accept": "application/json"
}

response = requests.get(url, params=params, headers=headers)

print(response.text)
```

--------------------------------
