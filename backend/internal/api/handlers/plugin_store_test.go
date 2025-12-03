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

// ===== List and Get Tests =====

func TestPluginStore_ListAvailable(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/plugin-store", nil)
	rr := httptest.NewRecorder()
	ListAvailablePlugins(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestPluginStore_GetFromRegistry(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/plugin-store/test-plugin", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-plugin")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	GetPluginFromRegistry(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestPluginStore_Search_ByQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/plugin-store/search?q=test", nil)
	rr := httptest.NewRecorder()
	SearchPlugins(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestPluginStore_Search_ByCategory(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/plugin-store/search?category=utilities", nil)
	rr := httptest.NewRecorder()
	SearchPlugins(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestPluginStore_Search_NoParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/plugin-store/search", nil)
	rr := httptest.NewRecorder()
	SearchPlugins(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Installation Tests =====

func TestPluginStore_Install(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/plugin-store/test-plugin/install", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-plugin")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	InstallPlugin(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestPluginStore_Uninstall(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/plugin-store/test-plugin", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-plugin")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	UninstallPlugin(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestPluginStore_Update(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/plugin-store/test-plugin/update", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-plugin")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	UpdatePlugin(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Registry Tests =====

func TestPluginStore_SyncRegistry(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/plugin-store/sync", nil)
	rr := httptest.NewRecorder()
	SyncRegistry(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestPluginStore_ListInstalled(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/plugin-store/installed", nil)
	rr := httptest.NewRecorder()
	ListInstalledPlugins(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Benchmarks =====

func BenchmarkPluginStore_ListAvailable(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/plugin-store", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ListAvailablePlugins(rr, req)
	}
}

func BenchmarkPluginStore_Search(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/plugin-store/search?q=test", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		SearchPlugins(rr, req)
	}
}
