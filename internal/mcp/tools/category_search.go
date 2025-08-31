package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"searxng-mcp/pkg/searxng"
)

// NewCategorySearchTool creates and registers a category search tool
func NewCategorySearchTool(server *mcp.Server, client searxng.Client) {
	availableCategories := strings.Join(getAllCategoryNames(), ", ")
	
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_category",
		Description: fmt.Sprintf("Perform a search in specific categories using SearXNG. Available categories: %s", availableCategories),
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"query": {
					Type:        "string",
					Description: "The search query to execute",
				},
				"categories": {
					Type:        "array",
					Description: fmt.Sprintf("Categories to search in. Available: %s", availableCategories),
					Items: &jsonschema.Schema{
						Type: "string",
						Enum: interfaceSlice(getAllCategoryNames()),
					},
					MinItems: intPtr(1),
				},
			},
			Required: []string{"query", "categories"},
		},
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
			validCategoriesFormatted := strings.Join(getAllCategoryNames(), "\", \"")
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Error: invalid categories found: \"%s\". Please use only valid categories from: [\"%s\"]. Example: {\"query\": \"machine learning\", \"categories\": [\"science\", \"it\"]}",
						strings.Join(invalidCategories, "\", \""),
						validCategoriesFormatted)},
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
		content := formatSearchResultsJSON(response)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: content},
			},
		}, nil, nil
	})
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

// interfaceSlice converts string slice to interface slice for JSON schema enum
func interfaceSlice(strings []string) []interface{} {
	result := make([]interface{}, len(strings))
	for i, s := range strings {
		result[i] = s
	}
	return result
}

// intPtr returns a pointer to an int value
func intPtr(i int) *int {
	return &i
}