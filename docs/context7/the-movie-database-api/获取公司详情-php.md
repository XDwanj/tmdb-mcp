# 获取公司详情 (PHP)

来源: https://developer.themoviedb.org/reference/company-details

使用 cURL 获取公司详情（按 ID）的示例 PHP 请求。需要 API 密钥。以 JSON 格式返回公司信息。

```php
<?php

$curl = curl_init();

curl_setopt($curl, CURLOPT_URL, 'https://api.themoviedb.org/3/company/company_id');
curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
curl_setopt($curl, CURLOPT_HTTPHEADER, array('accept: application/json'));

$response = curl_exec($curl);

if (curl_errno($curl)) {
    echo 'Error:' . curl_error($curl);
}

curl_close($curl);

echo $response;
?>
```

--------------------------------
