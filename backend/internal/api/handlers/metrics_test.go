// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetrics_GetCPU(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/metrics/cpu", nil)
	rr := httptest.NewRecorder()
	GetCPUMetrics(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestMetrics_GetMemory(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/metrics/memory", nil)
	rr := httptest.NewRecorder()
	GetMemoryMetrics(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestMetrics_GetDisk(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/metrics/disk", nil)
	rr := httptest.NewRecorder()
	GetDiskMetrics(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestMetrics_GetNetwork(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/metrics/network", nil)
	rr := httptest.NewRecorder()
	GetNetworkMetrics(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkMetrics_GetCPU(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/metrics/cpu", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		GetCPUMetrics(rr, req)
	}
}
