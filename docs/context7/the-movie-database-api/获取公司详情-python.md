# 获取公司详情 (Python)

来源: https://developer.themoviedb.org/reference/company-details

使用 'requests' 库获取公司详情（按 ID）的示例 Python 请求。需要 API 密钥。以 JSON 格式返回公司信息。

```python
import requests

url = "https://api.themoviedb.org/3/company/company_id"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)

print(response.json())
```

--------------------------------
