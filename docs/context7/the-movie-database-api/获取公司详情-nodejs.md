# 获取公司详情 (Node.js)

来源: https://developer.themoviedb.org/reference/company-details

使用 'axios' 库获取公司详情（按 ID）的示例 Node.js 请求。需要 API 密钥。以 JSON 格式返回公司信息。

```javascript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/company/company_id',
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
