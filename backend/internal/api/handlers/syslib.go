// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package handlers

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/sharing"
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/network"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/Stumpf-works/stumpfworks-nas/internal/system"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"go.uber.org/zap"
)

// ===== Helper Functions =====

// getSystemLib returns the system library and checks if it's initialized
func getSystemLib(w http.ResponseWriter) *system.SystemLibrary {
	lib := system.Get()
	if lib == nil {
		utils.RespondError(w, errors.InternalServerError("System library not initialized", nil))
		return nil
	}
	return lib
}

// ===== System Library Health =====

// SystemLibraryHealth returns health status of all system library subsystems
func SystemLibraryHealth(w http.ResponseWriter, r *http.Request) {
	lib := getSystemLib(w)
	if lib == nil {
		return
	}
	health, err := lib.HealthCheck()
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to get health status", err))
		return
	}
	utils.RespondSuccess(w, health)
}

// ===== ZFS Handlers =====

// ListZFSPools lists all ZFS pools
func ListZFSPools(w http.ResponseWriter, r *http.Request) {
	lib := getSystemLib(w)
	if lib == nil {
		return
	}
	if lib.Storage == nil || lib.Storage.ZFS == nil {
		utils.RespondError(w, errors.BadRequest("ZFS not available", nil))
		return
	}

	pools, err := lib.Storage.ZFS.ListPools()
	if err != nil {
		logger.Error("Failed to list ZFS pools", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list ZFS pools", err))
		return
	}

	utils.RespondSuccess(w, pools)
}

// GetZFSPool gets details of a specific ZFS pool
func GetZFSPool(w http.ResponseWriter, r *http.Request) {
	poolName := chi.URLParam(r, "name")
	lib := getSystemLib(w)
	if lib == nil {
		return
	}

	if lib.Storage == nil || lib.Storage.ZFS == nil {
		utils.RespondError(w, errors.BadRequest("ZFS not available", nil))
		return
	}

	pool, err := lib.Storage.ZFS.GetPool(poolName)
	if err != nil {
		logger.Error("Failed to get ZFS pool", zap.String("pool", poolName), zap.Error(err))
		utils.RespondError(w, errors.NotFound("Pool not found", err))
		return
	}

	utils.RespondSuccess(w, pool)
}

// CreateZFSPool creates a new ZFS pool
func CreateZFSPool(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string            `json:"name"`
		RaidType string            `json:"raid_type"` // mirror, raidz, raidz2, raidz3
		Devices  []string          `json:"devices"`
		Options  map[string]string `json:"options"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	lib := getSystemLib(w)
	if lib == nil {
		return
	}
	if lib.Storage == nil || lib.Storage.ZFS == nil {
		utils.RespondError(w, errors.BadRequest("ZFS not available", nil))
		return
	}

	if err := lib.Storage.ZFS.CreatePool(req.Name, req.RaidType, req.Devices, req.Options); err != nil {
		logger.Error("Failed to create ZFS pool", zap.String("pool", req.Name), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create pool", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "ZFS pool created successfully",
		"pool":    req.Name,
	})
}

// DestroyZFSPool destroys a ZFS pool
func DestroyZFSPool(w http.ResponseWriter, r *http.Request) {
	poolName := chi.URLParam(r, "name")
	lib := getSystemLib(w)
	if lib == nil {
		return
	}

	if lib.Storage == nil || lib.Storage.ZFS == nil {
		utils.RespondError(w, errors.BadRequest("ZFS not available", nil))
		return
	}

	force := r.URL.Query().Get("force") == "true"

	if err := lib.Storage.ZFS.DestroyPool(poolName, force); err != nil {
		logger.Error("Failed to destroy ZFS pool", zap.String("pool", poolName), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to destroy pool", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "ZFS pool destroyed successfully",
	})
}

// ScrubZFSPool starts a scrub operation on a ZFS pool
func ScrubZFSPool(w http.ResponseWriter, r *http.Request) {
	poolName := chi.URLParam(r, "name")
	lib := getSystemLib(w)
	if lib == nil {
		return
	}

	if lib.Storage == nil || lib.Storage.ZFS == nil {
		utils.RespondError(w, errors.BadRequest("ZFS not available", nil))
		return
	}

	if err := lib.Storage.ZFS.ScrubPool(poolName); err != nil {
		logger.Error("Failed to scrub ZFS pool", zap.String("pool", poolName), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to start scrub", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "ZFS pool scrub started",
	})
}

// ListZFSDatasets lists all datasets in a pool
func ListZFSDatasets(w http.ResponseWriter, r *http.Request) {
	poolName := chi.URLParam(r, "pool")
	lib := getSystemLib(w)
	if lib == nil {
		return
	}

	if lib.Storage == nil || lib.Storage.ZFS == nil {
		utils.RespondError(w, errors.BadRequest("ZFS not available", nil))
		return
	}

	datasets, err := lib.Storage.ZFS.ListDatasets(poolName)
	if err != nil {
		logger.Error("Failed to list ZFS datasets", zap.String("pool", poolName), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list datasets", err))
		return
	}

	utils.RespondSuccess(w, datasets)
}

// CreateZFSSnapshot creates a ZFS snapshot
func CreateZFSSnapshot(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Dataset  string `json:"dataset"`
		Snapshot string `json:"snapshot"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	lib := getSystemLib(w)
	if lib == nil {
		return
	}
	if lib.Storage == nil || lib.Storage.ZFS == nil {
		utils.RespondError(w, errors.BadRequest("ZFS not available", nil))
		return
	}

	if err := lib.Storage.ZFS.CreateSnapshot(req.Dataset, req.Snapshot); err != nil {
		logger.Error("Failed to create ZFS snapshot", zap.String("dataset", req.Dataset), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create snapshot", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Snapshot created successfully",
	})
}

// ListZFSSnapshots lists all snapshots for a dataset
func ListZFSSnapshots(w http.ResponseWriter, r *http.Request) {
	datasetName := chi.URLParam(r, "dataset")
	lib := getSystemLib(w)
	if lib == nil {
		return
	}

	if lib.Storage == nil || lib.Storage.ZFS == nil {
		utils.RespondError(w, errors.BadRequest("ZFS not available", nil))
		return
	}

	snapshots, err := lib.Storage.ZFS.ListSnapshots(datasetName)
	if err != nil {
		logger.Error("Failed to list ZFS snapshots", zap.String("dataset", datasetName), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list snapshots", err))
		return
	}

	utils.RespondSuccess(w, snapshots)
}

// ===== RAID Handlers =====

// ListRAIDArrays lists all mdadm RAID arrays
func ListRAIDArrays(w http.ResponseWriter, r *http.Request) {
	lib := getSystemLib(w)
	if lib == nil {
		return
	}
	if lib.Storage == nil || lib.Storage.RAID == nil {
		utils.RespondError(w, errors.BadRequest("RAID not available", nil))
		return
	}

	arrays, err := lib.Storage.RAID.ListArrays()
	if err != nil {
		logger.Error("Failed to list RAID arrays", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list RAID arrays", err))
		return
	}

	utils.RespondSuccess(w, arrays)
}

// GetRAIDArray gets details of a specific RAID array
func GetRAIDArray(w http.ResponseWriter, r *http.Request) {
	arrayName := chi.URLParam(r, "name")
	lib := getSystemLib(w)
	if lib == nil {
		return
	}

	if lib.Storage == nil || lib.Storage.RAID == nil {
		utils.RespondError(w, errors.BadRequest("RAID not available", nil))
		return
	}

	array, err := lib.Storage.RAID.GetArray(arrayName)
	if err != nil {
		logger.Error("Failed to get RAID array", zap.String("array", arrayName), zap.Error(err))
		utils.RespondError(w, errors.NotFound("RAID array not found", err))
		return
	}

	utils.RespondSuccess(w, array)
}

// CreateRAIDArray creates a new RAID array
func CreateRAIDArray(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string   `json:"name"`
		Level   string   `json:"level"` // 0, 1, 5, 6, 10
		Devices []string `json:"devices"`
		Spare   []string `json:"spare"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	lib := getSystemLib(w)
	if lib == nil {
		return
	}
	if lib.Storage == nil || lib.Storage.RAID == nil {
		utils.RespondError(w, errors.BadRequest("RAID not available", nil))
		return
	}

	if err := lib.Storage.RAID.CreateArray(req.Name, req.Level, req.Devices, req.Spare); err != nil {
		logger.Error("Failed to create RAID array", zap.String("array", req.Name), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create RAID array", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "RAID array created successfully",
		"array":   req.Name,
	})
}

// ===== SMART Handlers =====

// GetSMARTInfo gets SMART information for a disk
func GetSMARTInfo(w http.ResponseWriter, r *http.Request) {
	device := chi.URLParam(r, "device")
	lib := getSystemLib(w)
	if lib == nil {
		return
	}

	if lib.Storage == nil || lib.Storage.SMART == nil {
		utils.RespondError(w, errors.BadRequest("SMART not available", nil))
		return
	}

	info, err := lib.Storage.SMART.GetInfo(device)
	if err != nil {
		logger.Error("Failed to get SMART info", zap.String("device", device), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get SMART info", err))
		return
	}

	utils.RespondSuccess(w, info)
}

// RunSMARTTest runs a SMART self-test on a disk
func RunSMARTTest(w http.ResponseWriter, r *http.Request) {
	device := chi.URLParam(r, "device")
	testType := r.URL.Query().Get("type") // short, long, conveyance

	if testType == "" {
		testType = "short"
	}

	lib := getSystemLib(w)
	if lib == nil {
		return
	}
	if lib.Storage == nil || lib.Storage.SMART == nil {
		utils.RespondError(w, errors.BadRequest("SMART not available", nil))
		return
	}

	if err := lib.Storage.SMART.RunTest(device, testType); err != nil {
		logger.Error("Failed to run SMART test", zap.String("device", device), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to run SMART test", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "SMART test started",
		"type":    testType,
	})
}

// ===== Samba Handlers =====

// ListSambaShares lists all Samba shares
func ListSambaShares(w http.ResponseWriter, r *http.Request) {
	lib := system.Get()
	if lib == nil {
		utils.RespondError(w, errors.InternalServerError("System library not initialized", nil))
		return
	}
	if lib.Sharing == nil || lib.Sharing.Samba == nil {
		utils.RespondError(w, errors.BadRequest("Samba not available", nil))
		return
	}

	shares, err := lib.Sharing.Samba.ListShares()
	if err != nil {
		logger.Error("Failed to list Samba shares", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list Samba shares", err))
		return
	}

	utils.RespondSuccess(w, shares)
}

// GetSambaShare gets a specific Samba share
func GetSambaShare(w http.ResponseWriter, r *http.Request) {
	shareName := chi.URLParam(r, "name")
	lib := getSystemLib(w)
	if lib == nil {
		return
	}

	if lib.Sharing == nil || lib.Sharing.Samba == nil {
		utils.RespondError(w, errors.BadRequest("Samba not available", nil))
		return
	}

	share, err := lib.Sharing.Samba.GetShare(shareName)
	if err != nil {
		logger.Error("Failed to get Samba share", zap.String("share", shareName), zap.Error(err))
		utils.RespondError(w, errors.NotFound("Share not found", err))
		return
	}

	utils.RespondSuccess(w, share)
}

// CreateSambaShare creates a new Samba share
func CreateSambaShare(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string   `json:"name"`
		Path        string   `json:"path"`
		Comment     string   `json:"comment"`
		ValidUsers  []string `json:"valid_users"`
		ValidGroups []string `json:"valid_groups"`
		ReadOnly    bool     `json:"read_only"`
		Browseable  bool     `json:"browseable"`
		GuestOK     bool     `json:"guest_ok"`
		RecycleBin  bool     `json:"recycle_bin"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	lib := getSystemLib(w)
	if lib == nil {
		return
	}
	if lib.Sharing == nil || lib.Sharing.Samba == nil {
		utils.RespondError(w, errors.BadRequest("Samba not available", nil))
		return
	}

	share := sharing.SambaShare{
		Name:        req.Name,
		Path:        req.Path,
		Comment:     req.Comment,
		ValidUsers:  req.ValidUsers,
		ValidGroups: req.ValidGroups,
		ReadOnly:    req.ReadOnly,
		Browseable:  req.Browseable,
		GuestOK:     req.GuestOK,
		RecycleBin:  req.RecycleBin,
	}

	if err := lib.Sharing.Samba.CreateShare(share); err != nil {
		logger.Error("Failed to create Samba share", zap.String("share", req.Name), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create share", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Samba share created successfully",
		"share":   req.Name,
	})
}

// DeleteSambaShare deletes a Samba share
func DeleteSambaShare(w http.ResponseWriter, r *http.Request) {
	shareName := chi.URLParam(r, "name")
	lib := getSystemLib(w)
	if lib == nil {
		return
	}

	if lib.Sharing == nil || lib.Sharing.Samba == nil {
		utils.RespondError(w, errors.BadRequest("Samba not available", nil))
		return
	}

	if err := lib.Sharing.Samba.DeleteShare(shareName); err != nil {
		logger.Error("Failed to delete Samba share", zap.String("share", shareName), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to delete share", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Samba share deleted successfully",
	})
}

// GetSambaStatus gets Samba service status
func GetSambaStatus(w http.ResponseWriter, r *http.Request) {
	lib := getSystemLib(w)
	if lib == nil {
		return
	}
	if lib.Sharing == nil || lib.Sharing.Samba == nil {
		utils.RespondError(w, errors.BadRequest("Samba not available", nil))
		return
	}

	active, err := lib.Sharing.Samba.GetStatus()
	if err != nil {
		logger.Error("Failed to get Samba status", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get status", err))
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"active":  active,
		"enabled": lib.Sharing.Samba.IsEnabled(),
	})
}

// RestartSamba restarts the Samba service
func RestartSamba(w http.ResponseWriter, r *http.Request) {
	lib := getSystemLib(w)
	if lib == nil {
		return
	}
	if lib.Sharing == nil || lib.Sharing.Samba == nil {
		utils.RespondError(w, errors.BadRequest("Samba not available", nil))
		return
	}

	if err := lib.Sharing.Samba.Restart(); err != nil {
		logger.Error("Failed to restart Samba", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to restart Samba", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Samba restarted successfully",
	})
}

// ===== NFS Handlers =====

// ListNFSExports lists all NFS exports
func ListNFSExports(w http.ResponseWriter, r *http.Request) {
	lib := getSystemLib(w)
	if lib == nil {
		return
	}
	if lib.Sharing == nil || lib.Sharing.NFS == nil {
		utils.RespondError(w, errors.BadRequest("NFS not available", nil))
		return
	}

	exports, err := lib.Sharing.NFS.ListExports()
	if err != nil {
		logger.Error("Failed to list NFS exports", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list NFS exports", err))
		return
	}

	utils.RespondSuccess(w, exports)
}

// CreateNFSExport creates a new NFS export
func CreateNFSExport(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Path    string   `json:"path"`
		Clients []string `json:"clients"`
		Options []string `json:"options"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	lib := getSystemLib(w)
	if lib == nil {
		return
	}
	if lib.Sharing == nil || lib.Sharing.NFS == nil {
		utils.RespondError(w, errors.BadRequest("NFS not available", nil))
		return
	}

	export := sharing.NFSExport{
		Path:    req.Path,
		Clients: req.Clients,
		Options: req.Options,
	}

	if err := lib.Sharing.NFS.CreateExport(export); err != nil {
		logger.Error("Failed to create NFS export", zap.String("path", req.Path), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create export", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "NFS export created successfully",
		"path":    req.Path,
	})
}

// DeleteNFSExport deletes an NFS export
func DeleteNFSExport(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		utils.RespondError(w, errors.BadRequest("Path parameter required", nil))
		return
	}

	lib := getSystemLib(w)
	if lib == nil {
		return
	}
	if lib.Sharing == nil || lib.Sharing.NFS == nil {
		utils.RespondError(w, errors.BadRequest("NFS not available", nil))
		return
	}

	if err := lib.Sharing.NFS.DeleteExport(path); err != nil {
		logger.Error("Failed to delete NFS export", zap.String("path", path), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to delete export", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "NFS export deleted successfully",
	})
}

// RestartNFS restarts the NFS service
func RestartNFS(w http.ResponseWriter, r *http.Request) {
	lib := getSystemLib(w)
	if lib == nil {
		return
	}
	if lib.Sharing == nil || lib.Sharing.NFS == nil {
		utils.RespondError(w, errors.BadRequest("NFS not available", nil))
		return
	}

	if err := lib.Sharing.NFS.Restart(); err != nil {
		logger.Error("Failed to restart NFS", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to restart NFS", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "NFS restarted successfully",
	})
}

// ===== Network Interface Handlers =====

// CreateBondInterface creates a bonded network interface
func CreateBondInterface(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name       string   `json:"name"`
		Mode       string   `json:"mode"`
		Interfaces []string `json:"interfaces"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	lib := getSystemLib(w)
	if lib == nil {
		return
	}
	if lib.Network == nil || lib.Network.Interfaces == nil {
		utils.RespondError(w, errors.BadRequest("Network not available", nil))
		return
	}

	config := network.BondConfig{
		Name:       req.Name,
		Mode:       req.Mode,
		Slaves:     req.Interfaces,
	}

	if err := lib.Network.Interfaces.CreateBond(config); err != nil {
		logger.Error("Failed to create bond interface", zap.String("name", req.Name), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create bond", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Bond interface created successfully",
		"bond":    req.Name,
	})
}

// CreateVLANInterface creates a VLAN interface
func CreateVLANInterface(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Parent string `json:"parent"`
		VLANID int    `json:"vlan_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	lib := getSystemLib(w)
	if lib == nil {
		return
	}
	if lib.Network == nil || lib.Network.Interfaces == nil {
		utils.RespondError(w, errors.BadRequest("Network not available", nil))
		return
	}

	if err := lib.Network.Interfaces.CreateVLAN(req.Parent, req.VLANID); err != nil {
		logger.Error("Failed to create VLAN interface",
			zap.String("parent", req.Parent),
			zap.Int("vlan_id", req.VLANID),
			zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create VLAN", err))
		return
	}

	vlanName := req.Parent + "." + strconv.Itoa(req.VLANID)
	utils.RespondSuccess(w, map[string]string{
		"message": "VLAN interface created successfully",
		"vlan":    vlanName,
	})
}
