package app

import (
	"go-ecommerce/internal/config"
	"go-ecommerce/internal/middleware"
	"go-ecommerce/internal/modules/brand"
	"go-ecommerce/internal/modules/category"
	"go-ecommerce/internal/modules/product"
	"go-ecommerce/internal/modules/user"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter(cfg *config.Config, logger *zap.Logger, userHandler *user.Handler, categoryHandler *category.Handler, brandHandler *brand.Handler, productHandler *product.Handler) *gin.Engine {
	r := gin.Default()

	// 1. Global Middlewares
	r.Use(middleware.CorsConfig())
	r.Use(middleware.Logger(logger)) // Gắn Logger Zap vào
	r.Use(gin.Recovery())            // Chống crash server khi có panic

	api := r.Group("/api/v1")
	{
		// PUBLIC ROUTES (Ai cũng vào được)
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/refresh-token", userHandler.RefreshToken)
			auth.POST("/logout", userHandler.Logout)
		}

		// PRIVATE ROUTES (Phải đăng nhập)
		// Tạo một nhóm route có bảo vệ
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// Lấy thông tin cá nhân
			protected.GET("/me", userHandler.GetProfile)

			// ADMIN ROUTES (Phải đăng nhập + Là Admin)
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				// Category CRUD
				admin.POST("/categories", categoryHandler.Create)
				admin.GET("/categories", categoryHandler.GetAll)
				admin.GET("/categories/:id", categoryHandler.GetByID)
				admin.PUT("/categories/:id", categoryHandler.Update)
				admin.DELETE("/categories/:id", categoryHandler.Delete)

				// Brand CRUD
				admin.POST("/brands", brandHandler.Create)
				admin.GET("/brands", brandHandler.GetAll)
				admin.GET("/brands/:id", brandHandler.GetByID)
				admin.PUT("/brands/:id", brandHandler.Update)
				admin.DELETE("/brands/:id", brandHandler.Delete)

				// Product CRUD
				admin.POST("/products", productHandler.Create)
				admin.GET("/products", productHandler.GetAll)
				admin.GET("/products/:id", productHandler.GetByID)
				admin.PUT("/products/:id", productHandler.Update)
				admin.DELETE("/products/:id", productHandler.Delete)
			}
		}
	}

	return r
}
