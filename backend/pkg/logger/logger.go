package logger

import (
	"go.uber.org/zap"
)

// NewLogger khởi tạo zap logger
func NewLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}
