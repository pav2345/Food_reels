package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func SignToken(secret string, id uuid.UUID) (string, error) {
	claims := JWTClaims{
		ID: id.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(secret, tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	id, err := uuid.Parse(claims.ID)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
