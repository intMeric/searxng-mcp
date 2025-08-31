package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"searxng-mcp/pkg/searxng"
)

// NewAdvancedSearchTool creates and registers an advanced search tool
func NewAdvancedSearchTool(server *mcp.Server, client searxng.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_advanced",
		Description: "Perform an advanced search with language, time range, and pagination options using SearXNG",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args AdvancedSearchArgs) (*mcp.CallToolResult, any, error) {
		// Validate query
		if args.Query == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Error: query parameter is required and must be a non-empty string"},
				},
			}, nil, nil
		}

		// Build search options
		opts := searxng.SearchOptions{
			PageNo: 1, // default
		}

		if args.Page > 0 {
			if args.Page < 1 || args.Page > 50 {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{
						&mcp.TextContent{Text: "Error: page parameter must be between 1 and 50"},
					},
				}, nil, nil
			}
			opts.PageNo = args.Page
		}

		if args.Language != "" {
			opts.Language = args.Language
		}

		if args.TimeRange != "" {
			timeRange, err := searxng.ValidateTimeRange(args.TimeRange)
			if err != nil {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{
						&mcp.TextContent{Text: fmt.Sprintf("Error: invalid time_range '%s'. Valid options are: %s",
							args.TimeRange, strings.Join(getAllTimeRangeNames(), ", "))},
					},
				}, nil, nil
			}
			opts.TimeRange = timeRange
		}

		// Perform search
		response, err := searxng.SearchWithOptions(ctx, client, args.Query, opts)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Advanced search failed: %v", err)},
				},
			}, nil, nil
		}

		// Format results for MCP
		content := formatAdvancedSearchResults(response, args.Query, opts)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: content},
			},
		}, nil, nil
	})
}

// formatAdvancedSearchResults formats advanced search results for MCP display
func formatAdvancedSearchResults(response *searxng.SearchResponse, query string, opts searxng.SearchOptions) string {
	var optionsStr []string
	if opts.Language != "" {
		optionsStr = append(optionsStr, fmt.Sprintf("language: %s", opts.Language))
	}
	if opts.TimeRange != "" {
		optionsStr = append(optionsStr, fmt.Sprintf("time_range: %s", string(opts.TimeRange)))
	}
	if opts.PageNo > 1 {
		optionsStr = append(optionsStr, fmt.Sprintf("page: %d", opts.PageNo))
	}

	optionsDisplay := ""
	if len(optionsStr) > 0 {
		optionsDisplay = fmt.Sprintf(" [%s]", strings.Join(optionsStr, ", "))
	}

	if response.NumberOfResults == 0 {
		return fmt.Sprintf("No results found for query: \"%s\"%s", query, optionsDisplay)
	}

	result := fmt.Sprintf("Advanced search results for \"%s\"%s (Total: %d results)\n\n",
		response.Query, optionsDisplay, response.NumberOfResults)

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
		if searchResult.Score > 0 {
			result += fmt.Sprintf("   Score: %.1f\n", searchResult.Score)
		}
		if len(searchResult.Engines) > 0 {
			result += fmt.Sprintf("   Sources: %s\n", strings.Join(searchResult.Engines, ", "))
		}
		result += "\n"
	}

	if len(response.Results) > maxResults {
		result += fmt.Sprintf("... and %d more results\n", len(response.Results)-maxResults)
	}

	if opts.PageNo > 1 {
		result += fmt.Sprintf("\nShowing page %d of results\n", opts.PageNo)
	}

	return result
}

// getAllTimeRangeNames returns all valid time range names as strings
func getAllTimeRangeNames() []string {
	timeRanges := searxng.GetAllTimeRanges()
	names := make([]string, len(timeRanges))
	for i, tr := range timeRanges {
		names[i] = string(tr)
	}
	return names
}