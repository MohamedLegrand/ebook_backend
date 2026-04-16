package repository

import (
    "context"
    "ebook-backend/config"
    "ebook-backend/models"
)

// MonthlySales représente une vente mensuelle (pour les graphiques)
type MonthlySales struct {
    Month  string `json:"month"`
    Ventes int    `json:"ventes"`
}

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

// GetTotalRevenue retourne la somme de tous les montants des achats
func GetTotalRevenue(ctx context.Context) (int, error) {
    var total int
    query := `SELECT COALESCE(SUM(montant), 0) FROM achat`
    err := config.DB.QueryRow(ctx, query).Scan(&total)
    return total, err
}

// GetMonthlySales retourne les ventes mensuelles (somme des quantités) pour les 12 derniers mois
func GetMonthlySales(ctx context.Context) ([]MonthlySales, error) {
    query := `
        SELECT TO_CHAR(date_achat, 'Mon') as month, 
               EXTRACT(MONTH FROM date_achat) as month_num,
               COALESCE(SUM(quantite), 0) as ventes
        FROM achat
        WHERE date_achat >= NOW() - INTERVAL '11 months'
        GROUP BY month, month_num
        ORDER BY month_num
    `
    rows, err := config.DB.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var sales []MonthlySales
    for rows.Next() {
        var ms MonthlySales
        var monthNum int
        if err := rows.Scan(&ms.Month, &monthNum, &ms.Ventes); err != nil {
            return nil, err
        }
        sales = append(sales, ms)
    }
    return sales, nil
}