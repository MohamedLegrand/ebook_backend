package handlers

import (
    "ebook-backend/models"
    "ebook-backend/repository"
    "fmt"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

type InitierPaiementInput struct {
    MoyenPaiement  string `json:"moyen_paiement" binding:"required"`
    NumeroPaiement string `json:"numero_paiement" binding:"required"`
    MontantTotal   int    `json:"montant_total" binding:"required"`
    Items          []struct {
        LivreID  int `json:"livre_id"`
        Quantite int `json:"quantite"`
    } `json:"items" binding:"required"`
}

func InitierPaiement(c *gin.Context) {
    clientIDVal, exists := c.Get("client_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "non authentifié"})
        return
    }
    clientID := clientIDVal.(int)

    var input InitierPaiementInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Générer une référence unique
    reference := fmt.Sprintf("PAY-%d-%d", time.Now().UnixNano(), clientID)

    // Créer l'enregistrement paiement (statut completed pour simulation)
    paiement := &models.Paiement{
        ClientID:       clientID,
        Reference:      reference,
        MontantTotal:   input.MontantTotal,
        MoyenPaiement:  input.MoyenPaiement,
        NumeroPaiement: input.NumeroPaiement,
        Statut:         "completed",
    }
    if err := repository.CreatePaiement(c.Request.Context(), paiement); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur enregistrement paiement"})
        return
    }

    // Pour chaque article, créer un achat et mettre à jour le stock
    for _, item := range input.Items {
        livre, err := repository.GetBookByID(item.LivreID)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Livre %d introuvable", item.LivreID)})
            return
        }
        if livre.Stock < item.Quantite {
            c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Stock insuffisant pour le livre %s", livre.Titre)})
            return
        }

        montant := livre.PrixFCFA * item.Quantite
        achat := &models.Achat{
            ClientID: clientID,
            LivreID:  item.LivreID,
            Quantite: item.Quantite,
            Montant:  montant,
        }
        if err := repository.CreateAchat(c.Request.Context(), achat); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur enregistrement achat"})
            return
        }
        // Mise à jour du stock
        _ = repository.UpdateBookStock(item.LivreID, livre.Stock-item.Quantite)
    }

    c.JSON(http.StatusOK, gin.H{
        "success":   true,
        "reference": reference,
        "message":   "Paiement validé et achats enregistrés",
    })
}