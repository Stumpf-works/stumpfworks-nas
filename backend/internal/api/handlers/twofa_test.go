// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTwoFA_Enable_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/2fa/enable", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	EnableTwoFA(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestTwoFA_Verify_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/2fa/verify", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	VerifyTwoFA(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestTwoFA_Disable_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/2fa/disable", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	DisableTwoFA(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkTwoFA_Verify(b *testing.B) {
	reqBody := map[string]interface{}{"code": "123456"}
	body, _ := json.Marshal(reqBody)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/2fa/verify", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		VerifyTwoFA(rr, req)
	}
}
