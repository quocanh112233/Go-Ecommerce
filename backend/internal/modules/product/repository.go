package product

import (
	"context"

	"gorm.io/gorm"
)

// Repository interface
type Repository interface {
	// Product operations
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id uint) (*Product, error)
	GetAll(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uint) error

	// Variant operations
	CreateVariant(ctx context.Context, variant *ProductVariant) error
	UpdateVariantSKU(ctx context.Context, variantID uint, sku string) error
	GetVariantsByProductID(ctx context.Context, productID uint) ([]ProductVariant, error)

	// Image operations
	CreateImage(ctx context.Context, image *ProductImage) error
	GetImagesByProductID(ctx context.Context, productID uint) ([]ProductImage, error)
	DeleteImage(ctx context.Context, imageID uint) error

	// Transaction support
	WithTransaction(fn func(*gorm.DB) error) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new product repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// WithTransaction executes a function within a transaction
func (r *repository) WithTransaction(fn func(*gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *repository) Create(ctx context.Context, product *Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*Product, error) {
	var product Product
	err := r.db.WithContext(ctx).
		Preload("Variants").
		Preload("Images").
		First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *repository) GetAll(ctx context.Context) ([]Product, error) {
	var products []Product
	err := r.db.WithContext(ctx).
		Preload("Variants").
		Preload("Images").
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *repository) Update(ctx context.Context, product *Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Product{}, id).Error
}

// Variant operations
func (r *repository) CreateVariant(ctx context.Context, variant *ProductVariant) error {
	return r.db.WithContext(ctx).Create(variant).Error
}

func (r *repository) UpdateVariantSKU(ctx context.Context, variantID uint, sku string) error {
	return r.db.WithContext(ctx).Model(&ProductVariant{}).
		Where("id = ?", variantID).
		Update("sku", sku).Error
}

func (r *repository) GetVariantsByProductID(ctx context.Context, productID uint) ([]ProductVariant, error) {
	var variants []ProductVariant
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).Find(&variants).Error
	return variants, err
}

// Image operations
func (r *repository) CreateImage(ctx context.Context, image *ProductImage) error {
	return r.db.WithContext(ctx).Create(image).Error
}

func (r *repository) GetImagesByProductID(ctx context.Context, productID uint) ([]ProductImage, error) {
	var images []ProductImage
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("display_order ASC").
		Find(&images).Error
	return images, err
}

func (r *repository) DeleteImage(ctx context.Context, imageID uint) error {
	return r.db.WithContext(ctx).Delete(&ProductImage{}, imageID).Error
}
