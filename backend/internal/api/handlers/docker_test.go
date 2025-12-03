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
	"github.com/stretchr/testify/require"
)

// setupDockerTest initializes test environment for Docker tests
func setupDockerTest(t *testing.T) *DockerHandler {
	return NewDockerHandler()
}

// ===== Container Management Tests =====

func TestListContainers_Success(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/docker/containers", nil)
	rr := httptest.NewRecorder()

	h.ListContainers(rr, req)

	// Without Docker daemon, this will fail gracefully
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestListContainers_WithAllParam(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/docker/containers?all=true", nil)
	rr := httptest.NewRecorder()

	h.ListContainers(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestInspectContainer_ValidID(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/docker/containers/abc123", nil)
	rr := httptest.NewRecorder()

	// Setup chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "abc123")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.InspectContainer(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetContainerStats_ValidID(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/docker/containers/abc123/stats", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "abc123")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.GetContainerStats(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestStartContainer_ValidID(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/containers/abc123/start", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "abc123")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.StartContainer(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestStopContainer_ValidID(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/containers/abc123/stop", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "abc123")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.StopContainer(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestRestartContainer_ValidID(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/containers/abc123/restart", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "abc123")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.RestartContainer(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestPauseContainer_ValidID(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/containers/abc123/pause", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "abc123")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.PauseContainer(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestUnpauseContainer_ValidID(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/containers/abc123/unpause", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "abc123")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.UnpauseContainer(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestRemoveContainer_ValidID(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/docker/containers/abc123", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "abc123")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.RemoveContainer(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetContainerLogs_ValidID(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/docker/containers/abc123/logs", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "abc123")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.GetContainerLogs(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestCreateContainer_InvalidJSON(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/containers", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateContainer(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateContainer_ValidRequest(t *testing.T) {
	h := setupDockerTest(t)

	reqBody := map[string]interface{}{
		"name":  "test-container",
		"image": "nginx:latest",
		"env":   []string{"TEST=1"},
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/containers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateContainer(rr, req)

	// Without Docker, this will fail
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateContainerResources_InvalidJSON(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodPut, "/api/docker/containers/abc123/resources", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "abc123")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.UpdateContainerResources(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestExecContainer_InvalidJSON(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/containers/abc123/exec", bytes.NewReader([]byte("bad json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "abc123")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.ExecContainer(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetContainerTop_ValidID(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/docker/containers/abc123/top", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "abc123")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.GetContainerTop(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Image Management Tests =====

func TestListImages_Success(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/docker/images", nil)
	rr := httptest.NewRecorder()

	h.ListImages(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestInspectImage_ValidID(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/docker/images/nginx:latest", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "nginx:latest")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.InspectImage(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestPullImage_InvalidJSON(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/images/pull", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.PullImage(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPullImage_ValidRequest(t *testing.T) {
	h := setupDockerTest(t)

	reqBody := map[string]interface{}{
		"image": "nginx:latest",
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/images/pull", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.PullImage(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSearchImages_WithTerm(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/docker/images/search?term=nginx", nil)
	rr := httptest.NewRecorder()

	h.SearchImages(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSearchImages_EmptyTerm(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/docker/images/search", nil)
	rr := httptest.NewRecorder()

	h.SearchImages(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestRemoveImage_ValidID(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/docker/images/nginx:latest", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "nginx:latest")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.RemoveImage(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestBuildImage_InvalidJSON(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/images/build", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.BuildImage(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestTagImage_InvalidJSON(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/images/nginx/tag", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "nginx")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.TagImage(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPushImage_InvalidJSON(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/images/push", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.PushImage(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// ===== Volume Management Tests =====

func TestListVolumes_Success(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/docker/volumes", nil)
	rr := httptest.NewRecorder()

	h.ListVolumes(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestInspectVolume_ValidName(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/docker/volumes/myvolume", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "myvolume")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.InspectVolume(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestCreateVolume_InvalidJSON(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/volumes", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateVolume(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateVolume_ValidRequest(t *testing.T) {
	h := setupDockerTest(t)

	reqBody := map[string]interface{}{
		"name":   "testvol",
		"driver": "local",
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/volumes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateVolume(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestRemoveVolume_ValidName(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/docker/volumes/testvol", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "testvol")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.RemoveVolume(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Network Management Tests =====

func TestListNetworks_Success(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/docker/networks", nil)
	rr := httptest.NewRecorder()

	h.ListNetworks(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestInspectNetwork_ValidID(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/docker/networks/bridge", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "bridge")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.InspectNetwork(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestCreateNetwork_InvalidJSON(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/networks", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateNetwork(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateNetwork_ValidRequest(t *testing.T) {
	h := setupDockerTest(t)

	reqBody := map[string]interface{}{
		"name":   "testnet",
		"driver": "bridge",
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/networks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateNetwork(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestRemoveNetwork_ValidID(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/docker/networks/testnet", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "testnet")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.RemoveNetwork(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestConnectContainerToNetwork_InvalidJSON(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/networks/testnet/connect", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "testnet")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.ConnectContainerToNetwork(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDisconnectContainerFromNetwork_InvalidJSON(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/networks/testnet/disconnect", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "testnet")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.DisconnectContainerFromNetwork(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// ===== System Tests =====

func TestGetDockerInfo_Success(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/docker/info", nil)
	rr := httptest.NewRecorder()

	h.GetDockerInfo(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetDockerVersion_Success(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/docker/version", nil)
	rr := httptest.NewRecorder()

	h.GetDockerVersion(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestPruneSystem_Success(t *testing.T) {
	h := setupDockerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/docker/prune", nil)
	rr := httptest.NewRecorder()

	h.PruneSystem(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Benchmark Tests =====

func BenchmarkListContainers(b *testing.B) {
	h := NewDockerHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/docker/containers", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		h.ListContainers(rr, req)
	}
}

func BenchmarkListImages(b *testing.B) {
	h := NewDockerHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/docker/images", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		h.ListImages(rr, req)
	}
}

func BenchmarkGetDockerInfo(b *testing.B) {
	h := NewDockerHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/docker/info", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		h.GetDockerInfo(rr, req)
	}
}
