# 认证 API

来源: https://developer.themoviedb.org/reference/intro/getting-started

此端点允许您进行认证并检索请求令牌，这是许多其他 API 调用所必需的。

```APIDOC
# GET /authentication

# 描述
此端点用于与 TMDB API 进行认证。它返回一个可用于授予用户访问权限的请求令牌。

# 方法
GET

# 端点
https://api.themoviedb.org/3/authentication

# 参数
## 查询参数
- **api_key** (string) - 必需 - 您的 API 密钥。

# 请求示例
```
curl --request GET \
     --url 'https://api.themoviedb.org/3/authentication?api_key=YOUR_API_KEY' \
     --header 'accept: application/json'
```

# 响应
## 成功响应 (200)
- **success** (boolean) - 指示请求是否成功。
- **guest_session_id** (string) - 如果请求是针对访客会话，则返回访客会话 ID。
- **request_token** (string) - 用于认证的请求令牌。

## 响应示例
```json
{
  "success": true,
  "guest_session_id": "some_guest_session_id",
  "request_token": "some_request_token"
}
```
```

--------------------------------
