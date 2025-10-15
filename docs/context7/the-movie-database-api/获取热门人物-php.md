# 获取热门人物 (PHP)

来源: https://developer.themoviedb.org/reference/person-popular-list

使用 cURL 获取热门人物的 PHP 示例。此脚本演示了如何使用参数设置请求 URL 并发送 GET 请求。

```php
<?php

$curl = curl_init();

curl_setopt($curl, CURLOPT_URL, 'https://api.themoviedb.org/3/person/popular?language=en-US&page=1');
curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
curl_setopt($curl, CURLOPT_HTTPHEADER, array(
  'accept: application/json'
));

$response = curl_exec($curl);
$err = curl_error($curl);

curl_close($curl);

if ($err) {
  echo "cURL Error # " . $err;
} else {
  echo $response;
}

?>

```

--------------------------------
