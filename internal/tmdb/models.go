package tmdb

// SearchResult represents a single result from TMDB multi search
type SearchResult struct {
	ID           int     `json:"id"`
	MediaType    string  `json:"media_type"`     // "movie", "tv", "person"
	Title        string  `json:"title"`          // 电影标题
	Name         string  `json:"name"`           // 电视剧/人物名称
	ReleaseDate  string  `json:"release_date"`   // 上映日期
	FirstAirDate string  `json:"first_air_date"` // 首播日期
	VoteAverage  float64 `json:"vote_average"`   // 评分
	Overview     string  `json:"overview"`       // 简介
}

// SearchResponse represents the response from TMDB multi search API
type SearchResponse struct {
	Page         int            `json:"page"`
	Results      []SearchResult `json:"results"`
	TotalPages   int            `json:"total_pages"`
	TotalResults int            `json:"total_results"`
}
