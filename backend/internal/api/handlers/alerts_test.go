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

func setupAlertsTest(t *testing.T) *AlertHandler {
	return NewAlertHandler()
}

func TestAlerts_GetConfig(t *testing.T) {
	h := setupAlertsTest(t)
	req := httptest.NewRequest(http.MethodGet, "/api/alerts/config", nil)
	rr := httptest.NewRecorder()
	h.GetConfig(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestAlerts_UpdateConfig_InvalidJSON(t *testing.T) {
	h := setupAlertsTest(t)
	req := httptest.NewRequest(http.MethodPut, "/api/alerts/config", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.UpdateConfig(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAlerts_UpdateConfig_Valid(t *testing.T) {
	h := setupAlertsTest(t)
	reqBody := map[string]interface{}{"emailEnabled": true}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/alerts/config", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.UpdateConfig(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkAlerts_GetConfig(b *testing.B) {
	h := NewAlertHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/alerts/config", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		h.GetConfig(rr, req)
	}
}
