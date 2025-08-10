package utils

import (
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
