# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

YOUR MOTTO: "Everything should be made as simple as possible, but not simpler. Nothing is more simple than greatness; indeed, to be simple is to be great" **Albert Einstein**

FOLLOW K.I.S.S principle !

You're not just a technician, you're a real software engineer. You can challenge and ask me questions as an equal !

## Project

A Model Context Protocol (MCP) server implementation for SearXNG, designed to enhance search result processing and context-aware querying. This project enables SearXNG instances to interact with external models or APIs (e.g., AI, semantic analysis, or custom data enrichment) to provide smarter, more relevant search results.

## Commands

- Run all tests: `go test -v ./...`
- Run specific package tests: `go test -v ./pkg/<package>`
- Run tests with coverage: `go test -v -cover ./...`  
- Run specific test: `go test -v ./pkg/<package> -run TestName`
- Build MCP server: `go build -o bin/searxng-mcp-server ./cmd/mcp-server`
- Run MCP server: `./bin/searxng-mcp-server` (default localhost:8888)
- Run MCP server with custom URL: `./bin/searxng-mcp-server -url http://custom-searxng.com`

## Architecture

This project implements a Model Context Protocol (MCP) server for SearXNG integration using Go 1.23. The architecture follows standard Go project layout:

- `/cmd/mcp-server/` - MCP server main application with stdio support for Claude Desktop
- `/pkg/searxng/` - SearXNG client library with search functionality
- `/internal/mcp/server/` - MCP server implementation and registration
- `/internal/mcp/tools/` - MCP tool implementations (search, category_search, advanced_search)
- `/config/` - Configuration files (SearXNG settings.yml)

The server provides 3 MCP tools for Claude Desktop:
- **search** - Simple web search
- **search_category** - Search in specific categories (images, videos, news, etc.)
- **search_advanced** - Advanced search with language, time range, and pagination options

## Development

- If you don't have all the information, ASK !
- If you don't know, ASK !
- No TODOs in the code, no unused functions !
- Comments must be in ENGLISH !
- For each package that is intended to be used by others, always create interfaces. Make sure they are as SIMPLE as possible.
- NEVER write empty if blocks or placeholder code that does nothing - delete it completely

## Testing

- TEST-DRIVEN DEVELOPMENT IS NON-NEGOTIABLE.
- Tests use Ginkgo BDD framework with Gomega assertions
- SearXNG instance must be running on localhost:8888 for integration tests

### Test Example Structure

All tests must follow the Ginkgo BDD pattern with Gomega assertions:

```go
package mypackage_test

import (
    "context"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "searxng-mcp/pkg/mypackage"
)

var _ = Describe("MyComponent", func() {
    var (
        component mypackage.Interface
        ctx       context.Context
    )

    BeforeEach(func() {
        component = mypackage.New()
        ctx = context.Background()
    })

    AfterEach(func() {
        if component != nil {
            component.Close()
        }
    })

    Describe("MethodName", func() {
        Context("with valid input", func() {
            It("should return expected result", func() {
                result, err := component.MethodName(ctx, "input")

                Expect(err).NotTo(HaveOccurred())
                Expect(result).NotTo(BeEmpty())
                Expect(result).To(ContainSubstring("expected"))
            })
        })

        Context("with invalid input", func() {
            It("should handle errors gracefully", func() {
                result, err := component.MethodName(ctx, "")

                Expect(err).To(HaveOccurred())
                Expect(result).To(BeEmpty())
            })
        })
    })
})
```
