# Go Resty 请求方法重命名和更改

Source: https://github.com/go-resty/docs/blob/main/content/docs/upgrading-to-v3.md

本节详细介绍了 Go Resty Request 对象中方法的重命名和功能更改，以帮助迁移到 v3。

```APIDOC
Request Methods:

- `Request.QueryParams()`: 从 `QueryParam` 重命名。
- `Request.AuthToken()`: 从 `Token` 重命名。
- `Request.DoNotParseResponse()`: 从 `NotParseResponse` 重命名。
- `Request.SetExpectResponseContentType(contentType string)`: 从 `ExpectContentType` 重命名。
- `Request.SetForceResponseContentType(contentType string)`: 从 `ForceContentType` 重命名。
- `Request.SetOutputFileName(filename string)`: 从 `SetOutput` 重命名。
- `Request.EnableGenerateCurlCmd()`: 从 `EnableGenerateCurlOnDebug` 重命名。
- `Request.DisableGenerateCurlCmd()`: 从 `DisableGenerateCurlOnDebug` 重命名。
- `Request.CurlCmd()`: 从 `GenerateCurlCommand` 重命名。
- `Request.AddRetryConditions(conditions ...RetryConditionFunc)`: 从 `AddRetryCondition` 重命名。
```

--------------------------------
