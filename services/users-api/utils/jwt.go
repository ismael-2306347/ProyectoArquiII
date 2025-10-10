package utils

import (
	"users-api/domain"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint            `json:"user_id"`
	Username string          `json:"username"`
	Role     domain.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(user domain.User, secret string) (string, error) {
	// Implementar
	return "", nil
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
	// Implementar
	return nil, nil
}
