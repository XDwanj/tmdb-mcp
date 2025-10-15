# 获取评论详情 (cURL)

来源: https://developer.themoviedb.org/reference/review-details

此代码片段演示了如何使用 cURL 请求检索电影或电视剧集的评论详情。它包括 GET 方法、带有评论 ID 占位符的 API 端点以及必要的“Accept”标头。

```shell
curl --request GET \
     --url https://api.themoviedb.org/3/review/review_id \
     --header 'accept: application/json'
```

--------------------------------
