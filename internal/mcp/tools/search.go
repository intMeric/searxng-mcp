package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"searxng-mcp/pkg/searxng"
)

// NewSearchTool creates and registers a simple search tool
func NewSearchTool(server *mcp.Server, client searxng.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search",
		Description: "Perform a simple web search using SearXNG",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args SearchArgs) (*mcp.CallToolResult, any, error) {
		// Validate query
		if args.Query == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Error: query parameter is required and must be a non-empty string"},
				},
			}, nil, nil
		}

		// Perform search
		response, err := searxng.SimpleSearch(ctx, client, args.Query)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Search failed: %v", err)},
				},
			}, nil, nil
		}

		// Format results for MCP
		content := formatSearchResults(response, args.Query)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: content},
			},
		}, nil, nil
	})
}

// formatSearchResults formats search results for MCP display
func formatSearchResults(response *searxng.SearchResponse, query string) string {
	if response.NumberOfResults == 0 {
		return fmt.Sprintf("No results found for query: %s", query)
	}

	result := fmt.Sprintf("Search results for \"%s\" (Total: %d results)\n\n", 
		response.Query, response.NumberOfResults)

	// Limit to first 10 results for readability
	maxResults := 10
	if len(response.Results) < maxResults {
		maxResults = len(response.Results)
	}

	for i, searchResult := range response.Results[:maxResults] {
		result += fmt.Sprintf("%d. **%s**\n", i+1, searchResult.Title)
		result += fmt.Sprintf("   URL: %s\n", searchResult.URL)
		if searchResult.Content != "" {
			// Truncate content if too long
			content := searchResult.Content
			if len(content) > 200 {
				content = content[:200] + "..."
			}
			result += fmt.Sprintf("   Summary: %s\n", content)
		}
		result += "\n"
	}

	if len(response.Results) > maxResults {
		result += fmt.Sprintf("... and %d more results\n", len(response.Results)-maxResults)
	}

	return result
}