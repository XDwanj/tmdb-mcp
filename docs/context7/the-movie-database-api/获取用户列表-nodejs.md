# 获取用户列表 (Node.js)

来源: https://developer.themoviedb.org/reference/account-lists

此 Node.js 示例展示了如何使用 'axios' 库从 TMDb 检索用户的自定义列表。它发送带有指定参数和标头的 GET 请求，并以 JSON 格式返回列表数据。

```javascript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/account/{account_id}/lists',
  params: {page: '1'},
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
