// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeepalived_GetStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/ha/keepalived/status", nil)
	rr := httptest.NewRecorder()
	GetKeepalivedStatus(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestKeepalived_Configure_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/ha/keepalived/configure", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	ConfigureKeepalived(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkKeepalived_GetStatus(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/ha/keepalived/status", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		GetKeepalivedStatus(rr, req)
	}
}
