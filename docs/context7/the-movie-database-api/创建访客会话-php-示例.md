# 创建访客会话 - PHP 示例

来源: https://developer.themoviedb.org/reference/authentication-create-guest-session

一个 PHP 示例，演示如何使用 cURL 从 TMDb 获取访客会话。此脚本发送 GET 请求并输出 JSON 响应。

```php
<?php

$curl = curl_init();

curl_setopt($curl, CURLOPT_URL, 'https://api.themoviedb.org/3/authentication/guest_session/new');
curl_setopt($curl, CURLOPT_HTTPHEADER, [
  'accept: application/json'
]);

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
