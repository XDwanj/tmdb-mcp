# 获取电影翻译 (Node.js)

来源: https://developer.themoviedb.org/reference/movie-translations

提供使用 Node.js 和 'node-fetch' 库获取电影翻译的示例。确保您已安装该库并配置了 API 密钥。

```Node.js
const fetch = require('node-fetch');

const options = {
  method: 'GET',
  headers: {
    accept: 'application/json'
  }
};

fetch('https://api.themoviedb.org/3/movie/{movie_id}/translations', options)
  .then(response => response.json())
  .then(response => console.log(response))
  .catch(err => console.error(err));
```

--------------------------------
