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

// ===== System Library Health Tests =====

func TestSyslib_Health(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/syslib/health", nil)
	rr := httptest.NewRecorder()
	SystemLibraryHealth(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== ZFS Tests =====

func TestSyslib_ZFS_ListPools(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/syslib/zfs/pools", nil)
	rr := httptest.NewRecorder()
	ListZFSPools(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_ZFS_GetPool(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/syslib/zfs/pools/testpool", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "testpool")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	GetZFSPool(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_ZFS_CreatePool_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/syslib/zfs/pools", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateZFSPool(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_ZFS_CreatePool_ValidRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/syslib/zfs/pools", bytes.NewReader([]byte(`{"name":"testpool","raid_type":"mirror","devices":["/dev/sda","/dev/sdb"],"options":{}}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateZFSPool(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_ZFS_DestroyPool(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/syslib/zfs/pools/testpool", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "testpool")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	DestroyZFSPool(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_ZFS_ScrubPool(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/syslib/zfs/pools/testpool/scrub", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "testpool")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	ScrubZFSPool(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_ZFS_ListDatasets(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/syslib/zfs/pools/testpool/datasets", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("pool", "testpool")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	ListZFSDatasets(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_ZFS_CreateSnapshot_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/syslib/zfs/snapshots", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateZFSSnapshot(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_ZFS_ListSnapshots(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/syslib/zfs/datasets/testpool/data/snapshots", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("dataset", "testpool/data")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	ListZFSSnapshots(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== RAID Tests =====

func TestSyslib_RAID_ListArrays(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/syslib/raid/arrays", nil)
	rr := httptest.NewRecorder()
	ListRAIDArrays(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_RAID_GetArray(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/syslib/raid/arrays/md0", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "md0")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	GetRAIDArray(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_RAID_CreateArray_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/syslib/raid/arrays", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateRAIDArray(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_RAID_CreateArray_ValidRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/syslib/raid/arrays", bytes.NewReader([]byte(`{"name":"/dev/md0","level":"5","devices":["/dev/sda","/dev/sdb","/dev/sdc"],"spare":[]}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateRAIDArray(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== SMART Tests =====

func TestSyslib_SMART_GetInfo(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/syslib/smart/sda", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("device", "sda")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	GetSMARTInfo(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_SMART_RunTest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/syslib/smart/sda/test?type=short", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("device", "sda")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	RunSMARTTest(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_SMART_RunTest_DefaultType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/syslib/smart/sda/test", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("device", "sda")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	RunSMARTTest(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Samba Tests =====

func TestSyslib_Samba_ListShares(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/syslib/samba/shares", nil)
	rr := httptest.NewRecorder()
	ListSambaShares(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_Samba_GetShare(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/syslib/samba/shares/testshare", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "testshare")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	GetSambaShare(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_Samba_CreateShare_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/syslib/samba/shares", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateSambaShare(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_Samba_CreateShare_ValidRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/syslib/samba/shares", bytes.NewReader([]byte(`{"name":"testshare","path":"/mnt/test","comment":"Test Share","valid_users":[],"valid_groups":[],"read_only":false,"browseable":true,"guest_ok":false,"recycle_bin":false}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateSambaShare(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_Samba_UpdateShare_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/api/syslib/samba/shares/testshare", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "testshare")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	UpdateSambaShare(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_Samba_DeleteShare(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/syslib/samba/shares/testshare", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "testshare")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	DeleteSambaShare(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_Samba_GetStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/syslib/samba/status", nil)
	rr := httptest.NewRecorder()
	GetSambaStatus(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_Samba_Restart(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/syslib/samba/restart", nil)
	rr := httptest.NewRecorder()
	RestartSamba(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== NFS Tests =====

func TestSyslib_NFS_ListExports(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/syslib/nfs/exports", nil)
	rr := httptest.NewRecorder()
	ListNFSExports(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_NFS_CreateExport_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/syslib/nfs/exports", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateNFSExport(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_NFS_CreateExport_ValidRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/syslib/nfs/exports", bytes.NewReader([]byte(`{"path":"/mnt/test","clients":["*"],"options":["rw","sync"]}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateNFSExport(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_NFS_DeleteExport_MissingPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/syslib/nfs/exports", nil)
	rr := httptest.NewRecorder()
	DeleteNFSExport(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_NFS_DeleteExport_ValidRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/syslib/nfs/exports?path=/mnt/test", nil)
	rr := httptest.NewRecorder()
	DeleteNFSExport(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_NFS_Restart(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/syslib/nfs/restart", nil)
	rr := httptest.NewRecorder()
	RestartNFS(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Network Interface Tests =====

func TestSyslib_Network_CreateBond_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/syslib/network/bonds", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateBondInterface(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_Network_CreateBond_ValidRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/syslib/network/bonds", bytes.NewReader([]byte(`{"name":"bond0","mode":"balance-rr","interfaces":["eth0","eth1"]}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateBondInterface(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_Network_CreateVLAN_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/syslib/network/vlans", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateVLANInterface(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_Network_CreateVLAN_ValidRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/syslib/network/vlans", bytes.NewReader([]byte(`{"parent":"eth0","vlan_id":100}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CreateVLANInterface(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_Network_DeleteBond_EmptyName(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/syslib/network/bonds/", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	DeleteBondInterface(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_Network_DeleteBond_ValidRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/syslib/network/bonds/bond0", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "bond0")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	DeleteBondInterface(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_Network_DeleteVLAN_EmptyParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/syslib/network/vlans/", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("parent", "")
	rctx.URLParams.Add("vlanid", "")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	DeleteVLANInterface(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSyslib_Network_DeleteVLAN_ValidRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/syslib/network/vlans/eth0/100", nil)
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("parent", "eth0")
	rctx.URLParams.Add("vlanid", "100")
	req = req.WithContext(chi.NewRouteContext().WithValue(req.Context(), chi.RouteCtxKey, rctx))
	DeleteVLANInterface(rr, req)
	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
}

// ===== Benchmarks =====

func BenchmarkSyslib_Health(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/syslib/health", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		SystemLibraryHealth(rr, req)
	}
}

func BenchmarkSyslib_ZFS_ListPools(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/syslib/zfs/pools", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ListZFSPools(rr, req)
	}
}

func BenchmarkSyslib_Samba_ListShares(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/syslib/samba/shares", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ListSambaShares(rr, req)
	}
}
