// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/Stumpf-works/stumpfworks-nas/internal/config"
)

func setupHealthTest(t *testing.T) {
	config.GlobalConfig = &config.Config{
		App: config.AppConfig{
			Name:        "Stumpf.Works NAS",
			Version:     "1.3.0-test",
			Environment: "test",
		},
	}
}

func TestHealthCheck_Success(t *testing.T) {
	setupHealthTest(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	HealthCheck(rr, req)

	// Check status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check Content-Type
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify response structure
	assert.Contains(t, response, "status")
	assert.Contains(t, response, "service")
	assert.Contains(t, response, "version")

	// Verify values
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "Stumpf.Works NAS", response["service"])
	assert.Equal(t, "1.3.0-test", response["version"])
}

func TestHealthCheck_MultipleRequests(t *testing.T) {
	setupHealthTest(t)

	// Test that health check is consistent across multiple requests
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rr := httptest.NewRecorder()

		HealthCheck(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "Request %d should return 200", i+1)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "ok", response["status"], "Status should be 'ok' for request %d", i+1)
	}
}

func TestIndexHandler_Success(t *testing.T) {
	setupHealthTest(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	IndexHandler(rr, req)

	// Check status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check Content-Type
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify response structure
	assert.Contains(t, response, "name")
	assert.Contains(t, response, "version")
	assert.Contains(t, response, "environment")
	assert.Contains(t, response, "api_version")
	assert.Contains(t, response, "endpoints")

	// Verify values
	assert.Equal(t, "Stumpf.Works NAS", response["name"])
	assert.Equal(t, "1.3.0-test", response["version"])
	assert.Equal(t, "test", response["environment"])
	assert.Equal(t, "v1", response["api_version"])

	// Verify endpoints structure
	endpoints, ok := response["endpoints"].(map[string]interface{})
	require.True(t, ok, "endpoints should be a map")
	assert.Contains(t, endpoints, "health")
	assert.Contains(t, endpoints, "api")
	assert.Contains(t, endpoints, "ws")
	assert.Contains(t, endpoints, "docs")

	assert.Equal(t, "/health", endpoints["health"])
	assert.Equal(t, "/api/v1", endpoints["api"])
	assert.Equal(t, "/ws", endpoints["ws"])
}

func TestIndexHandler_DifferentEnvironments(t *testing.T) {
	environments := []string{"development", "production", "staging", "test"}

	for _, env := range environments {
		t.Run("Environment_"+env, func(t *testing.T) {
			config.GlobalConfig = &config.Config{
				App: config.AppConfig{
					Name:        "Stumpf.Works NAS",
					Version:     "1.3.0",
					Environment: env,
				},
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()

			IndexHandler(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)

			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, env, response["environment"])
		})
	}
}

func TestHealthCheck_WithDifferentHTTPMethods(t *testing.T) {
	setupHealthTest(t)

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

	for _, method := range methods {
		t.Run("Method_"+method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/health", nil)
			rr := httptest.NewRecorder()

			HealthCheck(rr, req)

			// Health check should respond to all methods with 200
			assert.Equal(t, http.StatusOK, rr.Code)

			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, "ok", response["status"])
		})
	}
}

func TestIndexHandler_WithDifferentHTTPMethods(t *testing.T) {
	setupHealthTest(t)

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

	for _, method := range methods {
		t.Run("Method_"+method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/", nil)
			rr := httptest.NewRecorder()

			IndexHandler(rr, req)

			// Index handler should respond to all methods with 200
			assert.Equal(t, http.StatusOK, rr.Code)
		})
	}
}

func TestHealthCheck_ConcurrentRequests(t *testing.T) {
	setupHealthTest(t)

	// Test concurrent access
	concurrency := 10
	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rr := httptest.NewRecorder()

			HealthCheck(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < concurrency; i++ {
		<-done
	}
}

// Benchmark tests
func BenchmarkHealthCheck(b *testing.B) {
	setupHealthTest(&testing.T{})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		HealthCheck(rr, req)
	}
}

func BenchmarkIndexHandler(b *testing.B) {
	setupHealthTest(&testing.T{})

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		IndexHandler(rr, req)
	}
}

func BenchmarkHealthCheckConcurrent(b *testing.B) {
	setupHealthTest(&testing.T{})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rr := httptest.NewRecorder()
			HealthCheck(rr, req)
		}
	})
}
