package main

import (
	"log"

	"go-ecommerce/internal/config"
	"go-ecommerce/internal/database"
	"go-ecommerce/internal/modules/user"
	"go-ecommerce/pkg/crypto"
)

func main() {
	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Load config failed: %v", err)
	}

	// Connect Database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Connect database failed: %v", err)
	}

	// Check if admin already exists
	var existingUser user.User
	result := db.Where("email = ?", "admin123@gmail.com").First(&existingUser)
	if result.RowsAffected > 0 {
		log.Println("Admin account already exists!")
		return
	}

	// Create Admin User
	admin := &user.User{
		Email:        "admin123@gmail.com",
		Username:     "Admin",
		Phone:        "0123456789",
		PasswordHash: crypto.HashPassword("Admin!123"),
		Role:         user.RoleAdmin,
		IsActive:     true,
	}

	if err := db.Create(admin).Error; err != nil {
		log.Fatalf("Failed to create admin: %v", err)
	}
	log.Println("Admin account created successfully!")

	// Create Address for Admin
	address := &user.Address{
		UserID:         admin.ID,
		RecipientName:  "Admin",
		RecipientPhone: "0123456789",
		Street:         "123 abc",
		City:           "xyz",
		District:       "mni",
		Ward:           "",
		IsDefault:      true,
	}

	if err := db.Create(address).Error; err != nil {
		log.Fatalf("Failed to create address: %v", err)
	}
	log.Println("Admin address created successfully!")

	log.Println("Seed completed!")
}
