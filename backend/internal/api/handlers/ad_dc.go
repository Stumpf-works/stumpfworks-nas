// Revision: 2025-11-23 | Author: Claude | Version: 1.0.0
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/internal/ad"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// ADDCHandler handles AD Domain Controller HTTP requests
type ADDCHandler struct {
	service *ad.DCService
}

// NewADDCHandler creates a new AD DC handler
func NewADDCHandler() *ADDCHandler {
	return &ADDCHandler{
		service: ad.GetDCService(),
	}
}

// ===== Domain Controller Management =====

// GetDCStatus returns the DC status
func (h *ADDCHandler) GetDCStatus(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	status := map[string]interface{}{
		"provisioned": h.service.IsProvisioned(),
		"config":      h.service.GetConfig(),
	}

	// Get service status
	if serviceStatus, err := h.service.GetSambaServiceStatus(); err == nil {
		status["service_status"] = serviceStatus
	}

	// Get domain info if provisioned
	if h.service.IsProvisioned() {
		if domainInfo, err := h.service.GetDomainInfo(); err == nil {
			status["domain_info"] = domainInfo
		}
	}

	utils.RespondSuccess(w, status)
}

// GetDCConfig returns the DC configuration
func (h *ADDCHandler) GetDCConfig(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	config := h.service.GetConfig()
	utils.RespondSuccess(w, config)
}

// UpdateDCConfig updates the DC configuration
func (h *ADDCHandler) UpdateDCConfig(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	var config ad.DCConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := h.service.UpdateConfig(&config); err != nil {
		logger.Error("Failed to update DC config", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to update configuration", err))
		return
	}

	logger.Info("DC configuration updated")
	utils.RespondSuccess(w, config)
}

// ProvisionDomain provisions a new AD domain
func (h *ADDCHandler) ProvisionDomain(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	var opts ad.ProvisionOptions
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	// Validate required fields
	if opts.Realm == "" || opts.Domain == "" || opts.AdminPassword == "" {
		utils.RespondError(w, errors.BadRequest("Realm, Domain, and AdminPassword are required", nil))
		return
	}

	// Set defaults
	if opts.ServerRole == "" {
		opts.ServerRole = "dc"
	}
	if opts.DNSBackend == "" {
		opts.DNSBackend = "SAMBA_INTERNAL"
	}

	logger.Info("Provisioning AD domain", zap.String("realm", opts.Realm), zap.String("domain", opts.Domain))

	if err := h.service.Provision(opts); err != nil {
		logger.Error("Failed to provision domain", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to provision domain", err))
		return
	}

	logger.Info("Domain provisioned successfully")
	utils.RespondSuccess(w, map[string]string{
		"message": "Domain provisioned successfully",
		"realm":   opts.Realm,
		"domain":  opts.Domain,
	})
}

// DemoteDomain demotes the domain controller
func (h *ADDCHandler) DemoteDomain(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	logger.Info("Demoting AD domain controller")

	if err := h.service.Demote(); err != nil {
		logger.Error("Failed to demote domain", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to demote domain", err))
		return
	}

	logger.Info("Domain controller demoted successfully")
	utils.RespondSuccess(w, map[string]string{
		"message": "Domain controller demoted successfully",
	})
}

// GetDomainInfo returns domain information
func (h *ADDCHandler) GetDomainInfo(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	info, err := h.service.GetDomainInfo()
	if err != nil {
		utils.RespondError(w, errors.NotFound("Domain info not available", err))
		return
	}

	utils.RespondSuccess(w, info)
}

// GetDomainLevel returns the domain functional level
func (h *ADDCHandler) GetDomainLevel(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	level, err := h.service.GetDomainLevel()
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to get domain level", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"level": level,
	})
}

// RaiseDomainLevel raises the domain functional level
func (h *ADDCHandler) RaiseDomainLevel(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	var req struct {
		Level string `json:"level"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := h.service.RaiseDomainLevel(req.Level); err != nil {
		logger.Error("Failed to raise domain level", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to raise domain level", err))
		return
	}

	logger.Info("Domain level raised", zap.String("level", req.Level))
	utils.RespondSuccess(w, map[string]string{
		"message": "Domain level raised successfully",
		"level":   req.Level,
	})
}

// RestartService restarts the Samba AD DC service
func (h *ADDCHandler) RestartService(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	if err := h.service.RestartSambaService(); err != nil {
		logger.Error("Failed to restart service", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to restart service", err))
		return
	}

	logger.Info("Samba AD DC service restarted")
	utils.RespondSuccess(w, map[string]string{
		"message": "Service restarted successfully",
	})
}

// ===== User Management =====

// ListUsers lists all AD users
func (h *ADDCHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	users, err := h.service.ListUsers()
	if err != nil {
		logger.Error("Failed to list users", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list users", err))
		return
	}

	utils.RespondSuccess(w, users)
}

// CreateUser creates a new AD user
func (h *ADDCHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	var req struct {
		User     ad.ADDCUser `json:"user"`
		Password string      `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.User.Username == "" || req.Password == "" {
		utils.RespondError(w, errors.BadRequest("Username and password are required", nil))
		return
	}

	if err := h.service.CreateUser(req.User, req.Password); err != nil {
		logger.Error("Failed to create user", zap.Error(err), zap.String("username", req.User.Username))
		utils.RespondError(w, errors.InternalServerError("Failed to create user", err))
		return
	}

	logger.Info("User created", zap.String("username", req.User.Username))
	utils.RespondSuccess(w, map[string]string{
		"message":  "User created successfully",
		"username": req.User.Username,
	})
}

// DeleteUser deletes an AD user
func (h *ADDCHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	username := chi.URLParam(r, "username")
	if username == "" {
		utils.RespondError(w, errors.BadRequest("Username is required", nil))
		return
	}

	if err := h.service.DeleteUser(username); err != nil {
		logger.Error("Failed to delete user", zap.Error(err), zap.String("username", username))
		utils.RespondError(w, errors.InternalServerError("Failed to delete user", err))
		return
	}

	logger.Info("User deleted", zap.String("username", username))
	utils.RespondSuccess(w, map[string]string{
		"message":  "User deleted successfully",
		"username": username,
	})
}

// EnableUser enables an AD user
func (h *ADDCHandler) EnableUser(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	username := chi.URLParam(r, "username")
	if username == "" {
		utils.RespondError(w, errors.BadRequest("Username is required", nil))
		return
	}

	if err := h.service.EnableUser(username); err != nil {
		logger.Error("Failed to enable user", zap.Error(err), zap.String("username", username))
		utils.RespondError(w, errors.InternalServerError("Failed to enable user", err))
		return
	}

	logger.Info("User enabled", zap.String("username", username))
	utils.RespondSuccess(w, map[string]string{
		"message":  "User enabled successfully",
		"username": username,
	})
}

// DisableUser disables an AD user
func (h *ADDCHandler) DisableUser(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	username := chi.URLParam(r, "username")
	if username == "" {
		utils.RespondError(w, errors.BadRequest("Username is required", nil))
		return
	}

	if err := h.service.DisableUser(username); err != nil {
		logger.Error("Failed to disable user", zap.Error(err), zap.String("username", username))
		utils.RespondError(w, errors.InternalServerError("Failed to disable user", err))
		return
	}

	logger.Info("User disabled", zap.String("username", username))
	utils.RespondSuccess(w, map[string]string{
		"message":  "User disabled successfully",
		"username": username,
	})
}

// SetUserPassword sets a user's password
func (h *ADDCHandler) SetUserPassword(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	username := chi.URLParam(r, "username")
	if username == "" {
		utils.RespondError(w, errors.BadRequest("Username is required", nil))
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Password == "" {
		utils.RespondError(w, errors.BadRequest("Password is required", nil))
		return
	}

	if err := h.service.SetUserPassword(username, req.Password); err != nil {
		logger.Error("Failed to set user password", zap.Error(err), zap.String("username", username))
		utils.RespondError(w, errors.InternalServerError("Failed to set password", err))
		return
	}

	logger.Info("User password set", zap.String("username", username))
	utils.RespondSuccess(w, map[string]string{
		"message":  "Password set successfully",
		"username": username,
	})
}

// SetUserExpiry sets user account expiry
func (h *ADDCHandler) SetUserExpiry(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	username := chi.URLParam(r, "username")
	if username == "" {
		utils.RespondError(w, errors.BadRequest("Username is required", nil))
		return
	}

	var req struct {
		Days int `json:"days"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := h.service.SetUserExpiry(username, req.Days); err != nil {
		logger.Error("Failed to set user expiry", zap.Error(err), zap.String("username", username))
		utils.RespondError(w, errors.InternalServerError("Failed to set expiry", err))
		return
	}

	logger.Info("User expiry set", zap.String("username", username), zap.Int("days", req.Days))
	utils.RespondSuccess(w, map[string]interface{}{
		"message":  "Expiry set successfully",
		"username": username,
		"days":     req.Days,
	})
}

// ===== Group Management =====

// ListGroups lists all AD groups
func (h *ADDCHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	groups, err := h.service.ListGroups()
	if err != nil {
		logger.Error("Failed to list groups", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list groups", err))
		return
	}

	utils.RespondSuccess(w, groups)
}

// CreateGroup creates a new AD group
func (h *ADDCHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	var group ad.ADGroup
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if group.Name == "" {
		utils.RespondError(w, errors.BadRequest("Group name is required", nil))
		return
	}

	if err := h.service.CreateGroup(group); err != nil {
		logger.Error("Failed to create group", zap.Error(err), zap.String("group", group.Name))
		utils.RespondError(w, errors.InternalServerError("Failed to create group", err))
		return
	}

	logger.Info("Group created", zap.String("group", group.Name))
	utils.RespondSuccess(w, map[string]string{
		"message": "Group created successfully",
		"name":    group.Name,
	})
}

// DeleteGroup deletes an AD group
func (h *ADDCHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	groupName := chi.URLParam(r, "name")
	if groupName == "" {
		utils.RespondError(w, errors.BadRequest("Group name is required", nil))
		return
	}

	if err := h.service.DeleteGroup(groupName); err != nil {
		logger.Error("Failed to delete group", zap.Error(err), zap.String("group", groupName))
		utils.RespondError(w, errors.InternalServerError("Failed to delete group", err))
		return
	}

	logger.Info("Group deleted", zap.String("group", groupName))
	utils.RespondSuccess(w, map[string]string{
		"message": "Group deleted successfully",
		"name":    groupName,
	})
}

// ListGroupMembers lists members of a group
func (h *ADDCHandler) ListGroupMembers(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	groupName := chi.URLParam(r, "name")
	if groupName == "" {
		utils.RespondError(w, errors.BadRequest("Group name is required", nil))
		return
	}

	members, err := h.service.ListGroupMembers(groupName)
	if err != nil {
		logger.Error("Failed to list group members", zap.Error(err), zap.String("group", groupName))
		utils.RespondError(w, errors.InternalServerError("Failed to list group members", err))
		return
	}

	utils.RespondSuccess(w, members)
}

// AddGroupMember adds a user to a group
func (h *ADDCHandler) AddGroupMember(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	groupName := chi.URLParam(r, "name")
	if groupName == "" {
		utils.RespondError(w, errors.BadRequest("Group name is required", nil))
		return
	}

	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Username == "" {
		utils.RespondError(w, errors.BadRequest("Username is required", nil))
		return
	}

	if err := h.service.AddGroupMember(groupName, req.Username); err != nil {
		logger.Error("Failed to add group member", zap.Error(err), zap.String("group", groupName), zap.String("username", req.Username))
		utils.RespondError(w, errors.InternalServerError("Failed to add group member", err))
		return
	}

	logger.Info("User added to group", zap.String("group", groupName), zap.String("username", req.Username))
	utils.RespondSuccess(w, map[string]string{
		"message":  "User added to group successfully",
		"group":    groupName,
		"username": req.Username,
	})
}

// RemoveGroupMember removes a user from a group
func (h *ADDCHandler) RemoveGroupMember(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	groupName := chi.URLParam(r, "name")
	username := chi.URLParam(r, "username")

	if groupName == "" || username == "" {
		utils.RespondError(w, errors.BadRequest("Group name and username are required", nil))
		return
	}

	if err := h.service.RemoveGroupMember(groupName, username); err != nil {
		logger.Error("Failed to remove group member", zap.Error(err), zap.String("group", groupName), zap.String("username", username))
		utils.RespondError(w, errors.InternalServerError("Failed to remove group member", err))
		return
	}

	logger.Info("User removed from group", zap.String("group", groupName), zap.String("username", username))
	utils.RespondSuccess(w, map[string]string{
		"message":  "User removed from group successfully",
		"group":    groupName,
		"username": username,
	})
}

// ===== Computer Management =====

// ListComputers lists all AD computers
func (h *ADDCHandler) ListComputers(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	computers, err := h.service.ListComputers()
	if err != nil {
		logger.Error("Failed to list computers", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list computers", err))
		return
	}

	utils.RespondSuccess(w, computers)
}

// CreateComputer creates a new AD computer
func (h *ADDCHandler) CreateComputer(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	var computer ad.ADComputer
	if err := json.NewDecoder(r.Body).Decode(&computer); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if computer.Name == "" {
		utils.RespondError(w, errors.BadRequest("Computer name is required", nil))
		return
	}

	if err := h.service.CreateComputer(computer); err != nil {
		logger.Error("Failed to create computer", zap.Error(err), zap.String("computer", computer.Name))
		utils.RespondError(w, errors.InternalServerError("Failed to create computer", err))
		return
	}

	logger.Info("Computer created", zap.String("computer", computer.Name))
	utils.RespondSuccess(w, map[string]string{
		"message": "Computer created successfully",
		"name":    computer.Name,
	})
}

// DeleteComputer deletes an AD computer
func (h *ADDCHandler) DeleteComputer(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	computerName := chi.URLParam(r, "name")
	if computerName == "" {
		utils.RespondError(w, errors.BadRequest("Computer name is required", nil))
		return
	}

	if err := h.service.DeleteComputer(computerName); err != nil {
		logger.Error("Failed to delete computer", zap.Error(err), zap.String("computer", computerName))
		utils.RespondError(w, errors.InternalServerError("Failed to delete computer", err))
		return
	}

	logger.Info("Computer deleted", zap.String("computer", computerName))
	utils.RespondSuccess(w, map[string]string{
		"message": "Computer deleted successfully",
		"name":    computerName,
	})
}

// ===== Organizational Unit Management =====

// ListOUs lists all OUs
func (h *ADDCHandler) ListOUs(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	ous, err := h.service.ListOUs()
	if err != nil {
		logger.Error("Failed to list OUs", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list OUs", err))
		return
	}

	utils.RespondSuccess(w, ous)
}

// CreateOU creates a new OU
func (h *ADDCHandler) CreateOU(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	var ou ad.ADOU
	if err := json.NewDecoder(r.Body).Decode(&ou); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if ou.Name == "" {
		utils.RespondError(w, errors.BadRequest("OU name is required", nil))
		return
	}

	if err := h.service.CreateOU(ou); err != nil {
		logger.Error("Failed to create OU", zap.Error(err), zap.String("ou", ou.Name))
		utils.RespondError(w, errors.InternalServerError("Failed to create OU", err))
		return
	}

	logger.Info("OU created", zap.String("ou", ou.Name))
	utils.RespondSuccess(w, map[string]string{
		"message": "OU created successfully",
		"name":    ou.Name,
	})
}

// DeleteOU deletes an OU
func (h *ADDCHandler) DeleteOU(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	var req struct {
		DN string `json:"dn"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.DN == "" {
		utils.RespondError(w, errors.BadRequest("OU DN is required", nil))
		return
	}

	if err := h.service.DeleteOU(req.DN); err != nil {
		logger.Error("Failed to delete OU", zap.Error(err), zap.String("dn", req.DN))
		utils.RespondError(w, errors.InternalServerError("Failed to delete OU", err))
		return
	}

	logger.Info("OU deleted", zap.String("dn", req.DN))
	utils.RespondSuccess(w, map[string]string{
		"message": "OU deleted successfully",
		"dn":      req.DN,
	})
}

// ===== Group Policy Management =====

// ListGPOs lists all GPOs
func (h *ADDCHandler) ListGPOs(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	gpos, err := h.service.ListGPOs()
	if err != nil {
		logger.Error("Failed to list GPOs", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list GPOs", err))
		return
	}

	utils.RespondSuccess(w, gpos)
}

// CreateGPO creates a new GPO
func (h *ADDCHandler) CreateGPO(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	var gpo ad.ADGPO
	if err := json.NewDecoder(r.Body).Decode(&gpo); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if gpo.Name == "" {
		utils.RespondError(w, errors.BadRequest("GPO name is required", nil))
		return
	}

	if err := h.service.CreateGPO(gpo); err != nil {
		logger.Error("Failed to create GPO", zap.Error(err), zap.String("gpo", gpo.Name))
		utils.RespondError(w, errors.InternalServerError("Failed to create GPO", err))
		return
	}

	logger.Info("GPO created", zap.String("gpo", gpo.Name))
	utils.RespondSuccess(w, map[string]string{
		"message": "GPO created successfully",
		"name":    gpo.Name,
	})
}

// DeleteGPO deletes a GPO
func (h *ADDCHandler) DeleteGPO(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	gpoName := chi.URLParam(r, "name")
	if gpoName == "" {
		utils.RespondError(w, errors.BadRequest("GPO name is required", nil))
		return
	}

	if err := h.service.DeleteGPO(gpoName); err != nil {
		logger.Error("Failed to delete GPO", zap.Error(err), zap.String("gpo", gpoName))
		utils.RespondError(w, errors.InternalServerError("Failed to delete GPO", err))
		return
	}

	logger.Info("GPO deleted", zap.String("gpo", gpoName))
	utils.RespondSuccess(w, map[string]string{
		"message": "GPO deleted successfully",
		"name":    gpoName,
	})
}

// LinkGPO links a GPO to an OU
func (h *ADDCHandler) LinkGPO(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	gpoName := chi.URLParam(r, "name")
	if gpoName == "" {
		utils.RespondError(w, errors.BadRequest("GPO name is required", nil))
		return
	}

	var req struct {
		OUDN string `json:"ou_dn"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.OUDN == "" {
		utils.RespondError(w, errors.BadRequest("OU DN is required", nil))
		return
	}

	if err := h.service.LinkGPO(gpoName, req.OUDN); err != nil {
		logger.Error("Failed to link GPO", zap.Error(err), zap.String("gpo", gpoName), zap.String("ou", req.OUDN))
		utils.RespondError(w, errors.InternalServerError("Failed to link GPO", err))
		return
	}

	logger.Info("GPO linked", zap.String("gpo", gpoName), zap.String("ou", req.OUDN))
	utils.RespondSuccess(w, map[string]string{
		"message": "GPO linked successfully",
		"gpo":     gpoName,
		"ou":      req.OUDN,
	})
}

// UnlinkGPO unlinks a GPO from an OU
func (h *ADDCHandler) UnlinkGPO(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	gpoName := chi.URLParam(r, "name")
	if gpoName == "" {
		utils.RespondError(w, errors.BadRequest("GPO name is required", nil))
		return
	}

	var req struct {
		OUDN string `json:"ou_dn"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.OUDN == "" {
		utils.RespondError(w, errors.BadRequest("OU DN is required", nil))
		return
	}

	if err := h.service.UnlinkGPO(gpoName, req.OUDN); err != nil {
		logger.Error("Failed to unlink GPO", zap.Error(err), zap.String("gpo", gpoName), zap.String("ou", req.OUDN))
		utils.RespondError(w, errors.InternalServerError("Failed to unlink GPO", err))
		return
	}

	logger.Info("GPO unlinked", zap.String("gpo", gpoName), zap.String("ou", req.OUDN))
	utils.RespondSuccess(w, map[string]string{
		"message": "GPO unlinked successfully",
		"gpo":     gpoName,
		"ou":      req.OUDN,
	})
}

// ===== DNS Management =====

// ListDNSZones lists all DNS zones
func (h *ADDCHandler) ListDNSZones(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	zones, err := h.service.ListDNSZones()
	if err != nil {
		logger.Error("Failed to list DNS zones", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list DNS zones", err))
		return
	}

	utils.RespondSuccess(w, zones)
}

// CreateDNSZone creates a new DNS zone
func (h *ADDCHandler) CreateDNSZone(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	var req struct {
		ZoneName string `json:"zone_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.ZoneName == "" {
		utils.RespondError(w, errors.BadRequest("Zone name is required", nil))
		return
	}

	if err := h.service.CreateDNSZone(req.ZoneName); err != nil {
		logger.Error("Failed to create DNS zone", zap.Error(err), zap.String("zone", req.ZoneName))
		utils.RespondError(w, errors.InternalServerError("Failed to create DNS zone", err))
		return
	}

	logger.Info("DNS zone created", zap.String("zone", req.ZoneName))
	utils.RespondSuccess(w, map[string]string{
		"message": "DNS zone created successfully",
		"zone":    req.ZoneName,
	})
}

// DeleteDNSZone deletes a DNS zone
func (h *ADDCHandler) DeleteDNSZone(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	zoneName := chi.URLParam(r, "zone")
	if zoneName == "" {
		utils.RespondError(w, errors.BadRequest("Zone name is required", nil))
		return
	}

	if err := h.service.DeleteDNSZone(zoneName); err != nil {
		logger.Error("Failed to delete DNS zone", zap.Error(err), zap.String("zone", zoneName))
		utils.RespondError(w, errors.InternalServerError("Failed to delete DNS zone", err))
		return
	}

	logger.Info("DNS zone deleted", zap.String("zone", zoneName))
	utils.RespondSuccess(w, map[string]string{
		"message": "DNS zone deleted successfully",
		"zone":    zoneName,
	})
}

// ListDNSRecords lists DNS records in a zone
func (h *ADDCHandler) ListDNSRecords(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	zoneName := chi.URLParam(r, "zone")
	if zoneName == "" {
		utils.RespondError(w, errors.BadRequest("Zone name is required", nil))
		return
	}

	records, err := h.service.ListDNSRecords(zoneName)
	if err != nil {
		logger.Error("Failed to list DNS records", zap.Error(err), zap.String("zone", zoneName))
		utils.RespondError(w, errors.InternalServerError("Failed to list DNS records", err))
		return
	}

	utils.RespondSuccess(w, records)
}

// AddDNSRecord adds a DNS record
func (h *ADDCHandler) AddDNSRecord(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	zoneName := chi.URLParam(r, "zone")
	if zoneName == "" {
		utils.RespondError(w, errors.BadRequest("Zone name is required", nil))
		return
	}

	var record ad.ADDNSRecord
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if record.Name == "" || record.RecordType == "" || record.Value == "" {
		utils.RespondError(w, errors.BadRequest("Name, type, and value are required", nil))
		return
	}

	if err := h.service.AddDNSRecord(zoneName, record); err != nil {
		logger.Error("Failed to add DNS record", zap.Error(err), zap.String("zone", zoneName), zap.String("record", record.Name))
		utils.RespondError(w, errors.InternalServerError("Failed to add DNS record", err))
		return
	}

	logger.Info("DNS record added", zap.String("zone", zoneName), zap.String("record", record.Name))
	utils.RespondSuccess(w, map[string]string{
		"message": "DNS record added successfully",
		"zone":    zoneName,
		"record":  record.Name,
	})
}

// DeleteDNSRecord deletes a DNS record
func (h *ADDCHandler) DeleteDNSRecord(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	zoneName := chi.URLParam(r, "zone")
	recordName := chi.URLParam(r, "record")

	if zoneName == "" || recordName == "" {
		utils.RespondError(w, errors.BadRequest("Zone and record name are required", nil))
		return
	}

	var req struct {
		RecordType string `json:"type"`
		Value      string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.RecordType == "" || req.Value == "" {
		utils.RespondError(w, errors.BadRequest("Record type and value are required", nil))
		return
	}

	if err := h.service.DeleteDNSRecord(zoneName, recordName, req.RecordType, req.Value); err != nil {
		logger.Error("Failed to delete DNS record", zap.Error(err), zap.String("zone", zoneName), zap.String("record", recordName))
		utils.RespondError(w, errors.InternalServerError("Failed to delete DNS record", err))
		return
	}

	logger.Info("DNS record deleted", zap.String("zone", zoneName), zap.String("record", recordName))
	utils.RespondSuccess(w, map[string]string{
		"message": "DNS record deleted successfully",
		"zone":    zoneName,
		"record":  recordName,
	})
}

// ===== FSMO Roles =====

// ShowFSMORoles shows FSMO roles
func (h *ADDCHandler) ShowFSMORoles(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	roles, err := h.service.ShowFSMORoles()
	if err != nil {
		logger.Error("Failed to show FSMO roles", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to show FSMO roles", err))
		return
	}

	utils.RespondSuccess(w, roles)
}

// TransferFSMORoles transfers FSMO roles
func (h *ADDCHandler) TransferFSMORoles(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	var req struct {
		Role     string `json:"role"`
		TargetDC string `json:"target_dc"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Role == "" || req.TargetDC == "" {
		utils.RespondError(w, errors.BadRequest("Role and target DC are required", nil))
		return
	}

	if err := h.service.TransferFSMORoles(req.Role, req.TargetDC); err != nil {
		logger.Error("Failed to transfer FSMO role", zap.Error(err), zap.String("role", req.Role))
		utils.RespondError(w, errors.InternalServerError("Failed to transfer FSMO role", err))
		return
	}

	logger.Info("FSMO role transferred", zap.String("role", req.Role), zap.String("target", req.TargetDC))
	utils.RespondSuccess(w, map[string]string{
		"message":   "FSMO role transferred successfully",
		"role":      req.Role,
		"target_dc": req.TargetDC,
	})
}

// SeizeFSMORoles seizes FSMO roles
func (h *ADDCHandler) SeizeFSMORoles(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	var req struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Role == "" {
		utils.RespondError(w, errors.BadRequest("Role is required", nil))
		return
	}

	if err := h.service.SeizeFSMORoles(req.Role); err != nil {
		logger.Error("Failed to seize FSMO role", zap.Error(err), zap.String("role", req.Role))
		utils.RespondError(w, errors.InternalServerError("Failed to seize FSMO role", err))
		return
	}

	logger.Info("FSMO role seized", zap.String("role", req.Role))
	utils.RespondSuccess(w, map[string]string{
		"message": "FSMO role seized successfully",
		"role":    req.Role,
	})
}

// ===== Utility Functions =====

// TestConfiguration tests the Samba configuration
func (h *ADDCHandler) TestConfiguration(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	if err := h.service.TestConfiguration(); err != nil {
		utils.RespondError(w, errors.InternalServerError("Configuration test failed", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Configuration is valid",
	})
}

// ShowDBCheck runs database check
func (h *ADDCHandler) ShowDBCheck(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	result, err := h.service.ShowDBCheck()
	if err != nil {
		logger.Error("Database check failed", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Database check failed", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"result": result,
	})
}

// BackupOnline performs an online backup
func (h *ADDCHandler) BackupOnline(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"AD DC service not available",
			nil,
		))
		return
	}

	var req struct {
		TargetDir string `json:"target_dir"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.TargetDir == "" {
		utils.RespondError(w, errors.BadRequest("Target directory is required", nil))
		return
	}

	if err := h.service.BackupOnline(req.TargetDir); err != nil {
		logger.Error("Backup failed", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Backup failed", err))
		return
	}

	logger.Info("Backup completed", zap.String("target", req.TargetDir))
	utils.RespondSuccess(w, map[string]string{
		"message":    "Backup completed successfully",
		"target_dir": req.TargetDir,
	})
}
