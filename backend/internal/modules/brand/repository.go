package brand

import (
	"context"

	"gorm.io/gorm"
)

// Repository interface
type Repository interface {
	Create(ctx context.Context, brand *Brand) error
	GetByID(ctx context.Context, id uint) (*Brand, error)
	GetAll(ctx context.Context) ([]Brand, error)
	Update(ctx context.Context, brand *Brand) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new brand repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, brand *Brand) error {
	return r.db.WithContext(ctx).Create(brand).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*Brand, error) {
	var brand Brand
	err := r.db.WithContext(ctx).First(&brand, id).Error
	if err != nil {
		return nil, err
	}
	return &brand, nil
}

func (r *repository) GetAll(ctx context.Context) ([]Brand, error) {
	var brands []Brand
	err := r.db.WithContext(ctx).Find(&brands).Error
	if err != nil {
		return nil, err
	}
	return brands, nil
}

func (r *repository) Update(ctx context.Context, brand *Brand) error {
	return r.db.WithContext(ctx).Save(brand).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Brand{}, id).Error
}
