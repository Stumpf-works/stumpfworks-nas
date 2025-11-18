package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
	"github.com/stumpf-works/stumpfworks-nas/plugins/asterisk-voip/ami"
)

// Server represents the API server
type Server struct {
	router    *chi.Mux
	amiClient *ami.Client
	port      string
}

// NewServer creates a new API server
func NewServer(amiClient *ami.Client, port string) *Server {
	s := &Server{
		router:    chi.NewRouter(),
		amiClient: amiClient,
		port:      port,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

// setupMiddleware configures middleware
func (s *Server) setupMiddleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))

	// CORS
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:8080", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

// setupRoutes configures API routes
func (s *Server) setupRoutes() {
	s.router.Get("/health", s.healthCheck)
	s.router.Get("/version", s.getVersion)

	// API v1
	s.router.Route("/api/v1", func(r chi.Router) {
		// Status & Info
		r.Get("/status", s.getStatus)
		r.Get("/ami/status", s.getAMIStatus)

		// Extensions
		r.Route("/extensions", func(r chi.Router) {
			r.Get("/", s.listExtensions)
			r.Post("/", s.createExtension)
			r.Get("/{id}", s.getExtension)
			r.Put("/{id}", s.updateExtension)
			r.Delete("/{id}", s.deleteExtension)
		})

		// Trunks
		r.Route("/trunks", func(r chi.Router) {
			r.Get("/", s.listTrunks)
			r.Post("/", s.createTrunk)
			r.Get("/{id}", s.getTrunk)
			r.Put("/{id}", s.updateTrunk)
			r.Delete("/{id}", s.deleteTrunk)
		})

		// Active Calls
		r.Route("/calls", func(r chi.Router) {
			r.Get("/", s.listActiveCalls)
			r.Post("/{id}/hangup", s.hangupCall)
			r.Post("/originate", s.originateCall)
		})

		// Voicemail
		r.Route("/voicemail", func(r chi.Router) {
			r.Get("/boxes", s.listVoicemailBoxes)
			r.Get("/boxes/{id}/messages", s.getVoicemailMessages)
		})

		// Recordings
		r.Route("/recordings", func(r chi.Router) {
			r.Get("/", s.listRecordings)
			r.Get("/{id}", s.getRecording)
			r.Delete("/{id}", s.deleteRecording)
		})

		// Configuration
		r.Route("/config", func(r chi.Router) {
			r.Get("/", s.getConfig)
			r.Put("/", s.updateConfig)
			r.Post("/reload", s.reloadConfig)
		})
	})
}

// Start starts the API server
func (s *Server) Start() error {
	log.Info().Str("port", s.port).Msg("Starting API server")
	return http.ListenAndServe(":"+s.port, s.router)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down API server")
	return nil
}

// Helper functions

// respondJSON sends JSON response
func (s *Server) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// respondError sends error response
func (s *Server) respondError(w http.ResponseWriter, status int, message string) {
	s.respondJSON(w, status, map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

// respondSuccess sends success response
func (s *Server) respondSuccess(w http.ResponseWriter, data interface{}) {
	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

// Health check
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"ami": map[string]interface{}{
			"connected": s.amiClient.IsConnected(),
		},
	}

	s.respondJSON(w, http.StatusOK, health)
}

// Get version
func (s *Server) getVersion(w http.ResponseWriter, r *http.Request) {
	version, err := s.amiClient.CoreShowVersion()
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to get Asterisk version")
		return
	}

	s.respondSuccess(w, map[string]string{
		"version": version,
	})
}
