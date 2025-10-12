package tools

import "github.com/XDwanj/tmdb-mcp/internal/tmdb"

// SearchParams represents the parameters for the search tool
type SearchParams struct {
	Query string `json:"query"` // 搜索关键词（必需）
	Page  int    `json:"page"`  // 页码（可选，默认 1）
}

// SearchResponse represents the response from the search tool
type SearchResponse struct {
	Results []tmdb.SearchResult `json:"results"`
}
