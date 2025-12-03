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

func TestVPN_GetConfig(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/vpn/config", nil)
	rr := httptest.NewRecorder()
	GetVPNConfig(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestVPN_UpdateConfig_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/api/vpn/config", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	UpdateVPNConfig(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestVPN_GetStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/vpn/status", nil)
	rr := httptest.NewRecorder()
	GetVPNStatus(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestVPN_ListClients(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/vpn/clients", nil)
	rr := httptest.NewRecorder()
	ListVPNClients(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkVPN_GetStatus(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/vpn/status", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		GetVPNStatus(rr, req)
	}
}
