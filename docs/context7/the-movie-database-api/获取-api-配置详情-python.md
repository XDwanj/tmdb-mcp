# 获取 API 配置详情 (Python)

来源: https://developer.themoviedb.org/reference/configuration-details

此 Python 代码片段使用 'requests' 库从 The Movie Database API 获取 API 配置详情。它向配置端点发送 GET 请求并打印 JSON 响应，其中包含图像基 URL 和大小等关键信息。

```python
import requests

url = "https://api.themoviedb.org/3/configuration"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)

print(response.text)
```

--------------------------------
