# 获取已评分电视剧集 (Node.js)

来源: https://developer.themoviedb.org/reference/account-rated-tv

使用 Node.js 和 TMDb API 获取用户已评分电视剧集的示例。演示了如何使用必要的标头和查询参数发出 GET 请求。

```javascript
const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/account/null/rated/tv',
  params: {
    language: 'en-US',
    page: '1',
    sort_by: 'created_at.asc'
  },
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
