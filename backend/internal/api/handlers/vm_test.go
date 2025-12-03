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

func TestVM_List(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/vm", nil)
	rr := httptest.NewRecorder()
	ListVMs(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestVM_Create_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/vm", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateVM(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestVM_Start(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/vm/1/start", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	StartVM(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestVM_Stop(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/vm/1/stop", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	StopVM(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestVM_Delete(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/vm/1", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	DeleteVM(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkVM_List(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/vm", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ListVMs(rr, req)
	}
}
