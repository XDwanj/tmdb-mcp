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

// Genre represents a movie/TV genre
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CastMember represents a cast member in credits
type CastMember struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Character string `json:"character"`
}

// CrewMember represents a crew member in credits
type CrewMember struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Job        string `json:"job"`
	Department string `json:"department"`
}

// Credits represents cast and crew information
type Credits struct {
	Cast []CastMember `json:"cast"`
	Crew []CrewMember `json:"crew"`
}

// Video represents a video (trailer, teaser, etc.)
type Video struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
	Site string `json:"site"` // "YouTube"
	Type string `json:"type"` // "Trailer", "Teaser", etc.
}

// Videos represents a collection of videos
type Videos struct {
	Results []Video `json:"results"`
}

// MovieDetails represents detailed information about a movie
type MovieDetails struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	ReleaseDate string  `json:"release_date"`
	Runtime     int     `json:"runtime"`
	VoteAverage float64 `json:"vote_average"`
	Overview    string  `json:"overview"`
	Genres      []Genre `json:"genres"`
	Credits     Credits `json:"credits"` // 通过 append_to_response 获取
	Videos      Videos  `json:"videos"`  // 通过 append_to_response 获取
}

// TVDetails represents detailed information about a TV show
type TVDetails struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	FirstAirDate     string  `json:"first_air_date"`
	LastAirDate      string  `json:"last_air_date"`
	NumberOfSeasons  int     `json:"number_of_seasons"`
	NumberOfEpisodes int     `json:"number_of_episodes"`
	VoteAverage      float64 `json:"vote_average"`
	Overview         string  `json:"overview"`
	Genres           []Genre `json:"genres"`
	Credits          Credits `json:"credits"` // 通过 append_to_response 获取
	Videos           Videos  `json:"videos"`  // 通过 append_to_response 获取
}

// CombinedCastCredit represents a cast credit in combined credits
type CombinedCastCredit struct {
	ID           int    `json:"id"`
	MediaType    string `json:"media_type"` // "movie" or "tv"
	Title        string `json:"title"`      // 电影标题
	Name         string `json:"name"`       // 电视剧名称
	Character    string `json:"character"`
	ReleaseDate  string `json:"release_date"`
	FirstAirDate string `json:"first_air_date"`
}

// CombinedCrewCredit represents a crew credit in combined credits
type CombinedCrewCredit struct {
	ID           int    `json:"id"`
	MediaType    string `json:"media_type"`
	Title        string `json:"title"`
	Name         string `json:"name"`
	Job          string `json:"job"`
	Department   string `json:"department"`
	ReleaseDate  string `json:"release_date"`
	FirstAirDate string `json:"first_air_date"`
}

// CombinedCredits represents combined cast and crew credits
type CombinedCredits struct {
	Cast []CombinedCastCredit `json:"cast"`
	Crew []CombinedCrewCredit `json:"crew"`
}

// PersonDetails represents detailed information about a person
type PersonDetails struct {
	ID                 int             `json:"id"`
	Name               string          `json:"name"`
	Birthday           string          `json:"birthday"`
	Deathday           string          `json:"deathday"`
	Biography          string          `json:"biography"`
	PlaceOfBirth       string          `json:"place_of_birth"`
	KnownForDepartment string          `json:"known_for_department"`
	CombinedCredits    CombinedCredits `json:"combined_credits"` // 通过 append_to_response 获取
}

// DiscoverMovieResult represents a single result from TMDB discover movies
type DiscoverMovieResult struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	ReleaseDate string  `json:"release_date"`
	VoteAverage float64 `json:"vote_average"`
	Overview    string  `json:"overview"`
	GenreIDs    []int   `json:"genre_ids"`
	Popularity  float64 `json:"popularity"`
}

// DiscoverMoviesResponse represents the response from TMDB discover movies API
type DiscoverMoviesResponse struct {
	Page         int                   `json:"page"`
	Results      []DiscoverMovieResult `json:"results"`
	TotalPages   int                   `json:"total_pages"`
	TotalResults int                   `json:"total_results"`
}

// DiscoverTVResult represents a single result from TMDB discover TV
type DiscoverTVResult struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	FirstAirDate  string   `json:"first_air_date"`
	VoteAverage   float64  `json:"vote_average"`
	Overview      string   `json:"overview"`
	GenreIDs      []int    `json:"genre_ids"`
	OriginCountry []string `json:"origin_country"`
	Popularity    float64  `json:"popularity"`
}

// DiscoverTVResponse represents the response from TMDB discover TV API
type DiscoverTVResponse struct {
	Page         int                `json:"page"`
	Results      []DiscoverTVResult `json:"results"`
	TotalPages   int                `json:"total_pages"`
	TotalResults int                `json:"total_results"`
}

// TrendingResult represents a single result from TMDB trending endpoint
type TrendingResult struct {
	ID                 int     `json:"id"`
	MediaType          string  `json:"media_type"`          // "movie", "tv", "person"
	Title              string  `json:"title"`               // 电影标题 (movie only)
	Name               string  `json:"name"`                // 电视剧/人物名称 (tv/person)
	ReleaseDate        string  `json:"release_date"`        // 上映日期 (movie only)
	FirstAirDate       string  `json:"first_air_date"`      // 首播日期 (tv only)
	VoteAverage        float64 `json:"vote_average"`        // 评分 (movie/tv)
	Overview           string  `json:"overview"`            // 简介 (movie/tv)
	Popularity         float64 `json:"popularity"`          // 流行度
	KnownForDepartment string  `json:"known_for_department"` // 职业 (person only)
}

// TrendingResponse represents the response from TMDB trending API
type TrendingResponse struct {
	Page         int              `json:"page"`
	Results      []TrendingResult `json:"results"`
	TotalPages   int              `json:"total_pages"`
	TotalResults int              `json:"total_results"`
}
