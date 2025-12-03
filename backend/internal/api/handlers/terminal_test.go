// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerminal_CreateSession(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/terminal/session", nil)
	rr := httptest.NewRecorder()
	CreateTerminalSession(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestTerminal_ListSessions(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/terminal/sessions", nil)
	rr := httptest.NewRecorder()
	ListTerminalSessions(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func BenchmarkTerminal_ListSessions(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/terminal/sessions", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ListTerminalSessions(rr, req)
	}
}
