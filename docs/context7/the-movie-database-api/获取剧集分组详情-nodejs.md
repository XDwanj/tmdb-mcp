# 获取剧集分组详情 (Node.js)

来源: https://developer.themoviedb.org/reference/tv-episode-group-details

使用 Node.js 获取剧集分组详情。此示例演示了如何发出 GET 请求并处理 JSON 响应。

```javascript
const options = {
  method: 'GET',
  headers: {
    accept: 'application/json'
  }
};

fetch('https://api.themoviedb.org/3/tv/episode_group/tv_episode_group_id', options)
  .then(response => response.json())
  .then(response => console.log(response))
  .catch(err => console.error(err));
```

--------------------------------
