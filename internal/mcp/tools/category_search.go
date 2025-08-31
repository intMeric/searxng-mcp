package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"searxng-mcp/pkg/searxng"
)

// NewCategorySearchTool creates and registers a category search tool
func NewCategorySearchTool(server *mcp.Server, client searxng.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_category",
		Description: "Perform a search in specific categories using SearXNG",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args CategorySearchArgs) (*mcp.CallToolResult, any, error) {
		// Validate query
		if args.Query == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Error: query parameter is required and must be a non-empty string"},
				},
			}, nil, nil
		}

		// Validate categories
		if len(args.Categories) == 0 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Error: categories parameter is required and must be a non-empty array"},
				},
			}, nil, nil
		}

		// Convert to searxng.Category slice and validate
		var categories []searxng.Category
		var invalidCategories []string

		for _, catStr := range args.Categories {
			category, err := searxng.ValidateCategory(catStr)
			if err != nil {
				invalidCategories = append(invalidCategories, catStr)
				continue
			}
			categories = append(categories, category)
		}

		if len(invalidCategories) > 0 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Error: invalid categories: %s. Valid categories are: %s",
						strings.Join(invalidCategories, ", "),
						strings.Join(getAllCategoryNames(), ", "))},
				},
			}, nil, nil
		}

		if len(categories) == 0 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Error: no valid categories provided"},
				},
			}, nil, nil
		}

		// Perform search
		response, err := searxng.SearchWithCategory(ctx, client, args.Query, categories...)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Category search failed: %v", err)},
				},
			}, nil, nil
		}

		// Format results for MCP
		content := formatCategorySearchResults(response, args.Query, args.Categories)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: content},
			},
		}, nil, nil
	})
}

// formatCategorySearchResults formats category search results for MCP display
func formatCategorySearchResults(response *searxng.SearchResponse, query string, categories []string) string {
	if response.NumberOfResults == 0 {
		return fmt.Sprintf("No results found for query: \"%s\" in categories: %s",
			query, strings.Join(categories, ", "))
	}

	result := fmt.Sprintf("Search results for \"%s\" in categories [%s] (Total: %d results)\n\n",
		response.Query, strings.Join(categories, ", "), response.NumberOfResults)

	// Limit to first 10 results for readability
	maxResults := 10
	if len(response.Results) < maxResults {
		maxResults = len(response.Results)
	}

	for i, searchResult := range response.Results[:maxResults] {
		result += fmt.Sprintf("%d. **%s** [%s]\n", i+1, searchResult.Title, searchResult.Category)
		result += fmt.Sprintf("   URL: %s\n", searchResult.URL)
		if searchResult.Content != "" {
			// Truncate content if too long
			content := searchResult.Content
			if len(content) > 200 {
				content = content[:200] + "..."
			}
			result += fmt.Sprintf("   Summary: %s\n", content)
		}
		if len(searchResult.Engines) > 0 {
			result += fmt.Sprintf("   Sources: %s\n", strings.Join(searchResult.Engines, ", "))
		}
		result += "\n"
	}

	if len(response.Results) > maxResults {
		result += fmt.Sprintf("... and %d more results\n", len(response.Results)-maxResults)
	}

	return result
}

// getAllCategoryNames returns all valid category names as strings
func getAllCategoryNames() []string {
	categories := searxng.GetAllCategories()
	names := make([]string, len(categories))
	for i, cat := range categories {
		names[i] = string(cat)
	}
	return names
}