package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Yarik7610/library-backend/user-service/config"
	"github.com/golang-jwt/jwt/v5"
)

const EXPIRATION_TIME_DAYS = 7

func CreateJWTToken(userID uint, isAdmin bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(int64(userID), 10),
		Audience:  jwt.ClaimStrings{strconv.FormatBool(isAdmin)},
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 24 * EXPIRATION_TIME_DAYS)},
	})

	tokenString, err := token.SignedString([]byte(config.Data.JWTSecret))
	return tokenString, err
}

func VerifyJWTToken(tokenString string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Data.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
