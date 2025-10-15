# 获取标记的图片 (PHP)

来源: https://developer.themoviedb.org/reference/person-tagged-images

使用 The Movie Database API 获取人物标记图片的 PHP 脚本。此示例使用 cURL 发出 GET 请求并检索 JSON 响应。

```php
<?php

$curl = curl_init();

curl_setopt($curl, CURLOPT_URL, 'https://api.themoviedb.org/3/person/person_id/tagged_images?page=1');
curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
curl_setopt($curl, CURLOPT_HTTPHEADER, array(
  'accept: application/json'
));

$response = curl_exec($curl);

if (curl_errno($curl)) {
    echo 'Error:' . curl_error($curl);
}

curl_close($curl);

echo $response;

?>
```