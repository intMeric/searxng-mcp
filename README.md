# searxng-mcp-server

A Model Context Protocol (MCP) server for SearXNG integration. Provides search tools for Claude Desktop.

## Requirements

- Go 1.23+
- Docker (for auto-launch)
- SearXNG instance (auto-launched or external)

## Installation

```bash
go build -o bin/searxng-mcp-server ./cmd/mcp-server
```

## Usage

```bash
# Default (expects SearXNG on localhost:8888)
./bin/searxng-mcp-server

# Custom SearXNG URL
./bin/searxng-mcp-server -url http://your-searxng.com

# Auto-launch SearXNG container (Unix/Linux only)
./bin/searxng-mcp-server -auto-launch
```

## MCP Tools

- `search` - Simple web search
- `search_category` - Search by category (images, videos, news, etc.)
- `search_advanced` - Advanced search with language, time range, pagination

## Claude Desktop Configuration

Add to your Claude Desktop config:

```json
{
  "mcpServers": {
    "searxng": {
      "command": "/path/to/bin/searxng-mcp-server",
      "args": ["-auto-launch"]
    }
  }
}
```

## Development

```bash
# Run tests
go test -v ./...

# Test specific package
go test -v ./pkg/searxng

# Test with coverage
go test -v -cover ./...
```

## Architecture

- `/cmd/mcp-server/` - Main application
- `/pkg/searxng/` - SearXNG client library
- `/internal/mcp/` - MCP server and tools implementation

## License

MIT