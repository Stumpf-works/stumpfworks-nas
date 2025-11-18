package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// Extension represents a SIP extension
type Extension struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Secret      string `json:"secret,omitempty"`
	Context     string `json:"context"`
	CallerID    string `json:"caller_id"`
	Mailbox     string `json:"mailbox"`
	Status      string `json:"status"` // online, offline, busy
	IPAddress   string `json:"ip_address,omitempty"`
	UserAgent   string `json:"user_agent,omitempty"`
	LastSeen    int64  `json:"last_seen,omitempty"`
	Description string `json:"description,omitempty"`
}

// Trunk represents a SIP trunk
type Trunk struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"` // peer, user, friend
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Username    string `json:"username,omitempty"`
	Secret      string `json:"secret,omitempty"`
	Context     string `json:"context"`
	Status      string `json:"status"` // registered, unreachable, lagged
	FromDomain  string `json:"from_domain,omitempty"`
	Description string `json:"description,omitempty"`
}

// Call represents an active call
type Call struct {
	ID          string `json:"id"`
	Channel     string `json:"channel"`
	State       string `json:"state"`
	CallerID    string `json:"caller_id"`
	CallerName  string `json:"caller_name"`
	Extension   string `json:"extension"`
	Context     string `json:"context"`
	Duration    int    `json:"duration"`
	StartTime   int64  `json:"start_time"`
	Application string `json:"application"`
	Data        string `json:"data"`
}

// VoicemailBox represents a voicemail box
type VoicemailBox struct {
	ID       string `json:"id"`
	Context  string `json:"context"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	NewCount int    `json:"new_count"`
	OldCount int    `json:"old_count"`
}

// VoicemailMessage represents a voicemail message
type VoicemailMessage struct {
	ID       string `json:"id"`
	Box      string `json:"box"`
	From     string `json:"from"`
	Duration int    `json:"duration"`
	Date     int64  `json:"date"`
	Read     bool   `json:"read"`
	FilePath string `json:"file_path"`
}

// Recording represents a call recording
type Recording struct {
	ID       string `json:"id"`
	FileName string `json:"file_name"`
	Channel  string `json:"channel"`
	Duration int    `json:"duration"`
	Date     int64  `json:"date"`
	Size     int64  `json:"size"`
	FilePath string `json:"file_path"`
}

// Status Handlers

func (s *Server) getStatus(w http.ResponseWriter, r *http.Request) {
	// Get system status from Asterisk
	status := map[string]interface{}{
		"asterisk": map[string]interface{}{
			"running":    true,
			"ami":        s.amiClient.IsConnected(),
			"uptime":     "24h 15m",
			"calls":      0,
			"channels":   0,
			"extensions": 0,
		},
		"services": map[string]interface{}{
			"sip":        "running",
			"voicemail":  "running",
			"recording":  "running",
			"conference": "running",
		},
	}

	s.respondSuccess(w, status)
}

func (s *Server) getAMIStatus(w http.ResponseWriter, r *http.Request) {
	s.respondSuccess(w, map[string]interface{}{
		"connected": s.amiClient.IsConnected(),
	})
}

// Extension Handlers

func (s *Server) listExtensions(w http.ResponseWriter, r *http.Request) {
	// TODO: Read from sip.conf and get status from AMI
	extensions := []Extension{
		{
			ID:       "1000",
			Name:     "Administrator",
			Context:  "internal",
			CallerID: "Administrator <1000>",
			Mailbox:  "1000@default",
			Status:   "offline",
		},
	}

	s.respondSuccess(w, extensions)
}

func (s *Server) createExtension(w http.ResponseWriter, r *http.Request) {
	var ext Extension
	if err := json.NewDecoder(r.Body).Decode(&ext); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	log.Info().Str("id", ext.ID).Msg("Creating extension")

	// TODO: Write to sip.conf and reload
	// For now, just return success
	s.respondSuccess(w, ext)
}

func (s *Server) getExtension(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// TODO: Read from sip.conf
	ext := Extension{
		ID:       id,
		Name:     "Test Extension",
		Context:  "internal",
		CallerID: fmt.Sprintf("Extension %s", id),
		Mailbox:  fmt.Sprintf("%s@default", id),
		Status:   "offline",
	}

	s.respondSuccess(w, ext)
}

func (s *Server) updateExtension(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var ext Extension
	if err := json.NewDecoder(r.Body).Decode(&ext); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ext.ID = id
	log.Info().Str("id", id).Msg("Updating extension")

	// TODO: Update sip.conf and reload
	s.respondSuccess(w, ext)
}

func (s *Server) deleteExtension(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	log.Info().Str("id", id).Msg("Deleting extension")

	// TODO: Remove from sip.conf and reload
	s.respondSuccess(w, map[string]string{
		"message": "Extension deleted successfully",
	})
}

// Trunk Handlers

func (s *Server) listTrunks(w http.ResponseWriter, r *http.Request) {
	// TODO: Read from sip.conf
	trunks := []Trunk{}
	s.respondSuccess(w, trunks)
}

func (s *Server) createTrunk(w http.ResponseWriter, r *http.Request) {
	var trunk Trunk
	if err := json.NewDecoder(r.Body).Decode(&trunk); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	log.Info().Str("id", trunk.ID).Msg("Creating trunk")

	// TODO: Write to sip.conf and reload
	s.respondSuccess(w, trunk)
}

func (s *Server) getTrunk(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// TODO: Read from sip.conf
	trunk := Trunk{
		ID:      id,
		Name:    "Example Trunk",
		Type:    "peer",
		Host:    "sip.provider.com",
		Port:    5060,
		Context: "from-trunk",
		Status:  "unreachable",
	}

	s.respondSuccess(w, trunk)
}

func (s *Server) updateTrunk(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var trunk Trunk
	if err := json.NewDecoder(r.Body).Decode(&trunk); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	trunk.ID = id
	log.Info().Str("id", id).Msg("Updating trunk")

	// TODO: Update sip.conf and reload
	s.respondSuccess(w, trunk)
}

func (s *Server) deleteTrunk(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	log.Info().Str("id", id).Msg("Deleting trunk")

	// TODO: Remove from sip.conf and reload
	s.respondSuccess(w, map[string]string{
		"message": "Trunk deleted successfully",
	})
}

// Call Handlers

func (s *Server) listActiveCalls(w http.ResponseWriter, r *http.Request) {
	// TODO: Get from AMI
	calls := []Call{}
	s.respondSuccess(w, calls)
}

func (s *Server) hangupCall(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	log.Info().Str("channel", id).Msg("Hanging up call")

	err := s.amiClient.Hangup(id)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to hangup call")
		return
	}

	s.respondSuccess(w, map[string]string{
		"message": "Call hung up successfully",
	})
}

func (s *Server) originateCall(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Channel   string `json:"channel"`
		Extension string `json:"extension"`
		Context   string `json:"context"`
		Timeout   int    `json:"timeout"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	log.Info().
		Str("channel", req.Channel).
		Str("extension", req.Extension).
		Msg("Originating call")

	err := s.amiClient.Originate(req.Channel, req.Extension, req.Context, req.Timeout)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to originate call")
		return
	}

	s.respondSuccess(w, map[string]string{
		"message": "Call originated successfully",
	})
}

// Voicemail Handlers

func (s *Server) listVoicemailBoxes(w http.ResponseWriter, r *http.Request) {
	// TODO: Read from voicemail.conf
	boxes := []VoicemailBox{}
	s.respondSuccess(w, boxes)
}

func (s *Server) getVoicemailMessages(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	log.Info().Str("box", id).Msg("Getting voicemail messages")

	// TODO: Read from filesystem
	messages := []VoicemailMessage{}
	s.respondSuccess(w, messages)
}

// Recording Handlers

func (s *Server) listRecordings(w http.ResponseWriter, r *http.Request) {
	// TODO: Read from filesystem
	recordings := []Recording{}
	s.respondSuccess(w, recordings)
}

func (s *Server) getRecording(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	log.Info().Str("id", id).Msg("Getting recording")

	// TODO: Stream recording file
	http.ServeFile(w, r, "/recordings/"+id)
}

func (s *Server) deleteRecording(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	log.Info().Str("id", id).Msg("Deleting recording")

	// TODO: Delete file
	s.respondSuccess(w, map[string]string{
		"message": "Recording deleted successfully",
	})
}

// Config Handlers

func (s *Server) getConfig(w http.ResponseWriter, r *http.Request) {
	// TODO: Read plugin.json config
	config := map[string]interface{}{
		"ami": map[string]interface{}{
			"host": "localhost",
			"port": 5038,
		},
		"sip": map[string]interface{}{
			"port": 5060,
		},
	}

	s.respondSuccess(w, config)
}

func (s *Server) updateConfig(w http.ResponseWriter, r *http.Request) {
	var config map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	log.Info().Msg("Updating configuration")

	// TODO: Write to plugin.json
	s.respondSuccess(w, config)
}

func (s *Server) reloadConfig(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Reloading Asterisk configuration")

	err := s.amiClient.Reload("")
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to reload configuration")
		return
	}

	s.respondSuccess(w, map[string]string{
		"message": "Configuration reloaded successfully",
	})
}
