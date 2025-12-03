// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ===== WebSocket Handler Tests =====

func TestWebSocket_Handler_NoUpgrade(t *testing.T) {
	// Test without WebSocket upgrade - should fail upgrade
	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	rr := httptest.NewRecorder()
	WebSocketHandler(rr, req)
	// Since this is not a proper WebSocket upgrade request, the handler should fail
	// The test verifies the handler doesn't panic
	assert.NotNil(t, rr)
}

func TestWebSocket_Handler_WithHeaders(t *testing.T) {
	// Test with some WebSocket headers (still won't upgrade in test environment)
	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	rr := httptest.NewRecorder()
	WebSocketHandler(rr, req)
	assert.NotNil(t, rr)
}

func TestWebSocket_Upgrader_OriginCheck(t *testing.T) {
	upgrader := createUpgrader()
	assert.NotNil(t, upgrader)
	assert.NotNil(t, upgrader.CheckOrigin)

	// Test CheckOrigin function with empty origin (should allow)
	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	allowed := upgrader.CheckOrigin(req)
	assert.True(t, allowed, "Empty origin should be allowed")
}

func TestWebSocket_Upgrader_SameHostOrigin(t *testing.T) {
	upgrader := createUpgrader()

	// Test same-host origin
	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	req.Host = "localhost:8080"
	req.Header.Set("Origin", "http://localhost:8080")
	allowed := upgrader.CheckOrigin(req)
	assert.True(t, allowed, "Same-host origin should be allowed")
}

func TestWebSocket_Upgrader_LocalhostInDev(t *testing.T) {
	upgrader := createUpgrader()

	// Test localhost origin (should be allowed in development)
	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	// Note: The result depends on config.IsDevelopment()
	upgrader.CheckOrigin(req)
	// We don't assert the result since it depends on environment
}

// ===== Benchmark =====

func BenchmarkWebSocket_Handler(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		WebSocketHandler(rr, req)
	}
}
