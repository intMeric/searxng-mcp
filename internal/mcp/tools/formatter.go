package tools

import (
	"encoding/json"
	"fmt"
	"searxng-mcp/pkg/searxng"
)

// SimplifiedResult represents a simplified search result
type SimplifiedResult struct {
	Title   string  `json:"title"`
	URL     string  `json:"url"`
	Summary string  `json:"summary"`
	Rank    int     `json:"rank"`
	Date    *string `json:"date,omitempty"`
}

// FormatSearchResultsJSON formats search results as simplified JSON (exported for testing)
func FormatSearchResultsJSON(response *searxng.SearchResponse) string {
	return formatSearchResultsJSON(response)
}

// formatSearchResultsJSON formats search results as simplified JSON
func formatSearchResultsJSON(response *searxng.SearchResponse) string {
	if response.NumberOfResults == 0 {
		return `{"results": [], "total": 0, "message": "No results found"}`
	}

	// Limit to first 10 results for readability
	maxResults := 10
	if len(response.Results) < maxResults {
		maxResults = len(response.Results)
	}

	var simplifiedResults []SimplifiedResult
	for i, searchResult := range response.Results[:maxResults] {
		// Truncate content if too long for summary
		summary := searchResult.Content
		if len(summary) > 500 {
			summary = summary[:500] + "..."
		}

		// Handle date if available
		var date *string
		if searchResult.PublishedDate != nil {
			if dateStr, ok := searchResult.PublishedDate.(string); ok && dateStr != "" {
				date = &dateStr
			}
		}

		simplifiedResult := SimplifiedResult{
			Title:   searchResult.Title,
			URL:     searchResult.URL,
			Summary: summary,
			Rank:    i + 1,
			Date:    date,
		}
		simplifiedResults = append(simplifiedResults, simplifiedResult)
	}

	// Create response structure
	response_data := map[string]interface{}{
		"results": simplifiedResults,
		"total":   response.NumberOfResults,
		"query":   response.Query,
	}

	jsonData, err := json.MarshalIndent(response_data, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "Failed to format results: %v"}`, err)
	}

	return string(jsonData)
}
