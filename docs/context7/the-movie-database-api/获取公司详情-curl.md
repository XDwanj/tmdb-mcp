# 获取公司详情 (cURL)

来源: https://developer.themoviedb.org/reference/company-details

获取公司详情（按 ID）的示例 cURL 请求。需要标头中的 API 密钥才能实现完整功能。以 JSON 格式返回公司信息。

```shell
curl --request GET \
     --url https://api.themoviedb.org/3/company/company_id \
     --header 'accept: application/json'
```

--------------------------------
