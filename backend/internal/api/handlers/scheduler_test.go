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

func TestScheduler_ListJobs(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/scheduler/jobs", nil)
	rr := httptest.NewRecorder()
	ListScheduledJobs(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestScheduler_CreateJob_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/scheduler/jobs", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateScheduledJob(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestScheduler_DeleteJob(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/scheduler/jobs/1", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	DeleteScheduledJob(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkScheduler_ListJobs(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/scheduler/jobs", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ListScheduledJobs(rr, req)
	}
}
