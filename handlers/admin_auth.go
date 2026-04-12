package handlers

import (
	"ebook-backend/models"
	"ebook-backend/repository"
	"ebook-backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminLogin godoc
// @Summary      Connexion administrateur
// @Description  Authentifie un administrateur et retourne un token JWT avec le rôle "admin"
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials body models.LoginRequest true "Identifiants admin (email, mot de passe)"
// @Success      200  {object}  map[string]interface{}  "token et infos admin"
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /admin/login [post]
func AdminLogin(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	admin, err := repository.GetAdminByEmail(req.Email)
	if err != nil || admin == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	// Vérifie le mot de passe (réutilise la fonction du repository client)
	if !repository.CheckPassword(admin.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	// Génère un token avec le rôle "admin"
	token, err := utils.GenerateTokenWithRole(admin.ID, admin.Email, "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"admin": gin.H{
			"id":        admin.ID,
			"email":     admin.Email,
			"full_name": admin.FullName,
		},
	})
}