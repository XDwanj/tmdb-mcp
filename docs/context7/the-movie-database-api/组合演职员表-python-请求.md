# 组合演职员表 Python 请求

来源: https://developer.themoviedb.org/reference/person-combined-credits

使用 'requests' 库调用 TMDB API 获取组合演职员表的 Python 示例。它展示了如何设置 GET 请求的 URL、参数和标头。

```Python
import requests

url = "https://api.themoviedb.org/3/person/person_id/combined_credits"

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
