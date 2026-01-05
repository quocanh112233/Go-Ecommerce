package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database   DatabaseConfig
	JWT        JWTConfig
	Cloudinary CloudinaryConfig
}
type JWTConfig struct {
	Secret            string
	AccessExpiration  time.Duration
	RefreshExpiration time.Duration
}
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
}

// LoadConfig đọc file .env và map vào struct
func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config

	cfg.Database.Host = viper.GetString("DB_HOST")
	cfg.Database.Port = viper.GetString("DB_PORT")
	cfg.Database.User = viper.GetString("DB_USER")
	cfg.Database.Password = viper.GetString("DB_PASSWORD")
	cfg.Database.Name = viper.GetString("DB_NAME")
	cfg.Database.SSLMode = viper.GetString("DB_SSLMODE")

	cfg.JWT.Secret = viper.GetString("JWT_SECRET")
	cfg.JWT.AccessExpiration = viper.GetDuration("JWT_ACCESS_EXPIRATION")
	cfg.JWT.RefreshExpiration = viper.GetDuration("JWT_REFRESH_EXPIRATION")

	// Set defaults if not specified
	if cfg.JWT.AccessExpiration == 0 {
		cfg.JWT.AccessExpiration = 15 * time.Minute
	}
	if cfg.JWT.RefreshExpiration == 0 {
		cfg.JWT.RefreshExpiration = 7 * 24 * time.Hour
	}

	// Cloudinary
	cfg.Cloudinary.CloudName = viper.GetString("CLOUDINARY_CLOUD_NAME")
	cfg.Cloudinary.APIKey = viper.GetString("CLOUDINARY_API_KEY")
	cfg.Cloudinary.APISecret = viper.GetString("CLOUDINARY_API_SECRET")

	return &cfg, nil
}
