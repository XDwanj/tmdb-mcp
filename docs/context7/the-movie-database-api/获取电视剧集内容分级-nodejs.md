# 获取电视剧集内容分级 (Node.js)

来源: https://developer.themoviedb.org/reference/tv-series-content-ratings

提供获取电视剧集内容分级的 Node.js 示例。它使用 'axios' 库执行 HTTP GET 请求。

```javascript
const axios = require('axios');

const url = 'https://api.themoviedb.org/3/tv/series_id/content_ratings';

const config = {
  headers: {
    'accept': 'application/json'
  }
};

axios.get(url, config)
  .then(response => {
    console.log(response.data);
  })
  .catch(error => {
    console.error(error);
  });
```

--------------------------------
