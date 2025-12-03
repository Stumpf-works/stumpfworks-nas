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

func TestCompose_List(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/compose", nil)
	rr := httptest.NewRecorder()
	ListComposeProjects(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestCompose_Create_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/compose", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateComposeProject(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCompose_Start(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/compose/1/start", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	StartComposeProject(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestCompose_Stop(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/compose/1/stop", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	StopComposeProject(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkCompose_List(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/compose", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ListComposeProjects(rr, req)
	}
}
