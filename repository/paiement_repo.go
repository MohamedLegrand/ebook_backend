package repository

import (
    "context"
    "ebook-backend/config"
    "ebook-backend/models"
)

func CreatePaiement(ctx context.Context, p *models.Paiement) error {
    query := `INSERT INTO paiements (client_id, reference, montant_total, moyen_paiement, numero_paiement, statut, date_creation, date_mise_a_jour)
              VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
              RETURNING id, date_creation, date_mise_a_jour`
    return config.DB.QueryRow(ctx, query,
        p.ClientID, p.Reference, p.MontantTotal, p.MoyenPaiement, p.NumeroPaiement, p.Statut,
    ).Scan(&p.ID, &p.DateCreation, &p.DateMiseAJour)
}

func UpdatePaiementStatus(ctx context.Context, reference, statut string) error {
    _, err := config.DB.Exec(ctx, "UPDATE paiements SET statut = $1, date_mise_a_jour = NOW() WHERE reference = $2", statut, reference)
    return err
}