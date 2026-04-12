package models

import "time"

type Paiement struct {
    ID            int       `json:"id"`
    ClientID      int       `json:"client_id"`
    Reference     string    `json:"reference"`
    MontantTotal  int       `json:"montant_total"`
    MoyenPaiement string    `json:"moyen_paiement"`
    NumeroPaiement string   `json:"numero_paiement"`
    Statut        string    `json:"statut"`
    DateCreation  time.Time `json:"date_creation"`
    DateMiseAJour time.Time `json:"date_mise_a_jour"`
}