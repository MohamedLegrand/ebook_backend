package middleware

import (
	"ebook-backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthRequired est le middleware qui protège les routes privées
// Il vérifie que la requête contient un token JWT valide
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. On cherche le token dans l'en-tête Authorization
		//    Le frontend React doit envoyer : Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token manquant, veuillez vous connecter",
			})
			c.Abort() // Stoppe la requête, le handler ne sera pas appelé
			return
		}

		// 2. L'en-tête doit commencer par "Bearer "
		//    On sépare le mot "Bearer" du token lui-même
		//    Exemple : "Bearer eyJhbGciOiJIUzI1..." → on garde juste "eyJhbGciOiJIUzI1..."
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Format du token invalide, utilisez : Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. On valide le token avec notre fonction utils
		//    Si le token est expiré ou falsifié, ValidateToken retourne une erreur
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token invalide ou expiré, veuillez vous reconnecter",
			})
			c.Abort()
			return
		}

		// 4. Token valide — on stocke les infos du client dans le contexte Gin
		//    Les handlers qui suivent pourront récupérer ces infos avec c.Get()
		c.Set("client_id", claims.ClientID)
		c.Set("client_email", claims.Email)

		// 5. Tout est bon, on laisse passer vers le handler
		c.Next()
	}
}