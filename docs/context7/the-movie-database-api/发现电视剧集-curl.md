# 发现电视剧集 (cURL)

来源: https://developer.themoviedb.org/reference/discover-tv

用于发现电视剧集的示例 cURL 请求。它演示了设置请求方法、带有用于过滤（成人内容、空首播日期、语言、页面和排序顺序）的查询参数的 URL，以及指定预期的响应格式（application/json）。

```Shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/discover/tv?include_adult=false&include_null_first_air_dates=false&language=en-US&page=1&sort_by=popularity.desc' \
     --header 'accept: application/json'
```

--------------------------------
