# 获取电视剧集演职员表 (Node.js)

来源: https://developer.themoviedb.org/reference/tv-series-credits

使用 'node-fetch' 获取电视剧集演职员表的 Node.js 示例。它演示了向 API 发出 GET 请求、设置 'accept' 标头以及处理 JSON 响应。

```javascript
import fetch from 'node-fetch';

const options = {
  method: 'GET',
  headers: {
    accept: 'application/json'
  }
};

fetch('https://api.themoviedb.org/3/tv/series_id/credits?language=en-US', options)
  .then(response => response.json())
  .then(response => console.log(response))
  .catch(err => console.error(err));
```

--------------------------------
