# 获取电视剧集演职员表 (Python)

来源: https://developer.themoviedb.org/reference/tv-series-credits

使用 'requests' 库获取电视剧集演职员表的示例 Python 代码。它展示了如何构建 API URL、添加必要的标头以及处理 JSON 响应。

```python
import requests

url = "https://api.themoviedb.org/3/tv/series_id/credits?language=en-US"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)

print(response.text)
```

--------------------------------
