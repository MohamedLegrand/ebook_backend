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

	// Route publique pour l'admin
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

		// Routes pour les achats (client)
		protected.POST("/achat", handlers.CreateAchat)          // créer un achat
		protected.GET("/client/achats", handlers.GetMyAchats)   // voir ses propres achats
	}

	// Routes admin protégées (nécessitent token + rôle admin)
	adminGroup := r.Group("/api/admin")
	adminGroup.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		adminGroup.GET("/dashboard", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Bienvenue administrateur"})
		})
		// Gestion des livres
		adminGroup.GET("/books", handlers.GetAllBooks)
		// Gestion des achats (admin)
		adminGroup.GET("/achats", handlers.GetAllAchatsAdmin)
	}
}