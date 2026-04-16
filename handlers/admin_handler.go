package handlers

import (
	"ebook-backend/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RevenueResponse représente la réponse du chiffre d'affaires
type RevenueResponse struct {
	Revenue int `json:"revenue"`
}

// MonthlySalesResponse représente une vente mensuelle
type MonthlySalesResponse struct {
	Month  string `json:"month"`
	Ventes int    `json:"ventes"`
}

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

// GetRevenue godoc
// @Summary      Chiffre d'affaires total
// @Description  Retourne le montant total des ventes (somme des montants des achats)
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  RevenueResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /admin/revenue [get]
func GetRevenue(c *gin.Context) {
	total, err := repository.GetTotalRevenue(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, RevenueResponse{Revenue: total})
}

// GetMonthlySales godoc
// @Summary      Ventes mensuelles
// @Description  Retourne le nombre de ventes (quantité totale) par mois
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  MonthlySalesResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /admin/sales/monthly [get]
func GetMonthlySales(c *gin.Context) {
	sales, err := repository.GetMonthlySales(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sales)
}