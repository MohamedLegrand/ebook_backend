package repository

import (
	"context"
	"ebook-backend/config"
	"ebook-backend/models"

	"golang.org/x/crypto/bcrypt"
)

// CreateClient crée un nouveau client
func CreateClient(req models.RegisterRequest) (*models.Client, error) {
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

// GetClientByEmail trouve un client par son email
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

// CheckPassword vérifie le mot de passe
func CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

// GetAllClients retourne tous les clients (admin)
func GetAllClients() ([]models.Client, error) {
	query := `SELECT id, full_name, email, created_at, updated_at FROM clients ORDER BY id`
	rows, err := config.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var c models.Client
		err := rows.Scan(&c.ID, &c.FullName, &c.Email, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, nil
}