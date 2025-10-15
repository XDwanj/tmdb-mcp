# 在 Go 中添加和删除资源/模板

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

演示在服务器上添加和删除资源及资源模板的方法。它展示了 AddResource、AddResourceTemplate、RemoveResources 和 RemoveResourceTemplates 的签名。

```Go
func (*Server) AddResource(*Resource, ResourceHandler)
func (*Server) AddResourceTemplate(*ResourceTemplate, ResourceHandler)

func (s *Server) RemoveResources(uris ...string)
func (s *Server) RemoveResourceTemplates(uriTemplates ...string)
```

--------------------------------
