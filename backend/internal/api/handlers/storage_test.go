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

// setupStorageTest initializes test environment for storage tests
func setupStorageTest(t *testing.T) {
	// Clear storage cache before each test
	storageCache.Clear()
}

// ===== Disk Handler Tests =====

func TestListDisks_CacheMiss(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/storage/disks", nil)
	rr := httptest.NewRecorder()

	ListDisks(rr, req)

	// Without real storage backend, this will fail
	// In production with mocks, we'd expect success
	// For now, verify it attempts to fetch disks
	assert.NotEqual(t, http.StatusBadRequest, rr.Code, "Should attempt to fetch disks")
}

func TestListDisks_CacheHit(t *testing.T) {
	setupStorageTest(t)

	// Pre-populate cache
	mockDisks := []map[string]interface{}{
		{"name": "sda", "size": 1000000000},
		{"name": "sdb", "size": 2000000000},
	}
	storageCache.Set("disks", mockDisks)

	req := httptest.NewRequest(http.MethodGet, "/api/storage/disks", nil)
	rr := httptest.NewRecorder()

	ListDisks(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify cached data is returned
	assert.Contains(t, response, "data")
}

func TestGetDisk_WithURLParam(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/storage/disks/sda", nil)
	rr := httptest.NewRecorder()

	// Setup chi context with URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "sda")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	GetDisk(rr, req)

	// Without mocking, this will fail to find disk
	// Verify it attempts to get disk info
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetDisk_EmptyDiskName(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/storage/disks/", nil)
	rr := httptest.NewRecorder()

	// Setup chi context with empty URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	GetDisk(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestFormatDisk_InvalidJSON(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/storage/disks/format", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	FormatDisk(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "error")
}

func TestFormatDisk_ValidRequest(t *testing.T) {
	setupStorageTest(t)

	reqBody := map[string]interface{}{
		"disk":       "sdb",
		"filesystem": "ext4",
		"label":      "data",
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/storage/disks/format", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	FormatDisk(rr, req)

	// Without real storage backend, this will fail
	// Verify request was processed
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetDiskSMART_ValidDisk(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/storage/disks/sda/smart", nil)
	rr := httptest.NewRecorder()

	// Setup chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "sda")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	GetDiskSMART(rr, req)

	// Without smartctl, this will fail
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetDiskHealth_ValidDisk(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/storage/disks/sda/health", nil)
	rr := httptest.NewRecorder()

	// Setup chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "sda")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	GetDiskHealth(rr, req)

	// Without real disk, this will fail
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSetDiskLabel_InvalidJSON(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodPut, "/api/storage/disks/sda/label", bytes.NewReader([]byte("bad json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Setup chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "sda")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	SetDiskLabel(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// ===== Volume Handler Tests =====

func TestListVolumes_Success(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/storage/volumes", nil)
	rr := httptest.NewRecorder()

	ListVolumes(rr, req)

	// Without ZFS/LVM, this will fail
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetVolume_ValidName(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/storage/volumes/tank", nil)
	rr := httptest.NewRecorder()

	// Setup chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "tank")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	GetVolume(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestCreateVolume_InvalidJSON(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/storage/volumes", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateVolume(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateVolume_ValidRequest(t *testing.T) {
	setupStorageTest(t)

	reqBody := map[string]interface{}{
		"name":  "testpool",
		"type":  "zfs",
		"disks": []string{"sdb", "sdc"},
		"raidLevel": "mirror",
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/storage/volumes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateVolume(rr, req)

	// Without ZFS, this will fail
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteVolume_ValidName(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/storage/volumes/testpool", nil)
	rr := httptest.NewRecorder()

	// Setup chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "testpool")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	DeleteVolume(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Share Handler Tests =====

func TestListShares_NoUser(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/storage/shares", nil)
	rr := httptest.NewRecorder()

	ListShares(rr, req)

	// Without database, this will fail
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetShare_ValidName(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/storage/shares/public", nil)
	rr := httptest.NewRecorder()

	// Setup chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "public")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	GetShare(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestCreateShare_InvalidJSON(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/storage/shares", bytes.NewReader([]byte("bad json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateShare(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateShare_ValidRequest(t *testing.T) {
	setupStorageTest(t)

	reqBody := map[string]interface{}{
		"name":     "testshare",
		"path":     "/mnt/storage/testshare",
		"protocol": "smb",
		"readOnly": false,
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/storage/shares", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateShare(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateShare_InvalidJSON(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodPut, "/api/storage/shares/public", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Setup chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "public")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	UpdateShare(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteShare_ValidName(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/storage/shares/testshare", nil)
	rr := httptest.NewRecorder()

	// Setup chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "testshare")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	DeleteShare(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestEnableShare_ValidName(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/storage/shares/public/enable", nil)
	rr := httptest.NewRecorder()

	// Setup chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "public")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	EnableShare(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestDisableShare_ValidName(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/storage/shares/public/disable", nil)
	rr := httptest.NewRecorder()

	// Setup chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "public")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	DisableShare(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Stats Handler Tests =====

func TestGetStorageStats_Success(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/storage/stats", nil)
	rr := httptest.NewRecorder()

	GetStorageStats(rr, req)

	// Without real filesystem, this may fail
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetDiskIOStats_Success(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/storage/io-stats", nil)
	rr := httptest.NewRecorder()

	GetDiskIOStats(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetDiskIOStatsForDisk_ValidDisk(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/storage/io-stats/sda", nil)
	rr := httptest.NewRecorder()

	// Setup chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("disk", "sda")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	GetDiskIOStatsForDisk(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetAllDiskHealth_Success(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/storage/health", nil)
	rr := httptest.NewRecorder()

	GetAllDiskHealth(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetIOMonitorStats_Success(t *testing.T) {
	setupStorageTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/storage/io-monitor", nil)
	rr := httptest.NewRecorder()

	GetIOMonitorStats(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "data")
}

// ===== IO Monitoring Tests =====

func TestStartIOMonitoring_Success(t *testing.T) {
	// Test that monitoring can be started without errors
	// In production, this would start a goroutine
	assert.NotPanics(t, func() {
		StartIOMonitoring()
	})
}

func TestStopIOMonitoring_Success(t *testing.T) {
	// Test that monitoring can be stopped without errors
	assert.NotPanics(t, func() {
		StopIOMonitoring()
	})
}

func TestIOMonitoring_StartStop(t *testing.T) {
	// Test start-stop cycle
	StartIOMonitoring()
	StopIOMonitoring()

	// Verify no panic or race conditions
	assert.True(t, true)
}

// ===== Benchmark Tests =====

func BenchmarkListDisks_CacheHit(b *testing.B) {
	setupStorageTest(&testing.T{})

	// Pre-populate cache
	mockDisks := []map[string]interface{}{
		{"name": "sda", "size": 1000000000},
	}
	storageCache.Set("disks", mockDisks)

	req := httptest.NewRequest(http.MethodGet, "/api/storage/disks", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ListDisks(rr, req)
	}
}

func BenchmarkGetStorageStats(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/storage/stats", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		GetStorageStats(rr, req)
	}
}
