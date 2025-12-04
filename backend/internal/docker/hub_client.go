// Package docker provides Docker management functionality
package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	// DefaultHubURL is the default Stumpfworks Hub URL
	DefaultHubURL = "https://hub.stumpfworks.de"

	// FallbackHubURL is used if the primary hub is unavailable
	FallbackHubURL = "http://localhost:8090"
)

// HubClient manages communication with Stumpfworks Hub
type HubClient struct {
	baseURL    string
	httpClient *http.Client
}

// HubResponse wraps Hub API responses
type HubResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   string          `json:"error,omitempty"`
}

// NewHubClient creates a new Hub client
func NewHubClient() *HubClient {
	// Check if custom hub URL is set
	hubURL := os.Getenv("STUMPFWORKS_HUB_URL")
	if hubURL == "" {
		hubURL = DefaultHubURL
	}

	return &HubClient{
		baseURL: hubURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ListTemplates fetches all templates from the Hub
func (c *HubClient) ListTemplates(ctx context.Context) ([]ComposeTemplate, error) {
	url := fmt.Sprintf("%s/api/v1/templates", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch templates from hub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hub returned status %d", resp.StatusCode)
	}

	var hubResp HubResponse
	if err := json.NewDecoder(resp.Body).Decode(&hubResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !hubResp.Success {
		return nil, fmt.Errorf("hub error: %s", hubResp.Error)
	}

	var templates []ComposeTemplate
	if err := json.Unmarshal(hubResp.Data, &templates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal templates: %w", err)
	}

	return templates, nil
}

// GetTemplate fetches a specific template from the Hub
func (c *HubClient) GetTemplate(ctx context.Context, id string) (*ComposeTemplate, error) {
	url := fmt.Sprintf("%s/api/v1/templates/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch template from hub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("template not found: %s", id)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hub returned status %d", resp.StatusCode)
	}

	var hubResp HubResponse
	if err := json.NewDecoder(resp.Body).Decode(&hubResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !hubResp.Success {
		return nil, fmt.Errorf("hub error: %s", hubResp.Error)
	}

	var template ComposeTemplate
	if err := json.Unmarshal(hubResp.Data, &template); err != nil {
		return nil, fmt.Errorf("failed to unmarshal template: %w", err)
	}

	return &template, nil
}

// GetTemplateCategories fetches all template categories from the Hub
func (c *HubClient) GetTemplateCategories(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("%s/api/v1/templates/categories", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories from hub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hub returned status %d", resp.StatusCode)
	}

	var hubResp HubResponse
	if err := json.NewDecoder(resp.Body).Decode(&hubResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !hubResp.Success {
		return nil, fmt.Errorf("hub error: %s", hubResp.Error)
	}

	var categories []string
	if err := json.Unmarshal(hubResp.Data, &categories); err != nil {
		return nil, fmt.Errorf("failed to unmarshal categories: %w", err)
	}

	return categories, nil
}

// SearchTemplates searches for templates on the Hub
func (c *HubClient) SearchTemplates(ctx context.Context, query string) ([]ComposeTemplate, error) {
	url := fmt.Sprintf("%s/api/v1/templates/search?q=%s", c.baseURL, query)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search templates on hub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hub returned status %d", resp.StatusCode)
	}

	var hubResp HubResponse
	if err := json.NewDecoder(resp.Body).Decode(&hubResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !hubResp.Success {
		return nil, fmt.Errorf("hub error: %s", hubResp.Error)
	}

	var templates []ComposeTemplate
	if err := json.Unmarshal(hubResp.Data, &templates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal templates: %w", err)
	}

	return templates, nil
}

// Global hub client instance
var globalHubClient *HubClient

// GetHubClient returns the global hub client instance
func GetHubClient() *HubClient {
	if globalHubClient == nil {
		globalHubClient = NewHubClient()
	}
	return globalHubClient
}


// GetBaseURL returns the Hub base URL
func (c *HubClient) GetBaseURL() string {
	return c.baseURL
}
