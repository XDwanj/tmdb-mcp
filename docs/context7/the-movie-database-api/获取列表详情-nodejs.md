# 获取列表详情 (Node.js)

来源: https://developer.themoviedb.org/reference/list-details

提供使用 Node.js 和 'axios' 库获取列表详情的示例。它包括设置请求 URL 和标头。

```JavaScript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/list/list_id?language=en-US&page=1',
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
