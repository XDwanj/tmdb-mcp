# 获取剧集视频 (cURL)

来源: https://developer.themoviedb.org/reference/tv-episode-videos

获取特定剧集视频的示例 cURL 请求。需要系列 ID、季号和集号。支持语言过滤。

```shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/tv/series_id/season/season_number/episode/episode_number/videos?language=en-US' \
     --header 'accept: application/json'
```

--------------------------------
