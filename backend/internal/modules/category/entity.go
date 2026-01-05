package category

import "time"

// Category entity
type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:text;not null" json:"name"`
	Slug        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Category) TableName() string {
	return "categories"
}
