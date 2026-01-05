package brand

import "time"

// Brand entity
type Brand struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"type:text;not null" json:"name"`
	Slug         string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Description  string    `gorm:"type:text" json:"description"`
	LogoURL      string    `gorm:"type:text" json:"logo_url"`
	LogoPublicID string    `gorm:"type:varchar(255)" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Brand) TableName() string {
	return "brands"
}
