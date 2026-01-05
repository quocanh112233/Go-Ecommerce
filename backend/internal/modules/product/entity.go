package product

import "time"

// Product entity
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:text;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Slug        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	CategoryID  uint      `gorm:"not null;index" json:"category_id"`
	BrandID     uint      `gorm:"not null;index" json:"brand_id"`
	TotalStock  int       `gorm:"default:0" json:"total_stock"`
	RatingAvg   float64   `gorm:"default:0" json:"rating_avg"`
	ReviewCount int       `gorm:"default:0" json:"review_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Variants []ProductVariant `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"variants,omitempty"`
	Images   []ProductImage   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"images,omitempty"`
}

func (Product) TableName() string {
	return "products"
}

// ProductVariant entity
type ProductVariant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"not null;index" json:"product_id"`
	Price     float64   `gorm:"not null;check:price >= 0" json:"price"`
	Stock     int       `gorm:"not null;check:stock >= 0" json:"stock"`
	Size      string    `gorm:"type:varchar(50);not null" json:"size"`
	SKU       string    `gorm:"type:varchar(100);uniqueIndex" json:"sku"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ProductVariant) TableName() string {
	return "product_variants"
}

// ProductImage entity
type ProductImage struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ProductID     uint      `gorm:"not null;index" json:"product_id"`
	ImageURL      string    `gorm:"type:text;not null" json:"image_url"`
	ImagePublicID string    `gorm:"type:varchar(255)" json:"-"`
	DisplayOrder  int       `gorm:"not null;check:display_order >= 1 AND display_order <= 5" json:"display_order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (ProductImage) TableName() string {
	return "product_images"
}
