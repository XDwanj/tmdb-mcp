# 协议数据类型（Go）

来源：https://github.com/modelcontextprotocol/go-sdk/blob/main/design/design.md

定义了模型上下文协议的关键数据结构，这些结构是从其 JSON Schema 生成的。它包括参数、结果和内容表示的类型，利用 Go 的接口和结构来实现灵活性。

```Go
package main

import "encoding/json"

// Meta includes arbitrary data and a progress token.
type Meta struct {
    Data        map[string]any `json:"data,omitempty"`
    ProgressToken string         `json:"progress_token,omitempty"`
}

// ReadResourceParams defines parameters for reading a resource.
type ReadResourceParams struct {
    URI string `json:"uri"`
}

// CallToolResult represents the outcome of a tool call.
type CallToolResult struct {
    Meta    Meta      `json:"_meta,omitempty"`
    Content []Content `json:"content"`
    IsError bool      `json:"isError,omitempty"`
}

// Content is an interface representing different types of content.
// It is implemented by types like TextContent, ImageContent, etc.
type Content interface {
    // (unexported methods)
}

// TextContent represents textual content.
type TextContent struct {
    Text string
}

// ImageContent represents image content (example).
type ImageContent struct {
    // Image data or reference
}

// AudioContent represents audio content (example).
type AudioContent struct {
    // Audio data or reference
}

// EmbeddedResource represents embedded resource content (example).
type EmbeddedResource struct {
    // Resource data or reference
}

// ResourceContents is a struct for representing multiple resource contents, using optional fields for union types.
type ResourceContents struct {
    TextContents     []TextContent     `json:"text_contents,omitempty"`
    ImageContents    []ImageContent    `json:"image_contents,omitempty"`
    AudioContents    []AudioContent    `json:"audio_contents,omitempty"`
    EmbeddedResource []EmbeddedResource `json:"embedded_resources,omitempty"`
}

func main() {
    // Example usage (conceptual):
    // params := ReadResourceParams{URI: "example.com/resource"}
    // text := TextContent{Text: "Hello, world!"}
    // result := CallToolResult{Content: []Content{&text}, IsError: false}
    // _ = json.Marshal(params)
    // _ = json.Marshal(result)
}
```

--------------------------------
