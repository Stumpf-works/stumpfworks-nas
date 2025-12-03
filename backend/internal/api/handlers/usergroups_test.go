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

func TestUserGroups_List(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/usergroups", nil)
	rr := httptest.NewRecorder()
	ListUserGroups(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestUserGroups_Create_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/usergroups", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateUserGroup(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUserGroups_Update_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/api/usergroups/1", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	UpdateUserGroup(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUserGroups_Delete(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/usergroups/1", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	DeleteUserGroup(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkUserGroups_List(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/usergroups", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ListUserGroups(rr, req)
	}
}
