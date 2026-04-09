package repository

import (
	"context"
	"ebook-backend/config"
	"ebook-backend/models"
	"errors"

	"github.com/jackc/pgx/v5"
)

// GetAdminByEmail récupère un administrateur par son email
func GetAdminByEmail(email string) (*models.Administrateur, error) {
	query := `SELECT id, full_name, email, password, created_at, updated_at FROM administrateurs WHERE email = $1`
	var admin models.Administrateur
	err := config.DB.QueryRow(context.Background(), query, email).Scan(
		&admin.ID,
		&admin.FullName,
		&admin.Email,
		&admin.Password,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &admin, nil
}