// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFiles_List(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/files?path=/", nil)
	rr := httptest.NewRecorder()
	ListFiles(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestFiles_CreateDirectory_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/files/mkdir", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateDirectory(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestFiles_Delete_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/files", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	DeleteFile(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestFiles_Copy_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/files/copy", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CopyFile(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestFiles_Move_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/files/move", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	MoveFile(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkFiles_List(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/files?path=/", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ListFiles(rr, req)
	}
}
