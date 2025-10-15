# 获取列表详情 (PHP)

来源: https://developer.themoviedb.org/reference/list-details

使用 PHP 和 cURL 检索列表详情的说明。此示例显示了进行 GET 请求和包含必要标头的设置。

```PHP
<?php

$curl = curl_init();

curl_setopt($curl, CURLOPT_URL, 'https://api.themoviedb.org/3/list/list_id?language=en-US&page=1');
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
