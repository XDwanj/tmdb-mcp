# 创建访客会话 - Node.js 示例

来源: https://developer.themoviedb.org/reference/authentication-create-guest-session

使用 Node.js 创建访客会话的示例。此代码片段使用 'axios' 库向 TMDb API 端点发送 GET 请求以创建访客会话。

```javascript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/authentication/guest_session/new',
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
