package api

import (
	"github.com/gin-gonic/gin"

	"github.com/stumpfworks/stumpfworks-nas/plugins/vpn-server/config"
	"github.com/stumpfworks/stumpfworks-nas/plugins/vpn-server/internal/core"
)

// SetupRouter creates and configures the Gin router
func SetupRouter(cfg *config.Config, vpnManager *core.VPNManager) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware(cfg))

	// Create handler
	handler := NewHandler(vpnManager)

	// API routes
	api := router.Group("/api/vpn")
	{
		// General status and control
		api.GET("/status", handler.GetStatus)
		api.POST("/start", handler.StartServer)
		api.POST("/stop", handler.StopServer)

		// Protocol control
		api.POST("/:protocol/start", handler.StartProtocol)
		api.POST("/:protocol/stop", handler.StopProtocol)

		// User management
		users := api.Group("/users")
		{
			users.GET("", handler.GetUsers)
			users.GET("/search", handler.SearchUsers)
			users.POST("", handler.CreateUser)
			users.GET("/:id", handler.GetUser)
			users.DELETE("/:id", handler.DeleteUser)
			users.PUT("/:id/protocols", handler.UpdateUserProtocols)
			users.POST("/:id/enable", handler.EnableUser)
			users.POST("/:id/disable", handler.DisableUser)
		}

		// WireGuard
		wireguard := api.Group("/wireguard")
		{
			wireguard.GET("/stats", handler.GetWireGuardStats)

			peers := wireguard.Group("/peers")
			{
				peers.GET("", handler.GetWireGuardPeers)
				peers.POST("", handler.CreateWireGuardPeer)
				peers.GET("/:id", handler.GetWireGuardPeer)
				peers.DELETE("/:id", handler.DeleteWireGuardPeer)
				peers.GET("/:id/config", handler.GetWireGuardPeerConfig)
				peers.GET("/:id/qrcode", handler.GetWireGuardPeerQRCode)
			}
		}

		// TODO: Add OpenVPN routes in Phase 2
		// openvpn := api.Group("/openvpn")
		// {
		//     openvpn.GET("/status", handler.GetOpenVPNStatus)
		//     openvpn.GET("/certificates", handler.GetCertificates)
		//     openvpn.POST("/certificates", handler.CreateCertificate)
		//     openvpn.DELETE("/certificates/:id", handler.RevokeCertificate)
		// }

		// TODO: Add PPTP routes in Phase 3
		// pptp := api.Group("/pptp")
		// {
		//     pptp.GET("/status", handler.GetPPTPStatus)
		//     pptp.GET("/connections", handler.GetPPTPConnections)
		// }

		// TODO: Add L2TP routes in Phase 3
		// l2tp := api.Group("/l2tp")
		// {
		//     l2tp.GET("/status", handler.GetL2TPStatus)
		//     l2tp.GET("/connections", handler.GetL2TPConnections)
		//     l2tp.PUT("/psk", handler.UpdateL2TPPSK)
		// }
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "vpn-server",
		})
	})

	return router
}
