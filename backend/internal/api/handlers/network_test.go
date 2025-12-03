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

func TestNetwork_GetInterfaces(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/network/interfaces", nil)
	rr := httptest.NewRecorder()
	GetNetworkInterfaces(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestNetwork_UpdateInterface_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/api/network/interfaces/eth0", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	UpdateNetworkInterface(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestNetwork_GetRoutes(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/network/routes", nil)
	rr := httptest.NewRecorder()
	GetRoutes(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestNetwork_GetDNS(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/network/dns", nil)
	rr := httptest.NewRecorder()
	GetDNSConfig(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkNetwork_GetInterfaces(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/network/interfaces", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		GetNetworkInterfaces(rr, req)
	}
}
