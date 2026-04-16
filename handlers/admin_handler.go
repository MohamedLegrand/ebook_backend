package handlers

import (
	"ebook-backend/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllClients retourne la liste de tous les clients (route admin)
// @Summary      Récupérer tous les clients
// @Description  Retourne la liste complète des clients (accès réservé aux administrateurs)
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Client
// @Failure      500  {object}  map[string]interface{}  "erreur interne"
// @Router       /admin/clients [get]
// @Security     BearerAuth
func GetAllClients(c *gin.Context) {
	clients, err := repository.GetAllClients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les clients",
		})
		return
	}
	c.JSON(http.StatusOK, clients)
}