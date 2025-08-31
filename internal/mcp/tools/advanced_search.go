package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"searxng-mcp/pkg/searxng"
)

// NewAdvancedSearchTool creates and registers an advanced search tool
func NewAdvancedSearchTool(server *mcp.Server, client searxng.Client) {
	availableTimeRanges := strings.Join(getAllTimeRangeNames(), ", ")
	
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_advanced",
		Description: fmt.Sprintf("Perform an advanced search with language, time range, and pagination options using SearXNG. Available time ranges: %s", availableTimeRanges),
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"query": {
					Type:        "string",
					Description: "The search query to execute",
				},
				"language": {
					Type:        "string",
					Description: "Language code for search results (e.g., 'en', 'fr', 'es')",
				},
				"time_range": {
					Type:        "string",
					Description: fmt.Sprintf("Time range for search results. Available: %s", availableTimeRanges),
					Enum:        interfaceSliceFromStringSlice(getAllTimeRangeNames()),
				},
				"page": {
					Type:        "integer",
					Description: "Page number for pagination (1-50)",
					Minimum:     floatPtr(1),
					Maximum:     floatPtr(50),
				},
			},
			Required: []string{"query"},
		},
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

// interfaceSliceFromStringSlice converts string slice to interface slice for JSON schema enum
func interfaceSliceFromStringSlice(strings []string) []interface{} {
	result := make([]interface{}, len(strings))
	for i, s := range strings {
		result[i] = s
	}
	return result
}

// floatPtr returns a pointer to a float64 value
func floatPtr(f float64) *float64 {
	return &f
}