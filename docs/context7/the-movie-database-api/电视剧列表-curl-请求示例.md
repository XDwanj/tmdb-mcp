# 电视剧列表 cURL 请求示例

来源: https://developer.themoviedb.org/reference/changes-tv-list

如何向电视剧列表端点发出 GET 请求的 cURL 示例。这演示了获取电视剧集更改所需的 URL 结构和标头。

```Shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/tv/changes?page=1' \
     --header 'accept: application/json'
```

--------------------------------
