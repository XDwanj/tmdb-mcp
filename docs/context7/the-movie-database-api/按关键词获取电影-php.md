# 按关键词获取电影 (PHP)

来源: https://developer.themoviedb.org/reference/keyword-movies

此 PHP 代码片段演示了如何使用 Guzzle HTTP 客户端通过关键词获取电影。它使用正确的 URL、查询参数和标头设置请求。确保通过 Composer 安装了 Guzzle。

```php
<?php

require 'vendor/autoload.php';

use GuzzleHttp\Client;
use GuzzleHttp\RequestOptions;

$client = new Client();

try {
    $response = $client->request('GET', 'https://api.themoviedb.org/3/keyword/keyword_id/movies', [
        RequestOptions::QUERY => [
            'include_adult' => 'false',
            'language' => 'en-US',
            'page' => '1',
        ],
        'headers' => [
            'accept' => 'application/json',
        ],
    ]);

    echo $response->getBody();

} catch (Exception $e) {
    echo 'Error: ' . $e->getMessage();
}
?>

```

--------------------------------
