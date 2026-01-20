package router

import (
	"sms-dashboard/internal/config"
	"sms-dashboard/internal/handler"
	"sms-dashboard/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// CORS for development (if frontend is separate)
	// In production (embed), same origin, so strictly not needed but good for safety if we dev separately
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Token")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	authHandler := handler.NewAuthHandler(cfg)
	smsHandler := handler.NewSMSHandler()

	// Initialize default user
	authHandler.InitDefaultUser()

	api := r.Group("/api")
	{
		api.POST("/login", authHandler.Login)

		// Public (Protected by API Token)
		// Assuming Android forwarder posts to /api/sms or /api/receive
		// Let's use /api/sms for receiving with Token check
		api.POST("/sms", middleware.APITokenMiddleware(cfg), smsHandler.Receive)

		// Private (Protected by JWT)
		authorized := api.Group("/")
		authorized.Use(middleware.AuthMiddleware(cfg))
		{
			authorized.POST("/change-password", authHandler.ChangePassword)
			authorized.GET("/sms/list", smsHandler.List)
			authorized.DELETE("/sms/:id", smsHandler.Delete)
		}
	}

	return r
}
