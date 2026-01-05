package middleware

import (
	"net/http"
	"strings"

	"go-ecommerce/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware kiểm tra JWT token trong header Authorization
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Lấy token từ Header: Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Tách chữ "Bearer " ra
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}

		tokenString := parts[1]

		// 2. Parse và Validate Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Kiểm tra thuật toán ký (thường là HMAC)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			// Trả về Secret Key (lấy từ Config)
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// 3. Lấy thông tin user từ Claims (payload của token)
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Lưu user_id và role vào context để các handler phía sau dùng lại
			c.Set("userID", claims["sub"]) // "sub" thường dùng lưu ID
			c.Set("role", claims["role"])

			c.Next() // Cho phép đi tiếp
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		}
	}
}
