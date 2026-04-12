package repository

import (
    "context"
    "ebook-backend/config"
    "ebook-backend/models"
)

// GetBookByID récupère un livre par son ID (inclut file_path)
func GetBookByID(id int) (*models.Book, error) {
    var b models.Book
    query := `SELECT id, titre, auteur, description, prix_fcfa, image, type, pages, stock, file_path, created_at, updated_at
              FROM livres WHERE id = $1`
    err := config.DB.QueryRow(context.Background(), query, id).Scan(
        &b.ID, &b.Titre, &b.Auteur, &b.Description, &b.PrixFCFA, &b.Image, &b.Type, &b.Pages, &b.Stock, &b.FilePath, &b.CreatedAt, &b.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &b, nil
}

// UpdateBookStock met à jour le stock d'un livre
func UpdateBookStock(bookID int, newStock int) error {
    _, err := config.DB.Exec(context.Background(), "UPDATE livres SET stock = $1, updated_at = NOW() WHERE id = $2", newStock, bookID)
    return err
}