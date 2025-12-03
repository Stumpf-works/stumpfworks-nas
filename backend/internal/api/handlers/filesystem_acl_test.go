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

func TestACL_Get_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/acl/get", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	GetACL(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestACL_Set_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/acl/set", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	SetACL(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestACL_Remove_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/acl/remove", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	RemoveACL(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkACL_Get(b *testing.B) {
	reqBody := map[string]interface{}{"path": "/test"}
	body, _ := json.Marshal(reqBody)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/acl/get", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		GetACL(rr, req)
	}
}
