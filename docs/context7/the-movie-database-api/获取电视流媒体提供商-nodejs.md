# 获取电视流媒体提供商 (Node.js)

来源: https://developer.themoviedb.org/reference/watch-provider-tv-list

提供使用 Node.js 从 TMDB API 检索电视流媒体提供商列表的示例。它包括请求的必要标头。

```javascript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/watch/providers/tv',
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
