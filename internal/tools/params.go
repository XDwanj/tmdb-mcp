package tools

import "github.com/XDwanj/tmdb-mcp/internal/tmdb"

// SearchParams represents the parameters for the search tool
type SearchParams struct {
	Query    string  `json:"query" jsonschema:"Search query for movies, TV shows, and people"`                                                  // 搜索关键词（必需）
	Page     int     `json:"page" jsonschema:"Page number (default: 1)"`                                                                        // 页码（可选，默认 1）
	Language *string `json:"language,omitempty" jsonschema:"ISO 639-1 language code (e.g., 'en', 'zh'). If not specified, uses config default"` // 语言参数（可选）
}

// SearchResponse represents the response from the search tool
type SearchResponse struct {
	Results []tmdb.SearchResult `json:"results" jsonschema:"List of search results"`
}

// GetDetailsParams represents the parameters for the get_details tool
type GetDetailsParams struct {
	MediaType string  `json:"media_type" jsonschema:"Media type (movie/tv/person)"`                                                              // 媒体类型（必需）
	ID        int     `json:"id" jsonschema:"TMDB ID of the content"`                                                                            // TMDB ID（必需）
	Language  *string `json:"language,omitempty" jsonschema:"ISO 639-1 language code (e.g., 'en', 'zh'). If not specified, uses config default"` // 语言参数（可选）
}

// DiscoverMoviesParams represents the parameters for the discover_movies tool
type DiscoverMoviesParams struct {
	WithGenres           *string  `json:"with_genres,omitempty" jsonschema:"Comma-separated genre IDs (e.g., '28,12' for Action and Adventure)"`
	PrimaryReleaseYear   *int     `json:"primary_release_year,omitempty" jsonschema:"Primary release year (e.g., 2020)"`
	VoteAverageGte       *float64 `json:"vote_average.gte,omitempty" jsonschema:"Minimum vote average (0-10)"`
	VoteAverageLte       *float64 `json:"vote_average.lte,omitempty" jsonschema:"Maximum vote average (0-10)"`
	WithOriginalLanguage *string  `json:"with_original_language,omitempty" jsonschema:"ISO 639-1 language code (e.g., 'en', 'zh')"`
	SortBy               *string  `json:"sort_by,omitempty" jsonschema:"Sort results by (e.g., 'popularity.desc', 'vote_average.desc', 'release_date.desc')"`
	Page                 *int     `json:"page,omitempty" jsonschema:"Page number (default: 1)"`
	Language             *string  `json:"language,omitempty" jsonschema:"ISO 639-1 language code (e.g., 'en', 'zh'). If not specified, uses config default"` // 语言参数（可选）
}

// DiscoverTVParams represents the parameters for the discover_tv tool
type DiscoverTVParams struct {
	WithGenres           *string  `json:"with_genres,omitempty" jsonschema:"Comma-separated genre IDs (e.g., '80,18' for Crime and Drama)"`
	FirstAirDateYear     *int     `json:"first_air_date_year,omitempty" jsonschema:"First air date year (e.g., 2020)"`
	VoteAverageGte       *float64 `json:"vote_average.gte,omitempty" jsonschema:"Minimum vote average (0-10)"`
	VoteAverageLte       *float64 `json:"vote_average.lte,omitempty" jsonschema:"Maximum vote average (0-10)"`
	WithOriginalLanguage *string  `json:"with_original_language,omitempty" jsonschema:"ISO 639-1 language code (e.g., 'en', 'zh')"`
	WithStatus           *string  `json:"with_status,omitempty" jsonschema:"TV show status (e.g., 'Returning Series', 'Ended', 'Canceled')"`
	SortBy               *string  `json:"sort_by,omitempty" jsonschema:"Sort results by (e.g., 'popularity.desc', 'vote_average.desc', 'first_air_date.desc')"`
	Page                 *int     `json:"page,omitempty" jsonschema:"Page number (default: 1)"`
	Language             *string  `json:"language,omitempty" jsonschema:"ISO 639-1 language code (e.g., 'en', 'zh'). If not specified, uses config default"` // 语言参数（可选）
}
