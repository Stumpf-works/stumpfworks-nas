// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/embedfs"
	"github.com/Stumpf-works/stumpfworks-nas/internal/api/handlers"
	mw "github.com/Stumpf-works/stumpfworks-nas/internal/api/middleware"
	"github.com/Stumpf-works/stumpfworks-nas/internal/config"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

// NewRouter creates and configures the HTTP router
func NewRouter(cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mw.LoggerMiddleware)
	r.Use(mw.RevisionMiddleware) // Add version headers to all responses
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS middleware - auto-detect origins in development
	var corsHandler *cors.Cors

	if cfg.IsDevelopment() {
		// Development mode: Allow all origins for local network access
		// This enables access from any device on the local network
		corsHandler = cors.New(cors.Options{
			AllowOriginFunc: func(r *http.Request, origin string) bool {
				// Always allow requests with no origin (same-origin or tools like curl)
				if origin == "" {
					return true
				}

				// Allow localhost in any form
				if strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1") {
					return true
				}

				// Allow private network ranges (RFC 1918)
				// 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
				if strings.Contains(origin, "192.168.") ||
					strings.Contains(origin, "10.") ||
					strings.Contains(origin, "172.16.") || strings.Contains(origin, "172.17.") ||
					strings.Contains(origin, "172.18.") || strings.Contains(origin, "172.19.") ||
					strings.Contains(origin, "172.20.") || strings.Contains(origin, "172.21.") ||
					strings.Contains(origin, "172.22.") || strings.Contains(origin, "172.23.") ||
					strings.Contains(origin, "172.24.") || strings.Contains(origin, "172.25.") ||
					strings.Contains(origin, "172.26.") || strings.Contains(origin, "172.27.") ||
					strings.Contains(origin, "172.28.") || strings.Contains(origin, "172.29.") ||
					strings.Contains(origin, "172.30.") || strings.Contains(origin, "172.31.") {
					return true
				}

				logger.Warn("CORS: Blocked origin in development mode", zap.String("origin", origin))
				return false
			},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		})
		logger.Info("CORS: Development mode - allowing all private network origins")
	} else {
		// Production mode: Use configured origins only
		allowedOrigins := cfg.Server.AllowedOrigins
		if len(allowedOrigins) == 0 {
			logger.Error("No CORS origins configured in production mode!")
			allowedOrigins = []string{} // Empty = block all
		}

		corsHandler = cors.New(cors.Options{
			AllowedOrigins:   allowedOrigins,
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		})
		logger.Info("CORS: Production mode - using configured origins", zap.Strings("origins", allowedOrigins))
	}

	r.Use(corsHandler.Handler)

	// Health check (no auth required)
	r.Get("/health", handlers.HealthCheck)

	// Prometheus metrics endpoint (no auth required for monitoring systems)
	r.Get("/metrics", handlers.PrometheusMetricsHandler)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Setup wizard routes (no auth required, always accessible)
		r.Group(func(r chi.Router) {
			r.Get("/setup/status", handlers.SetupStatus)
			r.Post("/setup/initialize", handlers.InitializeSetup)
		})

		// Public routes (no auth, but with IP blocking check)
		r.Group(func(r chi.Router) {
			r.Use(mw.IPBlockMiddleware)
			r.Post("/auth/login", handlers.Login)
			r.Post("/auth/login/2fa", handlers.LoginWith2FA)
			// r.Post("/auth/register", handlers.Register) // Will implement later
		})

		// Protected routes (auth required + setup check)
		r.Group(func(r chi.Router) {
			r.Use(mw.SetupRequired)
			r.Use(mw.AuthMiddleware)

			// Auth routes
			r.Post("/auth/logout", handlers.Logout)
			r.Post("/auth/refresh", handlers.RefreshToken)
			r.Get("/auth/me", handlers.GetCurrentUser)

			// System routes
			r.Get("/system/info", handlers.GetSystemInfo)
			r.Get("/system/metrics", handlers.GetSystemMetrics)

			// Update routes
			updateHandler := handlers.NewUpdateHandler()
			r.Get("/system/version", updateHandler.GetCurrentVersion)
			r.Get("/system/check-updates", updateHandler.CheckForUpdates)

			// Metrics and monitoring routes
			r.Route("/metrics", func(r chi.Router) {
				metricsHandler := handlers.NewMetricsHandler()

				r.Get("/history", metricsHandler.GetMetricsHistory)
				r.Get("/latest", metricsHandler.GetLatestMetric)
				r.Get("/trends", metricsHandler.GetTrends)
			})

			r.Route("/health", func(r chi.Router) {
				metricsHandler := handlers.NewMetricsHandler()

				r.Get("/scores", metricsHandler.GetHealthScores)
				r.Get("/score", metricsHandler.GetLatestHealthScore)
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

			// User Group routes (admin only)
			r.Route("/groups", func(r chi.Router) {
				r.Use(mw.AdminOnly)
				r.Get("/", handlers.ListGroups)
				r.Post("/", handlers.CreateGroup)
				r.Get("/{id}", handlers.GetGroup)
				r.Put("/{id}", handlers.UpdateGroup)
				r.Delete("/{id}", handlers.DeleteGroup)

				// Group member management
				r.Post("/{id}/members", handlers.AddGroupMember)
				r.Delete("/{id}/members/{userId}", handlers.RemoveGroupMember)
				r.Get("/{id}/members", handlers.GetGroupMembers)
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
					r.Put("/disks/{name}/label", handlers.SetDiskLabel)

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

			// System Library routes (Phase 1 integration)
			r.Route("/syslib", func(r chi.Router) {
				r.Use(mw.AdminOnly) // All system library operations require admin

				// System Library Health
				r.Get("/health", handlers.SystemLibraryHealth)

				// ZFS operations
				r.Route("/zfs", func(r chi.Router) {
					r.Get("/pools", handlers.ListZFSPools)
					r.Get("/pools/{name}", handlers.GetZFSPool)
					r.Post("/pools", handlers.CreateZFSPool)
					r.Delete("/pools/{name}", handlers.DestroyZFSPool)
					r.Post("/pools/{name}/scrub", handlers.ScrubZFSPool)

					r.Get("/pools/{pool}/datasets", handlers.ListZFSDatasets)
					r.Post("/snapshots", handlers.CreateZFSSnapshot)
					r.Get("/datasets/{dataset}/snapshots", handlers.ListZFSSnapshots)
				})

				// RAID operations
				r.Route("/raid", func(r chi.Router) {
					r.Get("/arrays", handlers.ListRAIDArrays)
					r.Get("/arrays/{name}", handlers.GetRAIDArray)
					r.Post("/arrays", handlers.CreateRAIDArray)
				})

				// SMART operations
				r.Route("/smart", func(r chi.Router) {
					r.Get("/{device}", handlers.GetSMARTInfo)
					r.Post("/{device}/test", handlers.RunSMARTTest)
				})

				// Samba operations
				r.Route("/samba", func(r chi.Router) {
					r.Get("/status", handlers.GetSambaStatus)
					r.Post("/restart", handlers.RestartSamba)
					r.Get("/shares", handlers.ListSambaShares)
					r.Get("/shares/{name}", handlers.GetSambaShare)
					r.Post("/shares", handlers.CreateSambaShare)
					r.Put("/shares/{name}", handlers.UpdateSambaShare)
					r.Delete("/shares/{name}", handlers.DeleteSambaShare)
				})

				// NFS operations
				r.Route("/nfs", func(r chi.Router) {
					r.Post("/restart", handlers.RestartNFS)
					r.Get("/exports", handlers.ListNFSExports)
					r.Post("/exports", handlers.CreateNFSExport)
					r.Delete("/exports", handlers.DeleteNFSExport)
				})

				// Network operations
				r.Route("/network", func(r chi.Router) {
					r.Post("/bond", handlers.CreateBondInterface)
					r.Delete("/bond/{name}", handlers.DeleteBondInterface)
					r.Post("/vlan", handlers.CreateVLANInterface)
					r.Delete("/vlan/{parent}/{vlanid}", handlers.DeleteVLANInterface)
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
				r.Post("/routes", netHandler.AddRoute)
				r.Delete("/routes", netHandler.DeleteRoute)
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
				r.Get("/containers/{id}/top", dockerHandler.GetContainerTop)
				r.Post("/containers/{id}/start", dockerHandler.StartContainer)
				r.Post("/containers/{id}/stop", dockerHandler.StopContainer)
				r.Post("/containers/{id}/restart", dockerHandler.RestartContainer)
				r.Post("/containers/{id}/pause", dockerHandler.PauseContainer)
				r.Post("/containers/{id}/unpause", dockerHandler.UnpauseContainer)
				r.Post("/containers/{id}/exec", dockerHandler.ExecContainer)
				r.Put("/containers/{id}/resources", dockerHandler.UpdateContainerResources)
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

			// Backup routes
			r.Route("/backups", func(r chi.Router) {
				backupHandler := handlers.NewBackupHandler()
				r.Use(backupHandler.CheckAvailability)

				// Backup jobs
				r.Get("/jobs", backupHandler.ListJobs)
				r.Post("/jobs", backupHandler.CreateJob)
				r.Get("/jobs/{id}", backupHandler.GetJob)
				r.Put("/jobs/{id}", backupHandler.UpdateJob)
				r.Delete("/jobs/{id}", backupHandler.DeleteJob)
				r.Post("/jobs/{id}/run", backupHandler.RunJob)

				// Backup history
				r.Get("/history", backupHandler.GetHistory)

				// Snapshots
				r.Get("/snapshots", backupHandler.ListSnapshots)
				r.Post("/snapshots", backupHandler.CreateSnapshot)
				r.Delete("/snapshots/{id}", backupHandler.DeleteSnapshot)
				r.Post("/snapshots/{id}/restore", backupHandler.RestoreSnapshot)
			})

			// Active Directory routes
			r.Route("/ad", func(r chi.Router) {
				adHandler := handlers.NewADHandler()

				// AD configuration
				r.Get("/config", adHandler.GetConfig)
				r.Put("/config", adHandler.UpdateConfig)
				r.Post("/test", adHandler.TestConnection)
				r.Get("/status", adHandler.GetStatus)

				// AD users
				r.Post("/authenticate", adHandler.Authenticate)
				r.Get("/users", adHandler.ListUsers)
				r.Post("/users/sync", adHandler.SyncUser)
			})

			// Active Directory Domain Controller routes
			r.Route("/ad-dc", func(r chi.Router) {
				dcHandler := handlers.NewADDCHandler()

				// Domain Controller Management
				r.Get("/status", dcHandler.GetDCStatus)
				r.Get("/config", dcHandler.GetDCConfig)
				r.Put("/config", dcHandler.UpdateDCConfig)
				r.Post("/provision", dcHandler.ProvisionDomain)
				r.Post("/demote", dcHandler.DemoteDomain)
				r.Get("/info", dcHandler.GetDomainInfo)
				r.Get("/level", dcHandler.GetDomainLevel)
				r.Post("/level/raise", dcHandler.RaiseDomainLevel)
				r.Post("/service/restart", dcHandler.RestartService)

				// User Management
				r.Route("/users", func(r chi.Router) {
					r.Get("/", dcHandler.ListUsers)
					r.Post("/", dcHandler.CreateUser)
					r.Delete("/{username}", dcHandler.DeleteUser)
					r.Post("/{username}/enable", dcHandler.EnableUser)
					r.Post("/{username}/disable", dcHandler.DisableUser)
					r.Post("/{username}/password", dcHandler.SetUserPassword)
					r.Post("/{username}/expiry", dcHandler.SetUserExpiry)
				})

				// Group Management
				r.Route("/groups", func(r chi.Router) {
					r.Get("/", dcHandler.ListGroups)
					r.Post("/", dcHandler.CreateGroup)
					r.Delete("/{name}", dcHandler.DeleteGroup)
					r.Get("/{name}/members", dcHandler.ListGroupMembers)
					r.Post("/{name}/members", dcHandler.AddGroupMember)
					r.Delete("/{name}/members/{username}", dcHandler.RemoveGroupMember)
				})

				// Computer Management
				r.Route("/computers", func(r chi.Router) {
					r.Get("/", dcHandler.ListComputers)
					r.Post("/", dcHandler.CreateComputer)
					r.Delete("/{name}", dcHandler.DeleteComputer)
				})

				// Organizational Unit Management
				r.Route("/ou", func(r chi.Router) {
					r.Get("/", dcHandler.ListOUs)
					r.Post("/", dcHandler.CreateOU)
					r.Delete("/", dcHandler.DeleteOU) // Uses DN in body
				})

				// Group Policy Management
				r.Route("/gpo", func(r chi.Router) {
					r.Get("/", dcHandler.ListGPOs)
					r.Post("/", dcHandler.CreateGPO)
					r.Delete("/{name}", dcHandler.DeleteGPO)
					r.Post("/{name}/link", dcHandler.LinkGPO)
					r.Post("/{name}/unlink", dcHandler.UnlinkGPO)
				})

				// DNS Management
				r.Route("/dns", func(r chi.Router) {
					r.Get("/zones", dcHandler.ListDNSZones)
					r.Post("/zones", dcHandler.CreateDNSZone)
					r.Delete("/zones/{zone}", dcHandler.DeleteDNSZone)
					r.Get("/zones/{zone}/records", dcHandler.ListDNSRecords)
					r.Post("/zones/{zone}/records", dcHandler.AddDNSRecord)
					r.Delete("/zones/{zone}/records/{record}", dcHandler.DeleteDNSRecord)
				})

				// FSMO Roles
				r.Route("/fsmo", func(r chi.Router) {
					r.Get("/", dcHandler.ShowFSMORoles)
					r.Post("/transfer", dcHandler.TransferFSMORoles)
					r.Post("/seize", dcHandler.SeizeFSMORoles)
				})

				// Utility
				r.Post("/test-config", dcHandler.TestConfiguration)
				r.Get("/dbcheck", dcHandler.ShowDBCheck)
				r.Post("/backup", dcHandler.BackupOnline)
			})

			// Audit Log routes
			r.Route("/audit", func(r chi.Router) {
				auditHandler := handlers.NewAuditHandler()

				// Audit log retrieval (admin only)
				r.Use(mw.AdminOnly)
				r.Get("/logs", auditHandler.ListAuditLogs)
				r.Get("/logs/recent", auditHandler.GetRecentAuditLogs)
				r.Get("/logs/{id}", auditHandler.GetAuditLog)
				r.Get("/stats", auditHandler.GetAuditStats)
			})

			// Failed Login Tracking routes
			r.Route("/security", func(r chi.Router) {
				failedLoginHandler := handlers.NewFailedLoginHandler()

				// Security management (admin only)
				r.Use(mw.AdminOnly)
				r.Get("/failed-logins", failedLoginHandler.ListFailedAttempts)
				r.Get("/blocked-ips", failedLoginHandler.GetBlockedIPs)
				r.Post("/unblock-ip", failedLoginHandler.UnblockIP)
				r.Get("/failed-logins/stats", failedLoginHandler.GetStats)
			})

			// Alert/Notification routes
			r.Route("/alerts", func(r chi.Router) {
				alertHandler := handlers.NewAlertHandler()

				// Alert management (admin only)
				r.Use(mw.AdminOnly)
				r.Get("/config", alertHandler.GetConfig)
				r.Put("/config", alertHandler.UpdateConfig)
				r.Post("/test/email", alertHandler.TestEmail)
				r.Post("/test/webhook", alertHandler.TestWebhook)
				r.Get("/logs", alertHandler.GetAlertLogs)
			})

			// Monitoring configuration routes
			r.Route("/monitoring", func(r chi.Router) {
				// Monitoring config management (admin only)
				r.Use(mw.AdminOnly)
				r.Get("/config", handlers.GetMonitoringConfig)
				r.Put("/config", handlers.UpdateMonitoringConfig)
			})

			// Scheduler/Task routes
			r.Route("/tasks", func(r chi.Router) {
				schedulerHandler := handlers.NewSchedulerHandler()

				// Task management (admin only)
				r.Use(mw.AdminOnly)
				r.Get("/", schedulerHandler.ListTasks)
				r.Post("/", schedulerHandler.CreateTask)
				r.Get("/{id}", schedulerHandler.GetTask)
				r.Put("/{id}", schedulerHandler.UpdateTask)
				r.Delete("/{id}", schedulerHandler.DeleteTask)
				r.Post("/{id}/run", schedulerHandler.RunTaskNow)
				r.Get("/{id}/executions", schedulerHandler.GetTaskExecutions)
				r.Post("/validate-cron", schedulerHandler.ValidateCron)
			})

			// Two-Factor Authentication routes
			r.Route("/2fa", func(r chi.Router) {
				twofaHandler := handlers.NewTwoFAHandler()

				// User 2FA management (requires authentication)
				r.Get("/status", twofaHandler.GetStatus)
				r.Post("/setup", twofaHandler.SetupTwoFactor)
				r.Post("/enable", twofaHandler.EnableTwoFactor)
				r.Post("/disable", twofaHandler.DisableTwoFactor)
				r.Post("/backup-codes/regenerate", twofaHandler.RegenerateBackupCodes)
			})

			// Plugin routes
			r.Route("/plugins", func(r chi.Router) {
				pluginHandler := handlers.NewPluginHandler()
				r.Use(pluginHandler.CheckAvailability)

				// Plugin management
				r.Get("/", pluginHandler.ListPlugins)
				r.Get("/{id}", pluginHandler.GetPlugin)
				r.Post("/install", pluginHandler.InstallPlugin)
				r.Delete("/{id}", pluginHandler.UninstallPlugin)
				r.Post("/{id}/enable", pluginHandler.EnablePlugin)
				r.Post("/{id}/disable", pluginHandler.DisablePlugin)
				r.Put("/{id}/config", pluginHandler.UpdatePluginConfig)

				// Plugin runtime control
				r.Post("/{id}/start", pluginHandler.StartPlugin)
				r.Post("/{id}/stop", pluginHandler.StopPlugin)
				r.Post("/{id}/restart", pluginHandler.RestartPlugin)
				r.Get("/{id}/status", pluginHandler.GetPluginStatus)
				r.Get("/running", pluginHandler.ListRunningPlugins)
			})

			// Plugin Store routes (registry-based installation)
			r.Route("/store", func(r chi.Router) {
				// Public endpoints (browsing)
				r.Get("/plugins", handlers.ListAvailablePlugins)
				r.Get("/plugins/{id}", handlers.GetPluginFromRegistry)
				r.Get("/plugins/search", handlers.SearchPlugins)

				// Admin-only endpoints (installation)
				r.Group(func(r chi.Router) {
					r.Use(mw.AdminOnly)

					r.Post("/plugins/{id}/install", handlers.InstallPlugin)
					r.Delete("/plugins/{id}/uninstall", handlers.UninstallPlugin)
					r.Post("/plugins/{id}/update", handlers.UpdatePlugin)
					r.Post("/sync", handlers.SyncRegistry)
					r.Get("/installed", handlers.ListInstalledPlugins)
				})
			})

			// Terminal WebSocket endpoint
			r.Route("/terminal", func(r chi.Router) {
				r.Use(mw.AdminOnly) // Terminal access requires admin privileges
				r.Get("/ws", handlers.TerminalWebSocketHandler)
			})
		})
	})

	// WebSocket endpoint
	r.Get("/ws", handlers.WebSocketHandler)

	// Serve embedded frontend static files (must be last to act as catch-all)
	// This handles all routes not matched above and serves the React SPA
	spaHandler, err := embedfs.NewSPAHandler()
	if err != nil {
		logger.Warn("Failed to initialize SPA handler, frontend will not be served",
			zap.Error(err))
	} else {
		// Catch-all route for SPA (must be last)
		r.Get("/*", spaHandler.ServeHTTP)
		logger.Info("Embedded frontend static file server initialized")
	}

	return r
}
