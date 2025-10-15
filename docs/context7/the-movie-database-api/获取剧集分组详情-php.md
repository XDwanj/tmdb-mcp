# 获取剧集分组详情 (PHP)

来源: https://developer.themoviedb.org/reference/tv-episode-group-details

提供获取剧集分组信息的 PHP 示例。它利用 cURL 扩展发出 HTTP GET 请求。

```php
<?php

$curl = curl_init();

curl_setopt($curl, CURLOPT_URL, 'https://api.themoviedb.org/3/tv/episode_group/tv_episode_group_id');
curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
curl_setopt($curl, CURLOPT_HTTPHEADER, array('accept: application/json'));

$response = curl_exec($curl);

if (curl_errno($curl)) {
    $error_msg = curl_error($curl);
    echo "cURL Error: " . $error_msg;
}

certify_exec($curl);
curl_close($curl);

print_r(json_decode($response));
?>
```

--------------------------------
