package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
}

// ===== ACCESS TOKEN (15 PHÚT) =====
func GenerateAccessToken(userID int64, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ===== REFRESH TOKEN (1 GIỜ) =====
func GenerateRefreshToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ===== PARSE TOKEN =====
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	userIDFloat, ok := mapClaims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user_id in token")
	}

	role, ok := mapClaims["role"].(string)
	if !ok {
		return nil, errors.New("invalid role in token")
	}

	return &Claims{
		UserID: int64(userIDFloat),
		Role:   role,
	}, nil
}
