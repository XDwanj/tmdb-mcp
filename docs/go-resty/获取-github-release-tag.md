# 获取 GitHub Release Tag

Source: https://github.com/go-resty/docs/blob/main/layouts/shortcodes/restyrelease.html

此代码片段演示如何使用 Go 模板函数构建 GitHub release tag URL。它检索 GitHub 存储库路径和版本号以创建指向特定版本的链接。

```go
package main

import "fmt"

func main() {
	// 示例用法：
	// 假设 $gh_repo 是 "https://github.com/go-resty/resty"
	// 假设 $version 是 "v1.14.0"
	ghRepo := "https://github.com/go-resty/resty"
	version := "v1.14.0"
	releaseTagURL := fmt.Sprintf("%s/releases/tag/%s", ghRepo, version)
	fmt.Println(releaseTagURL)
}

```

--------------------------------
