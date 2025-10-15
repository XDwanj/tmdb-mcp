# Go Resty SetMultipartFields 示例

Source: https://github.com/go-resty/docs/blob/main/content/docs/multipart.md

此 Go 代码片段演示如何使用 `SetMultipartFields` 方法构建 multipart/form-data 请求。它展示了添加简单的表单字段、使用文件路径上传文件、包括进度回调、指定文件名和内容类型以及从 `io.Reader` 上传数据。

```go
myImageFile, _ := os.Open("/path/to/image-1.png")
myImageFileStat, _ := myImageFile.Stat()

// 使用各种组合和可能性进行演示
client.R().
    SetMultipartFields(
        []*resty.MultipartField{
            // 添加表单数据，顺序得以保留
            {
                Name:   "field1",
                Values: []string{"field1value1", "field1value2"},
            },
            {
                Name:   "field2",
                Values: []string{"field2value1", "field2value2"},
            },
            // 添加文件上传
            {
                Name:             "myfile_1",
                FilePath:         "/path/to/file-1.txt",
            },
            // 添加带有进度回调的文件上传
            {
                Name:             "myfile_1",
                FilePath:         "/path/to/file-1.txt",
                ProgressCallback: func(mp MultipartFieldProgress) {
    				// 使用进度详细信息
    				},
            },
            // 带有文件名和内容类型
            {
                Name:             "myimage_1",
                FileName:         "image-1.png",
                ContentType:      "image/png",
                FilePath:         "/path/to/image-1.png",
            },
            // 带有 io.Reader 和文件大小
            {
                Name:             "myimage_2",
                FileName:         "image-2.png",
                ContentType:      "image/png",
                Reader:           myImageFile,
                FileSize:         myImageFileStat.Size(),
            },
            // 带有 io.Reader
            {
                Name:        "uploadManifest1",
                FileName:    "upload-file-1.json",
                ContentType: "application/json",
                Reader:      strings.NewReader(`{"input": {"name": "Uploaded document 1", "_filename" : ["file1.txt"]}}`),
            },
            // 带有 io.Reader 和进度回调
            {
                Name:             "image-file1",
                FileName:         "image-file1.png",
                ContentType:      "image/png",
                Reader:           bytes.NewReader(fileBytes),
                ProgressCallback: func(mp MultipartFieldProgress) {
                    // 使用进度详细信息
                },
            },
        }...
    )

```

--------------------------------
