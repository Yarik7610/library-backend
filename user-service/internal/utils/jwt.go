package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const EXPIRATION_TIME_DAYS = 7

func CreateJWTToken(ID uint, isAdmin bool, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     ID,
		"isAdmin": isAdmin,
		"exp":     time.Now().Add(time.Hour * 24 * EXPIRATION_TIME_DAYS).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}
