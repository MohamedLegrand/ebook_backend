package routes

import (
	"ebook-backend/handlers"
	"ebook-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// Routes publiques — pas besoin de token
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Routes protégées — le middleware vérifie le token avant chaque requête
	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired())
	{
		// Cette route ne sera accessible qu'avec un token valide
		protected.GET("/client/profile", func(c *gin.Context) {

			// On récupère les infos du client stockées par le middleware
			clientID, _ := c.Get("client_id")
			clientEmail, _ := c.Get("client_email")

			c.JSON(200, gin.H{
				"client_id": clientID,
				"email":     clientEmail,
			})
		})
	}
}