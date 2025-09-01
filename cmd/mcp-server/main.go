package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"searxng-mcp/internal/mcp/server"
	"searxng-mcp/pkg/config"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// createDefaultSettings creates a default settings.yml file for SearXNG using embedded config
func createDefaultSettings(settingsPath string) error {
	log.Printf("Creating default settings file: %s", settingsPath)
	return os.WriteFile(settingsPath, []byte(config.DefaultSettings), 0644)
}


// ensureSearXNGContainer ensures that a SearXNG Docker container is running
func ensureSearXNGContainer(ctx context.Context) error {
	log.Println("Checking SearXNG Docker container status...")

	// Check if container exists and its status
	cmd := exec.CommandContext(ctx, "docker", "ps", "-a", "--filter", "name=searxng", "--format", "{{.Names}}\t{{.Status}}")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check Docker containers: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	containerExists := false
	containerRunning := false

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) >= 2 && parts[0] == "searxng" {
			containerExists = true
			status := parts[1]
			if strings.HasPrefix(status, "Up") {
				containerRunning = true
				log.Println("SearXNG container is already running")
				break
			}
		}
	}

	if containerExists && !containerRunning {
		// Container exists but is stopped, start it
		log.Println("Starting existing SearXNG container...")
		cmd = exec.CommandContext(ctx, "docker", "start", "searxng")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to start existing SearXNG container: %w", err)
		}
		log.Println("SearXNG container started successfully")
	} else if !containerExists {
		// Container doesn't exist, create and run it
		log.Println("Creating new SearXNG container...")

		// Get user home directory for config path
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		configPath := filepath.Join(homeDir, ".config", "searxng-mcp")

		// Ensure config directory exists, create if not
		if err := os.MkdirAll(configPath, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		// Check if settings.yml exists, create one if not
		settingsPath := filepath.Join(configPath, "settings.yml")
		if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
			// Create settings from embedded config
			if err := createDefaultSettings(settingsPath); err != nil {
				return fmt.Errorf("failed to create default settings: %w", err)
			}
			log.Printf("Created default configuration from embedded settings: %s", settingsPath)
		}

		// Remove any existing container with the same name first (in case of race condition)
		removeCmd := exec.CommandContext(ctx, "docker", "rm", "-f", "searxng")
		if err := removeCmd.Run(); err != nil {
			log.Printf("Info: attempted to remove existing container: %v", err)
		}

		cmd = exec.CommandContext(ctx, "docker", "run", "--name", "searxng", "-d",
			"-p", "8888:8080",
			"-v", fmt.Sprintf("%s:/etc/searxng/", configPath),
			"docker.io/searxng/searxng:latest")

		output, err := cmd.CombinedOutput()
		if err != nil {
			// If creation fails, check if container now exists (race condition)
			checkCmd := exec.CommandContext(ctx, "docker", "ps", "-q", "--filter", "name=searxng")
			if checkOutput, checkErr := checkCmd.Output(); checkErr == nil && len(strings.TrimSpace(string(checkOutput))) > 0 {
				log.Println("SearXNG container was created by another process, continuing...")
				return nil
			}
			return fmt.Errorf("failed to create SearXNG container: %v - Output: %s", err, string(output))
		}
		log.Println("SearXNG container created and started successfully")
	}

	return nil
}

func main() {
	// Define command line flags
	var searxngURL string
	var autoLaunch bool
	flag.StringVar(&searxngURL, "url", "http://localhost:8888", "SearXNG server URL")
	flag.BoolVar(&autoLaunch, "auto-launch", false, "Automatically start SearXNG container with Docker if needed")
	flag.Parse()

	// Create context that cancels on interrupt signal
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Auto-launch SearXNG container if requested
	if autoLaunch {
		if err := ensureSearXNGContainer(ctx); err != nil {
			log.Fatalf("Failed to ensure SearXNG container is running: %v", err)
		}
	}

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
