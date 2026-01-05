package main

import (
	"log"

	"go-ecommerce/internal/app"
	"go-ecommerce/internal/config"
	"go-ecommerce/internal/database"
	"go-ecommerce/internal/modules/brand"
	"go-ecommerce/internal/modules/category"
	"go-ecommerce/internal/modules/user"
	"go-ecommerce/pkg/cloudinary"
	"go-ecommerce/pkg/logger"
)

func main() {
	//Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Load config failed: %v", err)
	}

	//Initialize Logger
	zapLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Logger init failed: %v", err)
	}
	defer zapLogger.Sync()

	//Connect Database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Connect database failed: %v", err)
	}

	//Kiểm tra kết nối
	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;")

	err = db.AutoMigrate(
		&user.User{},
		&user.Address{},
		&user.Session{},
		&category.Category{},
		&brand.Brand{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Database migration completed!")

	// Initialize Cloudinary
	cloudinaryClient, err := cloudinary.NewClient(&cfg.Cloudinary)
	if err != nil {
		log.Fatalf("Cloudinary init failed: %v", err)
	}

	// Initialize User Module
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, cfg)
	userHandler := user.NewHandler(userService)

	// Initialize Category Module
	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService)

	// Initialize Brand Module
	brandRepo := brand.NewRepository(db)
	brandService := brand.NewService(brandRepo, cloudinaryClient)
	brandHandler := brand.NewHandler(brandService)

	// Setup Router
	router := app.SetupRouter(cfg, zapLogger, userHandler, categoryHandler, brandHandler)

	// Start Server
	log.Println("Server is starting on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
