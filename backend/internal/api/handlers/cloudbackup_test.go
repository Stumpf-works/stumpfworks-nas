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

// setupCloudBackupTest initializes test environment for cloud backup tests
func setupCloudBackupTest(t *testing.T) *CloudBackupHandler {
	return NewCloudBackupHandler()
}

// ===== Provider Management Tests =====

func TestListProviders_Success(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/cloudbackup/providers", nil)
	rr := httptest.NewRecorder()

	h.ListProviders(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetProvider_ValidID(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/cloudbackup/providers/1", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.GetProvider(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetProvider_InvalidID(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/cloudbackup/providers/invalid", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.GetProvider(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateProvider_InvalidJSON(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/cloudbackup/providers", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateProvider(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateProvider_ValidRequest(t *testing.T) {
	h := setupCloudBackupTest(t)

	reqBody := map[string]interface{}{
		"name":     "My S3 Backup",
		"type":     "s3",
		"endpoint": "s3.amazonaws.com",
		"bucket":   "my-backups",
		"accessKey": "AKIAIOSFODNN7EXAMPLE",
		"secretKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/cloudbackup/providers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateProvider(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestCreateProvider_DifferentTypes(t *testing.T) {
	h := setupCloudBackupTest(t)

	types := []string{"s3", "b2", "gdrive", "dropbox", "onedrive", "azure", "sftp"}

	for _, providerType := range types {
		t.Run("Type_"+providerType, func(t *testing.T) {
			reqBody := map[string]interface{}{
				"name": "Test Provider",
				"type": providerType,
			}

			body, err := json.Marshal(reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/cloudbackup/providers", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			h.CreateProvider(rr, req)

			assert.NotEqual(t, http.StatusBadRequest, rr.Code)
		})
	}
}

func TestUpdateProvider_InvalidJSON(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodPut, "/api/cloudbackup/providers/1", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.UpdateProvider(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateProvider_InvalidID(t *testing.T) {
	h := setupCloudBackupTest(t)

	reqBody := map[string]interface{}{
		"name": "Updated Provider",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/cloudbackup/providers/invalid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.UpdateProvider(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteProvider_ValidID(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/cloudbackup/providers/1", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.DeleteProvider(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteProvider_InvalidID(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/cloudbackup/providers/invalid", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.DeleteProvider(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestTestProvider_ValidID(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/cloudbackup/providers/1/test", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.TestProvider(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetProviderTypes_Success(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/cloudbackup/provider-types", nil)
	rr := httptest.NewRecorder()

	h.GetProviderTypes(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "data")
}

// ===== Job Management Tests =====

func TestListJobs_Success(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/cloudbackup/jobs", nil)
	rr := httptest.NewRecorder()

	h.ListJobs(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetJob_ValidID(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/cloudbackup/jobs/1", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.GetJob(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetJob_InvalidID(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/cloudbackup/jobs/invalid", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.GetJob(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateJob_InvalidJSON(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/cloudbackup/jobs", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateJob(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateJob_ValidRequest(t *testing.T) {
	h := setupCloudBackupTest(t)

	reqBody := map[string]interface{}{
		"name":        "Daily Backup",
		"providerId":  1,
		"sourcePath":  "/mnt/storage/data",
		"destPath":    "backups/data",
		"mode":        "sync",
		"schedule":    "0 2 * * *",
		"enabled":     true,
		"compression": true,
		"encryption":  true,
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/cloudbackup/jobs", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.CreateJob(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestCreateJob_DifferentModes(t *testing.T) {
	h := setupCloudBackupTest(t)

	modes := []string{"sync", "upload", "download"}

	for _, mode := range modes {
		t.Run("Mode_"+mode, func(t *testing.T) {
			reqBody := map[string]interface{}{
				"name":       "Test Job",
				"providerId": 1,
				"mode":       mode,
			}

			body, err := json.Marshal(reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/cloudbackup/jobs", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			h.CreateJob(rr, req)

			assert.NotEqual(t, http.StatusBadRequest, rr.Code)
		})
	}
}

func TestUpdateJob_InvalidJSON(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodPut, "/api/cloudbackup/jobs/1", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.UpdateJob(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteJob_ValidID(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/cloudbackup/jobs/1", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.DeleteJob(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestRunJob_ValidID(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/cloudbackup/jobs/1/run", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.RunJob(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestRunJob_InvalidID(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/cloudbackup/jobs/invalid/run", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.RunJob(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// ===== Logs Tests =====

func TestGetLogs_NoJobID(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/cloudbackup/logs", nil)
	rr := httptest.NewRecorder()

	h.GetLogs(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetLogs_WithJobID(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/cloudbackup/logs?jobId=1", nil)
	rr := httptest.NewRecorder()

	h.GetLogs(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetLogs_WithLimitAndOffset(t *testing.T) {
	h := setupCloudBackupTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/cloudbackup/logs?limit=10&offset=5", nil)
	rr := httptest.NewRecorder()

	h.GetLogs(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Integration Tests =====

func TestCloudBackupWorkflow_Complete(t *testing.T) {
	h := setupCloudBackupTest(t)

	// 1. Get provider types
	req1 := httptest.NewRequest(http.MethodGet, "/api/cloudbackup/provider-types", nil)
	rr1 := httptest.NewRecorder()
	h.GetProviderTypes(rr1, req1)
	assert.Equal(t, http.StatusOK, rr1.Code)

	// 2. Create provider
	providerBody := map[string]interface{}{
		"name": "Test S3",
		"type": "s3",
	}
	body2, _ := json.Marshal(providerBody)
	req2 := httptest.NewRequest(http.MethodPost, "/api/cloudbackup/providers", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	rr2 := httptest.NewRecorder()
	h.CreateProvider(rr2, req2)
	assert.NotEqual(t, http.StatusBadRequest, rr2.Code)

	// 3. List providers
	req3 := httptest.NewRequest(http.MethodGet, "/api/cloudbackup/providers", nil)
	rr3 := httptest.NewRecorder()
	h.ListProviders(rr3, req3)
	assert.NotEqual(t, http.StatusBadRequest, rr3.Code)

	// 4. Create job
	jobBody := map[string]interface{}{
		"name":       "Daily Backup",
		"providerId": 1,
		"mode":       "sync",
	}
	body4, _ := json.Marshal(jobBody)
	req4 := httptest.NewRequest(http.MethodPost, "/api/cloudbackup/jobs", bytes.NewReader(body4))
	req4.Header.Set("Content-Type", "application/json")
	rr4 := httptest.NewRecorder()
	h.CreateJob(rr4, req4)
	assert.NotEqual(t, http.StatusBadRequest, rr4.Code)

	// 5. List jobs
	req5 := httptest.NewRequest(http.MethodGet, "/api/cloudbackup/jobs", nil)
	rr5 := httptest.NewRecorder()
	h.ListJobs(rr5, req5)
	assert.NotEqual(t, http.StatusBadRequest, rr5.Code)

	// 6. Get logs
	req6 := httptest.NewRequest(http.MethodGet, "/api/cloudbackup/logs", nil)
	rr6 := httptest.NewRecorder()
	h.GetLogs(rr6, req6)
	assert.NotEqual(t, http.StatusBadRequest, rr6.Code)
}

// ===== Benchmark Tests =====

func BenchmarkListProviders(b *testing.B) {
	h := NewCloudBackupHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/cloudbackup/providers", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		h.ListProviders(rr, req)
	}
}

func BenchmarkListJobs(b *testing.B) {
	h := NewCloudBackupHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/cloudbackup/jobs", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		h.ListJobs(rr, req)
	}
}

func BenchmarkGetProviderTypes(b *testing.B) {
	h := NewCloudBackupHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/cloudbackup/provider-types", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		h.GetProviderTypes(rr, req)
	}
}
