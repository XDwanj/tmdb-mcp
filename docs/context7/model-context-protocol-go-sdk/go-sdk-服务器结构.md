# Go SDK 服务器结构

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了 Go SDK 中服务器的结构，包括创建新服务器、连接到传输以及检索活动会话的方法。它还概述了服务器选项和会话详细信息。

```Go
type Server struct { /* ... */ }
func NewServer(impl *Implementation, opts *ServerOptions) *Server
func (*Server) Connect(context.Context, Transport) (*ServerSession, error)
func (*Server) Sessions() iter.Seq[*ServerSession]
// Methods for adding/removing server features are described below.

type ServerOptions struct { /* ... */ } // described below

type ServerSession struct { /* ... */ }
func (*ServerSession) Server() *Server
func (*ServerSession) Close() error
func (*ServerSession) Wait() error
// Methods for calling through the ServerSession are described below.
// For example: ServerSession.ListRoots.
```

--------------------------------
