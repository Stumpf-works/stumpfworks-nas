package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/stumpfworks/nas/internal/api/handlers"
	mw "github.com/stumpfworks/nas/internal/api/middleware"
	"github.com/stumpfworks/nas/internal/config"
)

// NewRouter creates and configures the HTTP router
func NewRouter(cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mw.LoggerMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"}, // Vite dev server
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check (no auth required)
	r.Get("/health", handlers.HealthCheck)
	r.Get("/", handlers.IndexHandler)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes (no auth)
		r.Group(func(r chi.Router) {
			r.Post("/auth/login", handlers.Login)
			// r.Post("/auth/register", handlers.Register) // Will implement later
		})

		// Protected routes (auth required)
		r.Group(func(r chi.Router) {
			r.Use(mw.AuthMiddleware)

			// Auth routes
			r.Post("/auth/logout", handlers.Logout)
			r.Post("/auth/refresh", handlers.RefreshToken)
			r.Get("/auth/me", handlers.GetCurrentUser)

			// System routes
			r.Get("/system/info", handlers.GetSystemInfo)
			r.Get("/system/metrics", handlers.GetSystemMetrics)

			// User routes (admin only for now)
			r.Route("/users", func(r chi.Router) {
				r.Use(mw.AdminOnly)
				r.Get("/", handlers.ListUsers)
				r.Post("/", handlers.CreateUser)
				r.Get("/{id}", handlers.GetUser)
				r.Put("/{id}", handlers.UpdateUser)
				r.Delete("/{id}", handlers.DeleteUser)
			})

			// Storage routes (will implement in next phase)
			// r.Route("/storage", func(r chi.Router) {
			// 	r.Get("/disks", handlers.ListDisks)
			// 	r.Get("/volumes", handlers.ListVolumes)
			// })

			// Network routes (will implement in next phase)
			// r.Route("/network", func(r chi.Router) {
			// 	r.Get("/interfaces", handlers.ListInterfaces)
			// })

			// Plugin routes (will implement in next phase)
			// r.Route("/plugins", func(r chi.Router) {
			// 	r.Get("/", handlers.ListPlugins)
			// })
		})
	})

	// WebSocket endpoint
	r.Get("/ws", handlers.WebSocketHandler)

	return r
}
