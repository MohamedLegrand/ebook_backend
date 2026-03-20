package repository

import (
	"context"
	"ebook-backend/config"
	"ebook-backend/models"

	"golang.org/x/crypto/bcrypt"
)

// Créer un nouveau client
func CreateClient(req models.RegisterRequest) (*models.Client, error) {
	// Hasher le mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	client := &models.Client{}

	query := `
		INSERT INTO clients (full_name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, full_name, email, created_at, updated_at
	`

	err = config.DB.QueryRow(
		context.Background(),
		query,
		req.FullName,
		req.Email,
		string(hashedPassword),
	).Scan(
		&client.ID,
		&client.FullName,
		&client.Email,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return client, nil
}

// Trouver un client par email
func GetClientByEmail(email string) (*models.Client, error) {
	client := &models.Client{}

	query := `
		SELECT id, full_name, email, password, created_at, updated_at
		FROM clients
		WHERE email = $1
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		email,
	).Scan(
		&client.ID,
		&client.FullName,
		&client.Email,
		&client.Password,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return client, nil
}

// Vérifier le mot de passe
func CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}