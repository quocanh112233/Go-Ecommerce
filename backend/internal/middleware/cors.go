package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CorsConfig cấu hình CORS để Frontend có thể gọi API
func CorsConfig() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 1. Cho phép các domain cụ thể gọi API
		// Khi dev thường là localhost:3000 (React/Next) hoặc 5173 (Vite)
		// Khi deploy, bạn thay bằng domain thật: "https://your-shop.com"
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://localhost:8080", // Cho phép chính nó (nếu cần test tool)
		},

		// 2. Các Method được phép sử dụng
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},

		// 3. Các Header được phép gửi lên
		// "Authorization" là bắt buộc để gửi JWT Token
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},

		// 4. Các Header mà Frontend được phép đọc từ Response
		ExposeHeaders: []string{"Content-Length"},

		// 5. Cho phép gửi Cookie/Credentials (nếu sau này bạn dùng Cookie)
		AllowCredentials: true,

		// 6. Cache lại preflight request (OPTIONS) trong 12 giờ để đỡ tốn tài nguyên server
		MaxAge: 12 * time.Hour,
	})
}
