package tools

// SearchArgs represents arguments for simple search tool
type SearchArgs struct {
	Query string `json:"query" jsonschema:"the search query to execute"`
}

// CategorySearchArgs represents arguments for category search tool
type CategorySearchArgs struct {
	Query      string   `json:"query" jsonschema:"the search query to execute"`
	Categories []string `json:"categories" jsonschema:"categories to search in"`
}

// AdvancedSearchArgs represents arguments for advanced search tool
type AdvancedSearchArgs struct {
	Query     string `json:"query" jsonschema:"the search query to execute"`
	Language  string `json:"language,omitempty" jsonschema:"language code for search results"`
	TimeRange string `json:"time_range,omitempty" jsonschema:"time range for search results"`
	Page      int    `json:"page,omitempty" jsonschema:"page number for pagination"`
}