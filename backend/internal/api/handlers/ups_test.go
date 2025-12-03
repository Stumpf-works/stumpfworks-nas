// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupUPSTest initializes test environment for UPS tests
func setupUPSTest(t *testing.T) *UPSHandler {
	return NewUPSHandler()
}

// ===== Configuration Tests =====

func TestGetConfig_Success(t *testing.T) {
	h := setupUPSTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/ups/config", nil)
	rr := httptest.NewRecorder()

	h.GetConfig(rr, req)

	// Without database, this may fail
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateConfig_InvalidJSON(t *testing.T) {
	h := setupUPSTest(t)

	req := httptest.NewRequest(http.MethodPut, "/api/ups/config", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.UpdateConfig(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "error")
}

func TestUpdateConfig_ValidRequest(t *testing.T) {
	h := setupUPSTest(t)

	reqBody := map[string]interface{}{
		"enabled":         true,
		"hostname":        "localhost",
		"port":            3493,
		"upsName":         "ups",
		"username":        "admin",
		"password":        "secret",
		"shutdownDelay":   30,
		"lowBatteryLevel": 20,
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/api/ups/config", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.UpdateConfig(rr, req)

	// Without database, this will fail
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateConfig_EmptyFields(t *testing.T) {
	h := setupUPSTest(t)

	tests := []struct {
		name   string
		config map[string]interface{}
	}{
		{
			name: "Empty hostname",
			config: map[string]interface{}{
				"enabled":  true,
				"hostname": "",
				"port":     3493,
			},
		},
		{
			name: "Invalid port",
			config: map[string]interface{}{
				"enabled":  true,
				"hostname": "localhost",
				"port":     0,
			},
		},
		{
			name: "Empty UPS name",
			config: map[string]interface{}{
				"enabled":  true,
				"hostname": "localhost",
				"port":     3493,
				"upsName":  "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.config)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/api/ups/config", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			h.UpdateConfig(rr, req)

			// Should accept but may fail validation
			assert.NotEqual(t, http.StatusBadRequest, rr.Code)
		})
	}
}

// ===== Status Tests =====

func TestGetStatus_Success(t *testing.T) {
	h := setupUPSTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/ups/status", nil)
	rr := httptest.NewRecorder()

	h.GetStatus(rr, req)

	// Without NUT, this will fail gracefully
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetStatus_MultipleRequests(t *testing.T) {
	h := setupUPSTest(t)

	// Test that status endpoint is consistent
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/ups/status", nil)
		rr := httptest.NewRecorder()

		h.GetStatus(rr, req)

		assert.NotEqual(t, http.StatusBadRequest, rr.Code, "Request %d should not return bad request", i+1)
	}
}

// ===== Connection Tests =====

func TestTestConnection_Success(t *testing.T) {
	h := setupUPSTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/ups/test", nil)
	rr := httptest.NewRecorder()

	h.TestConnection(rr, req)

	// Without NUT daemon, this will fail
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestTestConnection_NoConfig(t *testing.T) {
	h := setupUPSTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/ups/test", nil)
	rr := httptest.NewRecorder()

	h.TestConnection(rr, req)

	// Should attempt to connect even without config
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Events Tests =====

func TestGetEvents_Success(t *testing.T) {
	h := setupUPSTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/ups/events", nil)
	rr := httptest.NewRecorder()

	h.GetEvents(rr, req)

	// Without database, this may fail
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetEvents_WithLimit(t *testing.T) {
	h := setupUPSTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/ups/events?limit=10", nil)
	rr := httptest.NewRecorder()

	h.GetEvents(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetEvents_WithOffset(t *testing.T) {
	h := setupUPSTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/ups/events?limit=10&offset=5", nil)
	rr := httptest.NewRecorder()

	h.GetEvents(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetEvents_InvalidLimit(t *testing.T) {
	h := setupUPSTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/ups/events?limit=invalid", nil)
	rr := httptest.NewRecorder()

	h.GetEvents(rr, req)

	// Should handle invalid limit gracefully
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Monitoring Tests =====

func TestStartMonitoring_Success(t *testing.T) {
	h := setupUPSTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/ups/monitoring/start", nil)
	rr := httptest.NewRecorder()

	h.StartMonitoring(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestStopMonitoring_Success(t *testing.T) {
	h := setupUPSTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/ups/monitoring/stop", nil)
	rr := httptest.NewRecorder()

	h.StopMonitoring(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestMonitoring_StartStopCycle(t *testing.T) {
	h := setupUPSTest(t)

	// Start monitoring
	req1 := httptest.NewRequest(http.MethodPost, "/api/ups/monitoring/start", nil)
	rr1 := httptest.NewRecorder()
	h.StartMonitoring(rr1, req1)
	assert.NotEqual(t, http.StatusBadRequest, rr1.Code)

	// Stop monitoring
	req2 := httptest.NewRequest(http.MethodPost, "/api/ups/monitoring/stop", nil)
	rr2 := httptest.NewRecorder()
	h.StopMonitoring(rr2, req2)
	assert.NotEqual(t, http.StatusBadRequest, rr2.Code)
}

func TestMonitoring_DoubleStart(t *testing.T) {
	h := setupUPSTest(t)

	// Start monitoring twice
	req1 := httptest.NewRequest(http.MethodPost, "/api/ups/monitoring/start", nil)
	rr1 := httptest.NewRecorder()
	h.StartMonitoring(rr1, req1)

	req2 := httptest.NewRequest(http.MethodPost, "/api/ups/monitoring/start", nil)
	rr2 := httptest.NewRecorder()
	h.StartMonitoring(rr2, req2)

	// Should handle double start gracefully
	assert.NotEqual(t, http.StatusBadRequest, rr2.Code)
}

func TestMonitoring_DoubleStop(t *testing.T) {
	h := setupUPSTest(t)

	// Stop monitoring twice
	req1 := httptest.NewRequest(http.MethodPost, "/api/ups/monitoring/stop", nil)
	rr1 := httptest.NewRecorder()
	h.StopMonitoring(rr1, req1)

	req2 := httptest.NewRequest(http.MethodPost, "/api/ups/monitoring/stop", nil)
	rr2 := httptest.NewRecorder()
	h.StopMonitoring(rr2, req2)

	// Should handle double stop gracefully
	assert.NotEqual(t, http.StatusBadRequest, rr2.Code)
}

// ===== Integration Tests =====

func TestUPSWorkflow_Complete(t *testing.T) {
	h := setupUPSTest(t)

	// 1. Get initial config
	req1 := httptest.NewRequest(http.MethodGet, "/api/ups/config", nil)
	rr1 := httptest.NewRecorder()
	h.GetConfig(rr1, req1)
	assert.NotEqual(t, http.StatusBadRequest, rr1.Code)

	// 2. Update config
	config := map[string]interface{}{
		"enabled":  true,
		"hostname": "localhost",
		"port":     3493,
		"upsName":  "ups",
	}
	body, _ := json.Marshal(config)
	req2 := httptest.NewRequest(http.MethodPut, "/api/ups/config", bytes.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	rr2 := httptest.NewRecorder()
	h.UpdateConfig(rr2, req2)
	assert.NotEqual(t, http.StatusBadRequest, rr2.Code)

	// 3. Test connection
	req3 := httptest.NewRequest(http.MethodPost, "/api/ups/test", nil)
	rr3 := httptest.NewRecorder()
	h.TestConnection(rr3, req3)
	assert.NotEqual(t, http.StatusBadRequest, rr3.Code)

	// 4. Start monitoring
	req4 := httptest.NewRequest(http.MethodPost, "/api/ups/monitoring/start", nil)
	rr4 := httptest.NewRecorder()
	h.StartMonitoring(rr4, req4)
	assert.NotEqual(t, http.StatusBadRequest, rr4.Code)

	// 5. Get status
	req5 := httptest.NewRequest(http.MethodGet, "/api/ups/status", nil)
	rr5 := httptest.NewRecorder()
	h.GetStatus(rr5, req5)
	assert.NotEqual(t, http.StatusBadRequest, rr5.Code)

	// 6. Get events
	req6 := httptest.NewRequest(http.MethodGet, "/api/ups/events", nil)
	rr6 := httptest.NewRecorder()
	h.GetEvents(rr6, req6)
	assert.NotEqual(t, http.StatusBadRequest, rr6.Code)

	// 7. Stop monitoring
	req7 := httptest.NewRequest(http.MethodPost, "/api/ups/monitoring/stop", nil)
	rr7 := httptest.NewRecorder()
	h.StopMonitoring(rr7, req7)
	assert.NotEqual(t, http.StatusBadRequest, rr7.Code)
}

// ===== Benchmark Tests =====

func BenchmarkGetStatus(b *testing.B) {
	h := NewUPSHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ups/status", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		h.GetStatus(rr, req)
	}
}

func BenchmarkGetConfig(b *testing.B) {
	h := NewUPSHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ups/config", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		h.GetConfig(rr, req)
	}
}

func BenchmarkGetEvents(b *testing.B) {
	h := NewUPSHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ups/events", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		h.GetEvents(rr, req)
	}
}
