# 获取电影关键词 (PHP)

来源: https://developer.themoviedb.org/reference/movie-keywords

使用 PHP 和 cURL 获取电影关键词的说明。此示例演示了如何设置 cURL 请求，包括 URL 和标头。

```php
<?php

$curl = curl_init();

curl_setopt($curl, CURLOPT_URL, "https://api.themoviedb.org/3/movie/movie_id/keywords");
curl_setopt($curl, CURLOPT_HTTPHEADER, [
    "accept: application/json"
]);

$response = curl_exec($curl);

if (curl_errno($curl)) {
    echo 'Curl error: ' . curl_error($curl);
}

curl_close($curl);

echo $response;
?>
```

--------------------------------
