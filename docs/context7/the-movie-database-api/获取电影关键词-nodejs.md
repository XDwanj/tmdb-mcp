# 获取电影关键词 (Node.js)

来源: https://developer.themoviedb.org/reference/movie-keywords

提供了使用 Node.js 检索电影关键词的示例。它使用 'axios' 库向指定的 API 端点发出 GET 请求。

```javascript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/movie/movie_id/keywords',
  headers: {
    accept: 'application/json'
  }
};

axios.request(options).then(function (response) {
  console.log(response.data);
}).catch(function (error) {
  console.error(error);
});
```

--------------------------------
