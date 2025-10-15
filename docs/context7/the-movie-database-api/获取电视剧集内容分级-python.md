# 获取电视剧集内容分级 (Python)

来源: https://developer.themoviedb.org/reference/tv-series-content-ratings

展示了如何使用 Python 检索电视剧集的内容分级。此示例利用 'requests' 库发出 GET 请求。

```python
import requests

url = "https://api.themoviedb.org/3/tv/series_id/content_ratings"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)
print(response.text)
```

--------------------------------
