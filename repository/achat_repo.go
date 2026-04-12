package repository

import (
    "context"
    "ebook-backend/config"
    "ebook-backend/models"
)

func CreateAchat(ctx context.Context, achat *models.Achat) error {
    query := `INSERT INTO achat (client_id, livre_id, quantite, montant, date_achat, created_at, updated_at)
              VALUES ($1, $2, $3, $4, NOW(), NOW(), NOW())
              RETURNING id, created_at, updated_at`
    return config.DB.QueryRow(ctx, query,
        achat.ClientID, achat.LivreID, achat.Quantite, achat.Montant,
    ).Scan(&achat.ID, &achat.CreatedAt, &achat.UpdatedAt)
}

func GetAchatsByClient(ctx context.Context, clientID int) ([]models.Achat, error) {
    rows, err := config.DB.Query(ctx, `
        SELECT id, client_id, livre_id, quantite, montant, date_achat, created_at, updated_at
        FROM achat WHERE client_id = $1 ORDER BY date_achat DESC`, clientID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var achats []models.Achat
    for rows.Next() {
        var a models.Achat
        if err := rows.Scan(&a.ID, &a.ClientID, &a.LivreID, &a.Quantite, &a.Montant, &a.DateAchat, &a.CreatedAt, &a.UpdatedAt); err != nil {
            return nil, err
        }
        achats = append(achats, a)
    }
    return achats, nil
}

func GetAllAchats(ctx context.Context) ([]models.Achat, error) {
    rows, err := config.DB.Query(ctx, `
        SELECT id, client_id, livre_id, quantite, montant, date_achat, created_at, updated_at
        FROM achat ORDER BY date_achat DESC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var achats []models.Achat
    for rows.Next() {
        var a models.Achat
        if err := rows.Scan(&a.ID, &a.ClientID, &a.LivreID, &a.Quantite, &a.Montant, &a.DateAchat, &a.CreatedAt, &a.UpdatedAt); err != nil {
            return nil, err
        }
        achats = append(achats, a)
    }
    return achats, nil
}