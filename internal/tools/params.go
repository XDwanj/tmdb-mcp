package tools

import "github.com/XDwanj/tmdb-mcp/internal/tmdb"

// SearchParams represents the parameters for the search tool
type SearchParams struct {
	Query string `json:"query" jsonschema:"Search query for movies, TV shows, and people"` // 搜索关键词（必需）
	Page  int    `json:"page" jsonschema:"Page number (default: 1)"`                       // 页码（可选，默认 1）
}

// SearchResponse represents the response from the search tool
type SearchResponse struct {
	Results []tmdb.SearchResult `json:"results" jsonschema:"List of search results"`
}

// GetDetailsParams represents the parameters for the get_details tool
type GetDetailsParams struct {
	MediaType string `json:"media_type" jsonschema:"Media type (movie/tv/person)"` // 媒体类型（必需）
	ID        int    `json:"id" jsonschema:"TMDB ID of the content"`               // TMDB ID（必需）
}
