# 获取人物变更 (PHP)

来源: https://developer.themoviedb.org/reference/person-changes

此 PHP 代码提供了如何从 TMDB API 检索人物近期变更的示例。它使用 cURL 发出 HTTP GET 请求并显示 JSON 响应。

```php
<?php

$curl = curl_init();

curl_setopt($curl, CURLOPT_URL, "https://api.themoviedb.org/3/person/person_id/changes?page=1");
curl_setopt($curl, CURLOPT_HTTPHEADER, [
    "accept: application/json"
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
