# 获取电影演职员表 (Node.js)

来源: https://developer.themoviedb.org/reference/movie-credits

提供使用 `axios` 库获取电影演职员表的 Node.js 示例。它展示了如何使用指定的标头向 TMDB API 发出 GET 请求。

```javascript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/movie/movie_id/credits',
  params: {language: 'en-US'},
  headers: {
    accept: 'application/json'
  }
};

axios
  .request(options)
  .then(function (response) {
    console.log(response.data);
  })
  .catch(function (error) {
    console.error(error);
  });
```

--------------------------------
