// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAudit_ListLogs(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/audit", nil)
	rr := httptest.NewRecorder()
	ListAuditLogs(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestAudit_WithLimit(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/audit?limit=10", nil)
	rr := httptest.NewRecorder()
	ListAuditLogs(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestAudit_WithOffset(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/audit?limit=10&offset=5", nil)
	rr := httptest.NewRecorder()
	ListAuditLogs(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkAudit_ListLogs(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/audit", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ListAuditLogs(rr, req)
	}
}
