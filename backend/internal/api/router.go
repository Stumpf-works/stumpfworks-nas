package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/Stumpf-works/stumpfworks-nas/internal/api/handlers"
	mw "github.com/Stumpf-works/stumpfworks-nas/internal/api/middleware"
	"github.com/Stumpf-works/stumpfworks-nas/internal/config"
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
			r.Get("/system/updates", handlers.CheckForUpdates)

			// Admin-only system routes
			r.Route("/system", func(r chi.Router) {
				r.Use(mw.AdminOnly)
				r.Post("/updates", handlers.ApplyUpdates)
			})

			// User routes (admin only for now)
			r.Route("/users", func(r chi.Router) {
				r.Use(mw.AdminOnly)
				r.Get("/", handlers.ListUsers)
				r.Post("/", handlers.CreateUser)
				r.Get("/{id}", handlers.GetUser)
				r.Put("/{id}", handlers.UpdateUser)
				r.Delete("/{id}", handlers.DeleteUser)
			})

			// Storage routes
			r.Route("/storage", func(r chi.Router) {
				// Statistics and overview
				r.Get("/stats", handlers.GetStorageStats)
				r.Get("/health", handlers.GetAllDiskHealth)
				r.Get("/io", handlers.GetDiskIOStats)
				r.Get("/io/monitor", handlers.GetIOMonitorStats)

				// Disks
				r.Get("/disks", handlers.ListDisks)
				r.Get("/disks/{name}", handlers.GetDisk)
				r.Get("/disks/{name}/smart", handlers.GetDiskSMART)
				r.Get("/disks/{name}/health", handlers.GetDiskHealth)
				r.Get("/disks/{name}/io", handlers.GetDiskIOStatsForDisk)

				// Volumes
				r.Get("/volumes", handlers.ListVolumes)
				r.Get("/volumes/{id}", handlers.GetVolume)

				// Shares
				r.Get("/shares", handlers.ListShares)
				r.Get("/shares/{id}", handlers.GetShare)

				// Admin-only storage operations
				r.Group(func(r chi.Router) {
					r.Use(mw.AdminOnly)

					// Disk operations
					r.Post("/disks/format", handlers.FormatDisk)

					// Volume operations
					r.Post("/volumes", handlers.CreateVolume)
					r.Delete("/volumes/{id}", handlers.DeleteVolume)

					// Share operations
					r.Post("/shares", handlers.CreateShare)
					r.Put("/shares/{id}", handlers.UpdateShare)
					r.Delete("/shares/{id}", handlers.DeleteShare)
					r.Post("/shares/{id}/enable", handlers.EnableShare)
					r.Post("/shares/{id}/disable", handlers.DisableShare)
				})
			})

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
