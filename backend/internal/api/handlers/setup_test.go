// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ===== Setup Status Tests =====

func TestSetup_Status(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/setup/status", nil)
	rr := httptest.NewRecorder()
	SetupStatus(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Initialize Setup Tests =====

func TestSetup_Initialize_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/setup/initialize", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	InitializeSetup(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSetup_Initialize_EmptyUsername(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/setup/initialize", bytes.NewReader([]byte(`{"username":"","email":"test@example.com","password":"testpass123","fullName":"Test User"}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	InitializeSetup(rr, req)
	assert.NotEqual(t, http.StatusOK, rr.Code)
}

func TestSetup_Initialize_ShortUsername(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/setup/initialize", bytes.NewReader([]byte(`{"username":"ab","email":"test@example.com","password":"testpass123","fullName":"Test User"}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	InitializeSetup(rr, req)
	assert.NotEqual(t, http.StatusOK, rr.Code)
}

func TestSetup_Initialize_InvalidEmail(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/setup/initialize", bytes.NewReader([]byte(`{"username":"testuser","email":"invalid-email","password":"testpass123","fullName":"Test User"}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	InitializeSetup(rr, req)
	assert.NotEqual(t, http.StatusOK, rr.Code)
}

func TestSetup_Initialize_ShortPassword(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/setup/initialize", bytes.NewReader([]byte(`{"username":"testuser","email":"test@example.com","password":"short","fullName":"Test User"}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	InitializeSetup(rr, req)
	assert.NotEqual(t, http.StatusOK, rr.Code)
}

func TestSetup_Initialize_EmptyFullName(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/setup/initialize", bytes.NewReader([]byte(`{"username":"testuser","email":"test@example.com","password":"testpass123","fullName":"T"}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	InitializeSetup(rr, req)
	assert.NotEqual(t, http.StatusOK, rr.Code)
}

func TestSetup_Initialize_ValidRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/setup/initialize", bytes.NewReader([]byte(`{"username":"testuser","email":"test@example.com","password":"testpass123","fullName":"Test User"}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	InitializeSetup(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Benchmarks =====

func BenchmarkSetup_Status(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/setup/status", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		SetupStatus(rr, req)
	}
}

func BenchmarkSetup_Initialize(b *testing.B) {
	req := httptest.NewRequest(http.MethodPost, "/api/setup/initialize", bytes.NewReader([]byte(`{"username":"testuser","email":"test@example.com","password":"testpass123","fullName":"Test User"}`)))
	req.Header.Set("Content-Type", "application/json")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		InitializeSetup(rr, req)
	}
}
