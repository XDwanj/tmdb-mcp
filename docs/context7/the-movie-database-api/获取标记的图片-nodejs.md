# 获取标记的图片 (Node.js)

来源: https://developer.themoviedb.org/reference/person-tagged-images

演示如何使用 Node.js 获取人物的标记图片。这涉及到向 TMDb API 发出 HTTP GET 请求。

```javascript
const options = {
  method: 'GET',
  headers: {
    accept: 'application/json'
  }
};

fetch('https://api.themoviedb.org/3/person/person_id/tagged_images?page=1', options)
  .then(response => response.json())
  .then(response => console.log(response))
  .catch(err => console.error(err));
```

--------------------------------
