package searxng

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client defines the interface for SearXNG search operations
type Client interface {
	Search(ctx context.Context, req SearchRequest) (*SearchResponse, error)
}

// HTTPClient implements the Client interface using HTTP requests
type HTTPClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new SearXNG HTTP client
func NewClient(baseURL string) *HTTPClient {
	if baseURL == "" {
		baseURL = "http://localhost:8888"
	}

	return &HTTPClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Search performs a search query against the SearXNG instance
func (c *HTTPClient) Search(ctx context.Context, req SearchRequest) (*SearchResponse, error) {
	if req.Query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	searchURL := c.BaseURL + "/search"

	// Prepare form data
	formData := url.Values{}
	formData.Set("q", req.Query)
	formData.Set("format", "json")

	if req.Language != "" {
		formData.Set("language", req.Language)
	}

	if req.TimeRange != "" {
		formData.Set("time_range", string(req.TimeRange))
	}

	if len(req.Category) > 0 {
		categories := make([]string, len(req.Category))
		for i, cat := range req.Category {
			categories[i] = string(cat)
		}
		formData.Set("categories", strings.Join(categories, ","))
	}

	if req.PageNo > 0 {
		formData.Set("pageno", strconv.Itoa(req.PageNo))
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", searchURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("X-Forwarded-For", "127.0.0.1")
	httpReq.Header.Set("X-Real-IP", "127.0.0.1")

	// Execute request
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse JSON response
	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &searchResp, nil
}