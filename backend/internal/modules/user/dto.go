package user

import (
	"github.com/google/uuid"
	"time"
)

// RegisterRequest: Dữ liệu client gửi lên khi đăng ký
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=32"` // Mật khẩu từ 6-32 ký tự
	Username string `json:"full_name" binding:"required,min=2,max=100"`
	Phone    string `json:"phone" binding:"omitempty,e164"` // e164 là chuẩn số đt quốc tế
}

// UserResponse: Dữ liệu trả về cho client
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginRequest: Input đăng nhập
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse: Trả về client gồm 2 token
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"` // Giây
	User         UserResponse `json:"user"`
}
