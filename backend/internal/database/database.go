package database

import (
	"fmt"
	"log"
	"time"

	"go-ecommerce/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	// Tạo chuỗi kết nối
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Ho_Chi_Minh",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	// Mở kết nối qua GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Cấu hình Connection Pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	//Số lượng kết nối nhàn rỗi tối đa được giữ lại trong pool
	sqlDB.SetMaxIdleConns(10)

	// Số lượng kết nối tối đa được mở cùng lúc
	sqlDB.SetMaxOpenConns(100)

	//Thời gian tối đa một kết nối được tái sử dụng
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Connected to PostgreSQL successfully")
	return db, nil
}
