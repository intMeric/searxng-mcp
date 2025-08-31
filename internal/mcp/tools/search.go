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
		content := formatSearchResultsJSON(response)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: content},
			},
		}, nil, nil
	})
}

