# 获取列表详情 (Python)

来源: https://developer.themoviedb.org/reference/list-details

演示如何使用 Python 的 requests 库获取列表详情。它展示了如何使用 API 端点和标头设置请求。

```Python
import requests

url = "https://api.themoviedb.org/3/list/list_id?language=en-US&page=1"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)
print(response.text)
```

--------------------------------
