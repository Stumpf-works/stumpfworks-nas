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

func TestLXC_List(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/lxc", nil)
	rr := httptest.NewRecorder()
	ListLXCContainers(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestLXC_Create_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/lxc", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateLXCContainer(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLXC_Start(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/lxc/1/start", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	StartLXCContainer(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestLXC_Stop(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/lxc/1/stop", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	StopLXCContainer(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkLXC_List(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/lxc", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ListLXCContainers(rr, req)
	}
}
