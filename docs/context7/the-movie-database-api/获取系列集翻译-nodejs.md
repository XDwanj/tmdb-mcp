# 获取系列集翻译 (Node.js)

来源: https://developer.themoviedb.org/reference/collection-translations

提供使用 Node.js（可能使用 'axios' 或内置 'https' 模块等库）获取系列集翻译的示例。它展示了如何向 API 端点发出 GET 请求。

```javascript
const axios = require('axios');

axios
  .get('https://api.themoviedb.org/3/collection/collection_id/translations', {
    headers: {
      accept: 'application/json'
    }
  })
  .then(response => {
    console.log(response);
  })
  .catch(error => {
    console.error(error);
  });

```

--------------------------------
