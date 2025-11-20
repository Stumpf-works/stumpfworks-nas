package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the API client for StumpfWorks NAS
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

// Response represents a standard API response
type Response struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   *ErrorInfo      `json:"error,omitempty"`
}

// ErrorInfo contains error details
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewClient creates a new API client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Get performs a GET request
func (c *Client) Get(endpoint string, result interface{}) error {
	return c.request("GET", endpoint, nil, result)
}

// Post performs a POST request
func (c *Client) Post(endpoint string, body interface{}, result interface{}) error {
	return c.request("POST", endpoint, body, result)
}

// Put performs a PUT request
func (c *Client) Put(endpoint string, body interface{}, result interface{}) error {
	return c.request("PUT", endpoint, body, result)
}

// Delete performs a DELETE request
func (c *Client) Delete(endpoint string, result interface{}) error {
	return c.request("DELETE", endpoint, nil, result)
}

// request performs an HTTP request
func (c *Client) request(method, endpoint string, body interface{}, result interface{}) error {
	url := c.BaseURL + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !apiResp.Success {
		if apiResp.Error != nil {
			return fmt.Errorf("API error: %s", apiResp.Error.Message)
		}
		return fmt.Errorf("API request failed")
	}

	if result != nil && apiResp.Data != nil {
		if err := json.Unmarshal(apiResp.Data, result); err != nil {
			return fmt.Errorf("failed to unmarshal result: %w", err)
		}
	}

	return nil
}

// Health checks if the server is healthy
func (c *Client) Health() (map[string]interface{}, error) {
	var result map[string]interface{}
	err := c.Get("/api/health", &result)
	return result, err
}

// GetUsers retrieves all users
func (c *Client) GetUsers() ([]map[string]interface{}, error) {
	var users []map[string]interface{}
	err := c.Get("/api/users", &users)
	return users, err
}

// CreateUser creates a new user
func (c *Client) CreateUser(username, password, role string) error {
	body := map[string]interface{}{
		"username": username,
		"password": password,
		"role":     role,
	}
	return c.Post("/api/users", body, nil)
}

// DeleteUser deletes a user
func (c *Client) DeleteUser(username string) error {
	return c.Delete(fmt.Sprintf("/api/users/%s", username), nil)
}

// GetBackups retrieves all backups
func (c *Client) GetBackups() ([]map[string]interface{}, error) {
	var backups []map[string]interface{}
	err := c.Get("/api/backups", &backups)
	return backups, err
}

// CreateBackup creates a new backup
func (c *Client) CreateBackup() error {
	return c.Post("/api/backups", nil, nil)
}

// GetMetrics retrieves system metrics
func (c *Client) GetMetrics() (map[string]interface{}, error) {
	var metrics map[string]interface{}
	err := c.Get("/api/metrics", &metrics)
	return metrics, err
}
