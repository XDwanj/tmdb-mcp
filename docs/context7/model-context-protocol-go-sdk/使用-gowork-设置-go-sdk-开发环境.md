# 使用 go.work 设置 Go SDK 开发环境

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/CONTRIBUTING.md

演示如何初始化 Go 工作区，以便针对本地项目测试 SDK 更改。这对于多模块开发非常有用，其中 SDK 与使用它的项目一起被修改。

```bash
go work init ./project ./go-sdk
```

--------------------------------
