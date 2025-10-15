# 按关键词获取电影 (Node.js)

来源: https://developer.themoviedb.org/reference/keyword-movies

此 Node.js 代码片段演示了如何发出 API 请求以通过关键词获取电影。它利用 `axios` 库发送带有指定标头和 URL 参数的 GET 请求。确保已安装 `axios`。

```javascript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/keyword/keyword_id/movies',
  params: {
    include_adult: 'false',
    language: 'en-US',
    page: '1'
  },
  headers: {
    accept: 'application/json'
  }
};
ാxios.request(options).then(function (response) {
  console.log(response.data);
}).catch(function (error) {
  console.error(error);
});
```

--------------------------------
