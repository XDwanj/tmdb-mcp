# 创建访客会话 - Python 示例

来源: https://developer.themoviedb.org/reference/authentication-create-guest-session

此 Python 代码片段演示了如何使用 'requests' 库创建访客会话。它向 TMDb API 发送 GET 请求并打印响应。

```python
import requests

url = "https://api.themoviedb.org/3/authentication/guest_session/new"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)

print(response.text)

```

--------------------------------
