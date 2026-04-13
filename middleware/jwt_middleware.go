package middleware

import (
	"ebook-backend/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthRequired protège les routes avec token JWT
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Println("DEBUG - Authorization header:", authHeader) // ← log pour voir l'en-tête reçu
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token manquant"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format du token invalide"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide ou expiré"})
			c.Abort()
			return
		}

		c.Set("client_id", claims.ClientID)
		c.Set("client_email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// AdminRequired vérifie que l'utilisateur a le rôle "admin"
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Accès réservé aux administrateurs"})
			c.Abort()
			return
		}
		c.Next()
	}
}