package main

import (
	"log"
	"os"

	"ebook-backend/config"
	"ebook-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	// Importation de la documentation générée par swag init
	_ "ebook-backend/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           API E-Book
// @version         1.0
// @description     API de gestion de la librairie en ligne
// @termsOfService  http://swagger.io/terms/

// @contact.name   Support API
// @contact.email  support@ebook.com

// @license.name   Apache 2.0
// @license.url    http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Charger le fichier .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erreur chargement .env")
	}

	// Connexion à la base de données
	config.ConnectDB()
	defer config.DB.Close()

	// Initialiser Gin
	r := gin.Default()

	// Activer CORS pour que React puisse communiquer avec l'API
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Route pour servir la documentation Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Brancher toutes les routes
	routes.SetupRoutes(r)

	// Démarrer le serveur
	port := os.Getenv("PORT")
	log.Printf("Serveur démarré sur le port %s", port)
	r.Run(":" + port)
}