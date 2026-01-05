package product

import (
	"context"

	"go-ecommerce/internal/modules/brand"
	"go-ecommerce/internal/modules/category"
)

// CategoryRepoAdapter adapts category.Repository to CategoryGetter
type CategoryRepoAdapter struct {
	repo category.Repository
}

func NewCategoryRepoAdapter(repo category.Repository) *CategoryRepoAdapter {
	return &CategoryRepoAdapter{repo: repo}
}

func (a *CategoryRepoAdapter) GetByID(ctx context.Context, id uint) (string, error) {
	cat, err := a.repo.GetByID(ctx, id)
	if err != nil {
		return "", err
	}
	return cat.Name, nil
}

// BrandRepoAdapter adapts brand.Repository to BrandGetter
type BrandRepoAdapter struct {
	repo brand.Repository
}

func NewBrandRepoAdapter(repo brand.Repository) *BrandRepoAdapter {
	return &BrandRepoAdapter{repo: repo}
}

func (a *BrandRepoAdapter) GetByID(ctx context.Context, id uint) (string, error) {
	b, err := a.repo.GetByID(ctx, id)
	if err != nil {
		return "", err
	}
	return b.Name, nil
}
