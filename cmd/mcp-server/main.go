package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"searxng-mcp/internal/mcp/server"
)

func main() {
	// Define command line flags
	var searxngURL string
	flag.StringVar(&searxngURL, "url", "http://localhost:8888", "SearXNG server URL")
	flag.Parse()

	// Create context that cancels on interrupt signal
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Initialize our MCP server with custom URL
	mcpServer, err := server.NewSearXNGServer(searxngURL)
	if err != nil {
		log.Fatalf("Failed to create SearXNG MCP server: %v", err)
	}

	// Create stdio transport for Claude Desktop communication
	transport := &mcp.StdioTransport{}

	// Run the server
	if err := mcpServer.Run(ctx, transport); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	log.Println("SearXNG MCP Server shutdown gracefully")
}