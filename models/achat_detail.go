package models

import "time"

type AchatDetail struct {
    ID          int       `json:"id"`
    ClientID    int       `json:"client_id"`
    ClientNom   string    `json:"client_nom"`
    LivreID     int       `json:"livre_id"`
    LivreTitre  string    `json:"livre_titre"`
    LivreImage  string    `json:"livre_image"`
    LivreAuteur string    `json:"livre_auteur"`
    Quantite    int       `json:"quantite"`
    Montant     int       `json:"montant"`
    DateAchat   time.Time `json:"date_achat"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}