# 获取电视剧集演职员表 (PHP)

来源: https://developer.themoviedb.org/reference/tv-series-credits

使用 cURL 检索电视剧集演职员表的示例 PHP 代码。此代码片段演示了如何设置 cURL 会话、指定请求 URL 和标头，以及执行请求以获取 API 响应。

```php
<?php

$curl = curl_init();

curl_setopt($curl, CURLOPT_URL, 'https://api.themoviedb.org/3/tv/series_id/credits?language=en-US');
curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
curl_setopt($curl, CURLOPT_CUSTOMREQUEST, 'GET');


$headers = [
  'accept: application/json'
];

curl_setopt($curl, CURLOPT_HTTPHEADER, $headers);

$response = curl_exec($curl);

if (curl_errno($curl)) {
    echo 'Error:' . curl_error($curl);
}

curl_close($curl);

echo $response;
```

--------------------------------
