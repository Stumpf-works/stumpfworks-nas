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

func setupPluginHandler() *PluginHandler {
	return NewPluginHandler()
}

// ===== List and Get Tests =====

func TestPlugin_List(t *testing.T) {
	handler := setupPluginHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/plugins", nil)
	rr := httptest.NewRecorder()
	handler.ListPlugins(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestPlugin_Get(t *testing.T) {
	handler := setupPluginHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/plugins/test-plugin", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-plugin")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.GetPlugin(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Installation Tests =====

func TestPlugin_Install_InvalidJSON(t *testing.T) {
	handler := setupPluginHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/plugins/install", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.InstallPlugin(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPlugin_Install_EmptySourcePath(t *testing.T) {
	handler := setupPluginHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/plugins/install", bytes.NewReader([]byte(`{"sourcePath":""}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.InstallPlugin(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPlugin_Install_ValidRequest(t *testing.T) {
	handler := setupPluginHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/plugins/install", bytes.NewReader([]byte(`{"sourcePath":"/path/to/plugin"}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.InstallPlugin(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestPlugin_Uninstall(t *testing.T) {
	handler := setupPluginHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/plugins/test-plugin", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-plugin")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.UninstallPlugin(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Enable/Disable Tests =====

func TestPlugin_Enable(t *testing.T) {
	handler := setupPluginHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/plugins/test-plugin/enable", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-plugin")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.EnablePlugin(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestPlugin_Disable(t *testing.T) {
	handler := setupPluginHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/plugins/test-plugin/disable", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-plugin")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.DisablePlugin(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Configuration Tests =====

func TestPlugin_UpdateConfig_InvalidJSON(t *testing.T) {
	handler := setupPluginHandler()
	req := httptest.NewRequest(http.MethodPut, "/api/plugins/test-plugin/config", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-plugin")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.UpdatePluginConfig(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPlugin_UpdateConfig_ValidRequest(t *testing.T) {
	handler := setupPluginHandler()
	req := httptest.NewRequest(http.MethodPut, "/api/plugins/test-plugin/config", bytes.NewReader([]byte(`{"config":{"key":"value"}}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-plugin")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.UpdatePluginConfig(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Runtime Control Tests =====

func TestPlugin_Start(t *testing.T) {
	handler := setupPluginHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/plugins/test-plugin/start", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-plugin")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.StartPlugin(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestPlugin_Stop(t *testing.T) {
	handler := setupPluginHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/plugins/test-plugin/stop", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-plugin")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.StopPlugin(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestPlugin_Restart(t *testing.T) {
	handler := setupPluginHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/plugins/test-plugin/restart", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-plugin")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.RestartPlugin(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestPlugin_GetStatus(t *testing.T) {
	handler := setupPluginHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/plugins/test-plugin/status", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-plugin")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.GetPluginStatus(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestPlugin_ListRunningPlugins(t *testing.T) {
	handler := setupPluginHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/plugins/running", nil)
	rr := httptest.NewRecorder()
	handler.ListRunningPlugins(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Benchmarks =====

func BenchmarkPlugin_List(b *testing.B) {
	handler := setupPluginHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/plugins", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ListPlugins(rr, req)
	}
}

func BenchmarkPlugin_GetStatus(b *testing.B) {
	handler := setupPluginHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/plugins/test-plugin/status", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-plugin")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.GetPluginStatus(rr, req)
	}
}
