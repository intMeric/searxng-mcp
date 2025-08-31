package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"searxng-mcp/internal/mcp/server"
)

// createDefaultSettings creates a default settings.yml file for SearXNG
func createDefaultSettings(settingsPath string) error {
	log.Printf("Creating default settings file: %s", settingsPath)
	
	defaultSettings := `general:
  debug: false
  instance_name: "SearXNG MCP"
  privacypolicy_url: false
  donation_url: false
  contact_url: false
  enable_metrics: false

search:
  safe_search: 0
  autocomplete: ""
  autocomplete_min: 4
  favicon_resolver: ""
  default_lang: "auto"
  formats:
    - html
    - json

server:
  port: 8080
  bind_address: "0.0.0.0"
  secret_key: "changeme"
  image_proxy: false
  http_protocol_version: "1.0"
  method: "POST"

categories_as_tabs:
  general:
  images:
  videos:
  news:
  map:
  music:
  it:
  science:
  files:
  social media:

engines:
  - name: google
    engine: google
    shortcut: go
  
  - name: bing
    engine: bing
    shortcut: bi
  
  - name: duckduckgo
    engine: duckduckgo
    shortcut: ddg
  
  - name: wikipedia
    engine: wikipedia
    shortcut: wp
    display_type: ["infobox"]
    categories: [general]
`
	
	return os.WriteFile(settingsPath, []byte(defaultSettings), 0644)
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
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
			// Try to copy from project config first
			projectConfigPath := "config/settings.yml"
			if _, err := os.Stat(projectConfigPath); err == nil {
				if err := copyFile(projectConfigPath, settingsPath); err != nil {
					log.Printf("Warning: failed to copy project config, creating default: %v", err)
					if err := createDefaultSettings(settingsPath); err != nil {
						return fmt.Errorf("failed to create default settings: %w", err)
					}
				} else {
					log.Printf("Copied project configuration to: %s", settingsPath)
				}
			} else {
				if err := createDefaultSettings(settingsPath); err != nil {
					return fmt.Errorf("failed to create default settings: %w", err)
				}
			}
		}
		
		cmd = exec.CommandContext(ctx, "docker", "run", "--name", "searxng", "-d",
			"-p", "8888:8080",
			"-v", fmt.Sprintf("%s:/etc/searxng/", configPath),
			"docker.io/searxng/searxng:latest")
		
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create SearXNG container: %w", err)
		}
		log.Println("SearXNG container created and started successfully")
	}
	
	// Wait for SearXNG to be ready
	return waitForSearXNG(ctx, "http://localhost:8888")
}

// waitForSearXNG waits for SearXNG to be ready by checking its health
func waitForSearXNG(ctx context.Context, url string) error {
	log.Println("Waiting for SearXNG to be ready...")
	
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	timeout := time.After(60 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("timeout waiting for SearXNG to be ready")
		case <-ticker.C:
			resp, err := client.Get(url)
			if err == nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				log.Println("SearXNG is ready!")
				return nil
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
	}
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