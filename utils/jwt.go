package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims représente les données stockées dans le token
type Claims struct {
	ClientID int    `json:"client_id"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken crée un nouveau token JWT pour un client
// On lui passe l'id et l'email du client
// Le token expire après 24 heures
func GenerateToken(clientID int, email string) (string, error) {
	claims := Claims{
		ClientID: clientID,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// On crée le token avec l'algorithme HS256 et notre clé secrète du .env
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// ValidateToken vérifie qu'un token est valide et non expiré
// Retourne les données contenues dans le token si tout est bon
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			// Vérifie que l'algorithme utilisé est bien HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("algorithme de signature invalide")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token invalide")
	}

	return claims, nil
}