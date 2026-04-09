package models

import "time"

type Book struct {
    ID          int       `json:"id"`
    Titre       string    `json:"titre"`
    Auteur      string    `json:"auteur"`
    Description string    `json:"description"`
    PrixFCFA    int       `json:"prix_fcfa"`
    Image       string    `json:"image"`
    Type        string    `json:"type"`
    Pages       int       `json:"pages"`
    Stock       int       `json:"stock"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}