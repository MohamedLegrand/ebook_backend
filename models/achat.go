package models

import "time"

type Achat struct {
    ID        int       `json:"id"`
    ClientID  int       `json:"client_id"`
    LivreID   int       `json:"livre_id"`
    Quantite  int       `json:"quantite"`
    Montant   int       `json:"montant"`
    DateAchat time.Time `json:"date_achat"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}