# 获取系列集翻译 (Python)

来源: https://developer.themoviedb.org/reference/collection-translations

演示如何使用 Python 检索系列集翻译，通常使用 'requests' 库。该示例展示了向具有必要标头的指定 API 端点发出 GET 请求。

```python
import requests

url = "https://api.themoviedb.org/3/collection/collection_id/translations"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)

print(response.text)

```

--------------------------------
