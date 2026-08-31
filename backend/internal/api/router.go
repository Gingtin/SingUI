package api

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/singbox-ui/singbox-ui/internal/api/handlers"
	"github.com/singbox-ui/singbox-ui/internal/api/middleware"
)

func SetupRouter(webDistFS fs.FS) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.Cors())

	// Public Subscription Endpoints
	r.GET("/sub/:token", handlers.HandleSubscription)
	r.GET("/sub/view/:token", handlers.HandleSubView)

	// API Group
	api := r.Group("/api")
	{
		// Auth
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/logout", handlers.Logout)
			auth.GET("/info", middleware.JWTAuth(), handlers.GetUserInfo)
		}

		// Protected Routes
		protected := api.Group("")
		protected.Use(middleware.JWTAuth())
		{
			// Server & Core
			server := protected.Group("/server")
			{
				server.GET("/status", handlers.GetSystemStatus)
				server.GET("/core-status", handlers.GetCoreStatus)
				server.POST("/core/start", handlers.StartCore)
				server.POST("/core/stop", handlers.StopCore)
				server.POST("/core/restart", handlers.RestartCore)
				server.GET("/logs", handlers.GetLogs)
				server.GET("/logs/ws", handlers.StreamLogsWS)
				server.GET("/connections", handlers.GetActiveConnections)
			}

			// Inbounds & Clients
			inbounds := protected.Group("/inbounds")
			{
				inbounds.GET("", handlers.ListInbounds)
				inbounds.POST("", handlers.CreateInbound)
				inbounds.GET("/:id", handlers.GetInbound)
				inbounds.PUT("/:id", handlers.UpdateInbound)
				inbounds.DELETE("/:id", handlers.DeleteInbound)

				inbounds.POST("/:id/clients", handlers.AddClient)
				inbounds.PUT("/:id/clients/:clientId", handlers.UpdateClient)
				inbounds.DELETE("/:id/clients/:clientId", handlers.DeleteClient)
				inbounds.POST("/:id/clients/:clientId/reset", handlers.ResetClientTraffic)

				inbounds.POST("/reset-all", handlers.ResetAllTraffic)
				inbounds.GET("/reality-keypair", handlers.GenerateRealityKeypair)
				inbounds.GET("/random-uuid", handlers.GenerateRandomUUID)
			}

			// Settings
			settings := protected.Group("/settings")
			{
				settings.GET("", handlers.GetSettings)
				settings.POST("", handlers.UpdateSettings)
				settings.POST("/password", handlers.UpdatePassword)
				settings.GET("/backup", handlers.DownloadBackup)
				settings.POST("/restore", handlers.RestoreBackup)
			}
		}
	}

	// Static frontend assets serving
	if webDistFS != nil {
		fileServer := http.FileServer(http.FS(webDistFS))
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/sub") {
				c.JSON(http.StatusNotFound, gin.H{"error": "API route not found"})
				return
			}

			// Check if file exists in embed FS
			if _, err := fs.Stat(webDistFS, strings.TrimPrefix(path, "/")); err == nil && path != "/" {
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}

			// SPA Fallback: serve index.html
			indexFile, err := webDistFS.Open("index.html")
			if err != nil {
				c.String(http.StatusOK, "Sing-box UI Server is running. (Frontend dist not found)")
				return
			}
			defer indexFile.Close()

			stat, _ := indexFile.Stat()
			http.ServeContent(c.Writer, c.Request, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
		})
	}

	return r
}
