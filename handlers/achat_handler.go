package handlers

import (
	"ebook-backend/models"
	"ebook-backend/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AchatInput struct {
	LivreID  int `json:"livre_id" binding:"required"`
	Quantite int `json:"quantite" binding:"required,min=1"`
}

// CreateAchat godoc
// @Summary      Acheter un livre
// @Description  Permet à un client authentifié d'acheter un livre (vérifie le stock, crée un achat, met à jour le stock)
// @Tags         Achats
// @Accept       json
// @Produce      json
// @Param        achat body AchatInput true "Informations d'achat"
// @Success      201  {object}  models.Achat
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /achat [post]
func CreateAchat(c *gin.Context) {
	clientIDVal, exists := c.Get("client_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "non authentifié"})
		return
	}
	clientID := clientIDVal.(int)

	var input AchatInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Récupérer le livre et son prix
	livre, err := repository.GetBookByID(input.LivreID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Livre non trouvé"})
		return
	}
	if livre.Stock < input.Quantite {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock insuffisant"})
		return
	}

	montant := livre.PrixFCFA * input.Quantite
	achat := &models.Achat{
		ClientID:  clientID,
		LivreID:   input.LivreID,
		Quantite:  input.Quantite,
		Montant:   montant,
	}

	if err := repository.CreateAchat(c.Request.Context(), achat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'achat"})
		return
	}
	// Mise à jour du stock
	_ = repository.UpdateBookStock(input.LivreID, livre.Stock-input.Quantite)

	c.JSON(http.StatusCreated, achat)
}

// GetMyAchats godoc
// @Summary      Liste des achats du client connecté
// @Description  Retourne tous les achats effectués par le client authentifié
// @Tags         Achats
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Achat
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /client/achats [get]
func GetMyAchats(c *gin.Context) {
	clientIDVal, _ := c.Get("client_id")
	clientID := clientIDVal.(int)
	achats, err := repository.GetAchatsByClient(c.Request.Context(), clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, achats)
}

// GetAllAchatsAdmin godoc
// @Summary      (Admin) Liste tous les achats
// @Description  Récupère l'ensemble des achats de tous les clients (réservé aux administrateurs)
// @Tags         Achats
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Achat
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/achats [get]
func GetAllAchatsAdmin(c *gin.Context) {
	achats, err := repository.GetAllAchats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, achats)
}