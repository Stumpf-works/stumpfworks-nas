// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackup_ListBackups(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/backup", nil)
	rr := httptest.NewRecorder()
	ListBackups(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestBackup_CreateBackup_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/backup", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateBackup(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestBackup_CreateBackup_Valid(t *testing.T) {
	reqBody := map[string]interface{}{"name": "test-backup", "path": "/backup"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/backup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateBackup(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestBackup_RestoreBackup_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/backup/restore", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	RestoreBackup(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkBackup_ListBackups(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/backup", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ListBackups(rr, req)
	}
}
