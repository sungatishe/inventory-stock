package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"strconv"
	"time"
)

// Secret key для подписания JWT
var secretKey = []byte(os.Getenv("JWT_KEY"))

// Claims структура для хранения информации в токене (включая роль)
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// GenerateJWT генерирует JWT с ролью пользователя
func GenerateJWT(userID uint, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "auth-service",
			Subject:   strconv.Itoa(int(userID)), // Пользователь ID как строка
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ParseJwt парсит JWT и возвращает userID и роль
func ParseJwt(tokenString string) (uint, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return 0, "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, "", errors.New("invalid claims")
	}

	return claims.UserID, claims.Role, nil
}
