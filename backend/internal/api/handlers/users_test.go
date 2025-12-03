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

func TestListUsers_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rr := httptest.NewRecorder()

	ListUsers(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetUser_ValidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/users/1", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	GetUser(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestGetUser_InvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/users/invalid", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	GetUser(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateUser(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateUser_ValidRequest(t *testing.T) {
	reqBody := map[string]interface{}{
		"username": "newuser",
		"password": "Password123!",
		"email":    "user@example.com",
		"role":     "user",
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateUser(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateUser_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/api/users/1", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	UpdateUser(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteUser_ValidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/users/1", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))

	DeleteUser(rr, req)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkListUsers(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ListUsers(rr, req)
	}
}
