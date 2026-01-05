package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User Role Enum
type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleAdmin    UserRole = "admin"
)

// Bảng User
type User struct {
	//Dùng UUID làm khóa chính, tự động generate bởi DB
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Username     string    `gorm:"type:varchar(255);not null" json:"username"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Phone        string    `gorm:"type:varchar(20);index" json:"phone"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"` // Không trả về JSON

	//Role & Status
	Role     UserRole `gorm:"type:varchar(20);default:'customer'" json:"role"`
	IsActive bool     `gorm:"default:true" json:"is_active"`

	//Cloudinary info
	AvatarURL      string `gorm:"type:text" json:"avatar_url"`
	AvatarPublicID string `gorm:"type:varchar(255)" json:"-"`

	//Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	LastLogin *time.Time     `json:"last_login"` // Dùng pointer để cho phép null
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	//Relationships (Cho GORM biết để Preload)
	Addresses []Address `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
}

func (User) TableName() string {
	return "users"
}

// Bảng Address
type Address struct {
	ID     uint      `gorm:"primaryKey" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	// Thông tin người nhận hàng
	RecipientName  string `gorm:"type:varchar(100);not null" json:"recipient_name"`
	RecipientPhone string `gorm:"type:varchar(20);not null" json:"recipient_phone"`

	Street   string `gorm:"type:varchar(255);not null" json:"street"` // Số nhà, tên đường
	City     string `gorm:"type:varchar(100);not null" json:"city"`   // Tỉnh/Thành phố
	District string `gorm:"type:varchar(100)" json:"district"`        // Quận/Huyện
	Ward     string `gorm:"type:varchar(100)" json:"ward"`            // Phường/Xã

	IsDefault bool `gorm:"default:false" json:"is_default"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Address) TableName() string {
	return "addresses"
}

// Bảng Session
type Session struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	RefreshToken string    `gorm:"type:varchar(512);not null" json:"refresh_token"`
	UserAgent    string    `gorm:"type:varchar(255)" json:"user_agent"` // Để user biết đăng nhập từ đâu
	ClientIP     string    `gorm:"type:varchar(50)" json:"client_ip"`
	IsBlocked    bool      `gorm:"default:false" json:"is_blocked"`
	ExpiresAt    time.Time `gorm:"not null;index" json:"expires_at"` // Index field này để query dọn dẹp cho nhanh

	CreatedAt time.Time `json:"created_at"`
}

func (Session) TableName() string {
	return "sessions"
}
