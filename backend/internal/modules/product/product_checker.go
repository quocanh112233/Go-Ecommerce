package product

import (
	"context"

	"gorm.io/gorm"
)

// ProductChecker checks if products exist for a category or brand
type ProductChecker struct {
	db *gorm.DB
}

func NewProductChecker(db *gorm.DB) *ProductChecker {
	return &ProductChecker{db: db}
}

// HasProductsWithCategory checks if any products exist for a category
func (pc *ProductChecker) HasProductsWithCategory(ctx context.Context, categoryID uint) (bool, error) {
	var count int64
	err := pc.db.WithContext(ctx).Model(&Product{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count > 0, err
}

// HasProductsWithBrand checks if any products exist for a brand
func (pc *ProductChecker) HasProductsWithBrand(ctx context.Context, brandID uint) (bool, error) {
	var count int64
	err := pc.db.WithContext(ctx).Model(&Product{}).Where("brand_id = ?", brandID).Count(&count).Error
	return count > 0, err
}
