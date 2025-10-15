# Go Resty 包级类型更改

Source: https://github.com/go-resty/docs/blob/main/content/docs/upgrading-to-v3.md

本节详细介绍了 Go Resty 包中类型的更改。它解释了某些类型如何变得未导出、被新功能替换或集成到增强功能中，如负载均衡和多部分字段。

```APIDOC
User
  - 现在已设为未导出，称为 'credentials'。

SRVRecord
  - 已被支持 SRV 记录查找的新负载均衡器功能取代。

File
  - 已被增强的 MultipartField 功能取代。

RequestLog, ResponseLog
  - 已弃用：请改用 DebugLog。

RequestLogCallback, ResponseLogCallback
  - 已弃用：请改用 DebugLogCallbackFunc。
```

--------------------------------
