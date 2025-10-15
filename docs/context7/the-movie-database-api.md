================
代码片段
================
## 获取列表详情 (Ruby)

来源: https://developer.themoviedb.org/reference/list-details

使用 Ruby 和 'httparty' gem 获取列表详情的示例。它演示了如何使用 API 端点和标头配置 GET 请求。

```Ruby
require 'httparty'

url = 'https://api.themoviedb.org/3/list/list_id?language=en-US&page=1'
headers = { 'accept' => 'application/json' }

response = HTTParty.get(url, headers: headers)

puts response.body
```

--------------------------------

## 获取列表详情 (PHP)

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

## 创建访客会话 - PHP 示例

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

## 认证 API

来源: https://developer.themoviedb.org/reference/index

此端点允许您通过提供 API 密钥来认证您的 API 请求。

```APIDOC
## GET /authentication

## 描述
使用您的 API 密钥认证您的 API 请求。

## 方法
GET

## 端点
https://api.themoviedb.org/3/authentication

## 参数
### 标头参数
- **accept** (string) - 必需 - 指定所需的响应格式（例如，`application/json`）。

## 请求示例
```
curl --request GET \
     --url https://api.themoviedb.org/3/authentication \
     --header 'accept: application/json'
```

## 响应
### 成功响应 (200)
- **key** (string) - 认证的 API 密钥。
- **success** (boolean) - 指示认证是否成功。

### 响应示例
```json
{
  "key": "YOUR_API_KEY",
  "success": true
}
```
```

--------------------------------

## 获取热门人物 (PHP)

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

## 获取电影关键词 (Ruby)

来源: https://developer.themoviedb.org/reference/movie-keywords

演示使用 Ruby 获取电影关键词。此示例使用 'httparty' gem 向 API 端点发出 GET 请求。

```ruby
require 'httparty'

response = HTTParty.get('https://api.themoviedb.org/3/movie/movie_id/keywords',
  headers: { 'accept' => 'application/json' })

puts response.body
```

--------------------------------

## 获取电视剧集演职员表 (Node.js)

来源: https://developer.themoviedb.org/reference/tv-series-credits

使用 'node-fetch' 获取电视剧集演职员表的 Node.js 示例。它演示了向 API 发出 GET 请求、设置 'accept' 标头以及处理 JSON 响应。

```javascript
import fetch from 'node-fetch';

const options = {
  method: 'GET',
  headers: {
    accept: 'application/json'
  }
};

fetch('https://api.themoviedb.org/3/tv/series_id/credits?language=en-US', options)
  .then(response => response.json())
  .then(response => console.log(response))
  .catch(err => console.error(err));
```

--------------------------------

## 认证 API

来源: https://developer.themoviedb.org/reference/intro/getting-started

此端点允许您进行认证并检索请求令牌，这是许多其他 API 调用所必需的。

```APIDOC
## GET /authentication

## 描述
此端点用于与 TMDB API 进行认证。它返回一个可用于授予用户访问权限的请求令牌。

## 方法
GET

## 端点
https://api.themoviedb.org/3/authentication

## 参数
### 查询参数
- **api_key** (string) - 必需 - 您的 API 密钥。

## 请求示例
```
curl --request GET \
     --url 'https://api.themoviedb.org/3/authentication?api_key=YOUR_API_KEY' \
     --header 'accept: application/json'
```

## 响应
### 成功响应 (200)
- **success** (boolean) - 指示请求是否成功。
- **guest_session_id** (string) - 如果请求是针对访客会话，则返回访客会话 ID。
- **request_token** (string) - 用于认证的请求令牌。

### 响应示例
```json
{
  "success": true,
  "guest_session_id": "some_guest_session_id",
  "request_token": "some_request_token"
}
```
```

--------------------------------

## 创建访客会话 - Node.js 示例

来源: https://developer.themoviedb.org/reference/authentication-create-guest-session

使用 Node.js 创建访客会话的示例。此代码片段使用 'axios' 库向 TMDb API 端点发送 GET 请求以创建访客会话。

```javascript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/authentication/guest_session/new',
  headers: {
    accept: 'application/json'
  }
};

axios
  .request(options)
  .then(function (response) {
    console.log(response.data);
  })
  .catch(function (error) {
    console.error(error);
  });

```

--------------------------------

## 获取电影关键词 (PHP)

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

## 按关键词获取电影 (Python)

来源: https://developer.themoviedb.org/reference/keyword-movies

此 Python 代码片段展示了如何使用 `requests` 库获取与关键词相关的电影。它使用 API 端点、查询参数和必要的标头配置 GET 请求。请确保已安装 `requests` 库。

```python
import requests

url = "https://api.themoviedb.org/3/keyword/keyword_id/movies"

params = {
    "include_adult": "false",
    "language": "en-US",
    "page": "1"
}

headers = {
    "accept": "application/json"
}

response = requests.get(url, params=params, headers=headers)

print(response.text)
```

--------------------------------

## 获取电影关键词 (Node.js)

来源: https://developer.themoviedb.org/reference/movie-keywords

提供了使用 Node.js 检索电影关键词的示例。它使用 'axios' 库向指定的 API 端点发出 GET 请求。

```javascript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/movie/movie_id/keywords',
  headers: {
    accept: 'application/json'
  }
};

axios.request(options).then(function (response) {
  console.log(response.data);
}).catch(function (error) {
  console.error(error);
});
```

--------------------------------

## 获取列表详情 (Node.js)

来源: https://developer.themoviedb.org/reference/list-details

提供使用 Node.js 和 'axios' 库获取列表详情的示例。它包括设置请求 URL 和标头。

```JavaScript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/list/list_id?language=en-US&page=1',
  headers: {
    accept: 'application/json'
  }
};

axios
  .request(options)
  .then(function (response) {
    console.log(response.data);
  })
  .catch(function (error) {
    console.error(error);
  });
```

--------------------------------

## 获取标记的图片 (Node.js)

来源: https://developer.themoviedb.org/reference/person-tagged-images

演示如何使用 Node.js 获取人物的标记图片。这涉及到向 TMDb API 发出 HTTP GET 请求。

```javascript
const options = {
  method: 'GET',
  headers: {
    accept: 'application/json'
  }
};

fetch('https://api.themoviedb.org/3/person/person_id/tagged_images?page=1', options)
  .then(response => response.json())
  .then(response => console.log(response))
  .catch(err => console.error(err));
```

--------------------------------

## 获取列表详情 (cURL)

来源: https://developer.themoviedb.org/reference/list-details

如何使用 cURL 检索特定列表详情的示例。它指定了 GET 请求、带有列表 ID 占位符的 API 端点以及必需的标头。

```Shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/list/list_id?language=en-US&page=1' \
     --header 'accept: application/json'
```

--------------------------------

## 创建访客会话 - Python 示例

来源: https://developer.themoviedb.org/reference/authentication-create-guest-session

此 Python 代码片段演示了如何使用 'requests' 库创建访客会话。它向 TMDb API 发送 GET 请求并打印响应。

```python
import requests

url = "https://api.themoviedb.org/3/authentication/guest_session/new"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)

print(response.text)

```

--------------------------------

## 获取电影翻译 (Node.js)

来源: https://developer.themoviedb.org/reference/movie-translations

提供使用 Node.js 和 'node-fetch' 库获取电影翻译的示例。确保您已安装该库并配置了 API 密钥。

```Node.js
const fetch = require('node-fetch');

const options = {
  method: 'GET',
  headers: {
    accept: 'application/json'
  }
};

fetch('https://api.themoviedb.org/3/movie/{movie_id}/translations', options)
  .then(response => response.json())
  .then(response => console.log(response))
  .catch(err => console.error(err));
```

--------------------------------

## 按关键词获取电影 (PHP)

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

## 获取收藏电影 (PHP)

来源: https://developer.themoviedb.org/reference/account-get-favorites

展示了如何使用 PHP 获取用户收藏电影的示例。这涉及到向指定的 API 端点发出 HTTP GET 请求。

```PHP
// PHP example would go here, making a GET request to the API
```

--------------------------------

## 按关键词获取电影 (Node.js)

来源: https://developer.themoviedb.org/reference/keyword-movies

此 Node.js 代码片段演示了如何发出 API 请求以通过关键词获取电影。它利用 `axios` 库发送带有指定标头和 URL 参数的 GET 请求。确保已安装 `axios`。

```javascript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/keyword/keyword_id/movies',
  params: {
    include_adult: 'false',
    language: 'en-US',
    page: '1'
  },
  headers: {
    accept: 'application/json'
  }
};
ാxios.request(options).then(function (response) {
  console.log(response.data);
}).catch(function (error) {
  console.error(error);
});
```

--------------------------------

## 获取剧集分组详情 (Node.js)

来源: https://developer.themoviedb.org/reference/tv-episode-group-details

使用 Node.js 获取剧集分组详情。此示例演示了如何发出 GET 请求并处理 JSON 响应。

```javascript
const options = {
  method: 'GET',
  headers: {
    accept: 'application/json'
  }
};

fetch('https://api.themoviedb.org/3/tv/episode_group/tv_episode_group_id', options)
  .then(response => response.json())
  .then(response => console.log(response))
  .catch(err => console.error(err));
```

--------------------------------

## 获取电影演职员表 (Node.js)

来源: https://developer.themoviedb.org/reference/movie-credits

提供使用 `axios` 库获取电影演职员表的 Node.js 示例。它展示了如何使用指定的标头向 TMDB API 发出 GET 请求。

```javascript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/movie/movie_id/credits',
  params: {language: 'en-US'},
  headers: {
    accept: 'application/json'
  }
};

axios
  .request(options)
  .then(function (response) {
    console.log(response.data);
  })
  .catch(function (error) {
    console.error(error);
  });
```

--------------------------------

## cURL 请求示例

来源: https://developer.themoviedb.org/reference/authentication-create-session

如何使用 cURL 向创建会话端点发出 POST 请求的示例。它包括 URL、必需的标头以及指定的内容类型。

```shell
curl --request POST \
     --url https://api.themoviedb.org/3/authentication/session/new \
     --header 'accept: application/json' \
     --header 'content-type: application/json'
```

--------------------------------

## 获取电视剧集演职员表 (Ruby)

来源: https://developer.themoviedb.org/reference/tv-series-credits

使用 'httparty' gem 发出 GET 请求以获取电视剧集演职员表的示例 Ruby 代码。它展示了如何定义 API 端点、设置必需的标头以及处理响应。

```ruby
require 'httparty'

response = HTTParty.get('https://api.themoviedb.org/3/tv/series_id/credits?language=en-US',
  headers: { 'accept' => 'application/json' })

puts response.body
```

--------------------------------

## 获取电影演职员表 (Ruby)

来源: https://developer.themoviedb.org/reference/movie-credits

提供使用 `httparty` gem 发出 GET 请求以获取电影演职员表的 Ruby 示例。它展示了如何使用 URL、查询参数和标头配置请求。

```ruby
require 'httparty'

url = 'https://api.themoviedb.org/3/movie/movie_id/credits'
options = {
  query: {
    language: 'en-US'
  },
  headers: {
    'accept': 'application/json'
  }
}

response = HTTParty.get(url, options)

puts response.body
```

--------------------------------

## 电视剧列表 cURL 请求示例

来源: https://developer.themoviedb.org/reference/changes-tv-list

如何向电视剧列表端点发出 GET 请求的 cURL 示例。这演示了获取电视剧集更改所需的 URL 结构和标头。

```Shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/tv/changes?page=1' \
     --header 'accept: application/json'
```

--------------------------------

## 获取热门人物 (Node.js)

来源: https://developer.themoviedb.org/reference/person-popular-list

使用 Node.js 和 'axios' 库获取热门人物的示例。它向 API 端点发送 GET 请求，包括语言和页面参数。

```javascript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/person/popular',
  params: {
    language: 'en-US',
    page: '1'
  },
  headers: {
    accept: 'application/json'
  }
};

axios
  .request(options)
  .then(function (response) {
    console.log(response.data);
  })
  .catch(function (error) {
    console.error(error);
  });

```

--------------------------------

## 获取电视剧集内容分级 (Node.js)

来源: https://developer.themoviedb.org/reference/tv-series-content-ratings

提供获取电视剧集内容分级的 Node.js 示例。它使用 'axios' 库执行 HTTP GET 请求。

```javascript
const axios = require('axios');

const url = 'https://api.themoviedb.org/3/tv/series_id/content_ratings';

const config = {
  headers: {
    'accept': 'application/json'
  }
};

axios.get(url, config)
  .then(response => {
    console.log(response.data);
  })
  .catch(error => {
    console.error(error);
  });
```

--------------------------------

## 获取电视剧集演职员表 (PHP)

来源: https://developer.themoviedb.org/reference/tv-series-credits

使用 cURL 检索电视剧集演职员表的示例 PHP 代码。此代码片段演示了如何设置 cURL 会话、指定请求 URL 和标头，以及执行请求以获取 API 响应。

```php
<?php

$curl = curl_init();

curl_setopt($curl, CURLOPT_URL, 'https://api.themoviedb.org/3/tv/series_id/credits?language=en-US');
curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
curl_setopt($curl, CURLOPT_CUSTOMREQUEST, 'GET');


$headers = [
  'accept: application/json'
];

curl_setopt($curl, CURLOPT_HTTPHEADER, $headers);

$response = curl_exec($curl);

if (curl_errno($curl)) {
    echo 'Error:' . curl_error($curl);
}

curl_close($curl);

echo $response;
```

--------------------------------

## 电视剧集关键词 cURL 请求

来源: https://developer.themoviedb.org/reference/tv-series-keywords

此代码片段演示了如何向 TMDB API 发出 GET 请求以获取电视剧集的关键词。它包括必要的 URL 和标头。通常需要通过 API 密钥进行身份验证，但在本通用示例中省略了。

```bash
curl --request GET \
     --url https://api.themoviedb.org/3/tv/series_id/keywords \
     --header 'accept: application/json'
```

--------------------------------

## 获取公司详情 (Python)

来源: https://developer.themoviedb.org/reference/company-details

使用 'requests' 库获取公司详情（按 ID）的示例 Python 请求。需要 API 密钥。以 JSON 格式返回公司信息。

```python
import requests

url = "https://api.themoviedb.org/3/company/company_id"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)

print(response.json())
```

--------------------------------

## 获取电视剧集内容分级 (Python)

来源: https://developer.themoviedb.org/reference/tv-series-content-ratings

展示了如何使用 Python 检索电视剧集的内容分级。此示例利用 'requests' 库发出 GET 请求。

```python
import requests

url = "https://api.themoviedb.org/3/tv/series_id/content_ratings"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)
print(response.text)
```

--------------------------------

## 获取已评分电影 (cURL)

来源: https://developer.themoviedb.org/reference/guest-session-rated-movies

如何使用 cURL 获取访客会话已评分电影的示例。它演示了 GET 请求、用于语言、页面和排序的 URL 参数以及预期的 Accept 标头。

```Shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/guest_session/guest_session_id/rated/movies?language=en-US&page=1&sort_by=created_at.asc' \
     --header 'accept: application/json'
```

--------------------------------

## 获取公司详情 (Ruby)

来源: https://developer.themoviedb.org/reference/company-details

使用 'httparty' gem 获取公司详情（按 ID）的示例 Ruby 请求。需要 API 密钥。以 JSON 格式返回公司信息。

```ruby
require 'httparty'

url = 'https://api.themoviedb.org/3/company/company_id'

headers = { 'accept' => 'application/json' }

response = HTTParty.get(url, headers: headers)

puts response.parsed_response
```

--------------------------------

## GET /configuration

来源: https://developer.themoviedb.org/reference/configuration-details

获取 API 配置详情，其中包括图像的基 URL 和支持的图像大小。

```APIDOC
## GET /configuration

## 描述
检索 API 配置详情，包括图像基 URL 和支持的图像大小。

## 方法
GET

## 端点
https://api.themoviedb.org/3/configuration

## 参数
### 查询参数
无

### 请求正文
无

## 响应
### 成功响应 (200)
- **images** (object) - 包含图像配置详情。
  - **base_url** (string) - 图像的基 URL。
  - **secure_base_url** (string) - 图像的安全基 URL。
  - **backdrop_sizes** (string 数组) - 支持的背景图像大小。
  - **logo_sizes** (string 数组) - 支持的 Logo 图像大小。
  - **poster_sizes** (string 数组) - 支持的海报图像大小。
  - **profile_sizes** (string 数组) - 支持的个人资料图像大小。
  - **still_sizes** (string 数组) - 支持的剧照图像大小。
- **change_keys** (string 数组) - 触发 API 更改的键。

### 响应示例
```json
{
  "images": {
    "base_url": "http://image.tmdb.org/t/p/",
    "secure_base_url": "https://image.tmdb.org/t/p/",
    "backdrop_sizes": [
      "w300",
      "w780",
      "w1280",
      "original"
    ],
    "logo_sizes": [
      "w45",
      "w92",
      "w154",
      "w185",
      "w300",
      "w500"
    ],
    "poster_sizes": [
      "w92",
      "w154",
      "w185",
      "w342",
      "w500",
      "w780",
      "original"
    ],
    "profile_sizes": [
      "w45",
      "w185",
      "h632"
    ],
    "still_sizes": [
      "w92",
      "w185",
      "w300",
      "original"
    ]
  },
  "change_keys": [
    "adult",
    "air_date",
    "also_known_as",
    "backdrop_path",
    "biography",
    "birthday",
    "budget",
    "character",
    "created_by",
    "crew",
    ""}
```
```

--------------------------------

## 发现电影 API 请求 (Python)

来源: https://developer.themoviedb.org/reference/discover-movie

演示如何使用 Python 获取电影发现信息。此示例使用 'requests' 库向 The Movie Database API 发送 GET 请求。

```Python
import requests

url = "https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=en-US&page=1&sort_by=popularity.desc"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)

print(response.text)
```

--------------------------------

## 获取公司详情 (PHP)

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

## 发现电影 API 请求 (PHP)

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

## 获取公司详情 (cURL)

来源: https://developer.themoviedb.org/reference/company-details

获取公司详情（按 ID）的示例 cURL 请求。需要标头中的 API 密钥才能实现完整功能。以 JSON 格式返回公司信息。

```shell
curl --request GET \
     --url https://api.themoviedb.org/3/company/company_id \
     --header 'accept: application/json'
```

--------------------------------

## 获取剧集分组详情 (PHP)

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

## 获取电视剧集推荐 (cURL)

来源: https://developer.themoviedb.org/reference/tv-series-recommendations

演示如何使用 cURL 向 /tv/{series_id}/recommendations 端点发出 GET 请求。它包括设置 API 端点、语言、页码和 accept 标头。

```shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/tv/series_id/recommendations?language=en-US&page=1' \
     --header 'accept: application/json'
```

--------------------------------

## 获取电视剧集演职员表 (Python)

来源: https://developer.themoviedb.org/reference/tv-series-credits

使用 'requests' 库获取电视剧集演职员表的示例 Python 代码。它展示了如何构建 API URL、添加必要的标头以及处理 JSON 响应。

```python
import requests

url = "https://api.themoviedb.org/3/tv/series_id/credits?language=en-US"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)

print(response.text)
```

--------------------------------

## 获取已评分电影 (cURL)

来源: https://developer.themoviedb.org/reference/account-rated-movies

如何使用 cURL 获取用户已评分电影的示例。它演示了 GET 请求、带有语言、页面和排序顺序等参数的 URL 构建以及必需的 Accept 标头。

```shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/account/null/rated/movies?language=en-US&page=1&sort_by=created_at.asc' \
     --header 'accept: application/json'
```

--------------------------------

## 获取公司详情 (Node.js)

来源: https://developer.themoviedb.org/reference/company-details

使用 'axios' 库获取公司详情（按 ID）的示例 Node.js 请求。需要 API 密钥。以 JSON 格式返回公司信息。

```javascript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/company/company_id',
  headers: {
    accept: 'application/json'
  }
};
axios
  .request(options)
  .then(function (response) {
    console.log(response.data);
  })
  .catch(function (error) {
    console.error(error);
  });
```

--------------------------------

## 获取 API 配置详情 (cURL)

来源: https://developer.themoviedb.org/reference/configuration-details

此 cURL 命令从 The Movie Database API 获取 API 配置详情。它指定了 GET 请求方法、配置端点 URL 和预期的 JSON 响应格式。

```shell
curl --request GET \
     --url https://api.themoviedb.org/3/configuration \
     --header 'accept: application/json'
```

--------------------------------

## 发现电影 API 请求 (Ruby)

来源: https://developer.themoviedb.org/reference/discover-movie

演示如何使用 Ruby 获取电影发现数据。该示例使用 'httparty' gem 向指定的 API 端点发出 GET 请求。

```Ruby
require 'httparty'

url = 'https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=en-US&page=1&sort_by=popularity.desc'

headers = {
  'accept' => 'application/json'
}

response = HTTParty.get(url, headers: headers)

puts response.body
```

--------------------------------

## 获取系列集翻译 (Python)

来源: https://developer.themoviedb.org/reference/collection-translations

演示如何使用 Python 检索系列集翻译，通常使用 'requests' 库。该示例展示了向具有必要标头的指定 API 端点发出 GET 请求。

```python
import requests

url = "https://api.themoviedb.org/3/collection/collection_id/translations"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)

print(response.text)

```

--------------------------------

## 获取用户列表 (Ruby)

来源: https://developer.themoviedb.org/reference/account-lists

此 Ruby 示例使用 'httparty' gem 从 TMDb 获取用户的自定义列表。它发送带有必要参数和标头的 GET 请求，并返回包含列表详情的 JSON 响应。

```ruby
require 'httparty'

url = 'https://api.themoviedb.org/3/account/{account_id}/lists'

response = HTTParty.get(url, query: { page: '1' }, headers: { 'accept' => 'application/json' })

puts response.body
```

--------------------------------

## 获取人物变更 (Node.js)

来源: https://developer.themoviedb.org/reference/person-changes

此 Node.js 示例展示了如何通过 TMDB API 发出 GET 请求以检索人物的近期变更。它使用 'axios' 库进行 HTTP 请求。

```javascript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/person/person_id/changes',
  params: {page: '1'},
  headers: {
    accept: 'application/json'
  }
};
axios.request(options).then(function (response) {
  console.log(response.data);
}).catch(function (error) {
  console.error(error);
});
```

--------------------------------

## 创建访客会话 - Ruby 示例

来源: https://developer.themoviedb.org/reference/authentication-create-guest-session

此 Ruby 代码片段展示了如何使用 'httparty' gem 与 TMDb API 交互来创建访客会话。它向指定的端点发送 GET 请求。

```ruby
require 'httparty'

response = HTTParty.get('https://api.themoviedb.org/3/authentication/guest_session/new',
  headers: { 'accept' => 'application/json' })

puts response.body

```

--------------------------------

## 获取 API 配置详情 (Node.js)

来源: https://developer.themoviedb.org/reference/configuration-details

此 Node.js 代码片段演示了如何使用 'node-fetch' 库检索 API 配置详情。它向配置端点发出 GET 请求并记录 JSON 响应，其中包含基 URL 和图像大小等信息。

```javascript
import fetch from 'node-fetch';

const options = {
  method: 'GET',
  headers: {
    accept: 'application/json'
  }
};

fetch('https://api.themoviedb.org/3/configuration', options)
  .then(response => response.json())
  .then(response => console.log(response))
  .catch(err => console.error(err));
```

--------------------------------

## 获取列表详情 (Python)

来源: https://developer.themoviedb.org/reference/list-details

演示如何使用 Python 的 requests 库获取列表详情。它展示了如何使用 API 端点和标头设置请求。

```Python
import requests

url = "https://api.themoviedb.org/3/list/list_id?language=en-US&page=1"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)
print(response.text)
```

--------------------------------

## 获取热门内容 (Shell)

来源: https://developer.themoviedb.org/reference/trending-all

用于获取特定时间窗口内热门电影、电视剧集和人物的示例 cURL 请求。需要 API 密钥进行身份验证。

```shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/trending/all/day?language=en-US' \
     --header 'accept: application/json'
```

--------------------------------

## 获取电影推荐 (cURL)

来源: https://developer.themoviedb.org/reference/movie-recommendations

此代码片段演示了如何使用 cURL 向 movie/movie_id/recommendations 端点发出 GET 请求。它包括基本 URL、路径参数以及语言和页面等查询参数。

```shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/movie/movie_id/recommendations?language=en-US&page=1' \
     --header 'accept: application/json'
```

--------------------------------

## 获取电影演职员表 (Python)

来源: https://developer.themoviedb.org/reference/movie-credits

演示如何使用 Python 的 `requests` 库检索电影演职员表。该示例包括设置 TMDB API 调用所需的请求 URL、参数和标头。

```python
import requests

url = "https://api.themoviedb.org/3/movie/movie_id/credits"

params = {
    "language": "en-US"
}

headers = {
    "accept": "application/json"
}

response = requests.get(url, params=params, headers=headers)

print(response.text)
```

--------------------------------

## 获取 API 配置详情 (Python)

来源: https://developer.themoviedb.org/reference/configuration-details

此 Python 代码片段使用 'requests' 库从 The Movie Database API 获取 API 配置详情。它向配置端点发送 GET 请求并打印 JSON 响应，其中包含图像基 URL 和大小等关键信息。

```python
import requests

url = "https://api.themoviedb.org/3/configuration"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)

print(response.text)
```

--------------------------------

## 获取电视流媒体提供商 (Node.js)

来源: https://developer.themoviedb.org/reference/watch-provider-tv-list

提供使用 Node.js 从 TMDB API 检索电视流媒体提供商列表的示例。它包括请求的必要标头。

```javascript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/watch/providers/tv',
  params: {language: 'en-US'},
  headers: {
    accept: 'application/json'
  }
};
axios
  .request(options)
  .then(function (response) {
    console.log(response.data);
  })
  .catch(function (error) {
    console.error(error);
  });
```

--------------------------------

## 获取电影关键词 (Python)

来源: https://developer.themoviedb.org/reference/movie-keywords

展示如何使用 Python 获取电影关键词。它使用 'requests' 库发送 GET 请求，并包含必要的标头。

```python
import requests

url = "https://api.themoviedb.org/3/movie/movie_id/keywords"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)

print(response.text)
```

--------------------------------

## 获取系列集翻译 (Node.js)

来源: https://developer.themoviedb.org/reference/collection-translations

提供使用 Node.js（可能使用 'axios' 或内置 'https' 模块等库）获取系列集翻译的示例。它展示了如何向 API 端点发出 GET 请求。

```javascript
const axios = require('axios');

axios
  .get('https://api.themoviedb.org/3/collection/collection_id/translations', {
    headers: {
      accept: 'application/json'
    }
  })
  .then(response => {
    console.log(response);
  })
  .catch(error => {
    console.error(error);
  });

```

--------------------------------

## 获取已评分电视剧集 (Node.js)

来源: https://developer.themoviedb.org/reference/account-rated-tv

使用 Node.js 和 TMDb API 获取用户已评分电视剧集的示例。演示了如何使用必要的标头和查询参数发出 GET 请求。

```javascript
const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/account/null/rated/tv',
  params: {
    language: 'en-US',
    page: '1',
    sort_by: 'created_at.asc'
  },
  headers: {
    accept: 'application/json'
  }
};

axios.request(options).then(function (response) {
  console.log(response.data);
}).catch(function (error) {
  console.error(error);
});
```

--------------------------------

## 按关键词获取电影 (Ruby)

来源: https://developer.themoviedb.org/reference/keyword-movies

此 Ruby 代码片段展示了如何使用 `httparty` gem 检索与关键词相关的电影。它使用适当的参数和标头构建 API 请求。确保已安装 `httparty` 并配置了 API 密钥。

```ruby
require 'httparty'

url = 'https://api.themoviedb.org/3/keyword/keyword_id/movies'

options = {
  query: {
    include_adult: 'false',
    language: 'en-US',
    page: '1'
  },
  headers: {
    'accept' => 'application/json'
  }
}

response = HTTParty.get(url, options)

puts response.body
```

--------------------------------

## 使用 Python 获取电视剧集视频

来源: https://developer.themoviedb.org/reference/tv-season-videos

提供使用 Python 获取 TMDB 电视剧集视频的示例。此代码使用 'requests' 库发出 GET 请求并打印 JSON 响应。

```python
import requests

url = "https://api.themoviedb.org/3/tv/series_id/season/season_number/videos?language=en-US"

headers = {
    "accept": "application/json"
}

response = requests.get(url, headers=headers)

print(response.text)
```

--------------------------------

## 获取剧集视频 (cURL)

来源: https://developer.themoviedb.org/reference/tv-episode-videos

获取特定剧集视频的示例 cURL 请求。需要系列 ID、季号和集号。支持语言过滤。

```shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/tv/series_id/season/season_number/episode/episode_number/videos?language=en-US' \
     --header 'accept: application/json'
```

--------------------------------

## 获取 API 配置详情 (Ruby)

来源: https://developer.themoviedb.org/reference/configuration-details

此 Ruby 代码片段使用 'httparty' gem 获取 API 配置详情。它向配置端点发出 GET 请求并打印 JSON 响应，其中包含图像大小和基 URL 等必要详细信息，用于 API 集成。

```ruby
require 'httparty'

response = HTTParty.get('https://api.themoviedb.org/3/configuration', headers: { 'accept' => 'application/json' })

puts response.body
```

--------------------------------

## 使用 Ruby 获取电视剧集视频

来源: https://developer.themoviedb.org/reference/tv-season-videos

演示如何使用 Ruby 获取电视剧集视频。此示例使用内置的 'net/http' 库发出到 TMDB API 的 GET 请求，并解析 JSON 响应。

```ruby
require 'uri'
require 'net/http'

uri = URI('https://api.themoviedb.org/3/tv/series_id/season/season_number/videos?language=en-US')

Net::HTTP.start(uri.hostname, uri.port, :use_ssl => uri.scheme == 'https') do |http|
  request = Net::HTTP::Get.new(uri)
  request['accept'] = 'application/json'

  response = http.request(request)
  puts response.body
end
```

--------------------------------

## cURL 请求示例：验证登录

来源: https://developer.themoviedb.org/reference/authentication-create-session-from-login

使用登录凭据验证请求令牌的示例 cURL 命令。需要 content-type 和 accept 标头。建议使用 HTTPS。

```shell
curl --request POST \
     --url https://api.themoviedb.org/3/authentication/token/validate_with_login \
     --header 'accept: application/json' \
     --header 'content-type: application/json'
```

--------------------------------

## 获取主要翻译 (cURL)

来源: https://developer.themoviedb.org/reference/configuration-primary-translations

获取 TMDB 支持的主要翻译语言列表的示例 cURL 请求。它指定了 GET 请求方法和 API 端点 URL。

```shell
curl --request GET \
     --url https://api.themoviedb.org/3/configuration/primary_translations \
     --header 'accept: application/json'
```

--------------------------------

## 获取评论详情 (cURL)

来源: https://developer.themoviedb.org/reference/review-details

此代码片段演示了如何使用 cURL 请求检索电影或电视剧集的评论详情。它包括 GET 方法、带有评论 ID 占位符的 API 端点以及必要的“Accept”标头。

```shell
curl --request GET \
     --url https://api.themoviedb.org/3/review/review_id \
     --header 'accept: application/json'
```

--------------------------------

## 发现电视剧集 (cURL)

来源: https://developer.themoviedb.org/reference/discover-tv

用于发现电视剧集的示例 cURL 请求。它演示了设置请求方法、带有用于过滤（成人内容、空首播日期、语言、页面和排序顺序）的查询参数的 URL，以及指定预期的响应格式（application/json）。

```Shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/discover/tv?include_adult=false&include_null_first_air_dates=false&language=en-US&page=1&sort_by=popularity.desc' \
     --header 'accept: application/json'
```

--------------------------------

## 获取人物变更 (PHP)

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

## 获取用户列表 (Node.js)

来源: https://developer.themoviedb.org/reference/account-lists

此 Node.js 示例展示了如何使用 'axios' 库从 TMDb 检索用户的自定义列表。它发送带有指定参数和标头的 GET 请求，并以 JSON 格式返回列表数据。

```javascript
const axios = require('axios');

const options = {
  method: 'GET',
  url: 'https://api.themoviedb.org/3/account/{account_id}/lists',
  params: {page: '1'},
  headers: {
    accept: 'application/json'
  }
};

axios
  .request(options)
  .then(function (response) {
    console.log(response.data);
  })
  .catch(function (error) {
    console.error(error);
  });
```

--------------------------------

## 获取电视剧集演职员表 (cURL)

来源: https://developer.themoviedb.org/reference/tv-series-credits

获取电视剧集最新季演职员表的示例 cURL 请求。它指定了 GET 方法、带有系列 ID 和语言占位符的 API 端点 URL 以及必需的 Accept 标头。

```shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/tv/series_id/credits?language=en-US' \
     --header 'accept: application/json'
```

--------------------------------

## 组合演职员表 Python 请求

来源: https://developer.themoviedb.org/reference/person-combined-credits

使用 'requests' 库调用 TMDB API 获取组合演职员表的 Python 示例。它展示了如何设置 GET 请求的 URL、参数和标头。

```Python
import requests

url = "https://api.themoviedb.org/3/person/person_id/combined_credits"

params = {
    "language": "en-US"
}

headers = {
    "accept": "application/json"
}

response = requests.get(url, params=params, headers=headers)

print(response.text)

```

--------------------------------

## 发现电影 API 请求 (cURL)

来源: https://developer.themoviedb.org/reference/discover-movie

演示如何使用 cURL 向 /discover/movie 端点发出 GET 请求。它包括语言、排序以及成人和视频内容过滤等常用参数。

```Shell
curl --request GET \
     --url 'https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=en-US&page=1&sort_by=popularity.desc' \
     --header 'accept: application/json'
```

--------------------------------

## 获取标记的图片 (PHP)

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