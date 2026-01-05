package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

// GenerateAccessToken tạo ra JWT token chứa thông tin user
func GenerateAccessToken(userID uuid.UUID, role string, secret string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID.String(),
		"role": role,
		"exp":  time.Now().Add(duration).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken tạo ra chuỗi ngẫu nhiên lưu vào DB
func GenerateRefreshToken() string {
	return uuid.New().String()
}
