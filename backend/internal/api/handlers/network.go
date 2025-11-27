// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Stumpf-works/stumpfworks-nas/internal/network"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"

	"github.com/go-chi/chi/v5"
)

// NetworkHandler handles network-related requests
type NetworkHandler struct{}

// NewNetworkHandler creates a new network handler
func NewNetworkHandler() *NetworkHandler {
	return &NetworkHandler{}
}

// ListInterfaces handles GET /api/network/interfaces
func (h *NetworkHandler) ListInterfaces(w http.ResponseWriter, r *http.Request) {
	interfaces, err := network.ListInterfaces()
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to list interfaces", err))
		return
	}

	utils.RespondSuccess(w, interfaces)
}

// GetInterfaceStats handles GET /api/network/interfaces/stats
func (h *NetworkHandler) GetInterfaceStats(w http.ResponseWriter, r *http.Request) {
	stats, err := network.GetInterfaceStats()
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to get interface stats", err))
		return
	}

	utils.RespondSuccess(w, stats)
}

// SetInterfaceState handles POST /api/network/interfaces/{name}/state
func (h *NetworkHandler) SetInterfaceState(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var req struct {
		State string `json:"state"` // "up" or "down"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	var err error
	if req.State == "up" {
		err = network.SetInterfaceUp(name)
	} else if req.State == "down" {
		err = network.SetInterfaceDown(name)
	} else {
		utils.RespondError(w, errors.BadRequest("Invalid state", nil))
		return
	}

	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to set interface state", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "Interface state updated"})
}

// ConfigureInterface handles POST /api/network/interfaces/{name}/configure
func (h *NetworkHandler) ConfigureInterface(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var req struct {
		Mode    string `json:"mode"`    // "static" or "dhcp"
		Address string `json:"address"` // for static
		Netmask string `json:"netmask"` // for static
		Gateway string `json:"gateway"` // for static
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	var err error
	if req.Mode == "static" {
		if req.Address == "" || req.Netmask == "" {
			utils.RespondError(w, errors.BadRequest("Missing required fields", nil))
			return
		}
		err = network.ConfigureStaticIP(name, req.Address, req.Netmask, req.Gateway)
	} else if req.Mode == "dhcp" {
		err = network.ConfigureDHCP(name)
	} else {
		utils.RespondError(w, errors.BadRequest("Invalid mode", nil))
		return
	}

	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to configure interface", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "Interface configured successfully"})
}

// GetRoutes handles GET /api/network/routes
func (h *NetworkHandler) GetRoutes(w http.ResponseWriter, r *http.Request) {
	routes, err := network.GetRoutes()
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to get routes", err))
		return
	}

	utils.RespondSuccess(w, routes)
}

// AddRoute handles POST /api/network/routes
func (h *NetworkHandler) AddRoute(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Destination string `json:"destination"`
		Gateway     string `json:"gateway"`
		Interface   string `json:"interface"`
		Metric      int    `json:"metric"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	// Validate destination
	if req.Destination == "" {
		utils.RespondError(w, errors.BadRequest("Destination is required", nil))
		return
	}

	// Either gateway or interface must be provided
	if req.Gateway == "" && req.Interface == "" {
		utils.RespondError(w, errors.BadRequest("Either gateway or interface must be provided", nil))
		return
	}

	if err := network.AddRoute(req.Destination, req.Gateway, req.Interface, req.Metric); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to add route", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "Route added successfully"})
}

// DeleteRoute handles DELETE /api/network/routes
func (h *NetworkHandler) DeleteRoute(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Destination string `json:"destination"`
		Gateway     string `json:"gateway"`
		Interface   string `json:"interface"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	// Validate destination
	if req.Destination == "" {
		utils.RespondError(w, errors.BadRequest("Destination is required", nil))
		return
	}

	if err := network.DeleteRoute(req.Destination, req.Gateway, req.Interface); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to delete route", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "Route deleted successfully"})
}

// GetDNS handles GET /api/network/dns
func (h *NetworkHandler) GetDNS(w http.ResponseWriter, r *http.Request) {
	config, err := network.GetDNSConfig()
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to get DNS config", err))
		return
	}

	utils.RespondSuccess(w, config)
}

// SetDNS handles POST /api/network/dns
func (h *NetworkHandler) SetDNS(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Nameservers   []string `json:"nameservers"`
		SearchDomains []string `json:"searchDomains"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	if err := network.SetDNSConfig(req.Nameservers, req.SearchDomains); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to set DNS config", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "DNS configuration updated"})
}

// GetFirewallStatus handles GET /api/network/firewall
func (h *NetworkHandler) GetFirewallStatus(w http.ResponseWriter, r *http.Request) {
	status, err := network.GetFirewallStatus()
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to get firewall status", err))
		return
	}

	utils.RespondSuccess(w, status)
}

// SetFirewallState handles POST /api/network/firewall/state
func (h *NetworkHandler) SetFirewallState(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Enabled bool `json:"enabled"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	var err error
	if req.Enabled {
		err = network.EnableFirewall()
	} else {
		err = network.DisableFirewall()
	}

	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to set firewall state", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "Firewall state updated"})
}

// AddFirewallRule handles POST /api/network/firewall/rules
func (h *NetworkHandler) AddFirewallRule(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Action   string `json:"action"`   // allow, deny, reject
		Port     string `json:"port"`
		Protocol string `json:"protocol"` // tcp, udp
		From     string `json:"from"`
		To       string `json:"to"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	if err := network.AddFirewallRule(req.Action, req.Port, req.Protocol, req.From, req.To); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to add firewall rule", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "Firewall rule added"})
}

// DeleteFirewallRule handles DELETE /api/network/firewall/rules/{number}
func (h *NetworkHandler) DeleteFirewallRule(w http.ResponseWriter, r *http.Request) {
	numberStr := chi.URLParam(r, "number")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid rule number", err))
		return
	}

	if err := network.DeleteFirewallRule(number); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to delete firewall rule", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "Firewall rule deleted"})
}

// SetDefaultPolicy handles POST /api/network/firewall/default
func (h *NetworkHandler) SetDefaultPolicy(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Direction string `json:"direction"` // incoming, outgoing, routed
		Policy    string `json:"policy"`    // allow, deny, reject
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	if err := network.SetDefaultPolicy(req.Direction, req.Policy); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to set default policy", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "Default policy updated"})
}

// ResetFirewall handles POST /api/network/firewall/reset
func (h *NetworkHandler) ResetFirewall(w http.ResponseWriter, r *http.Request) {
	if err := network.ResetFirewall(); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to reset firewall", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "Firewall reset successfully"})
}

// Ping handles POST /api/network/diagnostics/ping
func (h *NetworkHandler) Ping(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Host  string `json:"host"`
		Count int    `json:"count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	if req.Count <= 0 {
		req.Count = 4
	}

	result, err := network.Ping(req.Host, req.Count)
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to ping", err))
		return
	}

	utils.RespondSuccess(w, result)
}

// Traceroute handles POST /api/network/diagnostics/traceroute
func (h *NetworkHandler) Traceroute(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Host string `json:"host"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	result, err := network.Traceroute(req.Host)
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to traceroute", err))
		return
	}

	utils.RespondSuccess(w, result)
}

// Netstat handles POST /api/network/diagnostics/netstat
func (h *NetworkHandler) Netstat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Options string `json:"options"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	result, err := network.Netstat(req.Options)
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to run netstat", err))
		return
	}

	utils.RespondSuccess(w, result)
}

// WakeOnLAN handles POST /api/network/wol
func (h *NetworkHandler) WakeOnLAN(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MacAddress string `json:"macAddress"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	if err := network.WakeOnLAN(req.MacAddress); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to send WOL packet", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "Wake-on-LAN packet sent"})
}
