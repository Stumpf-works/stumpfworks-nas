// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFailedLogin_List(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/failed-logins", nil)
	rr := httptest.NewRecorder()
	ListFailedLogins(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestFailedLogin_WithLimit(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/failed-logins?limit=10", nil)
	rr := httptest.NewRecorder()
	ListFailedLogins(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestFailedLogin_Clear(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/failed-logins", nil)
	rr := httptest.NewRecorder()
	ClearFailedLogins(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkFailedLogin_List(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/failed-logins", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ListFailedLogins(rr, req)
	}
}
