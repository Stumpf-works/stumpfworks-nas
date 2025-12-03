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
)

func TestQuota_List(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/quotas", nil)
	rr := httptest.NewRecorder()
	ListQuotas(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestQuota_Get(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/quotas/1", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	GetQuota(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestQuota_Set_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/quotas", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	SetQuota(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestQuota_Delete(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/quotas/1", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	DeleteQuota(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkQuota_List(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/quotas", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ListQuotas(rr, req)
	}
}
