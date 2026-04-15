package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Impossible de se connecter à la base de données : %v", err)
	}

	// Vérifier que la connexion fonctionne
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Base de données inaccessible : %v", err)
	}

	DB = pool
	log.Println("Connexion à PostgreSQL réussie !")
}