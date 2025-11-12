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

			// File Management routes
			r.Route("/files", func(r chi.Router) {
				// File browsing and info
				r.Get("/browse", handlers.BrowseFiles)
				r.Get("/info", handlers.GetFileInfo)
				r.Get("/download", handlers.DownloadFile)
				r.Get("/usage", handlers.GetDiskUsage)

				// File operations (write access required)
				r.Post("/upload", handlers.UploadFile)
				r.Post("/mkdir", handlers.CreateDirectory)
				r.Post("/rename", handlers.RenameFile)
				r.Post("/copy", handlers.CopyFiles)
				r.Post("/move", handlers.MoveFiles)
				r.Delete("/delete", handlers.DeleteFiles)

				// Chunked upload
				r.Post("/upload/start", handlers.StartChunkedUpload)
				r.Post("/upload/{sessionId}/chunk/{chunkIndex}", handlers.UploadChunk)
				r.Post("/upload/finalize", handlers.FinalizeUpload)
				r.Delete("/upload/{sessionId}", handlers.CancelUpload)
				r.Get("/upload/{sessionId}", handlers.GetUploadSession)

				// Archives
				r.Post("/archive/create", handlers.CreateArchive)
				r.Post("/archive/extract", handlers.ExtractArchive)

				// Permissions (admin only)
				r.Group(func(r chi.Router) {
					r.Use(mw.AdminOnly)
					r.Get("/permissions", handlers.GetFilePermissions)
					r.Post("/permissions", handlers.ChangeFilePermissions)
				})
			})

			// Network routes
			r.Route("/network", func(r chi.Router) {
				netHandler := handlers.NewNetworkHandler()

				// Interface management
				r.Get("/interfaces", netHandler.ListInterfaces)
				r.Get("/interfaces/stats", netHandler.GetInterfaceStats)

				// Routes and DNS
				r.Get("/routes", netHandler.GetRoutes)
				r.Get("/dns", netHandler.GetDNS)

				// Firewall (read-only)
				r.Get("/firewall", netHandler.GetFirewallStatus)

				// Diagnostics
				r.Post("/diagnostics/ping", netHandler.Ping)
				r.Post("/diagnostics/traceroute", netHandler.Traceroute)
				r.Post("/diagnostics/netstat", netHandler.Netstat)

				// Admin-only network operations
				r.Group(func(r chi.Router) {
					r.Use(mw.AdminOnly)

					// Interface configuration
					r.Post("/interfaces/{name}/state", netHandler.SetInterfaceState)
					r.Post("/interfaces/{name}/configure", netHandler.ConfigureInterface)

					// DNS configuration
					r.Post("/dns", netHandler.SetDNS)

					// Firewall management
					r.Post("/firewall/state", netHandler.SetFirewallState)
					r.Post("/firewall/rules", netHandler.AddFirewallRule)
					r.Delete("/firewall/rules/{number}", netHandler.DeleteFirewallRule)
					r.Post("/firewall/default", netHandler.SetDefaultPolicy)
					r.Post("/firewall/reset", netHandler.ResetFirewall)

					// Wake-on-LAN
					r.Post("/wol", netHandler.WakeOnLAN)
				})
			})

			// Docker routes
			r.Route("/docker", func(r chi.Router) {
				dockerHandler := handlers.NewDockerHandler()
				r.Use(dockerHandler.CheckAvailability)

				// Container routes
				r.Get("/containers", dockerHandler.ListContainers)
				r.Post("/containers", dockerHandler.CreateContainer)
				r.Get("/containers/{id}", dockerHandler.InspectContainer)
				r.Get("/containers/{id}/stats", dockerHandler.GetContainerStats)
				r.Get("/containers/{id}/logs", dockerHandler.GetContainerLogs)
				r.Post("/containers/{id}/start", dockerHandler.StartContainer)
				r.Post("/containers/{id}/stop", dockerHandler.StopContainer)
				r.Post("/containers/{id}/restart", dockerHandler.RestartContainer)
				r.Post("/containers/{id}/pause", dockerHandler.PauseContainer)
				r.Post("/containers/{id}/unpause", dockerHandler.UnpauseContainer)
				r.Delete("/containers/{id}", dockerHandler.RemoveContainer)

				// Image routes
				r.Get("/images", dockerHandler.ListImages)
				r.Get("/images/search", dockerHandler.SearchImages)
				r.Post("/images/pull", dockerHandler.PullImage)
				r.Get("/images/{id}", dockerHandler.InspectImage)
				r.Delete("/images/{id}", dockerHandler.RemoveImage)

				// Volume routes
				r.Get("/volumes", dockerHandler.ListVolumes)
				r.Post("/volumes", dockerHandler.CreateVolume)
				r.Get("/volumes/{name}", dockerHandler.InspectVolume)
				r.Delete("/volumes/{name}", dockerHandler.RemoveVolume)

				// Network routes
				r.Get("/networks", dockerHandler.ListNetworks)
				r.Post("/networks", dockerHandler.CreateNetwork)
				r.Get("/networks/{id}", dockerHandler.InspectNetwork)
				r.Delete("/networks/{id}", dockerHandler.RemoveNetwork)
				r.Post("/networks/{id}/connect", dockerHandler.ConnectContainerToNetwork)
				r.Post("/networks/{id}/disconnect", dockerHandler.DisconnectContainerFromNetwork)

				// System routes
				r.Get("/info", dockerHandler.GetDockerInfo)
				r.Get("/version", dockerHandler.GetDockerVersion)
				r.Post("/system/prune", dockerHandler.PruneSystem)

				// Docker Compose Stack routes
				composeHandler := handlers.NewComposeHandler("")
				r.Get("/stacks", composeHandler.ListStacks)
				r.Post("/stacks", composeHandler.CreateStack)
				r.Get("/stacks/{name}", composeHandler.GetStack)
				r.Put("/stacks/{name}", composeHandler.UpdateStack)
				r.Delete("/stacks/{name}", composeHandler.DeleteStack)
				r.Post("/stacks/{name}/deploy", composeHandler.DeployStack)
				r.Post("/stacks/{name}/stop", composeHandler.StopStack)
				r.Post("/stacks/{name}/restart", composeHandler.RestartStack)
				r.Post("/stacks/{name}/remove", composeHandler.RemoveStack)
				r.Get("/stacks/{name}/logs", composeHandler.GetStackLogs)
				r.Get("/stacks/{name}/compose", composeHandler.GetComposeFile)
			})

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
