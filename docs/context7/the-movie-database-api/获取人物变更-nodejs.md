# 获取人物变更 (Node.js)

来源: https://developer.themoviedb.org/reference/person-changes

此 Node.js 示例展示了如何通过 TMDB API 发出 GET 请求以检索人物的近期变更。它使用 'axios' 库进行 HTTP 请求。

```javascript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/person/person_id/changes',
  params: {page: '1'},
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
