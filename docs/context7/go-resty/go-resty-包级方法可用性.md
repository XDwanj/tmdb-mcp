# Go Resty 包级方法可用性

Source: https://github.com/go-resty/docs/blob/main/content/docs/upgrading-to-v3.md

本节列出了 Go Resty 中包级别的实用方法。这些方法涵盖了各种功能，例如检查字符串是否为空、内容类型检测和反序列化。

```APIDOC
IsStringEmpty
  - 检查字符串是否为空。

IsJSONType
  - 检测内容是否为 JSON 类型。

IsXMLType
  - 检测内容是否为 XML 类型。

DetectContentType
  - 检测数据的内容类型。

Unmarshalc
  - 带有内容类型检测的反序列化数据。

Backoff
  - 为重试实现退避策略。
```

--------------------------------
