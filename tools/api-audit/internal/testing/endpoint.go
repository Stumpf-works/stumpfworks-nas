// Revision: 2025-11-16 | Author: Claude | Version: 1.0.0
package testing

import (
	"encoding/json"
	"fmt"

	"api-audit/internal/client"
	"api-audit/internal/discovery"
	"api-audit/internal/report"
)

// EndpointTester tests API endpoints
type EndpointTester struct {
	client  *client.Client
	verbose bool
}

// NewEndpointTester creates a new endpoint tester
func NewEndpointTester(c *client.Client, verbose bool) *EndpointTester {
	return &EndpointTester{
		client:  c,
		verbose: verbose,
	}
}

// TestEndpoint tests a single endpoint
func (t *EndpointTester) TestEndpoint(ep discovery.Endpoint) report.EndpointResult {
	result := report.EndpointResult{
		Method:       ep.Method,
		Path:         ep.Path,
		Description:  ep.Description,
		AuthRequired: ep.AuthRequired,
		TestResults: report.TestResults{
			Status: "PASS",
			Errors: []string{},
		},
	}

	if t.verbose {
		fmt.Printf("Testing %s %s... ", ep.Method, ep.Path)
	}

	// Test the endpoint
	var resp *client.Response
	var err error

	switch ep.Method {
	case "GET":
		resp, err = t.client.Get(ep.Path)
	case "POST":
		// For POST, try with empty body first
		resp, err = t.client.Post(ep.Path, nil)
	case "PUT":
		resp, err = t.client.Put(ep.Path, nil)
	case "DELETE":
		resp, err = t.client.Delete(ep.Path)
	default:
		result.TestResults.Status = "SKIP"
		result.TestResults.Errors = append(result.TestResults.Errors, "unsupported method")
		return result
	}

	if err != nil {
		result.TestResults.Status = "FAIL"
		result.TestResults.Errors = append(result.TestResults.Errors, err.Error())
		if t.verbose {
			fmt.Println("FAIL")
		}
		return result
	}

	result.TestResults.ResponseCode = resp.StatusCode
	result.TestResults.ResponseTimeMs = resp.Duration.Milliseconds()

	// Parse headers
	result.TestResults.Headers = make(map[string]string)
	for k := range resp.Headers {
		result.TestResults.Headers[k] = resp.Headers.Get(k)
	}

	// Try to parse JSON response
	if len(resp.Body) > 0 {
		var body interface{}
		if err := json.Unmarshal(resp.Body, &body); err == nil {
			result.TestResults.ResponseBody = body
		}
	}

	// Check if response is acceptable
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.TestResults.Status = "PASS"
		if t.verbose {
			fmt.Printf("PASS (%dms)\n", resp.Duration.Milliseconds())
		}
	} else if resp.StatusCode == 401 && ep.AuthRequired {
		// 401 is expected for auth-required endpoints
		result.TestResults.Status = "PASS"
		if t.verbose {
			fmt.Printf("PASS (401 - auth required)\n")
		}
	} else if resp.StatusCode == 404 {
		result.TestResults.Status = "FAIL"
		result.TestResults.Errors = append(result.TestResults.Errors, "endpoint not found")
		if t.verbose {
			fmt.Println("FAIL (404)")
		}
	} else if resp.StatusCode >= 500 {
		result.TestResults.Status = "FAIL"
		result.TestResults.Errors = append(result.TestResults.Errors, fmt.Sprintf("server error: %d", resp.StatusCode))
		if t.verbose {
			fmt.Printf("FAIL (500)\n")
		}
	}

	return result
}

// TestAllEndpoints tests all endpoints
func (t *EndpointTester) TestAllEndpoints(endpoints []discovery.Endpoint) []report.EndpointResult {
	results := make([]report.EndpointResult, 0, len(endpoints))

	for _, ep := range endpoints {
		result := t.TestEndpoint(ep)
		results = append(results, result)
	}

	return results
}
