package routes

import (
	"ebook-backend/handlers"
	"ebook-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// Routes publiques
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	r.POST("/api/admin/login", handlers.AdminLogin)

	// Routes protégées (client)
	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/client/profile", func(c *gin.Context) {
			clientID, _ := c.Get("client_id")
			clientEmail, _ := c.Get("client_email")
			c.JSON(200, gin.H{
				"client_id": clientID,
				"email":     clientEmail,
			})
		})

		// Achats
		protected.POST("/achat", handlers.CreateAchat)
		protected.GET("/client/achats", handlers.GetMyAchats)
		protected.GET("/achat/:livre_id/download", handlers.DownloadBook) // nouvelle route
		protected.POST("/paiement/initier", handlers.InitierPaiement)
		
	}

	// Routes admin
	adminGroup := r.Group("/api/admin")
	adminGroup.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		adminGroup.GET("/dashboard", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Bienvenue administrateur"})
		})
		adminGroup.GET("/books", handlers.GetAllBooks)
		adminGroup.GET("/achats", handlers.GetAllAchatsAdmin)
	}
}