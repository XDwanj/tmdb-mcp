# 获取热门人物 (Node.js)

来源: https://developer.themoviedb.org/reference/person-popular-list

使用 Node.js 和 'axios' 库获取热门人物的示例。它向 API 端点发送 GET 请求，包括语言和页面参数。

```javascript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/person/popular',
  params: {
    language: 'en-US',
    page: '1'
  },
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
