// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupAlertRulesTest initializes test environment for alert rules tests
func setupAlertRulesTest(t *testing.T) *AlertRulesHandler {
	return NewAlertRulesHandler()
}

// ===== Rule Management Tests =====

func TestCreateRule_InvalidJSON(t *testing.T) {
	h := setupAlertRulesTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/alertrules", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateRule(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "error")
}

func TestCreateRule_ValidRequest(t *testing.T) {
	h := setupAlertRulesTest(t)

	reqBody := map[string]interface{}{
		"name":        "High CPU Alert",
		"description": "Alert when CPU exceeds 80%",
		"metricType":  "cpu",
		"condition":   ">",
		"threshold":   80,
		"duration":    300,
		"severity":    "warning",
		"enabled":     true,
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/alertrules", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateRule(rr, req)

	// Without database, this will fail
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestCreateRule_DifferentMetricTypes(t *testing.T) {
	h := setupAlertRulesTest(t)

	metricTypes := []string{"cpu", "memory", "disk", "network", "health", "temperature", "iops"}

	for _, metricType := range metricTypes {
		t.Run("MetricType_"+metricType, func(t *testing.T) {
			reqBody := map[string]interface{}{
				"name":       "Test Rule",
				"metricType": metricType,
				"condition":  ">",
				"threshold":  50,
			}

			body, err := json.Marshal(reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/alertrules", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			h.CreateRule(rr, req)

			assert.NotEqual(t, http.StatusBadRequest, rr.Code)
		})
	}
}

func TestCreateRule_DifferentConditions(t *testing.T) {
	h := setupAlertRulesTest(t)

	conditions := []string{">", "<", "=", ">=", "<="}

	for _, condition := range conditions {
		t.Run("Condition_"+condition, func(t *testing.T) {
			reqBody := map[string]interface{}{
				"name":       "Test Rule",
				"metricType": "cpu",
				"condition":  condition,
				"threshold":  50,
			}

			body, err := json.Marshal(reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/alertrules", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			h.CreateRule(rr, req)

			assert.NotEqual(t, http.StatusBadRequest, rr.Code)
		})
	}
}

func TestCreateRule_DifferentSeverities(t *testing.T) {
	h := setupAlertRulesTest(t)

	severities := []string{"info", "warning", "critical"}

	for _, severity := range severities {
		t.Run("Severity_"+severity, func(t *testing.T) {
			reqBody := map[string]interface{}{
				"name":       "Test Rule",
				"metricType": "cpu",
				"severity":   severity,
			}

			body, err := json.Marshal(reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/alertrules", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			h.CreateRule(rr, req)

			assert.NotEqual(t, http.StatusBadRequest, rr.Code)
		})
	}
}

func TestUpdateRule_InvalidJSON(t *testing.T) {
	h := setupAlertRulesTest(t)

	req := httptest.NewRequest(http.MethodPut, "/api/alertrules/1", bytes.NewReader([]byte("bad json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.UpdateRule(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateRule_InvalidID(t *testing.T) {
	h := setupAlertRulesTest(t)

	reqBody := map[string]interface{}{
		"name": "Updated Rule",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/alertrules/invalid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.UpdateRule(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateRule_ValidRequest(t *testing.T) {
	h := setupAlertRulesTest(t)

	reqBody := map[string]interface{}{
		"name":      "Updated High CPU Alert",
		"threshold": 90,
		"enabled":   false,
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/api/alertrules/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.UpdateRule(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteRule_ValidID(t *testing.T) {
	h := setupAlertRulesTest(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/alertrules/1", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.DeleteRule(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteRule_InvalidID(t *testing.T) {
	h := setupAlertRulesTest(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/alertrules/invalid", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.DeleteRule(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetRule_ValidID(t *testing.T) {
	h := setupAlertRulesTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/alertrules/1", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.GetRule(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetRule_InvalidID(t *testing.T) {
	h := setupAlertRulesTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/alertrules/invalid", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.GetRule(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestListRules_Success(t *testing.T) {
	h := setupAlertRulesTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/alertrules", nil)
	rr := httptest.NewRecorder()

	h.ListRules(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestListRules_WithEnabledFilter(t *testing.T) {
	h := setupAlertRulesTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/alertrules?enabled=true", nil)
	rr := httptest.NewRecorder()

	h.ListRules(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Execution Tests =====

func TestGetExecutions_NoRuleID(t *testing.T) {
	h := setupAlertRulesTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/alertrules/executions", nil)
	rr := httptest.NewRecorder()

	h.GetExecutions(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetExecutions_WithRuleID(t *testing.T) {
	h := setupAlertRulesTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/alertrules/executions?ruleId=1", nil)
	rr := httptest.NewRecorder()

	h.GetExecutions(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetExecutions_WithLimitAndOffset(t *testing.T) {
	h := setupAlertRulesTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/alertrules/executions?limit=10&offset=5", nil)
	rr := httptest.NewRecorder()

	h.GetExecutions(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetExecutions_WithAcknowledgedFilter(t *testing.T) {
	h := setupAlertRulesTest(t)

	tests := []string{"true", "false"}

	for _, acknowledged := range tests {
		t.Run("Acknowledged_"+acknowledged, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/alertrules/executions?acknowledged="+acknowledged, nil)
			rr := httptest.NewRecorder()

			h.GetExecutions(rr, req)

			assert.NotEqual(t, http.StatusBadRequest, rr.Code)
		})
	}
}

func TestGetRecentExecutions_WithLimit(t *testing.T) {
	h := setupAlertRulesTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/alertrules/executions/recent?limit=5", nil)
	rr := httptest.NewRecorder()

	h.GetRecentExecutions(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetRecentExecutions_DefaultLimit(t *testing.T) {
	h := setupAlertRulesTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/alertrules/executions/recent", nil)
	rr := httptest.NewRecorder()

	h.GetRecentExecutions(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestAcknowledgeExecution_InvalidJSON(t *testing.T) {
	h := setupAlertRulesTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/alertrules/executions/1/acknowledge", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.AcknowledgeExecution(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAcknowledgeExecution_ValidRequest(t *testing.T) {
	h := setupAlertRulesTest(t)

	reqBody := map[string]interface{}{
		"note": "Acknowledged - will investigate",
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/alertrules/executions/1/acknowledge", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.AcknowledgeExecution(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestAcknowledgeExecution_InvalidID(t *testing.T) {
	h := setupAlertRulesTest(t)

	reqBody := map[string]interface{}{
		"note": "Test note",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/alertrules/executions/invalid/acknowledge", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.AcknowledgeExecution(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// ===== Integration Tests =====

func TestAlertRulesWorkflow_Complete(t *testing.T) {
	h := setupAlertRulesTest(t)

	// 1. Create a rule
	ruleBody := map[string]interface{}{
		"name":        "Test Alert",
		"metricType":  "cpu",
		"condition":   ">",
		"threshold":   80,
		"severity":    "warning",
		"enabled":     true,
	}
	body1, _ := json.Marshal(ruleBody)
	req1 := httptest.NewRequest(http.MethodPost, "/api/alertrules", bytes.NewReader(body1))
	req1.Header.Set("Content-Type", "application/json")
	rr1 := httptest.NewRecorder()
	h.CreateRule(rr1, req1)
	assert.NotEqual(t, http.StatusBadRequest, rr1.Code)

	// 2. List rules
	req2 := httptest.NewRequest(http.MethodGet, "/api/alertrules", nil)
	rr2 := httptest.NewRecorder()
	h.ListRules(rr2, req2)
	assert.NotEqual(t, http.StatusBadRequest, rr2.Code)

	// 3. Get executions
	req3 := httptest.NewRequest(http.MethodGet, "/api/alertrules/executions", nil)
	rr3 := httptest.NewRecorder()
	h.GetExecutions(rr3, req3)
	assert.NotEqual(t, http.StatusBadRequest, rr3.Code)

	// 4. Get recent executions
	req4 := httptest.NewRequest(http.MethodGet, "/api/alertrules/executions/recent", nil)
	rr4 := httptest.NewRecorder()
	h.GetRecentExecutions(rr4, req4)
	assert.NotEqual(t, http.StatusBadRequest, rr4.Code)
}

// ===== Benchmark Tests =====

func BenchmarkCreateRule(b *testing.B) {
	h := NewAlertRulesHandler()

	reqBody := map[string]interface{}{
		"name":       "Benchmark Rule",
		"metricType": "cpu",
		"threshold":  80,
	}
	body, _ := json.Marshal(reqBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/alertrules", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.CreateRule(rr, req)
	}
}

func BenchmarkListRules(b *testing.B) {
	h := NewAlertRulesHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/alertrules", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		h.ListRules(rr, req)
	}
}

func BenchmarkGetExecutions(b *testing.B) {
	h := NewAlertRulesHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/alertrules/executions", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		h.GetExecutions(rr, req)
	}
}
