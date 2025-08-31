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
		content := formatSearchResultsJSON(response)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: content},
			},
		}, nil, nil
	})
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