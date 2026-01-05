package main

import (
	"log"

	"quocanh.com/go-ecommerce/internal/config"
	"quocanh.com/go-ecommerce/internal/database"
)

func main() {
	//Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Load config failed: %v", err)
	}

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

	// 3. (Tùy chọn) Auto Migration cho bảng User (để test thử)
	// Sau này nên dùng file migration riêng, nhưng lúc dev thì có thể dùng cái này cho nhanh
	// db.AutoMigrate(&entity.User{})

	log.Println("Server is starting...")
}
