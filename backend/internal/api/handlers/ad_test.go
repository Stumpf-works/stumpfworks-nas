// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupADHandler() *ADHandler {
	return NewADHandler()
}

// ===== Configuration Tests =====

func TestAD_GetConfig(t *testing.T) {
	handler := setupADHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad/config", nil)
	rr := httptest.NewRecorder()
	handler.GetConfig(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestAD_UpdateConfig_InvalidJSON(t *testing.T) {
	handler := setupADHandler()
	req := httptest.NewRequest(http.MethodPut, "/api/ad/config", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.UpdateConfig(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAD_UpdateConfig_ValidRequest(t *testing.T) {
	handler := setupADHandler()
	req := httptest.NewRequest(http.MethodPut, "/api/ad/config", bytes.NewReader([]byte(`{"enabled":true,"server":"ldap.example.com","domain":"example.com"}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.UpdateConfig(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Connection Tests =====

func TestAD_TestConnection(t *testing.T) {
	handler := setupADHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad/test", nil)
	rr := httptest.NewRecorder()
	handler.TestConnection(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Authentication Tests =====

func TestAD_Authenticate_InvalidJSON(t *testing.T) {
	handler := setupADHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad/authenticate", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.Authenticate(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAD_Authenticate_EmptyUsername(t *testing.T) {
	handler := setupADHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad/authenticate", bytes.NewReader([]byte(`{"username":"","password":"test"}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.Authenticate(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAD_Authenticate_EmptyPassword(t *testing.T) {
	handler := setupADHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad/authenticate", bytes.NewReader([]byte(`{"username":"test","password":""}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.Authenticate(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAD_Authenticate_ValidRequest(t *testing.T) {
	handler := setupADHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad/authenticate", bytes.NewReader([]byte(`{"username":"testuser","password":"testpass"}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.Authenticate(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== User Management Tests =====

func TestAD_ListUsers(t *testing.T) {
	handler := setupADHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad/users", nil)
	rr := httptest.NewRecorder()
	handler.ListUsers(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestAD_SyncUser_InvalidJSON(t *testing.T) {
	handler := setupADHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad/sync", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.SyncUser(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAD_SyncUser_EmptyUsername(t *testing.T) {
	handler := setupADHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad/sync", bytes.NewReader([]byte(`{"username":""}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.SyncUser(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAD_SyncUser_ValidRequest(t *testing.T) {
	handler := setupADHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad/sync", bytes.NewReader([]byte(`{"username":"testuser"}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.SyncUser(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Status Tests =====

func TestAD_GetStatus(t *testing.T) {
	handler := setupADHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad/status", nil)
	rr := httptest.NewRecorder()
	handler.GetStatus(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Benchmarks =====

func BenchmarkAD_GetConfig(b *testing.B) {
	handler := setupADHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad/config", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.GetConfig(rr, req)
	}
}

func BenchmarkAD_ListUsers(b *testing.B) {
	handler := setupADHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad/users", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ListUsers(rr, req)
	}
}

func BenchmarkAD_GetStatus(b *testing.B) {
	handler := setupADHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad/status", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.GetStatus(rr, req)
	}
}
