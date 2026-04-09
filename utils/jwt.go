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
	Role     string `json:"role"` // "client" ou "admin"
	jwt.RegisteredClaims
}

// GenerateToken crée un token pour un client (rôle "client")
func GenerateToken(clientID int, email string) (string, error) {
	return GenerateTokenWithRole(clientID, email, "client")
}

// GenerateTokenWithRole crée un token avec un rôle personnalisé
func GenerateTokenWithRole(userID int, email, role string) (string, error) {
	claims := Claims{
		ClientID: userID,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// ValidateToken vérifie un token et retourne les claims
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
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