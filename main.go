package main

import (
	"ebook-backend/config"
	"ebook-backend/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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
	// Sans ça, le navigateur bloquera les requêtes venant de localhost:3000
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Les requêtes OPTIONS sont envoyées par le navigateur avant chaque requête
		// C'est ce qu'on appelle le "preflight" — on doit y répondre 200
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Brancher toutes les routes
	routes.SetupRoutes(r)

	// Démarrer le serveur
	port := os.Getenv("PORT")
	log.Printf("Serveur démarré sur le port %s", port)
	r.Run(":" + port)
}