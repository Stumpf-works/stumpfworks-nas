// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdates_CheckForUpdates(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/updates/check", nil)
	rr := httptest.NewRecorder()
	CheckForUpdates(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestUpdates_ApplyUpdate_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/updates/apply", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	ApplyUpdate(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdates_GetHistory(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/updates/history", nil)
	rr := httptest.NewRecorder()
	GetUpdateHistory(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkUpdates_CheckForUpdates(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/updates/check", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		CheckForUpdates(rr, req)
	}
}
