package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole kiểm tra xem user có quyền truy cập không
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy role đã lưu từ AuthMiddleware
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		roleStr := userRole.(string)

		// Kiểm tra role của user có nằm trong danh sách cho phép không
		isValid := false
		for _, role := range allowedRoles {
			if roleStr == role {
				isValid = true
				break
			}
		}

		if !isValid {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}

		c.Next()
	}
}
