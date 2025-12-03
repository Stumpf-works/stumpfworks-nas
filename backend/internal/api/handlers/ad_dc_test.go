// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupADDCHandler() *ADDCHandler {
	return NewADDCHandler()
}

// ===== Domain Controller Management Tests =====

func TestADDC_GetStatus(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad-dc/status", nil)
	rr := httptest.NewRecorder()
	handler.GetDCStatus(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_GetConfig(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad-dc/config", nil)
	rr := httptest.NewRecorder()
	handler.GetDCConfig(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_UpdateConfig_InvalidJSON(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPut, "/api/ad-dc/config", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.UpdateDCConfig(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_ProvisionDomain_InvalidJSON(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/provision", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ProvisionDomain(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_ProvisionDomain_MissingFields(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/provision", bytes.NewReader([]byte(`{"realm":"","domain":"","adminPassword":""}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ProvisionDomain(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_ProvisionDomain_ValidRequest(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/provision", bytes.NewReader([]byte(`{"realm":"EXAMPLE.COM","domain":"EXAMPLE","adminPassword":"P@ssw0rd123"}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ProvisionDomain(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_DemoteDomain(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/demote", nil)
	rr := httptest.NewRecorder()
	handler.DemoteDomain(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_GetDomainInfo(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad-dc/domain/info", nil)
	rr := httptest.NewRecorder()
	handler.GetDomainInfo(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_RaiseDomainLevel_InvalidJSON(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/domain/raise-level", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.RaiseDomainLevel(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_RestartService(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/restart", nil)
	rr := httptest.NewRecorder()
	handler.RestartService(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== User Management Tests =====

func TestADDC_ListUsers(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad-dc/users", nil)
	rr := httptest.NewRecorder()
	handler.ListUsers(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_CreateUser_InvalidJSON(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/users", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.CreateUser(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_CreateUser_MissingFields(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/users", bytes.NewReader([]byte(`{"user":{"username":""},"password":""}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.CreateUser(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_DeleteUser(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/ad-dc/users/testuser", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("username", "testuser")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.DeleteUser(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_EnableUser(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/users/testuser/enable", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("username", "testuser")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.EnableUser(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_DisableUser(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/users/testuser/disable", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("username", "testuser")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.DisableUser(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_SetUserPassword_InvalidJSON(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/users/testuser/password", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("username", "testuser")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.SetUserPassword(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// ===== Group Management Tests =====

func TestADDC_ListGroups(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad-dc/groups", nil)
	rr := httptest.NewRecorder()
	handler.ListGroups(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_CreateGroup_InvalidJSON(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/groups", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.CreateGroup(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_CreateGroup_EmptyName(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/groups", bytes.NewReader([]byte(`{"name":""}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.CreateGroup(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_DeleteGroup(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/ad-dc/groups/testgroup", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "testgroup")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.DeleteGroup(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_ListGroupMembers(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad-dc/groups/testgroup/members", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "testgroup")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.ListGroupMembers(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_AddGroupMember_InvalidJSON(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/groups/testgroup/members", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "testgroup")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.AddGroupMember(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// ===== Computer Management Tests =====

func TestADDC_ListComputers(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad-dc/computers", nil)
	rr := httptest.NewRecorder()
	handler.ListComputers(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_CreateComputer_InvalidJSON(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/computers", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.CreateComputer(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_DeleteComputer(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/ad-dc/computers/testcomputer", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "testcomputer")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.DeleteComputer(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== OU Management Tests =====

func TestADDC_ListOUs(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad-dc/ous", nil)
	rr := httptest.NewRecorder()
	handler.ListOUs(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_CreateOU_InvalidJSON(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/ous", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.CreateOU(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_DeleteOU_InvalidJSON(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/ad-dc/ous", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.DeleteOU(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// ===== GPO Management Tests =====

func TestADDC_ListGPOs(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad-dc/gpos", nil)
	rr := httptest.NewRecorder()
	handler.ListGPOs(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_CreateGPO_InvalidJSON(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/gpos", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.CreateGPO(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_DeleteGPO(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/ad-dc/gpos/testgpo", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "testgpo")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.DeleteGPO(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_LinkGPO_InvalidJSON(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/gpos/testgpo/link", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "testgpo")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.LinkGPO(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// ===== DNS Management Tests =====

func TestADDC_ListDNSZones(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad-dc/dns/zones", nil)
	rr := httptest.NewRecorder()
	handler.ListDNSZones(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_CreateDNSZone_InvalidJSON(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/dns/zones", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.CreateDNSZone(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_DeleteDNSZone(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/ad-dc/dns/zones/testzone", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("zone", "testzone")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.DeleteDNSZone(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_ListDNSRecords(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad-dc/dns/zones/testzone/records", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("zone", "testzone")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.ListDNSRecords(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_AddDNSRecord_InvalidJSON(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/dns/zones/testzone/records", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("zone", "testzone")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	handler.AddDNSRecord(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// ===== FSMO Tests =====

func TestADDC_ShowFSMORoles(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad-dc/fsmo/roles", nil)
	rr := httptest.NewRecorder()
	handler.ShowFSMORoles(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_TransferFSMORoles_InvalidJSON(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/fsmo/transfer", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.TransferFSMORoles(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_SeizeFSMORoles_InvalidJSON(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/fsmo/seize", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.SeizeFSMORoles(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// ===== Utility Tests =====

func TestADDC_TestConfiguration(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/test-config", nil)
	rr := httptest.NewRecorder()
	handler.TestConfiguration(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_ShowDBCheck(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/dbcheck", nil)
	rr := httptest.NewRecorder()
	handler.ShowDBCheck(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestADDC_BackupOnline_InvalidJSON(t *testing.T) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/ad-dc/backup", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.BackupOnline(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// ===== Benchmarks =====

func BenchmarkADDC_GetStatus(b *testing.B) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad-dc/status", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.GetDCStatus(rr, req)
	}
}

func BenchmarkADDC_ListUsers(b *testing.B) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad-dc/users", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ListUsers(rr, req)
	}
}

func BenchmarkADDC_ListGroups(b *testing.B) {
	handler := setupADDCHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/ad-dc/groups", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ListGroups(rr, req)
	}
}
