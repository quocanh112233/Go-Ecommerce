package category

import (
	"context"

	"gorm.io/gorm"
)

// Repository interface
type Repository interface {
	Create(ctx context.Context, category *Category) error
	GetByID(ctx context.Context, id uint) (*Category, error)
	GetAll(ctx context.Context) ([]Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new category repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, category *Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*Category, error) {
	var category Category
	err := r.db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *repository) GetAll(ctx context.Context) ([]Category, error) {
	var categories []Category
	err := r.db.WithContext(ctx).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *repository) Update(ctx context.Context, category *Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Category{}, id).Error
}
