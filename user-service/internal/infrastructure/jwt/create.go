package jwt

import (
	"strconv"
	"time"

	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/config"
	"github.com/golang-jwt/jwt/v5"
)

func Create(userID uint, isAdmin bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(int64(userID), 10),
		Audience:  jwt.ClaimStrings{strconv.FormatBool(isAdmin)},
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Duration(config.Data.JWTExpirationSeconds))},
	})

	tokenString, err := token.SignedString([]byte(config.Data.JWTSecret))
	return tokenString, err
}
