// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestAddons_List(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/addons", nil)
	rr := httptest.NewRecorder()
	ListAddons(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestAddons_GetAddon(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/addons/test-addon", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-addon")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	GetAddon(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestAddons_GetAddon_EmptyID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/addons/", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	GetAddon(rr, req)
	assert.NotEqual(t, http.StatusOK, rr.Code)
}

func TestAddons_GetAddonStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/addons/test-addon/status", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-addon")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	GetAddonStatus(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestAddons_InstallAddon(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/addons/test-addon/install", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-addon")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	InstallAddon(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestAddons_InstallAddon_EmptyID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/addons//install", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	InstallAddon(rr, req)
	assert.NotEqual(t, http.StatusOK, rr.Code)
}

func TestAddons_UninstallAddon(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/addons/test-addon", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-addon")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	UninstallAddon(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestAddons_UninstallAddon_EmptyID(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/addons/", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	UninstallAddon(rr, req)
	assert.NotEqual(t, http.StatusOK, rr.Code)
}

func BenchmarkAddons_List(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/addons", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ListAddons(rr, req)
	}
}

func BenchmarkAddons_GetAddon(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/addons/test-addon", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-addon")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		GetAddon(rr, req)
	}
}
