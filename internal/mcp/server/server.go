package server

import (
	"context"
	"fmt"

	"searxng-mcp/internal/mcp/tools"
	"searxng-mcp/pkg/searxng"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// SearXNGServer wraps the MCP server with SearXNG functionality
type SearXNGServer struct {
	mcpServer     *mcp.Server
	searxngClient searxng.Client
}

// NewSearXNGServer creates a new MCP server with SearXNG tools
func NewSearXNGServer(searxngURL string) (*SearXNGServer, error) {
	// Default to localhost if empty URL provided
	if searxngURL == "" {
		searxngURL = "http://localhost:8888"
	}
	
	// Create SearXNG client
	searxngClient := searxng.NewClient(searxngURL)

	// Create MCP server with implementation details
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "searxng-mcp-server",
		Version: "1.0.0",
	}, nil)

	// Create our server wrapper
	server := &SearXNGServer{
		mcpServer:     mcpServer,
		searxngClient: searxngClient,
	}

	// Register SearXNG tools
	if err := server.registerTools(); err != nil {
		return nil, fmt.Errorf("failed to register tools: %w", err)
	}

	return server, nil
}

// registerTools registers all SearXNG tools with the MCP server
func (s *SearXNGServer) registerTools() error {
	// Register simple search tool
	tools.NewSearchTool(s.mcpServer, s.searxngClient)

	// Register category search tool
	tools.NewCategorySearchTool(s.mcpServer, s.searxngClient)

	// Register advanced search tool
	tools.NewAdvancedSearchTool(s.mcpServer, s.searxngClient)

	return nil
}

// Run starts the MCP server with the given transport
func (s *SearXNGServer) Run(ctx context.Context, transport mcp.Transport) error {
	return s.mcpServer.Run(ctx, transport)
}
