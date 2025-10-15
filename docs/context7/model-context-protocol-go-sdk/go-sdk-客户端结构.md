# Go SDK 客户端结构

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了 Go SDK 中客户端的结构，包括创建新客户端、连接到传输以及检索活动会话的方法。它还概述了客户端选项和会话详细信息。

```Go
type Client struct { /* ... */ }
func NewClient(impl *Implementation, opts *ClientOptions) *Client
func (*Client) Connect(context.Context, Transport) (*ClientSession, error)
func (*Client) Sessions() iter.Seq[*ClientSession]

type ClientOptions struct { /* ... */ } // described below

type ClientSession struct { /* ... */ }
func (*ClientSession) Client() *Client
func (*ClientSession) Close() error
func (*ClientSession) Wait() error
// Methods for calling through the ClientSession are described below.
// For example: ClientSession.ListTools.
```

--------------------------------
