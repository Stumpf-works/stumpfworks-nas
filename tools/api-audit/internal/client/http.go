// Revision: 2025-11-16 | Author: Claude | Version: 1.0.0
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is an HTTP client with retry and timeout capabilities
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	retries    int
}

// NewClient creates a new HTTP client
func NewClient(baseURL, token string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		retries: 3,
	}
}

// Request represents an HTTP request
type Request struct {
	Method  string
	Path    string
	Body    interface{}
	Headers map[string]string
}

// Response represents an HTTP response
type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
	Duration   time.Duration
}

// Do performs an HTTP request with retry logic
func (c *Client) Do(req Request) (*Response, error) {
	var lastErr error

	for attempt := 0; attempt <= c.retries; attempt++ {
		resp, err := c.doRequest(req)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// Don't retry on client errors (4xx)
		if resp != nil && resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return resp, nil
		}

		// Exponential backoff
		if attempt < c.retries {
			time.Sleep(time.Duration(1<<uint(attempt)) * time.Second)
		}
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", c.retries, lastErr)
}

func (c *Client) doRequest(req Request) (*Response, error) {
	url := c.baseURL + req.Path

	var bodyReader io.Reader
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	httpReq, err := http.NewRequest(req.Method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.token)
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	start := time.Now()
	httpResp, err := c.httpClient.Do(httpReq)
	duration := time.Since(start)

	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return &Response{
		StatusCode: httpResp.StatusCode,
		Body:       body,
		Headers:    httpResp.Header,
		Duration:   duration,
	}, nil
}

// Get performs a GET request
func (c *Client) Get(path string) (*Response, error) {
	return c.Do(Request{
		Method: "GET",
		Path:   path,
	})
}

// Post performs a POST request
func (c *Client) Post(path string, body interface{}) (*Response, error) {
	return c.Do(Request{
		Method: "POST",
		Path:   path,
		Body:   body,
	})
}

// Put performs a PUT request
func (c *Client) Put(path string, body interface{}) (*Response, error) {
	return c.Do(Request{
		Method: "PUT",
		Path:   path,
		Body:   body,
	})
}

// Delete performs a DELETE request
func (c *Client) Delete(path string) (*Response, error) {
	return c.Do(Request{
		Method: "DELETE",
		Path:   path,
	})
}

// ParseJSON parses JSON response body
func (r *Response) ParseJSON(v interface{}) error {
	if len(r.Body) == 0 {
		return nil
	}
	return json.Unmarshal(r.Body, v)
}
