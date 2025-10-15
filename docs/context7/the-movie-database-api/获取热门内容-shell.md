# 获取热门内容 (Shell)

来源: https://developer.themoviedb.org/reference/trending-all

用于获取特定时间窗口内热门电影、电视剧集和人物的示例 cURL 请求。需要 API 密钥进行身份验证。

```shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/trending/all/day?language=en-US' \
     --header 'accept: application/json'
```

--------------------------------
