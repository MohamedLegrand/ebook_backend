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

// GetAchatsByClient retourne les achats basiques (sans détails du livre)
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

// GetAchatsByClientWithDetails retourne les achats avec les infos du livre et du client
func GetAchatsByClientWithDetails(ctx context.Context, clientID int) ([]models.AchatDetail, error) {
    rows, err := config.DB.Query(ctx, `
        SELECT a.id, a.client_id, c.full_name, a.livre_id, l.titre, l.image, l.auteur,
               a.quantite, a.montant, a.date_achat, a.created_at, a.updated_at
        FROM achat a
        JOIN clients c ON a.client_id = c.id
        JOIN livres l ON a.livre_id = l.id
        WHERE a.client_id = $1
        ORDER BY a.date_achat DESC
    `, clientID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var details []models.AchatDetail
    for rows.Next() {
        var d models.AchatDetail
        if err := rows.Scan(
            &d.ID, &d.ClientID, &d.ClientNom,
            &d.LivreID, &d.LivreTitre, &d.LivreImage, &d.LivreAuteur,
            &d.Quantite, &d.Montant, &d.DateAchat, &d.CreatedAt, &d.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        details = append(details, d)
    }
    return details, nil
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