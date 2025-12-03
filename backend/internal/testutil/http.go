// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// HTTPTest provides utilities for HTTP handler testing
type HTTPTest struct {
	t *testing.T
}

// NewHTTPTest creates a new HTTP test helper
func NewHTTPTest(t *testing.T) *HTTPTest {
	return &HTTPTest{t: t}
}

// MakeRequest creates an HTTP request with the given parameters
func (h *HTTPTest) MakeRequest(method, path string, body interface{}) *http.Request {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			h.t.Fatalf("Failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// ExecuteRequest executes a request against a handler and returns the recorder
func (h *HTTPTest) ExecuteRequest(handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler(rr, req)
	return rr
}

// AssertStatusCode checks that the response has the expected status code
func (h *HTTPTest) AssertStatusCode(rr *httptest.ResponseRecorder, expectedCode int) {
	if rr.Code != expectedCode {
		h.t.Errorf("Expected status code %d, got %d. Body: %s",
			expectedCode, rr.Code, rr.Body.String())
	}
}

// AssertJSONResponse checks that the response is valid JSON and unmarshals it
func (h *HTTPTest) AssertJSONResponse(rr *httptest.ResponseRecorder, target interface{}) {
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" && contentType != "application/json; charset=utf-8" {
		h.t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	if err := json.Unmarshal(rr.Body.Bytes(), target); err != nil {
		h.t.Fatalf("Failed to unmarshal response: %v. Body: %s", err, rr.Body.String())
	}
}

// AssertErrorResponse checks that the response contains an error message
func (h *HTTPTest) AssertErrorResponse(rr *httptest.ResponseRecorder, expectedMsg string) {
	var response map[string]interface{}
	h.AssertJSONResponse(rr, &response)

	if err, ok := response["error"].(string); ok {
		if err != expectedMsg && expectedMsg != "" {
			h.t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err)
		}
	} else {
		h.t.Error("Response does not contain error field")
	}
}

// MakeAuthenticatedRequest adds an Authorization header to the request
func (h *HTTPTest) MakeAuthenticatedRequest(method, path, token string, body interface{}) *http.Request {
	req := h.MakeRequest(method, path, body)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}
