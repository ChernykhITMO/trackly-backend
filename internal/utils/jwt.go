package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"trackly-backend/internal/middleware"
)

func GenerateJwt(secret string, userId int) (*string, error) {
	claims := &middleware.JWTClaims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "your-app-name",
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	// Return the token
	return &tokenString, nil
}
