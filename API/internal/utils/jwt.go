package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tiago-prog/novels-api/internal/model"
)

func GenerateJWT(userID string, role model.Role, secret []byte) (string, error) {
	if len(secret) < 32 {
		return "", errors.New("secret must be at least 32 bytes")
	}

	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(3 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
		"nbf":  time.Now().Unix(),
		"iss":  "plataforma-de-novels",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
