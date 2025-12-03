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

func TestSystem_GetInfo(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/system/info", nil)
	rr := httptest.NewRecorder()
	GetSystemInfo(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSystem_GetHostname(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/system/hostname", nil)
	rr := httptest.NewRecorder()
	GetHostname(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSystem_SetHostname_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/api/system/hostname", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	SetHostname(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSystem_Reboot(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/system/reboot", nil)
	rr := httptest.NewRecorder()
	Reboot(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSystem_Shutdown(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/system/shutdown", nil)
	rr := httptest.NewRecorder()
	Shutdown(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkSystem_GetInfo(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/system/info", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		GetSystemInfo(rr, req)
	}
}
