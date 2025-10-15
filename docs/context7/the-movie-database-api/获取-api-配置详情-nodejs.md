# 获取 API 配置详情 (Node.js)

来源: https://developer.themoviedb.org/reference/configuration-details

此 Node.js 代码片段演示了如何使用 'node-fetch' 库检索 API 配置详情。它向配置端点发出 GET 请求并记录 JSON 响应，其中包含基 URL 和图像大小等信息。

```javascript
import fetch from 'node-fetch';

const options = {
  method: 'GET',
  headers: {
    accept: 'application/json'
  }
};

fetch('https://api.themoviedb.org/3/configuration', options)
  .then(response => response.json())
  .then(response => console.log(response))
  .catch(err => console.error(err));
```

--------------------------------
