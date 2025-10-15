# 发现电影 API 请求 (PHP)

来源: https://developer.themoviedb.org/reference/discover-movie

提供如何使用 PHP 检索电影发现数据的示例。它演示了向 API 端点发出带有常用发现参数的 GET 请求。

```PHP
<?php

$curl = curl_init();

curl_setopt($curl, CURLOPT_URL, 'https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=en-US&page=1&sort_by=popularity.desc');
curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
curl_setopt($curl, CURLOPT_HTTPHEADER, array(
  'accept: application/json'
));

$response = curl_exec($curl);

if (curl_errno($curl)) {
    $error_msg = curl_error($curl);
    print_r($error_msg);
}

curl_close($curl);

$data = json_decode($response, true);
print_r($data);

?>
```

--------------------------------
