// Package utils provides common utility functions used across the application,
// including hashing and JWT-related operations.
package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JWTUtil struct {
	secretKey       []byte
	expirationHours int
}

func NewJWTUtil(secret string, expirationHours int) *JWTUtil {
	return &JWTUtil{
		secretKey:       []byte(secret),
		expirationHours: expirationHours,
	}
}

func (j *JWTUtil) GenerateToken(userID int, email string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(j.expirationHours) * time.Hour)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWTUtil) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
