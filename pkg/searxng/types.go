package searxng


// Category represents search categories available in SearXNG
type Category string

const (
	CategoryGeneral     Category = "general"
	CategoryImages      Category = "images"
	CategoryVideos      Category = "videos"
	CategoryNews        Category = "news"
	CategoryMap         Category = "map"
	CategoryMusic       Category = "music"
	CategoryIT          Category = "it"
	CategoryScience     Category = "science"
	CategoryFiles       Category = "files"
	CategorySocialMedia Category = "social media"
)

// TimeRange represents time filtering options
type TimeRange string

const (
	TimeRangeDay   TimeRange = "day"
	TimeRangeMonth TimeRange = "month"
	TimeRangeYear  TimeRange = "year"
)

// SearchRequest represents the parameters for a search query
type SearchRequest struct {
	Query     string     `json:"q"`
	Language  string     `json:"language,omitempty"`
	TimeRange TimeRange  `json:"time_range,omitempty"`
	Category  []Category `json:"categories,omitempty"`
	PageNo    int        `json:"pageno,omitempty"`
}

// SearchResult represents a single search result
type SearchResult struct {
	URL           string    `json:"url"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Thumbnail     string    `json:"thumbnail"`
	Engine        string    `json:"engine"`
	Template      string    `json:"template"`
	ParsedURL     []string  `json:"parsed_url"`
	ImgSrc        string    `json:"img_src"`
	Priority      string    `json:"priority"`
	Engines       []string  `json:"engines"`
	Positions     []int     `json:"positions"`
	Score         float64   `json:"score"`
	Category      string    `json:"category"`
	PublishedDate any `json:"publishedDate"`
}

// SearchResponse represents the complete response from SearXNG API
type SearchResponse struct {
	Query           string         `json:"query"`
	NumberOfResults int            `json:"number_of_results"`
	Results         []SearchResult `json:"results"`
}